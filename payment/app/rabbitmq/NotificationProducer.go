package rabbitmq

import (
	"encoding/json"
	"log"
	"os"

	"github.com/coroo/go-starter/app/entity"
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
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		amqp.Table{
			"x-queue-type": "quorum",
		},
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

func (p *PaymentProducer) PublishPayment(request entity.WalletRechargeRequest) error {

	requestString, err := json.Marshal(request)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(requestString),
		},
	)
	if err != nil {
		return err
	}
	log.Printf("Payment message sent: %s", requestString)
	return nil
}

func (p *PaymentProducer) Close() {
	p.channel.Close()
	p.connection.Close()
}
