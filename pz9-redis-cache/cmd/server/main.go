package main

import (
	"context"
	"log"
	"net/http"

	"example.com/pz9-redis-cache/internal/cache"
	"example.com/pz9-redis-cache/internal/config"
	"example.com/pz9-redis-cache/internal/httpapi"
	"example.com/pz9-redis-cache/internal/service"
	"example.com/pz9-redis-cache/internal/task"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.New()
	repo := task.NewRepo()

	var redisClient *redis.Client
	client := cache.NewRedisClient(cfg)
	if err := cache.Ping(context.Background(), client); err != nil {
		log.Println("warning: redis is unavailable at startup:", err)
		log.Println("service will work without cache (fallback to repository)")
	} else {
		redisClient = client
		log.Println("redis connected:", cfg.RedisAddr)
	}

	taskService := service.NewTaskService(repo, redisClient, cfg)
	handler := httpapi.NewHandler(taskService)

	mux := http.NewServeMux()
	handler.Register(mux)

	log.Printf("server started on %s", cfg.ServerAddr)
	if err := http.ListenAndServe(cfg.ServerAddr, mux); err != nil {
		log.Fatal(err)
	}
}
