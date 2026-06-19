package config

import "os"

type Config struct {
	HTTPSAddr        string
	HTTPRedirectAddr string
	CertFile         string
	KeyFile          string
	DSN              string
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func New() Config {
	return Config{
		HTTPSAddr:        env("HTTPS_ADDR", ":8443"),
		HTTPRedirectAddr: env("HTTP_REDIRECT_ADDR", ":8085"),
		CertFile:         env("TLS_CERT_FILE", "certs/server.crt"),
		KeyFile:          env("TLS_KEY_FILE", "certs/server.key"),
		DSN: env("DB_DSN",
			"postgres://postgres:postgres@localhost:5432/study_security?sslmode=disable"),
	}
}
