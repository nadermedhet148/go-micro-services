package repositories

import (
	"os"
	"time"

	entity "github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/config"
	_ "github.com/joho/godotenv/autoload"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Save(entity entity.Payment) (int, error)
}

type PaymentDatabase struct {
	connection *gorm.DB
}

func NewPaymentRepository() PaymentRepository {
	db, err := config.ConnectDB()
	if err != nil {
		panic("Failed to connect database")
	}
	if os.Getenv("DB_HOST") != "" {
		db.AutoMigrate(&entity.Payment{})
	}
	return &PaymentDatabase{
		connection: db,
	}
}

func (db *PaymentDatabase) Save(entity entity.Payment) (int, error) {
	data := &entity
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	err := db.connection.Save(data)
	if err.Error != nil {
		return 0, err.Error
	}
	return data.ID, nil
}
