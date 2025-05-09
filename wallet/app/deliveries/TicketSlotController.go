package deliveries

import (
	"net/http"

	"github.com/coroo/go-starter/app/entity"
	usecases "github.com/coroo/go-starter/app/usecases"

	// "github.com/coroo/go-starter/app/middlewares"

	"github.com/gin-gonic/gin"
)

type TicketController interface {
	GetTickets(ctx *gin.Context)
	Save(ctx *gin.Context) error
}

type ticketController struct {
	usecases usecases.TicketService
}

func NewTicketController(router *gin.Engine, apiPrefix string, TicketService usecases.TicketService) {
	handlerTicket := &ticketController{
		usecases: TicketService,
	}
	TicketsGroup := router.Group(apiPrefix + "Ticket")
	{
		TicketsGroup.GET("", handlerTicket.TicketsIndex)
		TicketsGroup.POST("", handlerTicket.TicketCreate)
		TicketsGroup.POST("/lock", handlerTicket.TicketCreateWithLock)
		TicketsGroup.POST("/d-lock", handlerTicket.TicketCreateWithDLock)
		TicketsGroup.POST("/http-transaction", handlerTicket.TicketCreateWithHttpTransaction)
		TicketsGroup.POST("/http-mq", handlerTicket.TicketCreateWithRabbitMq)

	}
}

// GetTicketsIndex godoc
// @Security basicAuth
// @Summary Show all existing Tickets
// @Description Get all existing Tickets
// @Tags Tickets
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.Ticket
// @Failure 401 {object} dto.Response
// @Router /ticket/index [get]
func (deliveries *ticketController) TicketsIndex(c *gin.Context) {
	Tickets := deliveries.usecases.GetAllTickets()
	c.JSON(http.StatusOK, gin.H{"data": Tickets})
}

// CreateTickets godoc
// @Security basicAuth
// @Summary Create new Tickets
// @Description Create a new Tickets
// @Tags Tickets
// @Accept  json
// @Produce  json
// @Param Ticket body entity.Ticket true "Create Ticket"
// @Success 200 {object} entity.Ticket
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /ticket/create [post]
func (deliveries *ticketController) TicketCreate(c *gin.Context) {
	var TicketEntity entity.Ticket
	c.ShouldBindJSON(&TicketEntity)

	TicketPK, err := deliveries.usecases.SaveTicket(TicketEntity)
	if err != nil {
		c.JSON(http.StatusConflict, err)
	} else {
		TicketEntity.ID = TicketPK
		c.JSON(http.StatusOK, TicketEntity)
	}
}

func (deliveries *ticketController) TicketCreateWithLock(c *gin.Context) {

	var TicketEntity entity.Ticket
	c.ShouldBindJSON(&TicketEntity)

	TicketPK, err := deliveries.usecases.SaveTicketWithLock(TicketEntity)
	if err != nil {
		c.JSON(http.StatusConflict, err)
	} else {
		TicketEntity.ID = TicketPK
		c.JSON(http.StatusOK, TicketEntity)
	}
}

func (deliveries *ticketController) TicketCreateWithDLock(c *gin.Context) {

	var TicketEntity entity.Ticket
	c.ShouldBindJSON(&TicketEntity)

	TicketPK, err := deliveries.usecases.SaveTicketWithDLock(TicketEntity)
	if err != nil {
		c.JSON(http.StatusConflict, err)
	} else {
		TicketEntity.ID = TicketPK
		c.JSON(http.StatusOK, TicketEntity)
	}
}

func (deliveries *ticketController) TicketCreateWithHttpTransaction(c *gin.Context) {

	var TicketEntity entity.Ticket
	c.ShouldBindJSON(&TicketEntity)

	TicketPK, err := deliveries.usecases.TicketCreateWithHttpTransaction(TicketEntity)
	if err != nil {
		c.JSON(http.StatusConflict, err)
	} else {
		TicketEntity.ID = TicketPK
		c.JSON(http.StatusOK, TicketEntity)
	}
}

func (deliveries *ticketController) TicketCreateWithRabbitMq(c *gin.Context) {

	var TicketEntity entity.Ticket
	c.ShouldBindJSON(&TicketEntity)

	TicketPK, err := deliveries.usecases.TicketCreateWithRabbitMq(TicketEntity)
	if err != nil {
		c.JSON(http.StatusConflict, err)
	} else {
		TicketEntity.ID = TicketPK
		c.JSON(http.StatusOK, TicketEntity)
	}
}
