package activity

import (
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/apperror"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	activityService ActivityService
	cfg             *configs.Config
}

func NewActivityHandler(activityService ActivityService, cfg *configs.Config) *ActivityHandler {
	return &ActivityHandler{
		activityService: activityService,
		cfg:             cfg,
	}
}

func (h *ActivityHandler) InsertUpdateActivities(c *gin.Context) {
	ctx := c.Request.Context()

	var request ActivityRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, err)
		return
	}

	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		response.Error(c, apperror.NewBadRequest("Invalid Post Id", err))
		return
	}

	userId, exists := c.Get("id")
	if !exists {
		response.Error(c, apperror.NewUnauthorized("Unauthorized", nil))
		return
	}

	err = h.activityService.InsertUpdateActivities(ctx, postId, userId.(int), request)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success")
}

func (h *ActivityHandler) CountLikes(c *gin.Context) {
	ctx := c.Request.Context()

	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		response.Error(c, apperror.NewBadRequest("Invalid Post Id", err))
		return
	}

	count, err := h.activityService.CountLikes(ctx, postId)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, "success get like count", gin.H{
		"post_id":    postId,
		"like_count": count,
	})
}
