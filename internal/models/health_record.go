package models

import "time"

type HealthRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"-"`
	AnimalID  *uint     `json:"animal_id,omitempty"`
	Brinco    string    `json:"brinco,omitempty"`
	Type      string    `gorm:"not null" json:"type"` // Vacina, Vermífugo, etc.
	Product   string    `json:"product"`
	Date      time.Time `gorm:"not null" json:"date"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}
