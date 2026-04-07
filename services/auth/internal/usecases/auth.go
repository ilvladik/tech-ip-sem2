package usecases

import (
	"context"
	"tech-ip-sem2/services/auth/internal/domain"
)

type AuthenticationUsecase struct {
	users []domain.User
}

func NewAuthenticationUsecase() *AuthenticationUsecase {
	return &AuthenticationUsecase{
		users: []domain.User{
			{
				Login:    "student",
				Password: "student",
				Token:    "demo-token",
			},
		},
	}
}

func (u *AuthenticationUsecase) Login(ctx context.Context, in LoginInput) (*LoginOutput, error) {
	for _, user := range u.users {
		if user.Login == in.Login {
			return &LoginOutput{
				AccessToken: user.Token,
				TokenType:   "Bearer",
			}, nil
		}
	}
	return nil, domain.ErrorInvalidCredentials
}

func (u *AuthenticationUsecase) Verify(ctx context.Context, in VerifyInput) (*VerifyOutput, error) {
	for _, user := range u.users {
		if user.Token == in.Token {
			return &VerifyOutput{
				Subject: user.Login,
			}, nil
		}
	}
	return nil, domain.ErrorInvalidToken
}
