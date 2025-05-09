package deliveries

import (
	"net/http"

	"github.com/coroo/go-starter/app/entity"
	usecases "github.com/coroo/go-starter/app/usecases"
	"github.com/gin-gonic/gin"
)

type PaymentController interface {
	Save(ctx *gin.Context) error
}

type paymentController struct {
	usecases usecases.PaymentService
}

func NewPaymentController(router *gin.Engine, apiPrefix string, PaymentService usecases.PaymentService) {
	handlerPayment := &paymentController{
		usecases: PaymentService,
	}
	PaymentGroup := router.Group(apiPrefix + "Payment")
	{
		PaymentGroup.POST("", handlerPayment.PaymentCreate)

	}
}

func (deliveries *paymentController) PaymentCreate(c *gin.Context) {
	var PaymentEntity entity.Payment
	c.ShouldBindJSON(&PaymentEntity)

	PaymentPK, err := deliveries.usecases.SavePayment(PaymentEntity)
	if err != nil {
		c.JSON(http.StatusConflict, err)
	} else {
		PaymentEntity.ID = PaymentPK
		c.JSON(http.StatusOK, PaymentEntity)
	}
}
