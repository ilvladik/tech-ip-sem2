package http

import (
	"net/http"

	"tech-ip-sem2/services/tasks/internal/handlers/http/middlewares"
	"tech-ip-sem2/services/tasks/internal/usecases"
	sharedMiddleware "tech-ip-sem2/shared/middlewares"

	"go.uber.org/zap"
)

func RegisterRoutes(
	taskService *usecases.TaskUsecase,
	authClient middlewares.AuthenticationClient,
	log *zap.Logger,
) http.Handler {
	mux := http.NewServeMux()
	h := NewTaskHandler(taskService)

	mux.HandleFunc("POST /v1/tasks", h.handleCreateTask)
	mux.HandleFunc("GET /v1/tasks", h.handleGetTasks)
	mux.HandleFunc("GET /v1/tasks/{id}", h.handleGetTask)
	mux.HandleFunc("PATCH /v1/tasks/{id}", h.handleUpdateTask)
	mux.HandleFunc("DELETE /v1/tasks/{id}", h.handleDeleteTask)

	authMiddleware := middlewares.NewAuthenticationMiddleware(authClient)

	var handler http.Handler = mux
	handler = authMiddleware.Authenticate(handler)
	handler = sharedMiddleware.AccessLog(log)(handler)
	handler = sharedMiddleware.RequestId(handler)

	return handler
}
