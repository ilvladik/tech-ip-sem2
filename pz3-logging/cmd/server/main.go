package main

import (
	"log"
	"net/http"

	"example.com/pz3-logging/internal/httpapi"
	"example.com/pz3-logging/internal/student"
	applogger "example.com/pz3-logging/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	logger, err := applogger.New()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	repo := student.NewRepo(logger)
	handler := httpapi.NewHandler(repo, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /students/{id}", handler.GetStudentByID)
	mux.HandleFunc("POST /students", handler.CreateStudent)

	rootHandler := httpapi.LoggingMiddleware(logger, mux)

	logger.Info("server is starting", zap.String("addr", ":8083"))
	if err := http.ListenAndServe(":8083", rootHandler); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}
