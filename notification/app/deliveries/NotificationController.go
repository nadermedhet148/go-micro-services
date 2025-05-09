package deliveries

import (
	"net/http"

	"github.com/coroo/go-starter/app/entity"
	usecases "github.com/coroo/go-starter/app/usecases"
	"github.com/gin-gonic/gin"
)

type NotificationController interface {
	Save(ctx *gin.Context) error
}

type notificationController struct {
	usecases usecases.NotificationService
}

func NewNotificationController(router *gin.Engine, apiPrefix string, NotificationService usecases.NotificationService) {
	handlerNotification := &notificationController{
		usecases: NotificationService,
	}
	NotificationGroup := router.Group(apiPrefix + "Notification")
	{
		NotificationGroup.POST("", handlerNotification.NotificationCreate)

	}
}

func (deliveries *notificationController) NotificationCreate(c *gin.Context) {
	var NotificationEntity entity.Notification
	c.ShouldBindJSON(&NotificationEntity)

	NotificationPK, err := deliveries.usecases.SaveNotification(NotificationEntity)
	if err != nil {
		c.JSON(http.StatusConflict, err)
	} else {
		NotificationEntity.ID = NotificationPK
		c.JSON(http.StatusOK, NotificationEntity)
	}
}
