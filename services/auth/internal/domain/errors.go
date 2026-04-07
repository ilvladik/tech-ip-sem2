package domain

import "errors"

var (
	ErrorInvalidCredentials = errors.New("Invalid credentials")
	ErrorInvalidToken       = errors.New("Invalid token")
)
