package models

import "github.com/Dhruv9449/mou/pkg/database"

func InitializeModels() {

	MODELS := []interface{}{
		&User{},
		&BlogPost{},
		&Comment{},
	}

	database.DB.AutoMigrate(MODELS...)
}
