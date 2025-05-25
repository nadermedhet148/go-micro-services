package entity

import (
	"time"
)

type Transaction struct {
	ID         int       `gorm:"type:BIGINT UNSIGNED NOT NULL AUTO_INCREMENT" json:"id"`
	REF_NUMBER string    `gorm:"type:VARCHAR(255) NOT NULL" json:"ref_number"`
	AMOUNT     float64   `gorm:"type:DECIMAL(10,2) NOT NULL" json:"amount"`
	STATUS     string    `gorm:"type:ENUM('PENDING', 'COMPLETED', 'FAILED') NOT NULL" json:"status"`
	WALLET_ID  int       `gorm:"type:BIGINT UNSIGNED NOT NULL" json:"wallet_id"`
	REGION     string    `gorm:"type:VARCHAR(255) NOT NULL" json:"region"`
	TYPE       string    `gorm:"type:VARCHAR(255) NOT NULL" json:"type"` // "debit" or "credit"
	CreatedAt  time.Time `gorm:"type:TIMESTAMP DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:TIMESTAMP DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Transaction *Transaction) TableName() string {
	return "transactions"
}
