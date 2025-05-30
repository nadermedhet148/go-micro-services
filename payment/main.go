package main

import (
	"net/http"

	"github.com/coroo/go-starter/app/rabbitmq"
	"github.com/coroo/go-starter/routes"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	consumer, err := rabbitmq.NewPaymentConsumer()
	if err != nil {
		panic(err)
	}
	go consumer.ConsumePaymentEvents()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()
	routes.Api()

}
