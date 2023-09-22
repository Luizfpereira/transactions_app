package database

import (
	"log"
	"transactions_app/entity"

	"gorm.io/gorm"
)

// Migrate auto migrates the model Transaction to the database
func Migrate(instance *gorm.DB) {
	if err := instance.AutoMigrate(&entity.Transaction{}); err != nil {
		log.Fatalln("Could not migrate models to database!")
	}
	log.Println("Database migration completed!")
}
