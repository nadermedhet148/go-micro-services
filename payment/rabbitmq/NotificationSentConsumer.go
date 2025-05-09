package rabbitmq

import (
	"encoding/json"
	"log"
	"os"

	"github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/repositories"
	"github.com/streadway/amqp"
)

type NotificationSentConsumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func NewNotificationSentConsumer() (*NotificationSentConsumer, error) {
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

	return &NotificationSentConsumer{
		connection: conn,
		channel:    ch,
		queue:      q,
	}, nil
}

func (c *NotificationSentConsumer) ConsumeNotificationSentEvents() error {
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
			log.Printf("Received a notifcation message: %s", d.Body)
			var payment entity.Payment
			err := json.Unmarshal(d.Body, &payment)
			if err != nil {
				log.Printf("Failed to unmarshal payment: %v", err)
				continue
			}

			payment.Status = "success"
			PaymentRepository.Save(payment)
			producer, err := NewPaymentStatusProducer()

			if err != nil {
				log.Printf("Failed to create a producer: %v", err)
				return
			}
			paymentJson, err := json.Marshal(payment)
			if err != nil {
				log.Printf("Failed to marshal payment to JSON: %v", err)
				return
			}

			err = producer.PublishPaymentStatus(string(paymentJson))
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

func (c *NotificationSentConsumer) Close() {
	c.channel.Close()
	c.connection.Close()
}
