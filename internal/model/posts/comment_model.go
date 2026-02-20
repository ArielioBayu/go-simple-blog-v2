package posts

import "time"

type CommentRequest struct {
	CommentContent string `json:"comment_content"`
}

type CommentModel struct {
	ID             int       `json:"id"`
	PostId         int       `json:"post_id"`
	UserId         int       `json:"user_id"`
	CommentContent string    `json:"comment_content"`
	CreatedAt      time.Time `column:"created_at"`
	UpdatedAt      time.Time `column:"updated_at"`
	CreatedBy      string    `column:"created_by"`
	UpdatedBy      string    `column:"updated_by"`
}
