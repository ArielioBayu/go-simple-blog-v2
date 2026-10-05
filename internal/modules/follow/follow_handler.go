package follow

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	followService FollowService
	cfg           *configs.Config
}

func NewFollowHandler(followService FollowService, cfg *configs.Config) *FollowHandler {
	return &FollowHandler{
		followService: followService,
		cfg:           cfg,
	}
}

func (h *FollowHandler) FollowUser(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	targetUserID, err := strconv.Atoi(c.Param("targetUserId"))
	if err != nil || targetUserID <= 0 {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid target user id",
		})
		return
	}

	res, err := h.followService.FollowUser(ctx, currentUserID, targetUserID)
	if err != nil {
		if errors.Is(err, constants.ErrCannotFollowSelf) {
			c.JSON(http.StatusBadRequest, response.MessageResponse{
				Status:  http.StatusBadRequest,
				Message: "cannot follow yourself",
			})
			return
		}
		if errors.Is(err, constants.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: "target user not found",
			})
			return
		}
		if errors.Is(err, constants.ErrAlreadyFollowing) {
			c.JSON(http.StatusConflict, response.MessageResponse{
				Status:  http.StatusConflict,
				Message: "already following this user",
			})
			return
		}
		if errors.Is(err, constants.ErrFollowRequestPending) {
			c.JSON(http.StatusConflict, response.MessageResponse{
				Status:  http.StatusConflict,
				Message: "follow request already sent",
			})
			return
		}

		log.Printf("FollowUser error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	msg := "successfully followed user"
	if res.Status == utils.FollowStatusPending {
		msg = "follow request sent"
	}

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: msg,
		Data:    res,
	})
}

func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	targetUserID, err := strconv.Atoi(c.Param("targetUserId"))
	if err != nil || targetUserID <= 0 {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid target user id",
		})
		return
	}

	err = h.followService.UnfollowUser(ctx, currentUserID, targetUserID)
	if err != nil {
		if errors.Is(err, constants.ErrCannotFollowSelf) {
			c.JSON(http.StatusBadRequest, response.MessageResponse{
				Status:  http.StatusBadRequest,
				Message: "cannot unfollow yourself",
			})
			return
		}
		if errors.Is(err, constants.ErrFollowNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: "follow relationship not found",
			})
			return
		}

		log.Printf("UnfollowUser error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Status:  http.StatusOK,
		Message: "successfully unfollowed user",
	})
}

func (h *FollowHandler) RemoveFollower(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	followerUserID, err := strconv.Atoi(c.Param("followerUserId"))
	if err != nil || followerUserID <= 0 {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid follower user id",
		})
		return
	}

	err = h.followService.RemoveFollower(ctx, currentUserID, followerUserID)
	if err != nil {
		if errors.Is(err, constants.ErrCannotFollowSelf) {
			c.JSON(http.StatusBadRequest, response.MessageResponse{
				Status:  http.StatusBadRequest,
				Message: "cannot remove yourself",
			})
			return
		}
		if errors.Is(err, constants.ErrFollowNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: "follower not found",
			})
			return
		}

		log.Printf("RemoveFollower error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Status:  http.StatusOK,
		Message: "successfully removed follower",
	})
}

func (h *FollowHandler) GetPendingRequests(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", c.DefaultQuery("pageIndex", "1")))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", c.DefaultQuery("pageSize", "10")))

	res, err := h.followService.GetPendingRequests(ctx, currentUserID, page, limit)
	if err != nil {
		log.Printf("GetPendingRequests error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"message":    "success get follow requests",
		"data":       res.Data,
		"pagination": res.Pagination,
	})
}

func (h *FollowHandler) AcceptFollowRequest(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	followerUserID, err := strconv.Atoi(c.Param("followerUserId"))
	if err != nil || followerUserID <= 0 {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid follower user id",
		})
		return
	}

	err = h.followService.AcceptFollowRequest(ctx, currentUserID, followerUserID)
	if err != nil {
		if errors.Is(err, constants.ErrFollowRequestNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: "follow request not found",
			})
			return
		}

		log.Printf("AcceptFollowRequest error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Status:  http.StatusOK,
		Message: "follow request accepted",
	})
}

func (h *FollowHandler) RejectFollowRequest(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	followerUserID, err := strconv.Atoi(c.Param("followerUserId"))
	if err != nil || followerUserID <= 0 {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid follower user id",
		})
		return
	}

	err = h.followService.RejectFollowRequest(ctx, currentUserID, followerUserID)
	if err != nil {
		if errors.Is(err, constants.ErrFollowRequestNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: "follow request not found",
			})
			return
		}

		log.Printf("RejectFollowRequest error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Status:  http.StatusOK,
		Message: "follow request rejected",
	})
}

func (h *FollowHandler) GetFollowers(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	targetUserID, err := strconv.Atoi(c.Param("userId"))
	if err != nil || targetUserID <= 0 {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid user id",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", c.DefaultQuery("pageIndex", "1")))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", c.DefaultQuery("pageSize", "10")))
	search := c.Query("search")

	res, err := h.followService.GetFollowers(ctx, currentUserID, targetUserID, page, limit, search)
	if err != nil {
		if errors.Is(err, constants.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: "user not found",
			})
			return
		}
		if errors.Is(err, constants.ErrPrivateAccount) {
			c.JSON(http.StatusForbidden, response.MessageResponse{
				Status:  http.StatusForbidden,
				Message: constants.ErrPrivateAccount.Error(),
			})
			return
		}

		log.Printf("GetFollowers error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"message":    "success get followers",
		"data":       res.Data,
		"pagination": res.Pagination,
	})
}

func (h *FollowHandler) GetFollowing(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	targetUserID, err := strconv.Atoi(c.Param("userId"))
	if err != nil || targetUserID <= 0 {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid user id",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", c.DefaultQuery("pageIndex", "1")))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", c.DefaultQuery("pageSize", "10")))
	search := c.Query("search")

	res, err := h.followService.GetFollowing(ctx, currentUserID, targetUserID, page, limit, search)
	if err != nil {
		if errors.Is(err, constants.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: "user not found",
			})
			return
		}
		if errors.Is(err, constants.ErrPrivateAccount) {
			c.JSON(http.StatusForbidden, response.MessageResponse{
				Status:  http.StatusForbidden,
				Message: constants.ErrPrivateAccount.Error(),
			})
			return
		}

		log.Printf("GetFollowing error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"message":    "success get following",
		"data":       res.Data,
		"pagination": res.Pagination,
	})
}

func (h *FollowHandler) GetRelationshipStatus(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}

	targetUserID, err := strconv.Atoi(c.Param("userId"))
	if err != nil || targetUserID <= 0 {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid user id",
		})
		return
	}

	res, err := h.followService.GetRelationshipStatus(ctx, currentUserID, targetUserID)
	if err != nil {
		if errors.Is(err, constants.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: "user not found",
			})
			return
		}

		log.Printf("GetRelationshipStatus error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: "success get relationship status",
		Data:    res,
	})
}

