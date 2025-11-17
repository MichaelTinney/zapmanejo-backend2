package database

import (
	"zapmanejo-backend2/internal/models"
)

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Animal{},
		&models.HealthRecord{},
		&models.CostConfig{},
		&models.Payment{},
		&models.LifetimeSlot{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_animals_brinco ON animals(brinco)`)
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_animals_birth ON animals(birth_date)`)

	SeedLifetimeSlots() // creates the 200 empty Early Adopter slots
}
