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
	Save(key string, wallet entity.Wallet) (int, error)
	Update(key string, wallet entity.Wallet) error
	Delete(key string, wallet entity.Wallet) error
	Get(key string, id int) entity.Wallet
}

type WalletDatabase struct {
}

func NewWalletRepository() WalletRepository {
	return &WalletDatabase{}
}

func getDBConnection(key string) *gorm.DB {
	db, err := config.GetDBForKey(key)
	if err != nil {
		panic("Failed to connect database")
	}
	db.AutoMigrate(&entity.Wallet{})
	return db
}

func (db *WalletDatabase) Save(key string, Wallet entity.Wallet) (int, error) {
	conn := getDBConnection(key)

	data := &Wallet
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	err := conn.Create(data)
	if err.Error != nil {
		return 0, err.Error
	}
	return data.ID, nil
}

func (db *WalletDatabase) Update(key string, Wallet entity.Wallet) error {
	conn := getDBConnection(key)

	data := &Wallet
	data.UpdatedAt = time.Now()
	err := conn.Save(data)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (db *WalletDatabase) Delete(key string, Wallet entity.Wallet) error {
	conn := getDBConnection(key)

	data := &Wallet
	err := conn.Delete(data)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (db *WalletDatabase) Get(key string, id int) entity.Wallet {
	conn := getDBConnection(key)

	var Wallet entity.Wallet
	err := conn.Where("id = ?", id).First(&Wallet)
	if err.Error != nil {
		return Wallet
	}
	return Wallet
}
