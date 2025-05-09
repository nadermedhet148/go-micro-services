package rabbitmq

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

type PaymentStatusProducer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func NewPaymentStatusProducer() (*PaymentStatusProducer, error) {
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		os.Getenv("PAYMENT_STATUS_QUEUE"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &PaymentStatusProducer{
		connection: conn,
		channel:    ch,
		queue:      q,
	}, nil
}

func (p *PaymentStatusProducer) PublishPaymentStatus(PaymentStatus string) error {
	err := p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(PaymentStatus),
		},
	)
	if err != nil {
		return err
	}
	log.Printf("PaymentStatus message sent: %s", PaymentStatus)
	return nil
}

func (p *PaymentStatusProducer) Close() {
	p.channel.Close()
	p.connection.Close()
}
