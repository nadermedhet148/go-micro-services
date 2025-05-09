package rabbitmq

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

type NotificationSentProducer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func NewNotificationSentProducer() (*NotificationSentProducer, error) {
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		os.Getenv("PAYMENT_Notification_SENT_QUEUE"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &NotificationSentProducer{
		connection: conn,
		channel:    ch,
		queue:      q,
	}, nil
}

func (p *NotificationSentProducer) PublishNotificationSent(NotificationSentProducer string) error {
	err := p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(NotificationSentProducer),
		},
	)
	if err != nil {
		return err
	}
	log.Printf("NotificationSentProducer message sent: %s", NotificationSentProducer)
	return nil
}

func (p *NotificationSentProducer) Close() {
	p.channel.Close()
	p.connection.Close()
}
