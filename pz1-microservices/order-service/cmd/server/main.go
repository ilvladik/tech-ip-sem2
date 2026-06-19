package main

import (
	"log"
	"net/http"
	"os"

	"example.com/order-service/internal/order"
	"example.com/order-service/pkg/middleware"
)

func main() {
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://localhost:8081"
	}

	repo := order.NewRepo()
	client := order.NewUserServiceClient(userServiceURL)
	handler := order.NewHandler(repo, client)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /orders/by-user/{userID}", handler.GetOrdersByUser)
	mux.HandleFunc("GET /orders/{id}/full", handler.GetOrderWithUser)
	mux.HandleFunc("GET /orders/{id}", handler.GetOrderByID)

	addr := ":8082"
	log.Println("order-service started on", addr, "user-service:", userServiceURL)
	if err := http.ListenAndServe(addr, middleware.Logging(mux)); err != nil {
		log.Fatal(err)
	}
}
