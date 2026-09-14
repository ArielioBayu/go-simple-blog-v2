package activity

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
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
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Post Id",
		})
		return
	}

	userId, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}

	err = h.activityService.InsertUpdateActivities(ctx, postId, userId.(int), request)
	if err != nil {
		log.Printf("InsertUpdateActivities failed | error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Status:  http.StatusOK,
		Message: "success",
	})
}
