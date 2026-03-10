package http

import (
	"context"
	"net/http"
	"strings"

	authclient "tech-ip-sem2/services/tasks/internal/client"
	"tech-ip-sem2/shared/middleware"
)

const (
	userSubjectKey string = "subject"
)

type AuthMiddleware struct {
	authClient *authclient.Client
}

func NewAuthMiddleware(authClient *authclient.Client) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		requestID := middleware.GetRequestID(r.Context())
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		valid, subject, err := m.authClient.VerifyToken(r.Context(), token, requestID)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Authentication service unavailable")
			return
		}

		if !valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), userSubjectKey, subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
