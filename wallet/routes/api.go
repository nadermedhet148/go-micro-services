package routes

import (
	"os"

	"github.com/coroo/go-starter/app/deliveries"
	"github.com/coroo/go-starter/app/rabbitmq"
	"github.com/coroo/go-starter/app/repositories"
	"github.com/coroo/go-starter/app/services"
	"github.com/gin-gonic/gin"
)

func Api() {
	router := gin.Default()

	API_PREFIX := os.Getenv("API_PREFIX")

	router.GET("/", func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": os.Getenv("MAIN_DESCRIPTION"),
		})
	})
	PaymentProducer, err := rabbitmq.NewPaymentProducer()
	if err != nil {
		panic("Failed to initialize PaymentProducer: " + err.Error())
	}
	var (
		WalletRepository = repositories.NewWalletRepository()

		WalletService = services.NewWalletService(WalletRepository, PaymentProducer)
	)
	deliveries.NewWalletController(router, API_PREFIX, WalletService)

	router.Run(":" + os.Getenv("MAIN_PORT"))
}
