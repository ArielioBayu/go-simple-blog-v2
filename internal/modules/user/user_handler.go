package user

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
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
		Data: UserResponse{
			ID:        data.ID,
			Username:  data.Username,
			Email:     data.Email,
			Bio:       data.Bio,
			AvatarURL: data.AvatarURL,
			BannerURL: data.BannerURL,
			CreatedAt: data.CreatedAt,
		},
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	profile, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		log.Printf("GetProfile error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, response.MessageResponse{
			Status:  http.StatusNotFound,
			Message: "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: "success get user profile",
		Data:    profile,
	})
}

func (h *UserHandler) GetProfileByID(c *gin.Context) {
	ctx := c.Request.Context()

	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil || targetUserID <= 0 {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid user id",
		})
		return
	}

	profile, err := h.userService.GetProfileByID(ctx, targetUserID)
	if err != nil {
		log.Printf("GetProfileByID error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, response.MessageResponse{
			Status:  http.StatusNotFound,
			Message: "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: "success get user profile",
		Data:    profile,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	profile, err := h.userService.UpdateProfile(ctx, userId, req)
	if err != nil {
		if errors.Is(err, constants.ErrUsernameOrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, response.MessageResponse{
				Status:  http.StatusConflict,
				Message: "username already exists",
			})
			return
		}

		log.Printf("UpdateProfile error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: "profile updated successfully",
		Data:    profile,
	})
}
