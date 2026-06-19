package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

type KeycloakMiddleware struct {
	verifier *oidc.IDTokenVerifier
}

func NewKeycloakMiddleware(ctx context.Context, issuerURL string, clientID string) (*KeycloakMiddleware, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, err
	}

	return &KeycloakMiddleware{
		verifier: provider.Verifier(&oidc.Config{ClientID: clientID}),
	}, nil
}

func (m *KeycloakMiddleware) RequireToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r.Header.Get("Authorization"))
		if token == "" {
			writeAuthError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}

		if _, err := m.verifier.Verify(r.Context(), token); err != nil {
			writeAuthError(w, http.StatusUnauthorized, "invalid bearer token")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
