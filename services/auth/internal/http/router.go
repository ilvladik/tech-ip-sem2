package http

import (
	"net/http"
	"tech-ip-sem2/services/auth/internal/service"
	"tech-ip-sem2/shared/middleware"
)

func RegisterRoutes(auth *service.AuthenticationService) http.Handler {
	mux := http.NewServeMux()
	h := NewAuthHandler(auth)

	mux.HandleFunc("POST /v1/auth/login", h.handleLogin)
	mux.HandleFunc("GET /v1/auth/verify", h.handleVerify)

	var handler http.Handler = mux
	handler = middleware.Logging(handler)
	handler = middleware.RequestId(handler)
	return handler
}
