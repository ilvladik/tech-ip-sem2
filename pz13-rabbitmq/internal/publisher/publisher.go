package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"example.com/pz13-rabbitmq/pkg/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch        *amqp.Channel
	queueName string
}

func New(ch *amqp.Channel, queueName string) (*Publisher, error) {
	_, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return &Publisher{ch: ch, queueName: queueName}, nil
}

func (p *Publisher) PublishTaskCreated(ctx context.Context, ev events.TaskEvent) error {
	body, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	return p.ch.PublishWithContext(ctx, "", p.queueName, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

func (p *Publisher) Close() error {
	if p.ch == nil {
		return nil
	}
	return p.ch.Close()
}

func DialPublisher(url, queueName string) (*Publisher, *amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, nil, fmt.Errorf("channel: %w", err)
	}
	pub, err := New(ch, queueName)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, nil, err
	}
	return pub, conn, nil
}
