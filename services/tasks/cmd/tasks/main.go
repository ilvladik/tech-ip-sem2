package main

import (
	"log"
	"net/http"
	"os"
	authclient "tech-ip-sem2/services/tasks/internal/client"
	taskHttp "tech-ip-sem2/services/tasks/internal/http"
	"tech-ip-sem2/services/tasks/internal/service"
	"time"
)

func main() {
	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
		log.Panicf("TASKS_PORT not set, using default: %s", port)
	}

	authBaseURL := os.Getenv("AUTH_BASE_URL")
	if authBaseURL == "" {
		authBaseURL = "http://localhost:8081"
		log.Panicf("AUTH_BASE_URL not set, using default: %s", authBaseURL)
	}

	taskService := service.NewTaskService()
	authClient := authclient.NewClient(authBaseURL, 3*time.Second)

	handler := taskHttp.RegisterRoutes(taskService, authClient)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
