package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Ycnik/suprise/internal/config"
	"github.com/Ycnik/suprise/internal/database"
	"github.com/Ycnik/suprise/internal/handler"
	"github.com/Ycnik/suprise/internal/model"
	"github.com/Ycnik/suprise/internal/repository"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func BenchmarkHealth(b *testing.B) {
	repo, _ := newBenchmarkRepository(b)
	router := newBenchmarkRouter(repo)

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				b.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
			}
		}
	})
}

func BenchmarkListSoldaten(b *testing.B) {
	repo, db := newBenchmarkRepository(b)
	seedSoldaten(b, db, repo, "bench-list", 25)
	router := newBenchmarkRouter(repo)

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodGet, "/rest", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				b.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
			}
		}
	})
}

func BenchmarkFindSoldatByID(b *testing.B) {
	repo, db := newBenchmarkRepository(b)
	ids := seedSoldaten(b, db, repo, "bench-find", 1)
	router := newBenchmarkRouter(repo)
	path := "/rest/" + fmt.Sprint(ids[0])

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				b.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
			}
			if rec.Header().Get("ETag") == "" {
				b.Fatal("expected ETag header")
			}
		}
	})
}

func BenchmarkCreateSoldat(b *testing.B) {
	repo, db := newBenchmarkRepository(b)
	cleanupSoldatenByUsernamePrefix(b, db, "bench-create")
	router := newBenchmarkRouter(repo)

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodPost, "/rest", strings.NewReader(benchmarkSoldatJSON("bench-create")))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusCreated {
				b.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rec.Code, rec.Body.String())
			}
		}
	})
}

func newBenchmarkRepository(b *testing.B) (repository.SoldatRepository, *gorm.DB) {
	b.Helper()

	db, err := database.ConnectPostgres(config.FromEnv().DatabaseURL)
	if err != nil {
		b.Fatalf("connect postgres: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := database.Ping(ctx, db); err != nil {
		b.Fatalf("ping postgres: %v", err)
	}

	return repository.NewGormSoldatRepository(db), db
}

func seedSoldaten(b *testing.B, db *gorm.DB, repo repository.SoldatRepository, prefix string, count int) []int {
	b.Helper()

	cleanupSoldatenByUsernamePrefix(b, db, prefix)

	ids := make([]int, 0, count)
	for i := 0; i < count; i++ {
		soldat := &model.Soldat{
			Vorname:  "Eren",
			Nachname: "Jaeger",
			Username: fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), i),
		}
		if err := repo.Create(context.Background(), soldat); err != nil {
			b.Fatalf("seed soldat: %v", err)
		}
		ids = append(ids, soldat.ID)
	}

	return ids
}

func cleanupSoldatenByUsernamePrefix(b *testing.B, db *gorm.DB, prefix string) {
	b.Helper()

	b.Cleanup(func() {
		if err := db.Exec("DELETE FROM soldat.soldat WHERE username LIKE ?", prefix+"-%").Error; err != nil {
			b.Fatalf("cleanup benchmark soldaten: %v", err)
		}
	})
}

var benchmarkSoldatCounter int64

func benchmarkSoldatJSON(prefix string) string {
	n := atomic.AddInt64(&benchmarkSoldatCounter, 1)
	suffix := fmt.Sprintf("%d-%d", time.Now().UnixNano(), n)

	return fmt.Sprintf(`{
	"vorname": "Eren",
	"nachname": "Jaeger",
	"geburtsdatum": "2000-01-01",
	"geschlecht": "MAENNLICH",
	"rang": "SOLDAT",
	"username": %q,
	"ausruestung": {
		"waffe": "ODM_GEAR",
		"seriennummer": %q
	}
}`, prefix+"-"+suffix, "AOT-BENCH-"+suffix)
}

func newBenchmarkRouter(repo repository.SoldatRepository) http.Handler {
	r := chi.NewRouter()
	soldaten := handler.NewSoldatHandler(repo)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	r.Get("/rest", soldaten.List)
	r.Get("/rest/{id}", soldaten.FindByID)
	r.Post("/rest", soldaten.Create)

	return r
}
