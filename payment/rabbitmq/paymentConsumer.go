package rabbitmq

import (
	"encoding/json"
	"log"
	"os"

	"github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/repositories"
	"github.com/streadway/amqp"
)

type PaymentConsumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func NewPaymentConsumer() (*PaymentConsumer, error) {
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

	return &PaymentConsumer{
		connection: conn,
		channel:    ch,
		queue:      q,
	}, nil
}

func (c *PaymentConsumer) ConsumePayments() error {
	msgs, err := c.channel.Consume(
		c.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	PaymentRepository := repositories.NewPaymentRepository()
	consumerChan := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a payment message: %s", d.Body)
			var ticket entity.Ticket
			err := json.Unmarshal(d.Body, &ticket)
			if err != nil {
				log.Printf("Failed to unmarshal ticket: %v", err)
				continue
			}

			Payment := entity.Payment{
				REF_NUMBER: ticket.REF_NUMBER,
				Status:     "pending",
			}

			id, err := PaymentRepository.Save(Payment)
			if err != nil {
				log.Printf("Failed to save payment: %v", err)
				continue
			}
			Payment.ID = id

			producer, err := NewNotificationProducer()

			if err != nil {
				log.Printf("Failed to create a producer: %v", err)
				return
			}
			paymentJson, err := json.Marshal(Payment)
			if err != nil {
				log.Printf("Failed to marshal payment to JSON: %v", err)
				return
			}

			err = producer.PublishNotification(string(paymentJson))
			if err != nil {
				log.Printf("Failed to publish payment: %v", err)
				return
			}
			return
		}
	}()

	log.Printf("Waiting for payment messages. To exit press CTRL+C")
	<-consumerChan

	return nil
}

func (c *PaymentConsumer) Close() {
	c.channel.Close()
	c.connection.Close()
}
