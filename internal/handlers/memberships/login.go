package memberships

import (
	"errors"
	"log"
	"net/http"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SignIn(c *gin.Context) {
	ctx := c.Request.Context()

	var request memberships.SignInRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	token, refreshToken, err := h.membershipsService.SignIn(ctx, request)
	if err != nil {
		if errors.Is(err, constants.ErrDataNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
			return
		}
		if errors.Is(err, constants.ErrInvalidPassword) {
			c.JSON(http.StatusBadRequest, response.MessageResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}

		log.Printf("handler sign in failed | error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	utils.SetAccessTokenCookie(c, token)

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: "Success Login",
		Data: gin.H{
			"access_token":  token,
			"refresh_token": refreshToken,
		},
	})
}

func (h *Handler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()

	// binding json
	var request memberships.SignUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	err := h.membershipsService.SignUp(ctx, request)
	if err != nil {
		if errors.Is(err, constants.ErrUsernameOrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, response.MessageResponse{
				Status:  http.StatusConflict,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	// jika sukses
	c.JSON(http.StatusCreated, response.MessageResponse{
		Status:  http.StatusCreated,
		Message: "success created data",
	})
}
