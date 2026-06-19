package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Ycnik/suprise/internal/config"
	"github.com/Ycnik/suprise/internal/database"
	"github.com/Ycnik/suprise/internal/model"
	"github.com/Ycnik/suprise/internal/repository"
	"gorm.io/gorm"
)

type denyingTokenMiddleware struct{}

func (denyingTokenMiddleware) RequireToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
}

func TestHealth(t *testing.T) {
	router, _ := newIntegrationRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if strings.TrimSpace(rec.Body.String()) != `{"status":"ok"}` {
		t.Fatalf("unexpected response body: %s", rec.Body.String())
	}
}

func TestCreateSoldat(t *testing.T) {
	router, db := newIntegrationRouter(t)
	body, _ := validIntegrationSoldatJSON(t)

	req := httptest.NewRequest(http.MethodPost, "/rest", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var soldat model.Soldat
	if err := json.NewDecoder(rec.Body).Decode(&soldat); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	cleanupSoldat(t, db, soldat.ID)

	if soldat.ID <= 0 {
		t.Fatalf("expected created soldat id, got %d", soldat.ID)
	}
	if soldat.Geburtsdatum == nil {
		t.Fatal("expected geburtsdatum to be parsed")
	}
}

func TestCreateSoldatValidation(t *testing.T) {
	router, _ := newIntegrationRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/rest", strings.NewReader(`{"vorname":"E"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateSoldatRequiresTokenWhenAuthIsEnabled(t *testing.T) {
	_, db := newIntegrationRouter(t)
	router := NewRouter(repository.NewGormSoldatRepository(db), denyingTokenMiddleware{})

	req := httptest.NewRequest(http.MethodPost, "/rest", strings.NewReader(validSoldatJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestFindCreatedSoldatWithETag(t *testing.T) {
	router, db := newIntegrationRouter(t)
	body, username := validIntegrationSoldatJSON(t)

	createReq := httptest.NewRequest(http.MethodPost, "/rest", strings.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d: %s", http.StatusCreated, createRec.Code, createRec.Body.String())
	}

	var created model.Soldat
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	cleanupSoldat(t, db, created.ID)

	findReq := httptest.NewRequest(http.MethodGet, "/rest/"+strconv.Itoa(created.ID), nil)
	findRec := httptest.NewRecorder()
	router.ServeHTTP(findRec, findReq)

	if findRec.Code != http.StatusOK {
		t.Fatalf("expected find status %d, got %d: %s", http.StatusOK, findRec.Code, findRec.Body.String())
	}
	if findRec.Header().Get("ETag") != `"0"` {
		t.Fatalf("expected ETag %q, got %q", `"0"`, findRec.Header().Get("ETag"))
	}

	var soldat model.Soldat
	if err := json.NewDecoder(findRec.Body).Decode(&soldat); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if soldat.ID != created.ID || soldat.Username != username {
		t.Fatalf("unexpected soldat response: %+v", soldat)
	}
}

func newIntegrationRouter(t *testing.T) (http.Handler, *gorm.DB) {
	t.Helper()

	db, err := database.ConnectPostgres(config.FromEnv().DatabaseURL)
	if err != nil {
		t.Fatalf("connect postgres: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := database.Ping(ctx, db); err != nil {
		t.Fatalf("ping postgres: %v", err)
	}

	return NewRouter(repository.NewGormSoldatRepository(db), nil), db
}

func cleanupSoldat(t *testing.T, db *gorm.DB, id int) {
	t.Helper()

	t.Cleanup(func() {
		if id <= 0 {
			return
		}
		if err := db.Exec("DELETE FROM soldat.soldat WHERE id = ?", id).Error; err != nil {
			t.Fatalf("cleanup soldat %d: %v", id, err)
		}
	})
}

func validIntegrationSoldatJSON(t *testing.T) (string, string) {
	t.Helper()

	suffix := strconv.FormatInt(time.Now().UnixNano(), 10)
	username := "eren-" + suffix
	seriennummer := "AOT-" + suffix

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
}`, username, seriennummer), username
}

const validSoldatJSON = `{
	"vorname": "Eren",
	"nachname": "Jaeger",
	"geburtsdatum": "2000-01-01",
	"geschlecht": "MAENNLICH",
	"rang": "SOLDAT",
	"username": "eren",
	"ausruestung": {
		"waffe": "ODM_GEAR",
		"seriennummer": "AOT-12345ABC"
	}
}`
