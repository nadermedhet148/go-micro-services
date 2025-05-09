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

type NotificationRepository interface {
	Save(entity entity.Notification) (int, error)
}

type NotificationDatabase struct {
	connection *gorm.DB
}

func NewNotificationRepository() NotificationRepository {
	db, err := config.ConnectDB()
	if err != nil {
		panic("Failed to connect database")
	}
	if os.Getenv("DB_HOST") != "" {
		db.AutoMigrate(&entity.Notification{})
	}
	return &NotificationDatabase{
		connection: db,
	}
}

func (db *NotificationDatabase) Save(Notification entity.Notification) (int, error) {
	data := &Notification
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	err := db.connection.Create(data)
	if err.Error != nil {
		return 0, err.Error
	}
	return data.ID, nil
}
