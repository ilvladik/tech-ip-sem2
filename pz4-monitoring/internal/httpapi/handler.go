package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"example.com/pz4-monitoring/internal/metrics"
	"example.com/pz4-monitoring/internal/student"
)

type Handler struct {
	repo *student.Repo
}

func NewHandler(repo *student.Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid student id", http.StatusBadRequest)
		return
	}

	idLabel := strconv.FormatInt(id, 10)
	start := time.Now()
	defer func() {
		metrics.StudentHandlerDuration.WithLabelValues(idLabel).Observe(time.Since(start).Seconds())
	}()

	st, err := h.repo.GetByID(id)
	metrics.StudentRequestsTotal.WithLabelValues(idLabel).Inc()

	if err != nil {
		if errors.Is(err, student.ErrStudentNotFound) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": "студент с таким id не найден в системе",
				"id":    id,
			})
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(st)
}
