package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"example.com/pz13-rabbitmq/internal/amqp"
	"example.com/pz13-rabbitmq/internal/config"
	"example.com/pz13-rabbitmq/internal/publisher"
	"example.com/pz13-rabbitmq/internal/task"
	"example.com/pz13-rabbitmq/pkg/events"
)

func main() {
	cfg := config.TasksConfig()
	store := task.NewStore()

	var pub *publisher.Publisher
	conn := amqp.MustConnect(cfg.RabbitURL)
	ch := amqp.MustChannel(conn)
	var err error
	pub, err = publisher.New(ch, cfg.QueueName)
	if err != nil {
		log.Fatalf("publisher init: %v", err)
	}
	defer pub.Close()
	defer conn.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "tasks"})
	})
	mux.HandleFunc("POST /v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		if !checkAuth(r, cfg.ValidToken) {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		var body struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" {
			http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
			return
		}

		t := store.Create(body.Title, body.Description)
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = "unknown"
		}

		ev := events.NewTaskCreated(t.ID, requestID, "tasks-service")
		if err := pub.PublishTaskCreated(r.Context(), ev); err != nil {
			log.Printf("publish error: %v", err)
			if cfg.PublishMode == "strict" {
				http.Error(w, `{"error":"failed to publish event"}`, http.StatusInternalServerError)
				return
			}
		}

		writeJSON(w, http.StatusCreated, t)
	})

	log.Printf("tasks service on %s (publish_mode=%s)", cfg.ServerAddr, cfg.PublishMode)
	if err := http.ListenAndServe(cfg.ServerAddr, mux); err != nil {
		log.Fatal(err)
	}
}

func checkAuth(r *http.Request, token string) bool {
	auth := r.Header.Get("Authorization")
	return strings.TrimPrefix(auth, "Bearer ") == token
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
