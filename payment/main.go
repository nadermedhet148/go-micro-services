package main

import (
	"github.com/coroo/go-starter/rabbitmq"
	"github.com/coroo/go-starter/routes"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	consumer, err := rabbitmq.NewPaymentConsumer()
	if err != nil {
		panic(err)
	}
	go consumer.ConsumePayments()

	notificationConsumer, err := rabbitmq.NewNotificationSentConsumer()
	if err != nil {
		panic(err)
	}

	go notificationConsumer.ConsumeNotificationSentEvents()
	routes.Api()

}
