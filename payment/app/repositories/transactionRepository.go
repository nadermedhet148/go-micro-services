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
	Save(Transaction entity.Transaction) (int, error)
	Update(Transaction entity.Transaction) error
	Delete(Transaction entity.Transaction) error
	GetByRefNumber(refNumber string) entity.Transaction
}

type TransactionDatabase struct {
	connection *gorm.DB
}

func NewTransactionRepository() TransactionRepository {
	db, err := config.ConnectDB()
	if err != nil {
		panic("Failed to connect database")
	}
	db.AutoMigrate(&entity.Transaction{})

	return &TransactionDatabase{
		connection: db,
	}
}

func (db *TransactionDatabase) Save(Transaction entity.Transaction) (int, error) {
	data := &Transaction
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	err := db.connection.Create(data)
	if err.Error != nil {
		return 0, err.Error
	}
	return data.ID, nil
}

func (db *TransactionDatabase) Update(Transaction entity.Transaction) error {
	data := &Transaction
	data.UpdatedAt = time.Now()
	err := db.connection.Save(data)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (db *TransactionDatabase) Delete(Transaction entity.Transaction) error {
	data := &Transaction
	err := db.connection.Delete(data)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (db *TransactionDatabase) GetByRefNumber(refNumber string) entity.Transaction {
	var transaction entity.Transaction
	db.connection.Where("ref_number = ?", refNumber).First(&transaction)
	return transaction
}
