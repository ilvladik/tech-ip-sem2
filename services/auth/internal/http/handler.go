package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"tech-ip-sem2/services/auth/internal/service"
)

type AuthenticationHandler struct {
	authenticationService service.AuthenticationService
}

func NewAuthHandler(auth *service.AuthenticationService) *AuthenticationHandler {
	return &AuthenticationHandler{
		authenticationService: *auth,
	}
}

func (h *AuthenticationHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req service.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	loginResponse, err := h.authenticationService.Login(req)
	if err != nil {
		h.handleError(w, err)
		return
	}
	respondJson(w, http.StatusOK, loginResponse)
}

func (h *AuthenticationHandler) handleVerify(w http.ResponseWriter, r *http.Request) {
	authenticationHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authenticationHeader, "Bearer ") {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	token := strings.TrimPrefix(authenticationHeader, "Bearer ")
	verifyResponse, err := h.authenticationService.Verify(token)
	if err != nil {
		h.handleError(w, err)
		return
	}
	if verifyResponse.Valid {
		respondJson(w, http.StatusOK, verifyResponse)
	} else {
		respondJson(w, http.StatusUnauthorized, verifyResponse)
	}
}

func (h *AuthenticationHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrorInvalidCredentials:
		respondWithError(w, http.StatusBadRequest, err.Error())
	default:
		respondWithError(w, http.StatusBadRequest, "Internal server error")
	}
}

func respondJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondJson(w, status, map[string]string{"error": message})
}
