package http

import (
	"net/http"
	"strings"

	"tech-ip-sem2/services/tasks/internal/handlers/http/middlewares"
	"tech-ip-sem2/services/tasks/internal/usecases"
	sharedMiddleware "tech-ip-sem2/shared/middlewares"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func RegisterRoutes(
	taskService *usecases.TaskUsecase,
	authClient middlewares.AuthenticationClient,
	log *zap.Logger,
) http.Handler {
	taskMux := http.NewServeMux()
	h := NewTaskHandler(taskService)

	taskMux.HandleFunc("POST /v1/tasks", h.handleCreateTask)
	taskMux.HandleFunc("GET /v1/tasks", h.handleGetTasks)
	taskMux.HandleFunc("GET /v1/tasks/{id}", h.handleGetTask)
	taskMux.HandleFunc("PATCH /v1/tasks/{id}", h.handleUpdateTask)
	taskMux.HandleFunc("DELETE /v1/tasks/{id}", h.handleDeleteTask)

	authMiddleware := middlewares.NewAuthenticationMiddleware(authClient)

	var apiHandler http.Handler = taskMux
	apiHandler = authMiddleware.Authenticate(apiHandler)
	apiHandler = sharedMiddleware.AccessLog(log)(apiHandler)
	apiHandler = sharedMiddleware.Metrics(apiHandler, routeName)
	apiHandler = sharedMiddleware.RequestId(apiHandler)

	rootMux := http.NewServeMux()
	rootMux.Handle("/metrics", promhttp.Handler())
	rootMux.Handle("/", apiHandler)

	return rootMux
}

func routeName(r *http.Request) string {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/v1/tasks":
		return "/v1/tasks"
	case r.Method == http.MethodGet && r.URL.Path == "/v1/tasks":
		return "/v1/tasks"
	case r.URL.Path == "/v1/tasks" && (r.Method == http.MethodPatch || r.Method == http.MethodDelete):
		return "/v1/tasks"
	default:
		if strings.HasPrefix(r.URL.Path, "/v1/tasks/") {
			return "/v1/tasks/:id"
		}
		return r.URL.Path
	}
}
