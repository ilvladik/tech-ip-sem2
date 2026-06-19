package main

import (
	"log"
	"net/http"

	"example.com/pz4-monitoring/internal/httpapi"
	"example.com/pz4-monitoring/internal/student"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	repo := student.NewRepo()
	handler := httpapi.NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /students/{id}", handler.GetStudentByID)
	mux.Handle("GET /metrics", promhttp.Handler())

	rootHandler := httpapi.MetricsMiddleware(mux)

	log.Println("server started on :8084")
	if err := http.ListenAndServe(":8084", rootHandler); err != nil {
		log.Fatal(err)
	}
}
