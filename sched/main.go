package main

import (
	"net/http"

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

	go func() {
		c := cron.New(cron.WithSeconds())
		c.AddFunc("0/5 * * * * *", func() {
			print("Hello, World! This message is printed every 5 seconds.\n")
		})
		c.Start()
	}()

	select { // Keep the main goroutine running
	}
}
