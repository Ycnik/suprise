package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ycnik/suprise/internal/model"
	"github.com/Ycnik/suprise/internal/repository"
)

type denyingTokenMiddleware struct{}

func (denyingTokenMiddleware) RequireToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
}

func TestHealth(t *testing.T) {
	router := NewRouter(repository.NewMemorySoldatRepository(), nil)

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
	router := NewRouter(repository.NewMemorySoldatRepository(), nil)

	req := httptest.NewRequest(http.MethodPost, "/rest", strings.NewReader(validSoldatJSON))
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
	if soldat.ID != 1000 {
		t.Fatalf("expected created soldat id 1000, got %d", soldat.ID)
	}
	if soldat.Geburtsdatum == nil {
		t.Fatal("expected geburtsdatum to be parsed")
	}
}

func TestCreateSoldatValidation(t *testing.T) {
	router := NewRouter(repository.NewMemorySoldatRepository(), nil)

	req := httptest.NewRequest(http.MethodPost, "/rest", strings.NewReader(`{"vorname":"E"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var response struct {
		Error   string   `json:"error"`
		Details []string `json:"details"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode validation response: %v", err)
	}
	if response.Error != "validierung fehlgeschlagen" {
		t.Fatalf("unexpected validation error: %q", response.Error)
	}
	if len(response.Details) == 0 {
		t.Fatal("expected validation details")
	}
}

func TestCreateSoldatRejectsInvalidEnumValue(t *testing.T) {
	router := NewRouter(repository.NewMemorySoldatRepository(), nil)
	body := strings.Replace(validSoldatJSON, `"rang": "SOLDAT"`, `"rang": "ELITE_SOLDAT"`, 1)

	req := httptest.NewRequest(http.MethodPost, "/rest", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestDevResetRouteCanBeRegistered(t *testing.T) {
	router := NewRouterWithDevReset(repository.NewMemorySoldatRepository(), nil, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/dev/reset-db", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestDevResetRouteIsDisabledByDefault(t *testing.T) {
	router := NewRouter(repository.NewMemorySoldatRepository(), nil)

	req := httptest.NewRequest(http.MethodPost, "/dev/reset-db", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestCreateSoldatRequiresTokenWhenAuthIsEnabled(t *testing.T) {
	router := NewRouter(repository.NewMemorySoldatRepository(), denyingTokenMiddleware{})

	req := httptest.NewRequest(http.MethodPost, "/rest", strings.NewReader(validSoldatJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestFindCreatedSoldatWithETag(t *testing.T) {
	router := NewRouter(repository.NewMemorySoldatRepository(), nil)

	createReq := httptest.NewRequest(http.MethodPost, "/rest", strings.NewReader(validSoldatJSON))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d: %s", http.StatusCreated, createRec.Code, createRec.Body.String())
	}

	findReq := httptest.NewRequest(http.MethodGet, "/rest/1000", nil)
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
	if soldat.ID != 1000 || soldat.Username != "eren" {
		t.Fatalf("unexpected soldat response: %+v", soldat)
	}
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
