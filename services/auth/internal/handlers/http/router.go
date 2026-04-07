package http

import (
	"net/http"
	"tech-ip-sem2/services/auth/internal/usecases"
	"tech-ip-sem2/shared/middlewares"

	"go.uber.org/zap"
)

func RegisterRoutes(
	usecase *usecases.AuthenticationUsecase,
	log *zap.Logger,
) http.Handler {
	mux := http.NewServeMux()
	h := NewHTTPAuthenticationHandler(usecase)

	mux.HandleFunc("POST /v1/auth/login", h.handleLogin)
	mux.HandleFunc("GET /v1/auth/verify", h.handleVerify)

	var handler http.Handler = mux
	handler = middlewares.AccessLog(log)(handler)
	handler = middlewares.RequestId(handler)

	return handler
}
