package routes

import (
	"os"

	"github.com/coroo/go-starter/app/deliveries"
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
	var (
		TransactionRepository = repositories.NewTransactionRepository()

		transactionService = services.NewTransactionService(TransactionRepository)
	)
	deliveries.NewTransactionsController(router, API_PREFIX, transactionService)

	router.Run(":" + os.Getenv("MAIN_PORT"))
}
