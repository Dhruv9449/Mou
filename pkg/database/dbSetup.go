package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	if os.Getenv("DEBUG") == "true" {
		DB, err = gorm.Open(sqlite.Open("mou.db"), &gorm.Config{})
	} else {
		DB, err = gorm.Open(postgres.Open(os.Getenv("POSTGRES_URL")), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}
}
