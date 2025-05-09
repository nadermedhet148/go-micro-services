package main

import (
	"github.com/coroo/go-starter/rabbitmq"
	"github.com/coroo/go-starter/routes"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	paymentConsumer, err := rabbitmq.NewPaymentConsumer()
	if err != nil {
		panic(err)
	}
	go paymentConsumer.ConsumePayments()

	routes.Api()
}
