package posts

import "time"

type CommentRequest struct {
	CommentContent string `json:"comment_content"`
}

//	Entity ke database
type CommentModel struct {
	ID             int       `column:"id"`
	PostId         int       `column:"post_id"`
	UserId         int       `column:"user_id"`
	CommentContent string    `column:"comment_content"`
	CreatedAt      time.Time `column:"created_at"`
	UpdatedAt      time.Time `column:"updated_at"`
	CreatedBy      string    `column:"created_by"`
	UpdatedBy      string    `column:"updated_by"`
}
