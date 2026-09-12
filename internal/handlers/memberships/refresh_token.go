package memberships

import (
	"errors"
	"log"
	"net/http"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()

	var request memberships.RefreshTokenRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	refreshToken, err := h.membershipsService.GetIdRefreshToken(ctx, request)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRefreshTokenNotFound), errors.Is(err, constants.ErrDataNotFound):
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
			return

		case errors.Is(err, constants.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, response.MessageResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
			})
			return

		default:
			log.Printf("handler refresh failed | error: %v", err)
			c.JSON(http.StatusInternalServerError, response.MessageResponse{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
			})
			return
		}
	}

	token, err := h.membershipsService.ValidateRefreshToken(ctx, refreshToken.UserId, request)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRefreshTokenNotFound), errors.Is(err, constants.ErrDataNotFound):
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
			return

		case errors.Is(err, constants.ErrTokenExpired):
			c.JSON(http.StatusUnauthorized, response.MessageResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
			})
			return

		case errors.Is(err, constants.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, response.MessageResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
			})
			return

		default:
			log.Printf("handler refresh failed | error: %v", err)
			c.JSON(http.StatusInternalServerError, response.MessageResponse{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: "success refresh token",
		Data: gin.H{
			"access_token": token,
		},
	})
}
