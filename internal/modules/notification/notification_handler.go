package notification

import (
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/apperror"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service NotificationService
	cfg     *configs.Config
}

func NewNotificationHandler(service NotificationService, cfg *configs.Config) *NotificationHandler {
	return &NotificationHandler{
		service: service,
		cfg:     cfg,
	}
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	ctx := c.Request.Context()

	userID := int64(c.GetInt("id"))
	if userID == 0 {
		response.Error(c, apperror.NewUnauthorized("Unauthorized", nil))
		return
	}

	res, err := h.service.GetUnreadCount(ctx, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, "success get unread notification count", res)
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	ctx := c.Request.Context()

	userID := int64(c.GetInt("id"))
	if userID == 0 {
		response.Error(c, apperror.NewUnauthorized("Unauthorized", nil))
		return
	}

	page := 1
	if pStr := c.Query("page"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	} else if pStr := c.Query("pageIndex"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if lStr := c.Query("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	} else if lStr := c.Query("pageSize"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	res, err := h.service.GetNotifications(ctx, userID, page, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, "success get notifications", res.Data, res.Pagination)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	ctx := c.Request.Context()

	userID := int64(c.GetInt("id"))
	if userID == 0 {
		response.Error(c, apperror.NewUnauthorized("Unauthorized", nil))
		return
	}

	if err := h.service.MarkAllAsRead(ctx, userID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success mark all notifications as read")
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	ctx := c.Request.Context()

	userID := int64(c.GetInt("id"))
	if userID == 0 {
		response.Error(c, apperror.NewUnauthorized("Unauthorized", nil))
		return
	}

	notifID, err := strconv.ParseInt(c.Param("notificationId"), 10, 64)
	if err != nil || notifID <= 0 {
		response.Error(c, apperror.NewBadRequest("Invalid notification ID", err))
		return
	}

	if err := h.service.MarkAsReadByID(ctx, userID, notifID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success mark notification as read")
}
