package user

import (
	"net/http"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService UserService
	cfg         *configs.Config
}

func NewUserHandler(userService UserService, cfg *configs.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		cfg:         cfg,
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	data, err := h.userService.GetUserById(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	if data == nil {
		c.JSON(http.StatusNotFound, response.MessageResponse{
			Status:  http.StatusNotFound,
			Message: "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: "success get data",
		Data: gin.H{
			"id":         data.ID,
			"username":   data.Username,
			"email":      data.Email,
			"created_at": data.CreatedAt,
		},
	})
}
