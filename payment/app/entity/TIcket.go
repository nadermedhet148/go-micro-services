package entity

import (
	"time"
)

type Ticket struct {
	ID         int       `gorm:"type:BIGINT UNSIGNED NOT NULL AUTO_INCREMENT" json:"id"`
	REF_NUMBER string    `gorm:"type:VARCHAR(191) NOT NULL" json:"ref_number"`
	SLOT_ID    int       `gorm:"type:BIGINT NOT NULL" json:"slot_id"`
	Status     string    `gorm:"type:VARCHAR(100) NOT NULL" json:"status"`
	CreatedAt  time.Time `gorm:"type:TIMESTAMP DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:TIMESTAMP DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
}
