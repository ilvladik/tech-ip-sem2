package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"example.com/pz3-logging/internal/student"
	"go.uber.org/zap"
)

type Handler struct {
	repo *student.Repo
	log  *zap.Logger
}

func NewHandler(repo *student.Repo, log *zap.Logger) *Handler {
	return &Handler{repo: repo, log: log}
}

func (h *Handler) reqID(r *http.Request) string {
	return RequestIDFromContext(r.Context())
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.log.Warn("method not allowed for health endpoint",
			zap.String("request_id", h.reqID(r)),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.log.Debug("health endpoint called", zap.String("request_id", h.reqID(r)))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.log.Warn("method not allowed for student endpoint",
			zap.String("request_id", h.reqID(r)),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.log.Warn("invalid student id",
			zap.String("request_id", h.reqID(r)),
			zap.String("raw_id", r.PathValue("id")),
			zap.Error(err),
		)
		http.Error(w, "invalid student id", http.StatusBadRequest)
		return
	}

	st, err := h.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, student.ErrStudentNotFound) {
			h.log.Error("student not found",
				zap.String("request_id", h.reqID(r)),
				zap.Int64("student_id", id),
				zap.Error(err),
			)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": "студент с таким id не найден в системе",
				"id":    id,
			})
			return
		}
		h.log.Error("get student failed",
			zap.String("request_id", h.reqID(r)),
			zap.Int64("student_id", id),
			zap.Error(err),
		)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.log.Info("student returned successfully",
		zap.String("request_id", h.reqID(r)),
		zap.Int64("student_id", st.ID),
		zap.String("group", st.Group),
	)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(st)
}

func (h *Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.log.Warn("method not allowed for create student",
			zap.String("request_id", h.reqID(r)),
			zap.String("method", r.Method),
		)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		h.log.Error("failed to read request body",
			zap.String("request_id", h.reqID(r)),
			zap.Error(err),
		)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	h.log.Debug("create student request body received",
		zap.String("request_id", h.reqID(r)),
		zap.Int("body_size", len(body)),
		zap.String("body", string(body)),
	)

	var in student.CreateInput
	if err := json.Unmarshal(body, &in); err != nil {
		h.log.Warn("invalid json body",
			zap.String("request_id", h.reqID(r)),
			zap.Error(err),
		)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(in.FullName) == "" || strings.TrimSpace(in.Group) == "" || strings.TrimSpace(in.Email) == "" {
		h.log.Warn("validation failed: required fields missing",
			zap.String("request_id", h.reqID(r)),
			zap.String("full_name", in.FullName),
			zap.String("group", in.Group),
			zap.String("email", in.Email),
		)
		http.Error(w, "full_name, group and email are required", http.StatusBadRequest)
		return
	}

	st, err := h.repo.Create(in)
	if err != nil {
		if errors.Is(err, student.ErrEmailExists) {
			h.log.Warn("student create failed: email exists",
				zap.String("request_id", h.reqID(r)),
				zap.String("email", in.Email),
			)
			http.Error(w, "email already exists", http.StatusConflict)
			return
		}
		h.log.Error("student create failed",
			zap.String("request_id", h.reqID(r)),
			zap.Error(err),
		)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h.log.Info("student created successfully",
		zap.String("request_id", h.reqID(r)),
		zap.Int64("student_id", st.ID),
		zap.String("email", st.Email),
	)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(st)
}
