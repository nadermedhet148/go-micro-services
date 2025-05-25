package repositories

import (
	"time"

	entity "github.com/coroo/go-starter/app/entity"
	"github.com/coroo/go-starter/config"
	_ "github.com/joho/godotenv/autoload"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Save(key string, Transaction entity.Transaction) (int, error)
	Update(key string, Transaction entity.Transaction) error
	Delete(key string, Transaction entity.Transaction) error
	GetByRefNumber(key string, refNumber string) entity.Transaction
}

type TransactionDatabase struct {
}

func NewTransactionRepository() TransactionRepository {
	return &TransactionDatabase{}
}

func getDBConnection(key string) *gorm.DB {
	db, err := config.GetDBForKey(key)
	if err != nil {
		panic("Failed to connect database")
	}
	db.AutoMigrate(&entity.Transaction{})
	return db
}

func (db *TransactionDatabase) Save(key string, Transaction entity.Transaction) (int, error) {
	conn := getDBConnection(key)

	data := &Transaction
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	result := conn.Create(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return data.ID, nil
}

func (db *TransactionDatabase) Update(key string, Transaction entity.Transaction) error {
	conn := getDBConnection(key)

	data := &Transaction
	data.UpdatedAt = time.Now()
	result := conn.Save(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *TransactionDatabase) Delete(key string, Transaction entity.Transaction) error {
	conn := getDBConnection(key)
	data := &Transaction
	result := conn.Delete(data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *TransactionDatabase) GetByRefNumber(key string, refNumber string) entity.Transaction {
	conn := getDBConnection(key)
	var transaction entity.Transaction
	conn.Where("ref_number = ?", refNumber).First(&transaction)
	return transaction
}
