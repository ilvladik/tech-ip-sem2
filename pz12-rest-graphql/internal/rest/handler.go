package rest

import (
	"encoding/json"
	"net/http"

	"example.com/pz12-rest-graphql/internal/service"
)

type Handler struct {
	svc *service.TaskService
}

func New(svc *service.TaskService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/tasks", h.listTasks)
	mux.HandleFunc("GET /v1/tasks/{id}", h.getTask)
	mux.HandleFunc("POST /v1/tasks", h.createTask)
	mux.HandleFunc("PATCH /v1/tasks/{id}", h.patchTask)
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	// REST всегда отдаёт полный объект (демонстрация over-fetching).
	writeJSON(w, http.StatusOK, h.svc.List())
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := h.svc.Get(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if t == nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "bad json")
		return
	}
	if body.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	t := h.svc.Create(body.Title, body.Description)
	writeJSON(w, http.StatusCreated, t)
}

func (h *Handler) patchTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var body struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Done        *bool   `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "bad json")
		return
	}
	t, err := h.svc.Update(id, body.Title, body.Description, body.Done)
	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
