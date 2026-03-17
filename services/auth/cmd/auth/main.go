package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"tech-ip-sem2/services/auth/internal/authgrpc"
	authHttp "tech-ip-sem2/services/auth/internal/http"
	"tech-ip-sem2/services/auth/internal/service"
	authpb "tech-ip-sem2/services/auth/pkg/authpb/proto"

	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	authService := service.NewAuthService()

	httpErr := make(chan error, 1)
	grpcErr := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startHTTPServer(ctx, authService); err != nil {
			httpErr <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := startGRPCServer(ctx, authService); err != nil {
			grpcErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutting down servers...")
	case err := <-httpErr:
		log.Printf("HTTP server error: %v", err)
	case err := <-grpcErr:
		log.Printf("gRPC server error: %v", err)
	}

	time.Sleep(2 * time.Second)
	wg.Wait()
	log.Println("All servers stopped")
}

func startHTTPServer(ctx context.Context, authService *service.AuthenticationService) error {
	port := getEnv("HTTP_PORT", "8081")

	handler := authHttp.RegisterRoutes(authService)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down HTTP server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP shutdown error: %v", err)
		}
	}()

	log.Printf("HTTP server starting on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func startGRPCServer(ctx context.Context, authService *service.AuthenticationService) error {
	port := getEnv("GRPC_PORT", "50051")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	authServer := authgrpc.NewAuthenticationServer(authService)
	authpb.RegisterAuthenticationServiceServer(s, authServer)

	go func() {
		<-ctx.Done()
		log.Println("Shutting down gRPC server...")
		s.GracefulStop()
	}()

	log.Printf("gRPC server starting on port %s", port)
	if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		return err
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
