package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	authclient "tech-ip-sem2/services/tasks/internal/adapters/auth"
	taskHttp "tech-ip-sem2/services/tasks/internal/handlers/http"
	"tech-ip-sem2/services/tasks/internal/usecases"
	"tech-ip-sem2/shared/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log, err := logger.New("tasks")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	g, ctx := errgroup.WithContext(ctx)

	port := getEnv("TASKS_PORT", "8082")
	authAddr := getEnv("AUTH_GRPC_URL", "localhost:50051")

	taskUsecase := usecases.NewTaskUsecase()

	authClient, err := authclient.NewGrpcAuthenticationClient(authAddr, log)
	if err != nil {
		log.Fatal("failed to connect to auth service",
			zap.String("component", "auth_client"),
			zap.String("auth_addr", authAddr),
			zap.Error(err),
		)
	}
	defer func() {
		if err := authClient.Close(); err != nil {
			log.Error("auth client close error",
				zap.String("component", "auth_client"),
				zap.Error(err),
			)
		}
	}()

	handler := taskHttp.RegisterRoutes(taskUsecase, authClient, log)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	g.Go(func() error {
		log.Info("tasks http server started",
			zap.String("component", "http_server"),
			zap.String("addr", server.Addr),
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("tasks http server failed",
				zap.String("component", "http_server"),
				zap.Error(err),
			)
			return err
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		log.Info("shutting down tasks http server",
			zap.String("component", "main"),
		)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("http shutdown error",
				zap.String("component", "http_server"),
				zap.Error(err),
			)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		log.Error("server error",
			zap.String("component", "main"),
			zap.Error(err),
		)
	}

	log.Info("tasks server stopped gracefully",
		zap.String("component", "main"),
	)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
