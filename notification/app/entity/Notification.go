package entity

import (
	"time"
)

type Notification struct {
	ID         int       `gorm:"type:BIGINT UNSIGNED NOT NULL AUTO_INCREMENT" json:"id"`
	REF_NUMBER string    `gorm:"type:VARCHAR(191) NOT NULL" json:"ref_number"`
	TEXT       string    `gorm:"type:VARCHAR(500) NOT NULL" json:"text"`
	Status     string    `gorm:"type:VARCHAR(100) NOT NULL" json:"status"`
	CreatedAt  time.Time `gorm:"type:TIMESTAMP DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:TIMESTAMP DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
}
