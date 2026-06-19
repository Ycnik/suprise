package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

type KeycloakMiddleware struct {
	verifier *oidc.IDTokenVerifier
	clientID string
}

type keycloakClaims struct {
	AuthorizedParty string `json:"azp"`
	ClientID        string `json:"client_id"`
}

func NewKeycloakMiddleware(ctx context.Context, issuerURL string, clientID string) (*KeycloakMiddleware, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, err
	}

	return &KeycloakMiddleware{
		verifier: provider.Verifier(&oidc.Config{SkipClientIDCheck: true}),
		clientID: clientID,
	}, nil
}

func (m *KeycloakMiddleware) RequireToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r.Header.Get("Authorization"))
		if token == "" {
			writeAuthError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}

		if err := m.verifyToken(r.Context(), token); err != nil {
			writeAuthError(w, http.StatusUnauthorized, "invalid bearer token")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *KeycloakMiddleware) verifyToken(ctx context.Context, rawToken string) error {
	token, err := m.verifier.Verify(ctx, rawToken)
	if err != nil {
		return err
	}

	for _, audience := range token.Audience {
		if audience == m.clientID {
			return nil
		}
	}

	var claims keycloakClaims
	if err := token.Claims(&claims); err != nil {
		return err
	}

	if claims.AuthorizedParty == m.clientID || claims.ClientID == m.clientID {
		return nil
	}

	return errors.New("token was not issued for configured client")
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
