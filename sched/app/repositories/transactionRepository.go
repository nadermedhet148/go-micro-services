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
	GePendingTrx() entity.Transaction
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
func (db *TransactionDatabase) GePendingTrx() entity.Transaction {
	conn := getDBConnection("transaction")

	var transaction entity.Transaction
	result := conn.Where("status = ? AND created_at >= ?", "pending", time.Now().Add(-10*time.Minute)).First(&transaction)
	if result.Error != nil {
		return entity.Transaction{}
	}

	transaction.UpdatedAt = time.Now()
	conn.Save(&transaction)

	return transaction
}
