package categorycontroller

import (
	"net/http"

	"github.com/ArielioBayu/go-simple-blog-v2/config"
	"github.com/ArielioBayu/go-simple-blog-v2/models"
	"github.com/ArielioBayu/go-simple-blog-v2/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Index lists all categories
func Index(c *gin.Context) {
	var categories []models.Category

	if err := config.DB.Order("name ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// Show returns a single category with its posts
func Show(c *gin.Context) {
	identifier := c.Param("id")

	var category models.Category
	if err := config.DB.Preload("Posts.Author").
		Where("id = ? OR slug = ?", identifier, identifier).
		First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

// Create creates a new category (admin only)
func Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug := utils.GenerateSlug(req.Name)

	// Check if slug already exists
	var existing models.Category
	if err := config.DB.Where("slug = ?", slug).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Category with this name already exists"})
		return
	}

	category := models.Category{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Category created successfully",
		"category": category,
	})
}

// Update updates a category (admin only)
func Update(c *gin.Context) {
	categoryID := c.Param("id")

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, "id = ?", categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if req.Name != "" {
		category.Name = req.Name
		category.Slug = utils.GenerateSlug(req.Name)
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := config.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Category updated successfully",
		"category": category,
	})
}

// Delete deletes a category (admin only)
func Delete(c *gin.Context) {
	categoryID := c.Param("id")

	var category models.Category
	if err := config.DB.First(&category, "id = ?", categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if category has posts
	var postCount int64
	config.DB.Model(&models.Post{}).Where("category_id = ?", categoryID).Count(&postCount)
	if postCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category with existing posts"})
		return
	}

	if err := config.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
