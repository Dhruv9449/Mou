package handlers

import (
	"time"

	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/models"
	"github.com/Dhruv9449/mou/pkg/serializers"
	"github.com/Dhruv9449/mou/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func getBlogPosts(c *fiber.Ctx) error {
	var blogposts []models.BlogPost
	database.DB.Preload(clause.Associations).Order("created_on desc").Find(&blogposts)
	return c.JSON(serializers.BlogListSerializer(blogposts))
}

func getBlogPost(c *fiber.Ctx) error {
	var blogpost models.BlogPost
	title := utils.CovertSlugToTitle(c.Params("title"))
	database.DB.Preload(clause.Associations).Where("title = ?", title).First(&blogpost)

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

	if database.DB.Where("title = ?", c.FormValue("title")).First(&blogpost).RowsAffected != 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Blog with this title already exists",
		})
	}

	blogpost.Title = c.FormValue("title")
	blogpost.Content = c.FormValue("content")
	blogpost.Author = user
	blogpost.Thumbnail = c.FormValue("thumbnail")
	blogpost.CreatedOn = time.Now()
	blogpost.UpdatedOn = time.Now()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid blogpost data",
		})
	}

	database.DB.Create(&blogpost)
	database.DB.Save(&blogpost)

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
	title := utils.CovertSlugToTitle(c.Params("title"))
	database.DB.First(&blogpost, "title = ?", title)

	if blogpost.Title == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Blog not found",
		})
	}

	title = c.FormValue("title")

	if database.DB.Where("title = ?", title).First(&blogpost).RowsAffected != 0 && blogpost.Title != title {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Blog with this title already exists",
		})
	}

	blogpost.Title = title
	blogpost.Content = c.FormValue("content")
	blogpost.Thumbnail = c.FormValue("thumbnail")
	blogpost.UpdatedOn = time.Now()

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
	title := utils.CovertSlugToTitle(c.Params("title"))
	database.DB.Where("title = ?", title).First(&blogpost)

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
	group.Get("/:title", getBlogPost)
	group.Post("/", createBlogPost)
	group.Put("/:title", updateBlogPost)
	group.Delete("/:title", deleteBlogPost)
	group.Get("/:title/comments", getComments)
	group.Post("/:title/comments/", createComment)
	group.Delete("/:title/comments/:id", deleteComment)
}
