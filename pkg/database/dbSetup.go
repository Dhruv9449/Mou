package database

import (
	"log"

	"github.com/Dhruv9449/mou/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("vitty.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	MODELS := []interface{}{
		&models.User{},
		&models.BlogPost{},
		&models.Comment{},
	}

	DB.AutoMigrate(MODELS...)

}
