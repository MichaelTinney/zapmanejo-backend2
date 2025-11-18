// internal/models/lifetime_slot.go
package models

type LifetimeSlot struct {
	ID        uint `gorm:"primaryKey"`
	Slot      int  `gorm:"unique;not null"`
	UserID    *uint
	Occupied  bool `gorm:"default:false"`
}
