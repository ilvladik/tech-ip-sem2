package interceptors

import (
	"context"
	"tech-ip-sem2/shared/requestctx"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func AccessLogUnaryInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		fields := []zap.Field{
			zap.String("component", "grpc_server"),
			zap.String("request_id", requestctx.RequestID(ctx)),
			zap.String("method", info.FullMethod),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		}

		if err == nil {
			log.Info("grpc request completed", fields...)
			return resp, nil
		}

		_, ok := status.FromError(err)
		if !ok {
			log.Error("grpc request failed (non-status error)",
				append(fields, zap.Error(err))...,
			)
			return resp, err
		}

		return resp, err
	}
}
