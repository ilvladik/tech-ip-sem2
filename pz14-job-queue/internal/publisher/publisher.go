package publisher

import (
	"context"
	"encoding/json"

	"example.com/pz14-job-queue/pkg/jobs"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJob(ch *amqp.Channel, queue string, job jobs.TaskJob) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return ch.PublishWithContext(context.Background(), "", queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}
