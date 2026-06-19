package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ycnik/suprise/internal/handler"
	"github.com/Ycnik/suprise/internal/model"
	"github.com/Ycnik/suprise/internal/repository"
	"github.com/go-chi/chi/v5"
)

func BenchmarkHealth(b *testing.B) {
	router := newBenchmarkRouter(repository.NewMemorySoldatRepository())

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
	repo := repository.NewMemorySoldatRepository()
	seedSoldaten(b, repo, 100)
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
	repo := repository.NewMemorySoldatRepository()
	seedSoldaten(b, repo, 1)
	router := newBenchmarkRouter(repo)

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodGet, "/rest/1000", nil)
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
	router := newBenchmarkRouter(repository.NewMemorySoldatRepository())

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodPost, "/rest", strings.NewReader(validSoldatJSON))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusCreated {
				b.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rec.Code, rec.Body.String())
			}
		}
	})
}

func seedSoldaten(b *testing.B, repo *repository.MemorySoldatRepository, count int) {
	b.Helper()

	for i := 0; i < count; i++ {
		soldat := &model.Soldat{
			Vorname:  "Eren",
			Nachname: "Jaeger",
			Username: "eren",
		}
		if err := repo.Create(context.Background(), soldat); err != nil {
			b.Fatalf("seed soldat: %v", err)
		}
	}
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
