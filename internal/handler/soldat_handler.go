package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	Error   string   `json:"error"`
	Details []string `json:"details,omitempty"`
}

type createSoldatRequest struct {
	Vorname      string  `json:"vorname" validate:"required,min=2"`
	Nachname     string  `json:"nachname" validate:"required,min=2"`
	Geburtsdatum string  `json:"geburtsdatum,omitempty"`
	Geschlecht   *string `json:"geschlecht,omitempty" validate:"omitempty,oneof=MAENNLICH WEIBLICH"`
	Rang         *string `json:"rang,omitempty" validate:"omitempty,oneof=REKRUT SOLDAT ELITE-SOLDAT CAPTAIN KOMMANDANT"`
	Username     string  `json:"username" validate:"required,min=3"`
	Ausruestung  *struct {
		Waffe        string `json:"waffe" validate:"required,oneof=ODM_GEAR Schrotflinte Klinge"`
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
		writeValidationError(w, err)
		return
	}

	var geburtsdatum *time.Time
	if req.Geburtsdatum != "" {
		parsed, err := time.Parse("2006-01-02", req.Geburtsdatum)
		if err != nil {
			writeError(w, http.StatusBadRequest, "ungueltiges geburtsdatum")
			return
		}
		geburtsdatum = &parsed
	}

	soldat := model.Soldat{
		Vorname:      req.Vorname,
		Nachname:     req.Nachname,
		Geburtsdatum: geburtsdatum,
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

func writeValidationError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusBadRequest, errorResponse{
		Error:   "validierung fehlgeschlagen",
		Details: validationErrorDetails(err),
	})
}

func validationErrorDetails(err error) []string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []string{"request body ist ungueltig"}
	}

	details := make([]string, 0, len(validationErrors))
	for _, fieldError := range validationErrors {
		field := jsonFieldName(fieldError)
		switch fieldError.Tag() {
		case "required":
			details = append(details, field+" ist erforderlich")
		case "min":
			details = append(details, field+" muss mindestens "+fieldError.Param()+" Zeichen lang sein")
		case "oneof":
			details = append(details, field+" hat keinen erlaubten Wert")
		default:
			details = append(details, field+" ist ungueltig")
		}
	}
	return details
}

func jsonFieldName(fieldError validator.FieldError) string {
	name := fieldError.Field()
	switch name {
	case "Vorname":
		return "vorname"
	case "Nachname":
		return "nachname"
	case "Geburtsdatum":
		return "geburtsdatum"
	case "Geschlecht":
		return "geschlecht"
	case "Rang":
		return "rang"
	case "Username":
		return "username"
	case "Waffe":
		return nestedJSONField(fieldError, "waffe")
	case "Seriennummer":
		return nestedJSONField(fieldError, "seriennummer")
	default:
		return strings.ToLower(name)
	}
}

func nestedJSONField(fieldError validator.FieldError, field string) string {
	if strings.Contains(fieldError.StructNamespace(), ".Ausruestung.") {
		return "ausruestung." + field
	}
	return field
}
