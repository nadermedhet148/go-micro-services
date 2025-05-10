package repositories

import (
	"time"

	entity "github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/config"
	_ "github.com/joho/godotenv/autoload"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type WalletRepository interface {
	Save(wallet entity.Wallet) (int, error)
	Update(wallet entity.Wallet) error
	Delete(wallet entity.Wallet) error
	Get(id int) entity.Wallet
}

type WalletDatabase struct {
	connection *gorm.DB
}

func NewWalletRepository() WalletRepository {
	db, err := config.ConnectDB()
	if err != nil {
		panic("Failed to connect database")
	}
	db.AutoMigrate(&entity.Wallet{})

	return &WalletDatabase{
		connection: db,
	}
}

func (db *WalletDatabase) Save(Wallet entity.Wallet) (int, error) {
	data := &Wallet
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	err := db.connection.Create(data)
	if err.Error != nil {
		return 0, err.Error
	}
	return data.ID, nil
}

func (db *WalletDatabase) Update(Wallet entity.Wallet) error {
	data := &Wallet
	data.UpdatedAt = time.Now()
	err := db.connection.Save(data)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (db *WalletDatabase) Delete(Wallet entity.Wallet) error {
	data := &Wallet
	err := db.connection.Delete(data)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (db *WalletDatabase) Get(id int) entity.Wallet {
	var Wallet entity.Wallet
	err := db.connection.Where("id = ?", id).First(&Wallet)
	if err.Error != nil {
		return Wallet
	}
	return Wallet
}
