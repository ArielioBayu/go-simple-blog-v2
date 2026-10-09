package follow

import (
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type FollowModel struct {
	ID          int64     `column:"id"`
	FollowerID  int64     `column:"follower_id"`
	FollowingID int64     `column:"following_id"`
	Status      string    `column:"status"`
	CreatedAt   time.Time `column:"created_at"`
	UpdatedAt   time.Time `column:"updated_at"`
}

type FollowActionResponse struct {
	TargetUserID int    `json:"target_user_id"`
	Status       string `json:"status"`
	IsFollowing  bool   `json:"is_following"`
	IsPending    bool   `json:"is_pending"`
}

type FollowRequestUser struct {
	ID          int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	Username    string         `json:"username"`
	AvatarURL   string         `json:"avatar_url"`
	Bio         string         `json:"bio"`
	RequestedAt utils.JsonTime `json:"requested_at"`
}

type FollowPagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}

type FollowRequestsResponse struct {
	Data       []FollowRequestUser `json:"data"`
	Pagination FollowPagination    `json:"pagination"`
}

type FollowUserData struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	AvatarURL    string `json:"avatar_url"`
	Bio          string `json:"bio"`
	IsFollowing  bool   `json:"is_following"`
	IsFollowedBy bool   `json:"is_followed_by"`
}

type FollowListResponse struct {
	Data       []FollowUserData `json:"data"`
	Pagination FollowPagination `json:"pagination"`
}

type RelationshipResponse struct {
	TargetUserID   int  `json:"target_user_id"`
	IsSelf         bool `json:"is_self"`
	IsPrivate      bool `json:"is_private"`
	IsFollowing    bool `json:"is_following"`
	IsPending      bool `json:"is_pending"`
	IsFollowedBy   bool `json:"is_followed_by"`
	CanViewContent bool `json:"can_view_content"`
}

