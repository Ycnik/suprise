package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Ycnik/suprise/internal/repository"
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

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
