package rabbitmq

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

type PaymentProducer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func NewPaymentProducer() (*PaymentProducer, error) {
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		os.Getenv("PAYMENT_QUEUE"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &PaymentProducer{
		connection: conn,
		channel:    ch,
		queue:      q,
	}, nil
}

func (p *PaymentProducer) PublishPayment(payment string) error {
	err := p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(payment),
		},
	)
	if err != nil {
		return err
	}
	log.Printf("Payment message sent: %s", payment)
	return nil
}

func (p *PaymentProducer) Close() {
	p.channel.Close()
	p.connection.Close()
}
