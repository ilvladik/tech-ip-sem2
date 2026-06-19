package main

import (
	"log"
	"net/http"

	"example.com/user-service/internal/user"
	"example.com/user-service/pkg/middleware"
)

func main() {
	repo := user.NewRepo()
	handler := user.NewHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", handler.ListUsers)
	mux.HandleFunc("GET /users/{id}", handler.GetUserByID)

	addr := ":8081"
	log.Println("user-service started on", addr)
	if err := http.ListenAndServe(addr, middleware.Logging(mux)); err != nil {
		log.Fatal(err)
	}
}
