package serializers

import "github.com/Dhruv9449/mou/pkg/models"

func BlogPostSerializer(blogPost models.BlogPost) map[string]interface{} {
	return map[string]interface{}{
		"id":         blogPost.ID,
		"title":      blogPost.Title,
		"content":    blogPost.Content,
		"author":     blogPost.Author,
		"created_on": blogPost.CreatedOn,
		"updated_on": blogPost.UpdatedOn,
	}
}

func BlogBlockSerializer(blogPost models.BlogPost) map[string]interface{} {
	return map[string]interface{}{
		"id":         blogPost.ID,
		"title":      blogPost.Title,
		"author":     blogPost.Author,
		"created_on": blogPost.CreatedOn,
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
	return map[string]interface{}{
		"id":         comment.ID,
		"content":    comment.Content,
		"author":     comment.Author,
		"blog_post":  comment.BlogPost,
		"created_on": comment.CreatedOn,
		"replies":    comment.Replies,
	}
}

func CommentsListSerializer(comments []models.Comment) []map[string]interface{} {
	var commentsList []map[string]interface{}
	for _, comment := range comments {
		commentsList = append(commentsList, CommentSerializer(comment))
	}
	return commentsList
}
