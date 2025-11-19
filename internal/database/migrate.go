package database

import (
	"log"
	"zapmanejo-cleanbackend/internal/models"
)

func AutoMigrate() {
	log.Println("Starting database migration...")

	// Run GORM AutoMigrate with all models
	err := DB.AutoMigrate(
		&models.User{},
		&models.Animal{},
		&models.HealthRecord{},
		&models.CostConfig{},
		&models.Payment{},
		&models.LifetimeSlot{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
	log.Println("✓ Database schema migrated successfully")

	// Create indexes (idempotent - IF NOT EXISTS)
	log.Println("Creating indexes...")
	result := DB.Exec(`CREATE INDEX IF NOT EXISTS idx_animals_brinco ON animals(brinco)`)
	if result.Error != nil {
		log.Printf("Warning: Failed to create idx_animals_brinco: %v", result.Error)
	}

	result = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_animals_birth ON animals(birth_date)`)
	if result.Error != nil {
		log.Printf("Warning: Failed to create idx_animals_birth: %v", result.Error)
	}
	log.Println("✓ Indexes created successfully")

	// Seed lifetime slots (idempotent - only creates if empty)
	log.Println("Seeding lifetime slots...")
	SeedLifetimeSlots()
	log.Println("✓ Migration completed successfully")
}
