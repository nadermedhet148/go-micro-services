package rabbitmq

import (
	"encoding/json"
	"log"
	"os"

	"github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/repositories"
	"github.com/coroo/go-starter/app/services"
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

func (c *PaymentConsumer) ConsumePaymentEvents() error {
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
	transactionRepo := repositories.NewTransactionRepository()
	transactionService := services.NewTransactionService(transactionRepo)
	consumerChan := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a payment message: %s", d.Body)
			var trxRequest entity.WalletRechargeRequest
			err := json.Unmarshal(d.Body, &trxRequest)
			if err != nil {
				log.Printf("Failed to unmarshal payment: %v", err)
				continue
			}

			transactionService.CerateRechargeTransaction(trxRequest)
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
