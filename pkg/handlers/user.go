package handlers

import (
	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/models"
	"github.com/Dhruv9449/mou/pkg/serializers"
	"github.com/Dhruv9449/mou/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func getAllUsers(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to perform this action",
		})
	}

	var users []models.User
	database.DB.Find(&users)
	return c.JSON(serializers.UserListSerializer(users))
}

func getUser(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to perform this action",
		})
	}

	var userToGet models.User
	database.DB.First(&userToGet, c.Params("id"))

	if userToGet.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(serializers.UserSerializer(userToGet))
}

func deleteUser(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to perform this action",
		})
	}

	var userToDelete models.User
	database.DB.First(&userToDelete, c.Params("id"))

	if userToDelete.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	database.DB.Delete(&userToDelete)
	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func UserRouter(app fiber.Router) {
	group := app.Group("/users")
	group.Get("/", getAllUsers)
	group.Get("/:id", getUser)
	group.Delete("/:id", deleteUser)
}
