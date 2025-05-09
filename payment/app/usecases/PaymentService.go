package usecases

import (
	"github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/app/externalServices"
	"github.com/coroo/go-starter/app/repositories"
)

type PaymentService interface {
	SavePayment(Payment entity.Payment) (int, error)
}

type paymentService struct {
	repositories repositories.PaymentRepository
}

func NewTicketService(repository repositories.PaymentRepository) PaymentService {
	return &paymentService{
		repositories: repository,
	}
}

func (usecases *paymentService) SavePayment(Payment entity.Payment) (int, error) {
	Payment.Status = "pending"
	id, err := usecases.repositories.Save(Payment)
	Payment.ID = id
	if err != nil {
		return 0, err
	}
	notificationService := externalServices.NotificationService{}

	res, err := notificationService.SendNotification(externalServices.NotificationRequest{
		REF_NUMBER: Payment.REF_NUMBER,
	})

	if res.Status != "success" || err != nil {
		Payment.Status = "failed"
		usecases.repositories.Save(Payment)
		return 0, err
	}
	Payment.Status = "success"
	return usecases.repositories.Save(Payment)
}
