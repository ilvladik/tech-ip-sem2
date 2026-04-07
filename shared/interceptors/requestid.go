package interceptors

import (
	"context"
	"tech-ip-sem2/shared/requestctx"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestIDUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		md, _ := metadata.FromIncomingContext(ctx)

		var requestID string
		if values := md.Get("x-request-id"); len(values) > 0 {
			requestID = values[0]
		}

		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx = requestctx.WithRequestID(ctx, requestID)

		grpc.SetHeader(ctx, metadata.Pairs("x-request-id", requestID))

		return handler(ctx, req)
	}
}
