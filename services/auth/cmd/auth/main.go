package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	apiGrpc "tech-ip-sem2/services/auth/internal/handlers/grpc"
	apiHttp "tech-ip-sem2/services/auth/internal/handlers/http"
	"tech-ip-sem2/services/auth/internal/usecases"
	"tech-ip-sem2/services/auth/pkg/authpb"
	"tech-ip-sem2/shared/interceptors"
	"tech-ip-sem2/shared/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log, err := logger.New("auth")
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	authenticationUsecase := usecases.NewAuthenticationUsecase()

	g, ctx := errgroup.WithContext(ctx)

	httpServer := newHTTPServer(authenticationUsecase, log)

	g.Go(func() error {
		log.Info("http server started",
			zap.String("component", "http_server"),
			zap.String("addr", httpServer.Addr),
		)

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server failed",
				zap.String("component", "http_server"),
				zap.Error(err),
			)
			return err
		}

		return nil
	})

	grpcServer, lis := newGRPCServer(authenticationUsecase, log)

	g.Go(func() error {
		log.Info("grpc server started",
			zap.String("component", "grpc_server"),
			zap.String("addr", lis.Addr().String()),
		)

		if err := grpcServer.Serve(lis); err != nil {
			log.Error("grpc server failed",
				zap.String("component", "grpc_server"),
				zap.Error(err),
			)
			return err
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		log.Info("shutting down servers",
			zap.String("component", "main"),
		)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error("http shutdown error",
				zap.String("component", "http_server"),
				zap.Error(err),
			)
		}

		grpcServer.GracefulStop()

		return nil
	})

	if err := g.Wait(); err != nil {
		log.Error("server error",
			zap.String("component", "main"),
			zap.Error(err),
		)
	}

	log.Info("servers stopped gracefully",
		zap.String("component", "main"),
	)
}

func newHTTPServer(authenticationUsecase *usecases.AuthenticationUsecase, log *zap.Logger) *http.Server {
	port := getEnv("HTTP_PORT", "8081")

	handler := apiHttp.RegisterRoutes(authenticationUsecase, log)

	return &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
}

func newGRPCServer(authenticationUsecase *usecases.AuthenticationUsecase, log *zap.Logger) (*grpc.Server, net.Listener) {
	port := getEnv("GRPC_PORT", "50051")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("failed to listen",
			zap.String("component", "grpc_server"),
			zap.Error(err),
		)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RequestIDUnaryInterceptor(),
			interceptors.AccessLogUnaryInterceptor(log),
		),
	)

	handler := apiGrpc.NewGRPCAuthenticationHandler(authenticationUsecase)
	authpb.RegisterAuthenticationServiceServer(s, handler)

	return s, lis
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
