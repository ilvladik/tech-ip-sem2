package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"example.com/pz14-job-queue/internal/publisher"
	"example.com/pz14-job-queue/internal/rabbit"
	"example.com/pz14-job-queue/internal/store"
	"example.com/pz14-job-queue/pkg/jobs"
	amqp "github.com/rabbitmq/amqp091-go"
)

const maxAttempts = 3

func main() {
	rabbitURL := env("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	prefetch := envInt("PREFETCH", 1)

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("rabbit: %v", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel: %v", err)
	}
	defer ch.Close()
	if err := rabbit.DeclareQueues(ch); err != nil {
		log.Fatalf("declare: %v", err)
	}
	if err := ch.Qos(prefetch, 0, false); err != nil {
		log.Fatalf("qos: %v", err)
	}

	msgs, err := ch.Consume(rabbit.QueueJobs, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume: %v", err)
	}

	processed := store.NewProcessedStore()
	log.Printf("worker started queue=%s prefetch=%d", rabbit.QueueJobs, prefetch)

	for d := range msgs {
		var job jobs.TaskJob
		if err := json.Unmarshal(d.Body, &job); err != nil {
			log.Printf("bad json: %v", err)
			_ = d.Nack(false, false)
			continue
		}
		if processed.Exists(job.MessageID) {
			log.Printf("duplicate message_id=%s skip", job.MessageID)
			_ = d.Ack(false)
			continue
		}

		log.Printf("processing task_id=%s attempt=%d message_id=%s", job.TaskID, job.Attempt, job.MessageID)
		if err := processTask(job); err != nil {
			log.Printf("error: %v", err)
			job.Attempt++
			if job.Attempt <= maxAttempts {
				if pubErr := publisher.PublishJob(ch, rabbit.QueueJobs, job); pubErr != nil {
					log.Printf("retry publish error: %v", pubErr)
				} else {
					log.Printf("retry scheduled attempt=%d", job.Attempt)
				}
				_ = d.Ack(false)
				continue
			}
			if pubErr := publisher.PublishJob(ch, rabbit.QueueJobsDLQ, job); pubErr != nil {
				log.Printf("dlq publish error: %v", pubErr)
			} else {
				log.Printf("moved to DLQ task_id=%s", job.TaskID)
			}
			_ = d.Ack(false)
			continue
		}

		processed.MarkDone(job.MessageID)
		_ = d.Ack(false)
		log.Printf("done task_id=%s message_id=%s", job.TaskID, job.MessageID)
	}
}

func processTask(job jobs.TaskJob) error {
	time.Sleep(1 * time.Second)
	if job.TaskID == "t_fail" {
		return errors.New("simulated processing error")
	}
	return nil
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func envInt(k string, d int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return d
}
