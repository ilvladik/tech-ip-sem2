package rabbit

import amqp "github.com/rabbitmq/amqp091-go"

const (
	QueueJobs    = "task_jobs"
	QueueJobsDLQ = "task_jobs_dlq"
)

func DeclareQueues(ch *amqp.Channel) error {
	if _, err := ch.QueueDeclare(QueueJobsDLQ, true, false, false, false, nil); err != nil {
		return err
	}
	args := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": QueueJobsDLQ,
	}
	_, err := ch.QueueDeclare(QueueJobs, true, false, false, false, args)
	return err
}
