package repositories

import (
	"os"
	"time"

	entity "github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/config"
	"gorm.io/gorm/clause"

	// "github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AvailableSlotRepository interface {
	SaveAvailableSlot(paymentMethod entity.AvailableSlot) (int, error)
	UpdateAvailableSlot(paymentMethod entity.AvailableSlot) error
	DeleteAvailableSlot(paymentMethod entity.AvailableSlot) error
	GetAllAvailableSlots() []entity.AvailableSlot
	GetAvailableSlot(id string) []entity.AvailableSlot
}

type AvailableSlotDatabase struct {
	connection *gorm.DB
}

func NewAvailableSlotRepository() AvailableSlotRepository {
	db, err := config.ConnectDB()
	if err != nil {
		panic("Failed to connect database")
	}

	if os.Getenv("DB_HOST") != "" {
		db.AutoMigrate(&entity.AvailableSlot{})
	} else {
		db.AutoMigrate(&entity.AvailableSlotTesting{})
	}
	return &AvailableSlotDatabase{
		connection: db,
	}
}

func (db *AvailableSlotDatabase) SaveAvailableSlot(AvailableSlot entity.AvailableSlot) (int, error) {
	data := &AvailableSlot
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	err := db.connection.Create(data)
	if err.Error != nil {
		return 0, err.Error
	}
	return data.ID, nil
}

func (db *AvailableSlotDatabase) UpdateAvailableSlot(AvailableSlot entity.AvailableSlot) error {
	data := &AvailableSlot
	data.UpdatedAt = time.Now()
	db.connection.Save(data)
	return nil
}

func (db *AvailableSlotDatabase) DeleteAvailableSlot(AvailableSlot entity.AvailableSlot) error {
	db.connection.Delete(&AvailableSlot)
	return nil
}

func (db *AvailableSlotDatabase) GetAllAvailableSlots() []entity.AvailableSlot {
	var AvailableSlots []entity.AvailableSlot
	db.connection.Preload(clause.Associations).Find(&AvailableSlots)
	return AvailableSlots
}

func (db *AvailableSlotDatabase) GetAvailableSlot(id string) []entity.AvailableSlot {
	var AvailableSlot []entity.AvailableSlot
	db.connection.Preload(clause.Associations).Where("id = ?", id).First(&AvailableSlot)
	return AvailableSlot
}
