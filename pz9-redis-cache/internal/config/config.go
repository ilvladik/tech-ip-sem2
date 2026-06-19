package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerAddr        string
	RedisAddr         string
	RedisPassword     string
	CacheTTL          time.Duration
	CacheTTLJitter    time.Duration
	ListCacheTTL      time.Duration
	NegativeCacheTTL  time.Duration
	RedisDialTimeout  time.Duration
	RedisReadTimeout  time.Duration
	RedisWriteTimeout time.Duration
}

func New() Config {
	return Config{
		ServerAddr:        env("SERVER_ADDR", ":8089"),
		RedisAddr:         env("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     env("REDIS_PASSWORD", ""),
		CacheTTL:          durationEnv("CACHE_TTL_SEC", 120) * time.Second,
		CacheTTLJitter:    durationEnv("CACHE_TTL_JITTER_SEC", 30) * time.Second,
		ListCacheTTL:      durationEnv("LIST_CACHE_TTL_SEC", 60) * time.Second,
		NegativeCacheTTL:  durationEnv("NEGATIVE_CACHE_TTL_SEC", 30) * time.Second,
		RedisDialTimeout:  durationEnv("REDIS_DIAL_TIMEOUT_SEC", 2) * time.Second,
		RedisReadTimeout:  durationEnv("REDIS_READ_TIMEOUT_SEC", 2) * time.Second,
		RedisWriteTimeout: durationEnv("REDIS_WRITE_TIMEOUT_SEC", 2) * time.Second,
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func durationEnv(key string, fallback int64) time.Duration {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return time.Duration(n)
		}
	}
	return time.Duration(fallback)
}
