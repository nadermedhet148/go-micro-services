package deliveries

import (
	"net/http"

	"github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/services"
	"github.com/gin-gonic/gin"
)

type TransactionsController interface {
	TransactionCallBack(ctx *gin.Context) error
}

type transactionsController struct {
	service services.TransactionService
}

func NewTransactionsController(router *gin.Engine, apiPrefix string, transactionsService services.TransactionService) {
	handler := &transactionsController{
		service: transactionsService,
	}
	Group := router.Group(apiPrefix + "transactions")
	{
		Group.POST("", handler.TransactionCallBack)
	}
}

func (deliveries *transactionsController) TransactionCallBack(c *gin.Context) {
	var req entity.TransactionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := deliveries.service.UpdateTransaction(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transactions recharged successfully"})
}
