package serializers

import (
	"github.com/Dhruv9449/mou/pkg/database"
	"github.com/Dhruv9449/mou/pkg/models"
	"github.com/dustin/go-humanize"
	"gorm.io/gorm/clause"
)

func BlogPostSerializer(blogPost models.BlogPost) map[string]interface{} {
	return map[string]interface{}{
		"id":         blogPost.ID,
		"title":      blogPost.Title,
		"content":    blogPost.Content,
		"author":     blogPost.Author,
		"thumbnail":  blogPost.Thumbnail,
		"created_on": blogPost.CreatedOn.Format("January 2, 2006"),
		"updated_on": blogPost.UpdatedOn.Format("January 2, 2006"),
	}
}

func BlogBlockSerializer(blogPost models.BlogPost) map[string]interface{} {
	return map[string]interface{}{
		"id":         blogPost.ID,
		"title":      blogPost.Title,
		"author":     UserBlockSerializer(blogPost.Author),
		"created_on": blogPost.CreatedOn.Format("January 2, 2006"),
		"thumbnail":  blogPost.Thumbnail,
	}
}

func BlogListSerializer(blogPosts []models.BlogPost) []map[string]interface{} {
	var blogs []map[string]interface{}
	for _, blogPost := range blogPosts {
		blogs = append(blogs, BlogBlockSerializer(blogPost))
	}
	return blogs
}

func CommentSerializer(comment models.Comment) map[string]interface{} {
	var replies []models.Comment
	database.DB.Preload(clause.Associations).Where("parent_id = ?", comment.ID).Find(&replies)

	return map[string]interface{}{
		"id":         comment.ID,
		"content":    comment.Content,
		"author":     UserBlockSerializer(comment.Author),
		"created_on": humanize.Time(comment.CreatedOn),
		"replies":    CommentsListSerializer(replies),
	}
}

func CommentsListSerializer(comments []models.Comment) []map[string]interface{} {
	var commentsList []map[string]interface{}
	for _, comment := range comments {
		commentsList = append(commentsList, CommentSerializer(comment))
	}
	return commentsList
}
