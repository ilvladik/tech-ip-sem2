package httpapi

import (
	"encoding/json"
	"log"
	"net/http"

	"example.com/pz7-docker/tasks/internal/authclient"
	"example.com/pz7-docker/tasks/internal/task"
)

type Handler struct {
	auth *authclient.Client
}

func New(auth *authclient.Client) *Handler {
	return &Handler{auth: auth}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /v1/tasks", h.listTasks)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "tasks",
	})
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = "unknown"
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, `{"error":"missing authorization"}`, http.StatusUnauthorized)
		return
	}
	if err := h.auth.Validate(authHeader); err != nil {
		log.Printf("request_id=%s auth_error=%v", requestID, err)
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	log.Printf("request_id=%s list_tasks ok", requestID)
	writeJSON(w, http.StatusOK, map[string]any{
		"request_id": requestID,
		"tasks":      task.List(),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
