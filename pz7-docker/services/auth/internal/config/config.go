package config

import "os"

type Config struct {
	Addr          string
	ValidToken    string
	ListenAddress string
}

func New() Config {
	port := env("AUTH_PORT", "8087")
	return Config{
		Addr:          ":" + port,
		ValidToken:    env("AUTH_VALID_TOKEN", "demo-token"),
		ListenAddress: env("LISTEN_ADDRESS", "0.0.0.0"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
