package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Ycnik/suprise/internal/model"
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

type createSoldatRequest struct {
	Vorname      string     `json:"vorname" validate:"required,min=2"`
	Nachname     string     `json:"nachname" validate:"required,min=2"`
	Geburtsdatum *time.Time `json:"geburtsdatum,omitempty"`
	Geschlecht   *string    `json:"geschlecht,omitempty"`
	Rang         *string    `json:"rang,omitempty"`
	Username     string     `json:"username" validate:"required,min=3"`
	Ausruestung  *struct {
		Waffe        string `json:"waffe" validate:"required"`
		Seriennummer string `json:"seriennummer" validate:"required"`
	} `json:"ausruestung,omitempty" validate:"omitempty"`
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

func (h *SoldatHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createSoldatRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ungueltiges json")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, "validierung fehlgeschlagen")
		return
	}

	soldat := model.Soldat{
		Vorname:      req.Vorname,
		Nachname:     req.Nachname,
		Geburtsdatum: req.Geburtsdatum,
		Geschlecht:   req.Geschlecht,
		Rang:         req.Rang,
		Username:     req.Username,
	}
	if req.Ausruestung != nil {
		soldat.Ausruestung = &model.Ausruestung{
			Waffe:        req.Ausruestung.Waffe,
			Seriennummer: req.Ausruestung.Seriennummer,
		}
	}

	if err := h.repo.Create(r.Context(), &soldat); err != nil {
		writeError(w, http.StatusInternalServerError, "soldat konnte nicht angelegt werden")
		return
	}

	writeJSON(w, http.StatusCreated, soldat)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
