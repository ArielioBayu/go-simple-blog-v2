package memberships

import (
	"errors"
	"log"
	"net/http"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()

	var request memberships.RefreshTokenRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	refreshToken, err := h.membershipsService.GetIdRefreshToken(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
	}

	token, err := h.membershipsService.ValidateRefreshToken(ctx, refreshToken.UserId, request)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrTokenExpired):
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})

		case errors.Is(err, constants.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})

		default:
			log.Printf("handler refresh failed | error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}
