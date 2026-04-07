package grpc

import (
	"context"
	"tech-ip-sem2/services/auth/internal/usecases"
	"tech-ip-sem2/services/auth/pkg/authpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCAuthenticationHandler struct {
	authpb.UnimplementedAuthenticationServiceServer
	usecase *usecases.AuthenticationUsecase
}

func NewGRPCAuthenticationHandler(usecase *usecases.AuthenticationUsecase) *GRPCAuthenticationHandler {
	return &GRPCAuthenticationHandler{
		usecase: usecase,
	}
}

func (h *GRPCAuthenticationHandler) Verify(ctx context.Context, req *authpb.VerifyRequest) (*authpb.VerifyResponse, error) {
	out, err := h.usecase.Verify(ctx, usecases.VerifyInput{
		Token: req.GetToken(),
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &authpb.VerifyResponse{
		Subject: out.Subject,
	}, nil
}
