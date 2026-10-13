package notification

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NotificationRoutes(r *gin.RouterGroup, h *NotificationHandler) {
	group := r.Group("/notifications")
	group.Use(middleware.AuthMiddlewareToken())

	group.GET("/unread-count", h.GetUnreadCount)
	group.GET("", h.GetNotifications)
	group.PUT("/read-all", h.MarkAllAsRead)
	group.PUT("/:notificationId/read", h.MarkAsRead)
}
