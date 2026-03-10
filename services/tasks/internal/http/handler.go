package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"tech-ip-sem2/services/tasks/internal/service"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req service.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Title == "" {
		respondWithError(w, http.StatusBadRequest, "Title is required")
		return
	}

	if req.DueDate != nil {
		if _, err := time.Parse("2006-01-02", *req.DueDate); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid due_date format, expected YYYY-MM-DD")
			return
		}
	}

	task := h.taskService.Create(req)
	respondJson(w, http.StatusCreated, task)
}

func (h *TaskHandler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.taskService.GetAll()

	response := make([]service.Task, len(tasks))
	for i, task := range tasks {
		response[i] = *task
	}

	respondJson(w, http.StatusOK, response)
}
func (h *TaskHandler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	task, err := h.taskService.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "")
		return
	}

	respondJson(w, http.StatusOK, &task)
}

func (h *TaskHandler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task Id format")
		return
	}

	var req service.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.DueDate != nil {
		if _, err := time.Parse("2006-01-02", *req.DueDate); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid due_date format, expected YYYY-MM-DD")
			return
		}
	}

	task, err := h.taskService.Update(id, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	respondJson(w, http.StatusOK, task)
}

func (h *TaskHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID format")
		return
	}

	if err := h.taskService.Delete(id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case service.ErrorTaskNotFound:
		respondWithError(w, http.StatusNotFound, err.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func respondJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondJson(w, status, map[string]string{"error": message})
}
