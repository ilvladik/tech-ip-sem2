package config

import "os"

type Config struct {
	ServerAddr   string
	RabbitURL    string
	QueueName    string
	PublishMode  string // best_effort | strict
	ValidToken   string
}

func TasksConfig() Config {
	return Config{
		ServerAddr:  env("SERVER_ADDR", ":8096"),
		RabbitURL:   env("RABBIT_URL", "amqp://guest:guest@localhost:5672/"),
		QueueName:   env("QUEUE_NAME", "task_events"),
		PublishMode: env("PUBLISH_MODE", "best_effort"),
		ValidToken:  env("AUTH_TOKEN", "demo-token"),
	}
}

func WorkerConfig() Config {
	return Config{
		RabbitURL: env("RABBIT_URL", "amqp://guest:guest@localhost:5672/"),
		QueueName: env("QUEUE_NAME", "task_events"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
