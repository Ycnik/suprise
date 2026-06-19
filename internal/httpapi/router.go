package httpapi

import (
	"net/http"

	"github.com/Ycnik/suprise/internal/handler"
	"github.com/Ycnik/suprise/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type TokenMiddleware interface {
	RequireToken(http.Handler) http.Handler
}

func NewRouter(repo repository.SoldatRepository, keycloak TokenMiddleware) http.Handler {
	return NewRouterWithDevReset(repo, keycloak, nil)
}

func NewRouterWithDevReset(repo repository.SoldatRepository, keycloak TokenMiddleware, resetHandler http.HandlerFunc) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	soldaten := handler.NewSoldatHandler(repo)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Get("/rest", soldaten.List)
	r.Get("/rest/{id}", soldaten.FindByID)
	if keycloak == nil {
		r.Post("/rest", soldaten.Create)
	} else {
		r.With(keycloak.RequireToken).Post("/rest", soldaten.Create)
	}

	if resetHandler != nil {
		r.Post("/dev/reset-db", resetHandler)
	}

	return r
}
