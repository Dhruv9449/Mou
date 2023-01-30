package management

import (
	"fmt"

	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/models"
	"github.com/Dhruv9449/mou/pkg/utils/dbutils"
	"github.com/urfave/cli"
)

func CreateSuperuser(c *cli.Context) error {
	var email string

	fmt.Println("Enter email: ")
	fmt.Scan(&email)

	if !dbutils.CheckUserExists(email) {
		fmt.Println("First login to create a superuser")
		return nil
	}

	user := dbutils.GetUserByEmail(email)

	user.Role = "admin"

	database.DB.Save(&user)

	fmt.Println("Superuser created successfully!")
	return nil
}

func ViewAllUsers(c *cli.Context) error {
	var users []models.User

	database.DB.Find(&users)

	for _, user := range users {
		fmt.Println(user)
	}

	return nil
}
