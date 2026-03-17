package client

import (
	"context"
	"fmt"
	authpb "tech-ip-sem2/services/auth/pkg/authpb/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type GRPCClient struct {
	client authpb.AuthenticationServiceClient
	conn   *grpc.ClientConn
}

func NewGRPCClient(addr string) (*GRPCClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &GRPCClient{
		client: authpb.NewAuthenticationServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

func (c *GRPCClient) Verify(ctx context.Context, token string) (bool, string, error) {
	if token == "" {
		return false, "", fmt.Errorf("empty token")
	}

	resp, err := c.client.Verify(ctx, &authpb.VerifyRequest{Token: token})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unauthenticated:
				return false, "", nil
			case codes.DeadlineExceeded:
				return false, "", fmt.Errorf("auth service timeout: %w", err)
			case codes.Unavailable:
				return false, "", fmt.Errorf("auth service unavailable: %w", err)
			default:
				return false, "", fmt.Errorf("auth service error: %w", err)
			}
		}
		return false, "", fmt.Errorf("failed to verify token: %w", err)
	}

	return resp.Valid, resp.Subject, nil
}
