package posts

import "github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"

type ActivityRequest struct {
	IsLiked bool `json:"is_liked"`
}

type ActivityModel struct {
	ID        int            `column:"id"`
	PostId    int            `column:"post_id"`
	UserId    int            `column:"user_id"`
	IsLiked   bool           `column:"is_liked"`
	CreatedAt utils.JsonTime `column:"created_at"`
	UpdatedAt utils.JsonTime `column:"updated_at"`
	CreatedBy string         `column:"created_by"`
	UpdatedBy string         `column:"updated_by"`
}
