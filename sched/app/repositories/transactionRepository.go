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
	GePendingTrxs(bucket int) []entity.Transaction
}

type TransactionDatabase struct {
}

func NewTransactionRepository() TransactionRepository {
	return &TransactionDatabase{}
}

func getDBConnection(bucket int) *gorm.DB {
	db, err := config.GetDBForKey(bucket)
	if err != nil {
		panic("Failed to connect database")
	}
	db.AutoMigrate(&entity.Transaction{})
	return db
}
func (db *TransactionDatabase) GePendingTrxs(bucket int) []entity.Transaction {
	conn := getDBConnection(bucket)

	var transactions []entity.Transaction
	result := conn.Where("status = ? AND created_at >= ?", "pending", time.Now().Add(-10*time.Minute)).Find(&transactions)
	if result.Error != nil {
		return []entity.Transaction{}
	}

	return transactions
}
