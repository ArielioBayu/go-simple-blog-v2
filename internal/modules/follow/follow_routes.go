package follow

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/gin-gonic/gin"
)

func FollowRoutes(r *gin.RouterGroup, h *FollowHandler) {
	group := r.Group("/users")
	group.Use(middleware.AuthMiddlewareToken())

	// Core Follow & Unfollow
	group.POST("/follow/:targetUserId", h.FollowUser)
	group.DELETE("/follow/:targetUserId", h.UnfollowUser)
	group.DELETE("/followers/:followerUserId", h.RemoveFollower)

	// Follow Requests for Private Accounts
	group.GET("/follow/requests", h.GetPendingRequests)
	group.POST("/follow/requests/:followerUserId/accept", h.AcceptFollowRequest)
	group.POST("/follow/requests/:followerUserId/reject", h.RejectFollowRequest)

	// Social Graph Lists & Relationship
	group.GET("/:userId/followers", h.GetFollowers)
	group.GET("/:userId/following", h.GetFollowing)
	group.GET("/:userId/relationship", h.GetRelationshipStatus)
}

