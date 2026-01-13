package postcontroller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/config"
	"github.com/ArielioBayu/go-simple-blog-v2/models"
	"github.com/ArielioBayu/go-simple-blog-v2/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreatePostRequest struct {
	Title      string   `json:"title" binding:"required,min=3,max=255"`
	Content    string   `json:"content" binding:"required"`
	Excerpt    string   `json:"excerpt"`
	CategoryID string   `json:"category_id"`
	Tags       []string `json:"tags"`
	Status     string   `json:"status"` // draft, published
}

type UpdatePostRequest struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Excerpt    string   `json:"excerpt"`
	CategoryID string   `json:"category_id"`
	Tags       []string `json:"tags"`
	Status     string   `json:"status"`
}

// Index lists all published posts with pagination
func Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	categoryID := c.Query("category_id")
	status := c.DefaultQuery("status", "published")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	query := config.DB.Preload("Author").Preload("Category").Preload("Tags")

	// Apply filters
	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	query = query.Where("status = ?", status)

	// Get total count
	var total int64
	query.Model(&models.Post{}).Count(&total)

	// Get posts
	var posts []models.Post
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// Show returns a single post by ID or slug
func Show(c *gin.Context) {
	identifier := c.Param("id")

	var post models.Post
	query := config.DB.Preload("Author").Preload("Category").Preload("Tags").Preload("Comments.User")

	// Try to find by ID first, then by slug
	if err := query.Where("id = ? OR slug = ?", identifier, identifier).First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// Create creates a new post
func Create(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	roleName, _ := c.Get("roleName")

	// Check if user has permission to create posts
	if roleName != "admin" && roleName != "author" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only authors and admins can create posts"})
		return
	}

	// Generate slug from title
	slug := utils.GenerateSlug(req.Title)

	// Check if slug already exists and make it unique
	var count int64
	config.DB.Model(&models.Post{}).Where("slug LIKE ?", slug+"%").Count(&count)
	if count > 0 {
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	post := models.Post{
		Title:      req.Title,
		Slug:       slug,
		Content:    req.Content,
		Excerpt:    req.Excerpt,
		AuthorID:   userID.(string),
		CategoryID: req.CategoryID,
		Status:     req.Status,
	}

	if post.Status == "" {
		post.Status = "draft"
	}

	if post.Status == "published" {
		now := time.Now()
		post.PublishedAt = &now
	}

	// Start transaction
	tx := config.DB.Begin()

	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Handle tags
	if len(req.Tags) > 0 {
		var tags []models.Tag
		for _, tagName := range req.Tags {
			var tag models.Tag
			tagSlug := utils.GenerateSlug(tagName)

			// Find or create tag
			if err := tx.Where("slug = ?", tagSlug).First(&tag).Error; err == gorm.ErrRecordNotFound {
				tag = models.Tag{Name: tagName, Slug: tagSlug}
				if err := tx.Create(&tag).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
					return
				}
			}
			tags = append(tags, tag)
		}
		if err := tx.Model(&post).Association("Tags").Append(tags); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate tags"})
			return
		}
	}

	tx.Commit()

	// Load relationships
	config.DB.Preload("Author").Preload("Category").Preload("Tags").First(&post, "id = ?", post.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post":    post,
	})
}

// Update updates a post
func Update(c *gin.Context) {
	postID := c.Param("id")
	userID, _ := c.Get("userID")
	roleName, _ := c.Get("roleName")

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post models.Post
	if err := config.DB.First(&post, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check permission: user must be the author or an admin
	if post.AuthorID != userID.(string) && roleName != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this post"})
		return
	}

	// Update fields if provided
	if req.Title != "" {
		post.Title = req.Title
		post.Slug = utils.GenerateSlug(req.Title)
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Excerpt != "" {
		post.Excerpt = req.Excerpt
	}
	if req.CategoryID != "" {
		post.CategoryID = req.CategoryID
	}
	if req.Status != "" {
		// If changing to published, set published_at
		if req.Status == "published" && post.Status != "published" {
			now := time.Now()
			post.PublishedAt = &now
		}
		post.Status = req.Status
	}

	// Update post
	if err := config.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	// Update tags if provided
	if req.Tags != nil {
		var tags []models.Tag
		for _, tagName := range req.Tags {
			var tag models.Tag
			tagSlug := utils.GenerateSlug(tagName)

			if err := config.DB.Where("slug = ?", tagSlug).First(&tag).Error; err == gorm.ErrRecordNotFound {
				tag = models.Tag{Name: tagName, Slug: tagSlug}
				config.DB.Create(&tag)
			}
			tags = append(tags, tag)
		}
		config.DB.Model(&post).Association("Tags").Replace(tags)
	}

	// Load relationships
	config.DB.Preload("Author").Preload("Category").Preload("Tags").First(&post, "id = ?", post.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"post":    post,
	})
}

// Delete deletes a post
func Delete(c *gin.Context) {
	postID := c.Param("id")
	userID, _ := c.Get("userID")
	roleName, _ := c.Get("roleName")

	var post models.Post
	if err := config.DB.First(&post, "id = ?", postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check permission
	if post.AuthorID != userID.(string) && roleName != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this post"})
		return
	}

	// Delete post (this will cascade delete comments and associations)
	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// MyPosts returns posts created by the authenticated user
func MyPosts(c *gin.Context) {
	userID, _ := c.Get("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var total int64
	config.DB.Model(&models.Post{}).Where("author_id = ?", userID).Count(&total)

	var posts []models.Post
	if err := config.DB.Preload("Category").Preload("Tags").
		Where("author_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}
