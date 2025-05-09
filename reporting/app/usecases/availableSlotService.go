package usecases

import (
	entity "github.com/coroo/go-starter/app/entity"
	repositories "github.com/coroo/go-starter/app/repositories"
	_ "github.com/go-sql-driver/mysql"
)

type AvailableSlotService interface {
	SaveAvailableSlot(entity.AvailableSlot) (int, error)
	UpdateAvailableSlot(entity.AvailableSlot) error
	DeleteAvailableSlot(entity.AvailableSlot) error
	GetAllAvailableSlots() []entity.AvailableSlot
	GetAvailableSlot(id string) []entity.AvailableSlot
}

type availableSlotService struct {
	repositories repositories.AvailableSlotRepository
}

func NewAvailableSlotService(repository repositories.AvailableSlotRepository) AvailableSlotService {
	return &availableSlotService{
		repositories: repository,
	}
}

func (usecases *availableSlotService) GetAllAvailableSlots() []entity.AvailableSlot {
	return usecases.repositories.GetAllAvailableSlots()
}

func (usecases *availableSlotService) GetAvailableSlot(id string) []entity.AvailableSlot {
	return usecases.repositories.GetAvailableSlot(id)
}

func (usecases *availableSlotService) SaveAvailableSlot(AvailableSlot entity.AvailableSlot) (int, error) {
	return usecases.repositories.SaveAvailableSlot(AvailableSlot)
}

func (usecases *availableSlotService) UpdateAvailableSlot(AvailableSlot entity.AvailableSlot) error {
	return usecases.repositories.UpdateAvailableSlot(AvailableSlot)
}

func (usecases *availableSlotService) DeleteAvailableSlot(AvailableSlot entity.AvailableSlot) error {
	return usecases.repositories.DeleteAvailableSlot(AvailableSlot)
}
