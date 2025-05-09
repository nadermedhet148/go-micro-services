package routes

import (
	"os"

	"github.com/coroo/go-starter/app/deliveries"
	"github.com/coroo/go-starter/app/repositories"
	"github.com/coroo/go-starter/app/usecases"
	"github.com/gin-gonic/gin"
)

func Api() {
	router := gin.Default()
	// router.Use(middlewares.BasicAuth())

	API_PREFIX := os.Getenv("API_PREFIX")

	router.GET("/", func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": os.Getenv("MAIN_DESCRIPTION"),
		})
	})
	var (
		PaymentRepository repositories.PaymentRepository = repositories.NewPaymentRepository()
		PaymentService    usecases.PaymentService        = usecases.NewTicketService(PaymentRepository)
	)
	deliveries.NewPaymentController(router, API_PREFIX, PaymentService)

	router.Run(":" + os.Getenv("MAIN_PORT"))
}
