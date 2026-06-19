package handler

import (
	"net/http"

	"github.com/Ycnik/suprise/internal/database"
	"gorm.io/gorm"
)

type DevHandler struct {
	db *gorm.DB
}

func NewDevHandler(db *gorm.DB) *DevHandler {
	return &DevHandler{db: db}
}

func (h *DevHandler) ResetDatabase(w http.ResponseWriter, r *http.Request) {
	if err := database.ResetAndPopulate(h.db); err != nil {
		writeError(w, http.StatusInternalServerError, "datenbank konnte nicht zurueckgesetzt werden")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "datenbank wurde zurueckgesetzt und befuellt",
	})
}
