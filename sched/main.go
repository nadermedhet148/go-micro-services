package main

import (
	"net/http"

	"github.com/coroo/go-starter/app/rabbitmq"
	"github.com/coroo/go-starter/app/repositories"
	"github.com/coroo/go-starter/app/services"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	paymentProducer, err := rabbitmq.NewPaymentProducer()
	if err != nil {
		panic("Failed to initialize PaymentProducer: " + err.Error())
	}
	transactionRepo := repositories.NewTransactionRepository()
	transactionService := services.NewTransactionService(transactionRepo, *paymentProducer)

	go func() {
		c := cron.New(cron.WithSeconds())
		c.AddFunc("0/40 * * * * *", func() {
			transactionService.RunExpiredTransactionCleanup()
		})
		c.Start()
	}()

	select { // Keep the main goroutine running
	}
}
