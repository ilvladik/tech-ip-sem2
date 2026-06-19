package main

import (
	"log"
	"net/http"

	"example.com/pz7-docker/tasks/internal/authclient"
	"example.com/pz7-docker/tasks/internal/config"
	"example.com/pz7-docker/tasks/internal/httpapi"
)

func main() {
	cfg := config.New()
	mux := http.NewServeMux()
	httpapi.New(authclient.New(cfg.AuthBaseURL)).Register(mux)

	addr := cfg.ListenAddress + cfg.Addr
	log.Printf("tasks service started on %s (auth=%s)", addr, cfg.AuthBaseURL)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
