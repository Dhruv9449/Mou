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

func getComments(c *fiber.Ctx) error {
	blog_title := utils.CovertSlugToTitle(c.Params("title"))

	var blog models.BlogPost

	if database.DB.Where("title = ?", blog_title).First(&blog).RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Blog not found",
		})
	}

	var comments []models.Comment
	database.DB.Preload(clause.Associations).Where("blog_post_id = ? AND parent_id IS NULL", blog.ID).Order("created_on desc").Find(&comments)

	return c.Status(fiber.StatusOK).JSON(serializers.CommentsListSerializer(comments))
}

func createComment(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	blog_title := utils.CovertSlugToTitle(c.Params("title"))
	var blog models.BlogPost

	if database.DB.Where("title = ?", blog_title).First(&blog).RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Blog not found",
		})
	}

	var comment models.Comment

	comment.Content = c.FormValue("content")
	comment.Author = user
	comment.BlogPost = blog
	comment.CreatedOn = time.Now()

	var parent models.Comment

	parent_id := c.FormValue("parent_id")

	if parent_id != "" && parent_id != "null" {
		if database.DB.Where("id = ?", parent_id).First(&parent).RowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Parent comment not found",
			})
		}

		comment.ParentID = &parent.ID
		database.DB.Create(&comment)
	} else {
		database.DB.Create(&comment)
	}

	database.DB.Save(&comment)

	return c.Status(fiber.StatusCreated).JSON(serializers.CommentSerializer(comment))
}

func deleteComment(c *fiber.Ctx) error {
	user, err := utils.GetUserFromHeader(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	var comment models.Comment

	if database.DB.Where("id = ?", c.Params("id")).First(&comment).RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Comment not found",
		})
	}

	if comment.AuthorID != user.ID && user.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized to delete this comment",
		})
	}

	database.DB.Delete(&comment)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Comment deleted successfully",
	})
}
