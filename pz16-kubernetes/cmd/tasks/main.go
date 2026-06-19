package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	port := env("TASKS_PORT", "8082")
	authURL := env("AUTH_BASE_URL", "http://auth:8081")
	logLevel := env("LOG_LEVEL", "info")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	addr := ":" + port
	log.Printf("tasks listening on %s log_level=%s auth=%s", addr, logLevel, authURL)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
