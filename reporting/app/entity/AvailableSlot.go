package entity

import (
	"time"
)

type AvailableSlot struct {
	ID        int       `gorm:"type:BIGINT UNSIGNED NOT NULL AUTO_INCREMENT" json:"id"`
	LOCATION  string    `gorm:"type:VARCHAR(191) NOT NULL" json:"location"`
	EVENT_ID  int       `gorm:"type:BIGINT NOT NULL" json:"event_id"`
	CreatedAt time.Time `gorm:"type:TIMESTAMP DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:TIMESTAMP DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
}
type AvailableSlotTesting struct {
	ID        int       `gorm:"type:BIGINT UNSIGNED NOT NULL" json:"id"`
	LOCATION  string    `gorm:"type:VARCHAR(191) NOT NULL" json:"location"`
	EVENT_ID  int       `gorm:"type:BIGINT NOT NULL" json:"event_id"`
	CreatedAt time.Time `gorm:"type:TIMESTAMP DEFAULT CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:TIMESTAMP DEFAULT CURRENT_TIMESTAMP" json:"updated_at"`
}

func (AvailableSlot *AvailableSlotTesting) TableName() string {
	return "available_slots"
}
