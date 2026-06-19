package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"example.com/pz14-job-queue/internal/publisher"
	"example.com/pz14-job-queue/internal/rabbit"
	"example.com/pz14-job-queue/pkg/jobs"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	addr := env("SERVER_ADDR", ":8097")
	rabbitURL := env("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	token := env("AUTH_TOKEN", "demo-token")

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("rabbit: %v", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel: %v", err)
	}
	defer ch.Close()
	if err := rabbit.DeclareQueues(ch); err != nil {
		log.Fatalf("declare queues: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("POST /v1/jobs/process-task", func(w http.ResponseWriter, r *http.Request) {
		if !authOK(r, token) {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		var body struct {
			TaskID string `json:"task_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.TaskID == "" {
			http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
			return
		}
		job := jobs.TaskJob{
			Job:       "process_task",
			TaskID:    body.TaskID,
			Attempt:   1,
			MessageID: uuid.NewString(),
		}
		if err := publisher.PublishJob(ch, rabbit.QueueJobs, job); err != nil {
			log.Printf("publish error: %v", err)
			http.Error(w, `{"error":"publish failed"}`, http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusAccepted, map[string]string{
			"status":  "accepted",
			"task_id": body.TaskID,
		})
	})

	log.Printf("tasks jobs API on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func authOK(r *http.Request, token string) bool {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ") == token
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
