package rabbitmq

import (
	"encoding/json"
	"log"
	"os"

	"github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/kafka"
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

	return &PaymentConsumer{
		connection: conn,
		channel:    ch,
		queue:      q,
	}, nil
}

func (c *PaymentConsumer) ConsumePaymentEvents() error {
	err := c.channel.Qos(
		100,   // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		c.queue.Name,
		"",
		false, // autoAck set to false for manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	transactionRepo := repositories.NewTransactionRepository()
	eventProducer, err := kafka.NewPaymentEventProducer()
	transactionService := services.NewTransactionService(transactionRepo, eventProducer)
	consumerChan := make(chan bool)

	workerCount := 20 // or set to your desired parallelism, should not exceed prefetch count

	for i := 0; i < workerCount; i++ {
		go func() {
			for d := range msgs {
				log.Printf("Received a payment message: %s", d.Body)
				var trxRequest entity.WalletRechargeRequest
				err := json.Unmarshal(d.Body, &trxRequest)
				if err != nil {
					log.Printf("Failed to unmarshal payment: %v", err)
					d.Nack(false, false) // reject message, don't requeue
					continue
				}

				_, err = transactionService.CerateRechargeTransaction(trxRequest)
				if err != nil {
					log.Printf("Failed to process transaction: %v", err)
					d.Nack(false, true) // reject message, requeue
					continue
				}

				d.Ack(false) // manually acknowledge message
			}
		}()
	}

	log.Printf("Waiting for payment messages. To exit press CTRL+C")
	<-consumerChan

	return nil
}

func (c *PaymentConsumer) Close() {
	c.channel.Close()
	c.connection.Close()
}
