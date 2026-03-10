package main

import (
	"log"
	"net/http"
	"os"
	authHttp "tech-ip-sem2/services/auth/internal/http"
	"tech-ip-sem2/services/auth/internal/service"
)

func main() {
	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8081"
		log.Panicf("AUTH_PORT not set, using default: %s", port)
	}
	handler := authHttp.RegisterRoutes(service.NewAuthService())
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
