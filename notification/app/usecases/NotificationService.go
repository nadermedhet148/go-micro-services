package usecases

import (
	"github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/repositories"
)

type NotificationService interface {
	SaveNotification(Notification entity.Notification) (int, error)
}

type notificationService struct {
	repositories repositories.NotificationRepository
}

func NewTicketService(repository repositories.NotificationRepository) NotificationService {
	return &notificationService{
		repositories: repository,
	}
}

func (usecases *notificationService) SaveNotification(Notification entity.Notification) (int, error) {
	return usecases.repositories.Save(Notification)
}
