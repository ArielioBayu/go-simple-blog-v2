package user

import (
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/apperror"
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
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	data, err := h.userService.GetUserById(ctx, userId)
	if err != nil {
		response.Error(c, err)
		return
	}

	if data == nil {
		response.Error(c, constants.ErrUserNotFound)
		return
	}

	response.Data(c, http.StatusOK, "success get data", UserResponse{
		ID:        data.ID,
		Username:  data.Username,
		Email:     data.Email,
		Bio:       data.Bio,
		AvatarURL: data.AvatarURL,
		BannerURL: data.BannerURL,
		CreatedAt: data.CreatedAt,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.GetInt("id")
	if userId == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	profile, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		response.Error(c, err)
		return
	}

	if profile == nil {
		response.Error(c, constants.ErrUserNotFound)
		return
	}

	response.Data(c, http.StatusOK, "success get user profile", profile)
}

func (h *UserHandler) GetProfileByID(c *gin.Context) {
	ctx := c.Request.Context()

	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil || targetUserID <= 0 {
		response.Error(c, apperror.NewBadRequest("invalid user id", nil))
		return
	}

	profile, err := h.userService.GetProfileByID(ctx, targetUserID)
	if err != nil {
		response.Error(c, err)
		return
	}

	if profile == nil {
		response.Error(c, constants.ErrUserNotFound)
		return
	}

	response.Data(c, http.StatusOK, "success get user profile", profile)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.GetInt("id")
	if userId == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	profile, err := h.userService.UpdateProfile(ctx, userId, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, "profile updated successfully", profile)
}

func (h *UserHandler) UpdatePrivacy(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.GetInt("id")
	if userId == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	var req UpdatePrivacyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.NewBadRequest("invalid request body: is_private is required (true or false)", err))
		return
	}

	if err := h.userService.UpdatePrivacy(ctx, userId, req); err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, "account privacy updated", gin.H{
		"is_private": *req.IsPrivate,
	})
}
