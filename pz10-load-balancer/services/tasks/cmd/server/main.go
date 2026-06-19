package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Task struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

func main() {
	instanceID := os.Getenv("INSTANCE_ID")
	if instanceID == "" {
		instanceID = "tasks-unknown"
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8091"
	}

	tasks := []Task{
		{ID: 1, Title: "Изучить NGINX"},
		{ID: 2, Title: "Освоить load balancing"},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler(instanceID))
	mux.HandleFunc("GET /whoami", whoamiHandler(instanceID))
	mux.HandleFunc("GET /v1/tasks", tasksHandler(instanceID, tasks))

	handler := loggingMiddleware(instanceID, mux)
	addr := ":" + port
	log.Println("tasks service started on", addr, "instance =", instanceID)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(instanceID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, instanceID, map[string]string{
			"status":   "ok",
			"instance": instanceID,
		})
	}
}

func whoamiHandler(instanceID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, instanceID, map[string]string{
			"instance": instanceID,
		})
	}
}

func tasksHandler(instanceID string, tasks []Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, instanceID, tasks)
	}
}

func writeJSON(w http.ResponseWriter, instanceID string, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Instance-ID", instanceID)
	_ = json.NewEncoder(w).Encode(payload)
}

func loggingMiddleware(instanceID string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"instance=%s method=%s path=%s remote=%s duration=%s",
			instanceID,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}
