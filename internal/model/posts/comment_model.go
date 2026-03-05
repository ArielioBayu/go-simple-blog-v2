package posts

import "github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"

type CommentRequest struct {
	CommentContent string `json:"comment_content"`
}

//	Entity ke database
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
