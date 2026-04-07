package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"tech-ip-sem2/services/auth/internal/usecases"
	"tech-ip-sem2/shared/httpx"
)

type HTTPAuthenticationHandler struct {
	authenticationUsecase *usecases.AuthenticationUsecase
}

func NewHTTPAuthenticationHandler(usecase *usecases.AuthenticationUsecase) *HTTPAuthenticationHandler {
	return &HTTPAuthenticationHandler{
		authenticationUsecase: usecase,
	}
}

func (h *HTTPAuthenticationHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	out, err := h.authenticationUsecase.Login(r.Context(), usecases.LoginInput{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		httpx.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	httpx.RespondJson(w, http.StatusOK, LoginResponse{
		AccessToken: out.AccessToken,
		TokenType:   out.TokenType,
	})
}

func (h *HTTPAuthenticationHandler) handleVerify(w http.ResponseWriter, r *http.Request) {
	authenticationHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authenticationHeader, "Bearer ") {
		httpx.RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	token := strings.TrimPrefix(authenticationHeader, "Bearer ")
	out, err := h.authenticationUsecase.Verify(r.Context(), usecases.VerifyInput{
		Token: token,
	})
	if err != nil {
		msg := err.Error()
		httpx.RespondJson(w, http.StatusUnauthorized, VerifyResponse{
			Error: &msg,
			Valid: false,
		})
		return
	}
	httpx.RespondJson(w, http.StatusOK, VerifyResponse{
		Subject: &out.Subject,
		Valid:   true,
	})
}
