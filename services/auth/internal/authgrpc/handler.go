package authgrpc

import (
	"context"
	"log"
	"strings"
	"tech-ip-sem2/services/auth/internal/service"
	authpb "tech-ip-sem2/services/auth/pkg/authpb/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthenticationServer struct {
	authpb.UnimplementedAuthenticationServiceServer
	authService *service.AuthenticationService
}

func NewAuthenticationServer(authService *service.AuthenticationService) *AuthenticationServer {
	return &AuthenticationServer{
		authService: authService,
	}
}

func (s *AuthenticationServer) Verify(ctx context.Context, req *authpb.VerifyRequest) (*authpb.VerifyResponse, error) {
	token := req.GetToken()
	log.Printf("[gRPC IN] Verify called with token length: %d", len(token))
	if token == "" {
		log.Printf("[gRPC ERROR] Empty token")
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	if after, ok := strings.CutPrefix(token, "Bearer "); ok {
		token = after
	}

	verifyResponse, err := s.authService.Verify(token)
	if err != nil {
		log.Printf("[gRPC ERROR] Verify failed: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Printf("[gRPC OUT] Verify success - valid: %v, subject: %s",
		verifyResponse.Valid, verifyResponse.Subject)

	return &authpb.VerifyResponse{
		Valid:   verifyResponse.Valid,
		Subject: verifyResponse.Subject,
		Error:   verifyResponse.Error,
	}, nil
}
