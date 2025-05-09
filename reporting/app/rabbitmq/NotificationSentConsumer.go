package rabbitmq

import (
	"encoding/json"
	"log"
	"os"

	"github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/repositories"
	"github.com/streadway/amqp"
)

type PaymentStatusConsumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func NewPaymentStatusConsumer() (*PaymentStatusConsumer, error) {
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

	return &PaymentStatusConsumer{
		connection: conn,
		channel:    ch,
		queue:      q,
	}, nil
}

func (c *PaymentStatusConsumer) ConsumePaymentStatusEvents() error {
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
	ticketRepository := repositories.NewTicketRepository()
	consumerChan := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a notifcation message: %s", d.Body)
			var payment entity.Payment
			err := json.Unmarshal(d.Body, &payment)
			if err != nil {
				log.Printf("Failed to unmarshal payment: %v", err)
				continue
			}

			ticket := ticketRepository.GetTicket(payment.REF_NUMBER)

			if payment.Status == "success" {
				ticket.Status = "success"
			} else {
				ticket.Status = "failed"
			}

			ticketRepository.UpdateTicket(ticket)
			return
		}
	}()

	log.Printf("Waiting for payment messages. To exit press CTRL+C")
	<-consumerChan

	return nil
}

func (c *PaymentStatusConsumer) Close() {
	c.channel.Close()
	c.connection.Close()
}
