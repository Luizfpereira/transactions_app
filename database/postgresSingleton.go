package database

import (
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var instance *gorm.DB
var lock = &sync.Mutex{}

// ConnectSingleton creates a postgres gorm instance if it was not created before.
// If a connection alreay exists, it returns the instance
func ConnectSingleton(connection string) *gorm.DB {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			log.Println("Creating database instance...")
			var err error
			instance, err = gorm.Open(postgres.Open(connection), &gorm.Config{})
			if err != nil {
				log.Panic("Failed to connect to database!")
			}
		} else {
			log.Println("Connected to database!")
		}
	} else {
		log.Println("Connected to database!")
	}
	return instance
}
