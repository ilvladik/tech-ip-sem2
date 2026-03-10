package service

import "errors"

var (
	ErrorInvalidCredentials = errors.New("Invalid credentials")
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type VerifyResponse struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject,omitempty"`
	Error   string `json:"error"`
}

type AuthenticationService struct{}

func NewAuthService() *AuthenticationService {
	return &AuthenticationService{}
}

func (s *AuthenticationService) Login(request LoginRequest) (*LoginResponse, error) {
	if request.Username == "student" && request.Password == "student" {
		return &LoginResponse{
			AccessToken: "demo-token",
			TokenType:   "Bearer",
		}, nil
	}
	return nil, ErrorInvalidCredentials
}

func (s *AuthenticationService) Verify(token string) (*VerifyResponse, error) {
	if token == "demo-token" {
		return &VerifyResponse{
			Valid:   true,
			Subject: "student",
		}, nil
	}
	return &VerifyResponse{
		Valid: false,
		Error: "unauthorized",
	}, nil
}
