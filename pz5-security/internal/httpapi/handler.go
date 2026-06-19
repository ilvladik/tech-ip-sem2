package httpapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"example.com/pz5-security/internal/student"
)

type Handler struct {
	repo          *student.Repo
	stmtByID      *sql.Stmt
	stmtByEmail   *sql.Stmt
}

func NewHandler(repo *student.Repo, stmtByID, stmtByEmail *sql.Stmt) *Handler {
	return &Handler{
		repo:        repo,
		stmtByID:    stmtByID,
		stmtByEmail: stmtByEmail,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"scheme": "https",
	})
}

func (h *Handler) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rawID := r.URL.Query().Get("id")
	if rawID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, ok := validateStudentID(rawID)
	if !ok {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var st student.Student
	err := h.stmtByID.QueryRow(id).Scan(&st.ID, &st.FullName, &st.StudyGroup, &st.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": "студент с таким id не найден в базе",
				"id":    id,
			})
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(st)
}

func (h *Handler) GetStudentByEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}
	if !validateEmail(email) {
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}

	var st student.Student
	err := h.stmtByEmail.QueryRow(email).Scan(&st.ID, &st.FullName, &st.StudyGroup, &st.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(st)
}

// GetStudentUnsafe — только для демонстрации SQL-инъекции в учебных целях.
func (h *Handler) GetStudentUnsafe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rawID := r.URL.Query().Get("id")
	if rawID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	st, err := h.repo.UnsafeGetByID(rawID)
	if err != nil {
		if errors.Is(err, student.ErrStudentNotFound) {
			http.Error(w, "student not found", http.StatusNotFound)
			return
		}
		http.Error(w, "query error: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"warning": "unsafe SQL concatenation — do not use in production",
		"query":   h.repo.UnsafeGetByIDQuery(rawID),
		"student": st,
	})
}
