package middlewares

import (
	"context"
	"net/http"
	"strings"

	"tech-ip-sem2/shared/httpx"
	"tech-ip-sem2/shared/requestctx"
)

type AuthenticationClient interface {
	Verify(ctx context.Context, token string) (*string, error)
}

type AuthenticationMiddleware struct {
	authClient AuthenticationClient
}

func NewAuthenticationMiddleware(authClient AuthenticationClient) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		authClient: authClient,
	}
}

func (m *AuthenticationMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticationHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authenticationHeader, "Bearer ") {
			httpx.RespondWithError(w, http.StatusBadRequest, "Invalid request")
			return
		}
		token := strings.TrimPrefix(authenticationHeader, "Bearer ")
		subject, err := m.authClient.Verify(r.Context(), token)

		if err != nil {
			httpx.RespondWithError(w, http.StatusServiceUnavailable, "Authentication service unavailable")
			return
		}

		if subject == nil {
			httpx.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		ctx := requestctx.WithSubject(r.Context(), *subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
