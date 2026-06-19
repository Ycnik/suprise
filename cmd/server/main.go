package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Ycnik/suprise/internal/auth"
	"github.com/Ycnik/suprise/internal/config"
	"github.com/Ycnik/suprise/internal/database"
	"github.com/Ycnik/suprise/internal/handler"
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

	var keycloak *auth.KeycloakMiddleware
	if cfg.AuthEnabled {
		if cfg.OIDCIssuerURL == "" || cfg.OIDCClientID == "" {
			log.Fatal("AUTH_ENABLED=true benoetigt OIDC_ISSUER_URL und OIDC_CLIENT_ID")
		}

		keycloak, err = auth.NewKeycloakMiddleware(context.Background(), cfg.OIDCIssuerURL, cfg.OIDCClientID)
		if err != nil {
			log.Fatalf("keycloak middleware konnte nicht erstellt werden: %v", err)
		}
	}

	var resetHandler http.HandlerFunc
	if cfg.DBResetEnabled {
		log.Print("dev reset endpoint ist aktiv: POST /dev/reset-db")
		resetHandler = handler.NewDevHandler(db).ResetDatabase
	}

	router := httpapi.NewRouterWithDevReset(repo, keycloak, resetHandler)

	log.Printf("server startet auf %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		log.Fatalf("server beendet: %v", err)
	}
}
