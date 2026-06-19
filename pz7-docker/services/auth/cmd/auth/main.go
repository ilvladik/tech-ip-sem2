package main

import (
	"log"
	"net/http"

	"example.com/pz7-docker/auth/internal/config"
	"example.com/pz7-docker/auth/internal/httpapi"
)

func main() {
	cfg := config.New()
	mux := http.NewServeMux()
	httpapi.New(cfg.ValidToken).Register(mux)

	addr := cfg.ListenAddress + cfg.Addr
	log.Printf("auth service started on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
