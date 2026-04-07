package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"tech-ip-sem2/services/tasks/internal/domain"
	"tech-ip-sem2/services/tasks/internal/usecases"
	"tech-ip-sem2/shared/httpx"
)

type TaskHandler struct {
	taskUsecase *usecases.TaskUsecase
}

func NewTaskHandler(taskUsecase *usecases.TaskUsecase) *TaskHandler {
	return &TaskHandler{
		taskUsecase: taskUsecase,
	}
}

func (h *TaskHandler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if req.Title == "" {
		httpx.RespondWithError(w, http.StatusBadRequest, "Title is required")
		return
	}
	var dueTime *time.Time
	if req.DueDate != nil {
		parsedTime, err := time.Parse(time.DateOnly, *req.DueDate)
		if err != nil {
			httpx.RespondWithError(w, http.StatusBadRequest,
				fmt.Sprintf("Invalid due_date format, expected %s", time.DateOnly))
			return
		}
		dueTime = &parsedTime
	}

	task, err := h.taskUsecase.Add(r.Context(), usecases.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueTime,
	})
	if err != nil {
		h.handleError(w, err)
	}
	httpx.RespondJson(w, http.StatusCreated, task)
}

func (h *TaskHandler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.taskUsecase.GetAll(r.Context())

	list := make([]TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		list = append(list, h.mapTaskToTaskResponse(task))
	}

	httpx.RespondJson(w, http.StatusOK, list)
}
func (h *TaskHandler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondWithError(w, http.StatusBadRequest, "Invalid task Id format")
		return
	}

	task, err := h.taskUsecase.Get(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpx.RespondJson(w, http.StatusOK, h.mapTaskToTaskResponse(*task))
}

func (h *TaskHandler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondWithError(w, http.StatusBadRequest, "Invalid task Id format")
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.RespondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	var dueTime *time.Time
	if req.DueDate != nil {
		parsedTime, err := time.Parse(time.DateOnly, *req.DueDate)
		if err != nil {
			httpx.RespondWithError(w, http.StatusBadRequest,
				fmt.Sprintf("Invalid due_date format, expected %s", time.DateOnly))
			return
		}
		dueTime = &parsedTime
	}

	task, err := h.taskUsecase.Update(r.Context(), id, usecases.UpdateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueTime,
		Done:        req.Done,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	httpx.RespondJson(w, http.StatusOK, h.mapTaskToTaskResponse(*task))
}

func (h *TaskHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpx.RespondWithError(w, http.StatusBadRequest, "Invalid task Id format")
		return
	}

	if err := h.taskUsecase.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case usecases.ErrorTaskNotFound:
		httpx.RespondWithError(w, http.StatusNotFound, err.Error())
	default:
		httpx.RespondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func (h *TaskHandler) mapTaskToTaskResponse(task domain.Task) TaskResponse {
	tr := TaskResponse{
		Title:       task.Title,
		Description: task.Description,
		Done:        task.Done,
	}

	if dd := task.DueDate; dd != nil {
		formatted := (*dd).Format(time.DateTime)
		tr.DueDate = &formatted
	}
	return tr
}
