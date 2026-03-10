package memberships

import (
	"errors"
	"log"
	"net/http"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SignIn(c *gin.Context) {
	ctx := c.Request.Context()

	var request memberships.SignInRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	token, refreshToken, err := h.membershipsService.SignIn(ctx, request)
	if err != nil {
		if errors.Is(err, constants.ErrDataNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, constants.ErrInvalidPassword) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		log.Printf("handler sign in failed | error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Success Login",
		"access_token":  token,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()

	// binding json
	var request memberships.SignUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := h.membershipsService.SignUp(ctx, request)
	if err != nil {
		if errors.Is(err, constants.ErrUsernameOrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// jika sukses
	c.JSON(http.StatusCreated, gin.H{
		"Message": "success created data",
	})
}
