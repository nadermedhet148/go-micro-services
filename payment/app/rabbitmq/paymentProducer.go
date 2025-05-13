package rabbitmq

import (
	"encoding/json"
	"log"
	"os"

	"github.com/coroo/go-starter/app/entity"
	"github.com/streadway/amqp"
)

type NotificationProducer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func NewNotificationProducer() (*NotificationProducer, error) {
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		os.Getenv("Notification_QUEUE"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &NotificationProducer{
		connection: conn,
		channel:    ch,
		queue:      q,
	}, nil
}

func (p *NotificationProducer) PublishNotification(request entity.WalletRechargeRequest) error {

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
	log.Printf("Notification message sent: %s", requestString)
	return nil
}

func (p *NotificationProducer) Close() {
	p.channel.Close()
	p.connection.Close()
}
