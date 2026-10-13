package notification

import (
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type NotificationType string

const (
	NotificationTypeLike          NotificationType = "like"
	NotificationTypeComment       NotificationType = "comment"
	NotificationTypeReply         NotificationType = "reply"
	NotificationTypeFollow        NotificationType = "follow"
	NotificationTypeFollowRequest NotificationType = "follow_request"
)

type EntityType string

const (
	EntityTypePost    EntityType = "post"
	EntityTypeComment EntityType = "comment"
	EntityTypeUser    EntityType = "user"
)

type NotificationModel struct {
	ID         int64            `column:"id"`
	UserID     int64            `column:"user_id"`
	ActorID    int64            `column:"actor_id"`
	Type       NotificationType `column:"type"`
	EntityID   int64            `column:"entity_id"`
	EntityType EntityType       `column:"entity_type"`
	Message    string           `column:"message"`
	IsRead     bool             `column:"is_read"`
	CreatedAt  utils.JsonTime   `column:"created_at"`
	ReadAt     *utils.JsonTime  `column:"read_at"`
}

type ActorInfo struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type NotificationItemResponse struct {
	ID         int64            `json:"id"`
	Actor      ActorInfo        `json:"actor"`
	Type       NotificationType `json:"type"`
	EntityID   int64            `json:"entity_id"`
	EntityType EntityType       `json:"entity_type"`
	Message    string           `json:"message"`
	IsRead     bool             `json:"is_read"`
	CreatedAt  utils.JsonTime   `json:"created_at"`
	ReadAt     *utils.JsonTime  `json:"read_at,omitempty"`
}

type NotificationPagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}

type GetNotificationsResponse struct {
	Data       []NotificationItemResponse `json:"data"`
	Pagination NotificationPagination     `json:"pagination"`
}

type UnreadCountResponse struct {
	UnreadCount int `json:"unread_count"`
}

type CreateNotificationInput struct {
	UserID     int64            `json:"user_id"`
	ActorID    int64            `json:"actor_id"`
	Type       NotificationType `json:"type"`
	EntityID   int64            `json:"entity_id"`
	EntityType EntityType       `json:"entity_type"`
	Message    string           `json:"message"`
}
