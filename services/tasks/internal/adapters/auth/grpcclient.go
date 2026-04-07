package client

import (
	"context"
	"fmt"
	"tech-ip-sem2/services/auth/pkg/authpb"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type GrpcAuthenticationClient struct {
	client authpb.AuthenticationServiceClient
	conn   *grpc.ClientConn
	log    *zap.Logger
}

func NewGrpcAuthenticationClient(addr string, log *zap.Logger) (*GrpcAuthenticationClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("Authentication service connection failed: %w", err)
	}

	return &GrpcAuthenticationClient{
		client: authpb.NewAuthenticationServiceClient(conn),
		conn:   conn,
		log:    log,
	}, nil
}

func (c *GrpcAuthenticationClient) Close() error {
	return c.conn.Close()
}

func (c *GrpcAuthenticationClient) Verify(ctx context.Context, token string) (*string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.client.Verify(ctx, &authpb.VerifyRequest{Token: token})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return nil, nil
			case codes.DeadlineExceeded:
				return nil, fmt.Errorf("Authentication service timeout: %w", err)
			case codes.Unavailable:
				return nil, fmt.Errorf("Authentication service unavailable: %w", err)
			default:
				return nil, fmt.Errorf("auth service error: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	return &resp.Subject, nil
}
