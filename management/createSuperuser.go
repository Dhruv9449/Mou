package main

import (
	"fmt"

	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/utils/dbutils"
)

func createSuperuser() {
	var email string

	fmt.Println("Enter email: ")
	fmt.Scan(&email)

	if !dbutils.CheckUserExists(email) {
		fmt.Println("First login to create a superuser")
		return
	}

	user := dbutils.GetUserByEmail(email)

	user.Role = "admin"

	database.DB.Save(&user)

	fmt.Println("Superuser created successfully!")
}

func main() {
	database.Connect()
	createSuperuser()
}
