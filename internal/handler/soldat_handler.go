package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Ycnik/suprise/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type SoldatHandler struct {
	repo     repository.SoldatRepository
	validate *validator.Validate
}

func NewSoldatHandler(repo repository.SoldatRepository) *SoldatHandler {
	return &SoldatHandler{
		repo:     repo,
		validate: validator.New(),
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *SoldatHandler) List(w http.ResponseWriter, r *http.Request) {
	soldaten, err := h.repo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "soldaten konnten nicht geladen werden")
		return
	}

	writeJSON(w, http.StatusOK, soldaten)
}

func (h *SoldatHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "ungueltige soldat id")
		return
	}

	soldat, err := h.repo.FindByID(r.Context(), id)
	if errors.Is(err, repository.ErrSoldatNotFound) {
		writeError(w, http.StatusNotFound, "soldat nicht gefunden")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "soldat konnte nicht geladen werden")
		return
	}

	w.Header().Set("ETag", fmt.Sprintf(`"%d"`, soldat.Version))
	writeJSON(w, http.StatusOK, soldat)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
