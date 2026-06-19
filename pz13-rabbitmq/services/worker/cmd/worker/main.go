package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"example.com/pz13-rabbitmq/internal/amqp"
	"example.com/pz13-rabbitmq/internal/config"
	"example.com/pz13-rabbitmq/pkg/events"
	amqplib "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := config.WorkerConfig()
	prefetch := envInt("PREFETCH", 1)

	conn := amqp.MustConnect(cfg.RabbitURL)
	defer conn.Close()

	ch := amqp.MustChannel(conn)
	defer ch.Close()

	_, err := ch.QueueDeclare(cfg.QueueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("queue declare: %v", err)
	}
	if err := ch.Qos(prefetch, 0, false); err != nil {
		log.Fatalf("qos: %v", err)
	}

	msgs, err := ch.Consume(cfg.QueueName, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume: %v", err)
	}

	log.Printf("worker started queue=%s prefetch=%d", cfg.QueueName, prefetch)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sig:
			log.Println("worker stopped")
			return
		case d, ok := <-msgs:
			if !ok {
				return
			}
			handleDelivery(d)
		}
	}
}

func handleDelivery(d amqplib.Delivery) {
	var ev events.TaskEvent
	if err := json.Unmarshal(d.Body, &ev); err != nil {
		log.Printf("bad message: %v body=%s", err, string(d.Body))
		_ = d.Nack(false, false)
		return
	}
	log.Printf("received event=%s task_id=%s ts=%s request_id=%s producer=%s",
		ev.Event, ev.TaskID, ev.TS, ev.RequestID, ev.Producer)
	if err := d.Ack(false); err != nil {
		log.Printf("ack error: %v", err)
	}
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
