package handlers

import (
	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/models"
	"github.com/Dhruv9449/mou/pkg/serializers"
	"github.com/Dhruv9449/mou/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func getBlogPosts(c *fiber.Ctx) error {
	var blogposts []models.BlogPost
	database.DB.Find(&blogposts)
	return c.JSON(serializers.BlogListSerializer(blogposts))
}

func getBlogPost(c *fiber.Ctx) error {
	var blogpost models.BlogPost
	database.DB.First(&blogpost, c.Params("id"))

	if blogpost.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Blog not found",
		})
	}

	return c.JSON(serializers.BlogPostSerializer(blogpost))
}

func createBlogPost(c *fiber.Ctx) error {
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

	var blogpost models.BlogPost
	err = c.BodyParser(&blogpost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid blogpost data",
		})
	}

	database.DB.Create(&blogpost)

	return c.JSON(serializers.BlogPostSerializer(blogpost))
}

func updateBlogPost(c *fiber.Ctx) error {
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

	var blogpost models.BlogPost
	database.DB.First(&blogpost, c.Params("id"))

	if blogpost.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Blog not found",
		})
	}

	err = c.BodyParser(&blogpost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid blogpost data",
		})
	}

	database.DB.Save(&blogpost)

	return c.JSON(serializers.BlogPostSerializer(blogpost))
}

func deleteBlogPost(c *fiber.Ctx) error {
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

	var blogpost models.BlogPost
	database.DB.First(&blogpost, c.Params("id"))

	if blogpost.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Blog not found",
		})
	}

	database.DB.Delete(&blogpost)

	return c.SendStatus(fiber.StatusNoContent)
}

func BlogRouter(app fiber.Router) {
	group := app.Group("/blog")
	group.Get("/", getBlogPosts)
	group.Get("/:id", getBlogPost)
	group.Post("/", createBlogPost)
	group.Put("/:id", updateBlogPost)
	group.Delete("/:id", deleteBlogPost)
}
