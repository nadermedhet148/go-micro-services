package entity

import (
	"time"
)

type Wallet struct {
	ID        int       `gorm:"type:BIGINT UNSIGNED NOT NULL AUTO_INCREMENT" json:"id"`
	Region    string    `gorm:"type:VARCHAR(255) NOT NULL" json:"Region"`
	Name      string    `gorm:"type:VARCHAR(255) NOT NULL" json:"name"`
	Balance   float64   `gorm:"type:DECIMAL(10,2) NOT NULL" json:"balance"`
	UserID    int       `gorm:"type:BIGINT UNSIGNED NOT NULL" json:"user_id"`
	CreatedAt time.Time `gorm:"type:TIMESTAMP DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:TIMESTAMP DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
}

func (Wallet *Wallet) TableName() string {
	return "wallets"
}
