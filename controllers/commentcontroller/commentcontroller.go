package commentcontroller

import (
	"net/http"

	"github.com/ArielioBayu/go-simple-blog-v2/config"
	"github.com/ArielioBayu/go-simple-blog-v2/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateCommentRequest struct {
	Content  string  `json:"content" binding:"required,min=1"`
	PostID   string  `json:"post_id" binding:"required"`
	ParentID *string `json:"parent_id"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

// Index lists comments for a post
func Index(c *gin.Context) {
	postID := c.Query("post_id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}

	var comments []models.Comment
	if err := config.DB.Preload("User").Preload("Replies.User").
		Where("post_id = ? AND parent_id IS NULL", postID).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

// Create creates a new comment
func Create(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")

	// Check if post exists
	var post models.Post
	if err := config.DB.First(&post, "id = ?", req.PostID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// If replying to a comment, check if parent comment exists
	if req.ParentID != nil && *req.ParentID != "" {
		var parentComment models.Comment
		if err := config.DB.First(&parentComment, "id = ?", *req.ParentID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Parent comment not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	comment := models.Comment{
		Content:  req.Content,
		PostID:   req.PostID,
		UserID:   userID.(string),
		ParentID: req.ParentID,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Load user for response
	config.DB.Preload("User").First(&comment, "id = ?", comment.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": comment,
	})
}

// Update updates a comment
func Update(c *gin.Context) {
	commentID := c.Param("id")
	userID, _ := c.Get("userID")
	roleName, _ := c.Get("roleName")

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var comment models.Comment
	if err := config.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check permission: user must be the comment author or admin
	if comment.UserID != userID.(string) && roleName != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this comment"})
		return
	}

	comment.Content = req.Content

	if err := config.DB.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	// Load user for response
	config.DB.Preload("User").First(&comment, "id = ?", comment.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
		"comment": comment,
	})
}

// Delete deletes a comment
func Delete(c *gin.Context) {
	commentID := c.Param("id")
	userID, _ := c.Get("userID")
	roleName, _ := c.Get("roleName")

	var comment models.Comment
	if err := config.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check permission
	if comment.UserID != userID.(string) && roleName != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this comment"})
		return
	}

	// Delete comment and its replies
	if err := config.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
