package comment

import "github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"

type CommentRequest struct {
	CommentContent string `json:"comment_content"`
}

type UpdateCommentRequest struct {
	CommentContent string `json:"comment_content"`
}

type ReplyRequest struct {
	ReplyContent  string `json:"reply_content"`
	ReplyToUserID *int64 `json:"reply_to_user_id,omitempty"`
}

type UpdateReplyRequest struct {
	ReplyContent string `json:"reply_content"`
}

type CommentAuthInfo struct {
	CommentID       int
	CommentAuthorID int
	PostAuthorID    int
}

type ReplyAuthInfo struct {
	ReplyID       int
	ReplyAuthorID int
	PostAuthorID  int
}

type CommentModel struct {
	ID             int            `column:"id"`
	PostId         int            `column:"post_id"`
	UserId         int            `column:"user_id"`
	CommentContent string         `column:"comment_content"`
	CreatedAt      utils.JsonTime `column:"created_at"`
	UpdatedAt      utils.JsonTime `column:"updated_at"`
	CreatedBy      string         `column:"created_by"`
	UpdatedBy      string         `column:"updated_by"`
}

type CommentReplyModel struct {
	ID            int            `column:"id"`
	CommentID     int            `column:"comment_id"`
	UserID        int            `column:"user_id"`
	ReplyToUserID *int64         `column:"reply_to_user_id"`
	ReplyContent  string         `column:"reply_content"`
	CreatedAt     utils.JsonTime `column:"created_at"`
	UpdatedAt     utils.JsonTime `column:"updated_at"`
	CreatedBy     string         `column:"created_by"`
	UpdatedBy     string         `column:"updated_by"`
}

type GetCommentReply struct {
	ID              int            `json:"id"`
	CommentID       int            `json:"comment_id"`
	UserID          int            `json:"user_id"`
	Username        string         `json:"username"`
	AvatarURL       string         `json:"avatar_url"`
	ReplyToUserID   *int64         `json:"reply_to_user_id,omitempty"`
	ReplyToUsername string         `json:"reply_to_username,omitempty"`
	ReplyContent    string         `json:"reply_content"`
	CreatedAt       utils.JsonTime `json:"created_at"`
}

type GetComment struct {
	ID             int               `json:"id"`
	PostId         int               `json:"post_id,omitempty"`
	UserId         int               `json:"user_id"`
	Username       string            `json:"username"`
	AvatarURL      string            `json:"avatar_url"`
	CommentContent string            `json:"comment_content"`
	CreatedAt      utils.JsonTime    `json:"created_at"`
	RepliesCount   int               `json:"replies_count"`
	Replies        []GetCommentReply `json:"replies"`
}

type ReplyPagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}

type GetRepliesResponse struct {
	Data       []GetCommentReply `json:"data"`
	Pagination ReplyPagination   `json:"pagination"`
}

type CommentPagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}

type GetCommentsResponse struct {
	Data       []GetComment      `json:"data"`
	Pagination CommentPagination `json:"pagination"`
}
