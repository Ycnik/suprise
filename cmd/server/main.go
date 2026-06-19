package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Ycnik/suprise/internal/config"
	"github.com/Ycnik/suprise/internal/database"
	"github.com/Ycnik/suprise/internal/httpapi"
	"github.com/Ycnik/suprise/internal/repository"
)

func main() {
	cfg := config.FromEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.ConnectPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("datenbankverbindung fehlgeschlagen: %v", err)
	}
	if err := database.Ping(ctx, db); err != nil {
		log.Fatalf("datenbank ping fehlgeschlagen: %v", err)
	}

	repo := repository.NewGormSoldatRepository(db)
	router := httpapi.NewRouter(repo)

	log.Printf("server startet auf %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		log.Fatalf("server beendet: %v", err)
	}
}
