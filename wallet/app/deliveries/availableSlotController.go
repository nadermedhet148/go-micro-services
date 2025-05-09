package deliveries

import (
	"net/http"

	"github.com/coroo/go-starter/app/entity"
	usecases "github.com/coroo/go-starter/app/usecases"

	// "github.com/coroo/go-starter/app/middlewares"

	"github.com/gin-gonic/gin"
)

type AvailableSlotController interface {
	GetAvailableSlots(ctx *gin.Context)
	Save(ctx *gin.Context) error
}

type availableSlotController struct {
	usecases usecases.AvailableSlotService
}

func NewAvailableSlotController(router *gin.Engine, apiPrefix string, AvailableSlotService usecases.AvailableSlotService) {
	handlerAvailableSlot := &availableSlotController{
		usecases: AvailableSlotService,
	}
	AvailableSlotsGroup := router.Group(apiPrefix + "AvailableSlot")
	{
		AvailableSlotsGroup.GET("", handlerAvailableSlot.AvailableSlotsIndex)
		AvailableSlotsGroup.POST("", handlerAvailableSlot.AvailableSlotCreate)
	}
}

// GetAvailableSlotsIndex godoc
// @Security basicAuth
// @Summary Show all existing AvailableSlots
// @Description Get all existing AvailableSlots
// @Tags AvailableSlots
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.AvailableSlot
// @Failure 401 {object} dto.Response
// @Router /AvailableSlot/index [get]
func (deliveries *availableSlotController) AvailableSlotsIndex(c *gin.Context) {
	AvailableSlots := deliveries.usecases.GetAllAvailableSlots()
	c.JSON(http.StatusOK, gin.H{"data": AvailableSlots})
}

// CreateAvailableSlots godoc
// @Security basicAuth
// @Summary Create new AvailableSlots
// @Description Create a new AvailableSlots
// @Tags AvailableSlots
// @Accept  json
// @Produce  json
// @Param AvailableSlot body entity.AvailableSlot true "Create AvailableSlot"
// @Success 200 {object} entity.AvailableSlot
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /AvailableSlot/create [post]
func (deliveries *availableSlotController) AvailableSlotCreate(c *gin.Context) {
	var AvailableSlotEntity entity.AvailableSlot
	c.ShouldBindJSON(&AvailableSlotEntity)
	AvailableSlotPK, err := deliveries.usecases.SaveAvailableSlot(AvailableSlotEntity)
	if err != nil {
		c.JSON(http.StatusConflict, err)
	} else {
		AvailableSlotEntity.ID = AvailableSlotPK
		c.JSON(http.StatusOK, AvailableSlotEntity)
	}
}
