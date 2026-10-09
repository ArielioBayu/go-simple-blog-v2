package follow

import (
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/apperror"
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
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	targetUserID, err := strconv.Atoi(c.Param("targetUserId"))
	if err != nil || targetUserID <= 0 {
		response.Error(c, apperror.NewBadRequest("invalid target user id", nil))
		return
	}

	res, err := h.followService.FollowUser(ctx, currentUserID, targetUserID)
	if err != nil {
		response.Error(c, err)
		return
	}

	msg := "successfully followed user"
	if res.Status == utils.FollowStatusPending {
		msg = "follow request sent"
	}

	response.Data(c, http.StatusOK, msg, res)
}

func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	targetUserID, err := strconv.Atoi(c.Param("targetUserId"))
	if err != nil || targetUserID <= 0 {
		response.Error(c, apperror.NewBadRequest("invalid target user id", nil))
		return
	}

	err = h.followService.UnfollowUser(ctx, currentUserID, targetUserID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "successfully unfollowed user")
}

func (h *FollowHandler) RemoveFollower(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	followerUserID, err := strconv.Atoi(c.Param("followerUserId"))
	if err != nil || followerUserID <= 0 {
		response.Error(c, apperror.NewBadRequest("invalid follower user id", nil))
		return
	}

	err = h.followService.RemoveFollower(ctx, currentUserID, followerUserID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "successfully removed follower")
}

func (h *FollowHandler) GetPendingRequests(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", c.DefaultQuery("pageIndex", "1")))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", c.DefaultQuery("pageSize", "10")))

	res, err := h.followService.GetPendingRequests(ctx, currentUserID, page, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, "success get follow requests", res.Data, res.Pagination)
}

func (h *FollowHandler) AcceptFollowRequest(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	followerUserID, err := strconv.Atoi(c.Param("followerUserId"))
	if err != nil || followerUserID <= 0 {
		response.Error(c, apperror.NewBadRequest("invalid follower user id", nil))
		return
	}

	err = h.followService.AcceptFollowRequest(ctx, currentUserID, followerUserID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "follow request accepted")
}

func (h *FollowHandler) RejectFollowRequest(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	followerUserID, err := strconv.Atoi(c.Param("followerUserId"))
	if err != nil || followerUserID <= 0 {
		response.Error(c, apperror.NewBadRequest("invalid follower user id", nil))
		return
	}

	err = h.followService.RejectFollowRequest(ctx, currentUserID, followerUserID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "follow request rejected")
}

func (h *FollowHandler) GetFollowers(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	targetUserID, err := strconv.Atoi(c.Param("userId"))
	if err != nil || targetUserID <= 0 {
		response.Error(c, apperror.NewBadRequest("invalid user id", nil))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", c.DefaultQuery("pageIndex", "1")))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", c.DefaultQuery("pageSize", "10")))
	search := c.Query("search")

	res, err := h.followService.GetFollowers(ctx, currentUserID, targetUserID, page, limit, search)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, "success get followers", res.Data, res.Pagination)
}

func (h *FollowHandler) GetFollowing(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	targetUserID, err := strconv.Atoi(c.Param("userId"))
	if err != nil || targetUserID <= 0 {
		response.Error(c, apperror.NewBadRequest("invalid user id", nil))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", c.DefaultQuery("pageIndex", "1")))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", c.DefaultQuery("pageSize", "10")))
	search := c.Query("search")

	res, err := h.followService.GetFollowing(ctx, currentUserID, targetUserID, page, limit, search)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, "success get following", res.Data, res.Pagination)
}

func (h *FollowHandler) GetRelationshipStatus(c *gin.Context) {
	ctx := c.Request.Context()
	currentUserID := c.GetInt("id")
	if currentUserID == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	targetUserID, err := strconv.Atoi(c.Param("userId"))
	if err != nil || targetUserID <= 0 {
		response.Error(c, apperror.NewBadRequest("invalid user id", nil))
		return
	}

	res, err := h.followService.GetRelationshipStatus(ctx, currentUserID, targetUserID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, "success get relationship status", res)
}
