package management

import (
	"fmt"

	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/models"
)

func DeleteBlogs() {
	var id int

	fmt.Println("Enter id of blog to be deleted, -1 to delete all and 0 to cancel")
	fmt.Scan(&id)

	if id == -1 {
		database.DB.Unscoped().Delete(&models.BlogPost{})
		fmt.Println("Deleted all blogs")
	}

	if id == 0 {
		fmt.Println("Cancelled")
		return
	}

	var blog models.BlogPost
	database.DB.First(&blog, id)

	fmt.Println("Deleted blog!")

}
