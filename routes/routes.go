package routes

import (
	"github.com/ArielioBayu/go-simple-blog-v2/controllers/authcontroller"
	"github.com/ArielioBayu/go-simple-blog-v2/controllers/categorycontroller"
	"github.com/ArielioBayu/go-simple-blog-v2/controllers/commentcontroller"
	"github.com/ArielioBayu/go-simple-blog-v2/controllers/postcontroller"
	"github.com/ArielioBayu/go-simple-blog-v2/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(r *gin.Engine) {
	// API v1 group
	api := r.Group("/api/v1")
	{
		// Public routes
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "message": "Simple Blog V2 API"})
		})

		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authcontroller.Register)
			auth.POST("/login", authcontroller.Login)
		}

		// Public post routes
		posts := api.Group("/posts")
		{
			posts.GET("", postcontroller.Index)
			posts.GET("/:id", postcontroller.Show)
		}

		// Public category routes
		categories := api.Group("/categories")
		{
			categories.GET("", categorycontroller.Index)
			categories.GET("/:id", categorycontroller.Show)
		}

		// Public comment routes (read only)
		comments := api.Group("/comments")
		{
			comments.GET("", commentcontroller.Index)
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Auth routes (protected)
			protected.GET("/auth/me", authcontroller.Me)

			// Post routes (protected)
			protectedPosts := protected.Group("/posts")
			{
				protectedPosts.GET("/my-posts", postcontroller.MyPosts)
				protectedPosts.POST("", postcontroller.Create)
				protectedPosts.PUT("/:id", postcontroller.Update)
				protectedPosts.DELETE("/:id", postcontroller.Delete)
			}

			// Comment routes (protected)
			protectedComments := protected.Group("/comments")
			{
				protectedComments.POST("", commentcontroller.Create)
				protectedComments.PUT("/:id", commentcontroller.Update)
				protectedComments.DELETE("/:id", commentcontroller.Delete)
			}

			// Category routes (admin only)
			adminCategories := protected.Group("/categories")
			adminCategories.Use(middleware.RequireRole("admin"))
			{
				adminCategories.POST("", categorycontroller.Create)
				adminCategories.PUT("/:id", categorycontroller.Update)
				adminCategories.DELETE("/:id", categorycontroller.Delete)
			}
		}
	}
}
