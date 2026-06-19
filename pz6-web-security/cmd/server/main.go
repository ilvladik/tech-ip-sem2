package main

import (
	"log"
	"net/http"

	"example.com/pz6-web-security/internal/config"
	"example.com/pz6-web-security/internal/httpapi"
	"example.com/pz6-web-security/internal/store"
)

func main() {
	cfg := config.New()
	s := store.New()

	handler, err := httpapi.NewHandler(s, cfg.SecureCookie)
	if err != nil {
		log.Fatalf("init handler: %v", err)
	}

	mux := http.NewServeMux()
	handler.Register(mux)

	log.Printf("web security server started on %s (secure cookie=%v)", cfg.Addr, cfg.SecureCookie)
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
