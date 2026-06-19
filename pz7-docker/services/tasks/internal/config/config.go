package config

import "os"

type Config struct {
	Addr          string
	AuthBaseURL   string
	ListenAddress string
}

func New() Config {
	port := env("TASKS_PORT", "8088")
	return Config{
		Addr:          ":" + port,
		AuthBaseURL:   env("AUTH_BASE_URL", "http://localhost:8087"),
		ListenAddress: env("LISTEN_ADDRESS", "0.0.0.0"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
