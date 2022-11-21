package dbutils

import (
	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/models"
)

func CheckUserExists(email string) bool {
	var user models.User
	return database.DB.Where("email = ?", email).First(&user).RowsAffected != 0
}

func GetUserByEmail(email string) models.User {
	var user models.User
	database.DB.Where("email = ?", email).First(&user)
	return user
}
