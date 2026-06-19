package httpapi

import (
	"net/http"

	"github.com/Ycnik/suprise/internal/handler"
	"github.com/Ycnik/suprise/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(repo repository.SoldatRepository) http.Handler {
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
	r.Post("/rest", soldaten.Create)
	r.Get("/rest/{id}", soldaten.FindByID)

	return r
}
