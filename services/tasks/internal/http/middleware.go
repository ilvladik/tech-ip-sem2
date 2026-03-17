package http

import (
	"context"
	"net/http"
	"strings"

	"tech-ip-sem2/services/tasks/internal/client"
	"tech-ip-sem2/shared/middleware"
)

const (
	userSubjectKey contextKey = "subject"
)

type contextKey string

type AuthMiddleware struct {
	authClient *client.GRPCClient
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewAuthMiddleware(authClient *client.GRPCClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization format. Use Bearer token")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			respondWithError(w, http.StatusUnauthorized, "Empty token")
			return
		}

		requestID := middleware.GetRequestID(r.Context())

		ctx := context.WithValue(r.Context(), "request_id", requestID)

		valid, subject, err := m.authClient.Verify(ctx, token)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Authentication service unavailable")
			return
		}

		if !valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		ctx = context.WithValue(r.Context(), userSubjectKey, subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserSubject(ctx context.Context) (string, bool) {
	subject, ok := ctx.Value(userSubjectKey).(string)
	return subject, ok
}
