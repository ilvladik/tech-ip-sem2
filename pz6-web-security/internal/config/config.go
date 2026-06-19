package config

import "os"

type Config struct {
	Addr         string
	SecureCookie bool
}

func New() Config {
	secure := os.Getenv("SECURE_COOKIE") == "true"
	return Config{
		Addr:         env("SERVER_ADDR", ":8086"),
		SecureCookie: secure,
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
