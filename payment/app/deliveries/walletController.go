package deliveries

import (
	"net/http"

	"github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/services"
	"github.com/gin-gonic/gin"
)

type WalletController interface {
	GetWallet(ctx *gin.Context)
	Save(ctx *gin.Context) error
}

type walletController struct {
	service services.WalletService
}

func NewWalletController(router *gin.Engine, apiPrefix string, WalletService services.WalletService) {
	handler := &walletController{
		service: WalletService,
	}
	Group := router.Group(apiPrefix + "wallets")
	{
		Group.POST("", handler.CreateWallet)
	}
}

func (deliveries *walletController) CreateWallet(c *gin.Context) {
	var wallet entity.Wallet
	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := deliveries.service.CerateWallet(wallet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (deliveries *walletController) RechargeWallet(c *gin.Context) {
	var req entity.WalletRechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := deliveries.service.RechargeWallet(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet recharged successfully"})
}
