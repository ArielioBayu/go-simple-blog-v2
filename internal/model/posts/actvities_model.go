package posts

import "time"

type ActivityRequest struct {
	IsLiked bool `json:"is_liked"`
}

type ActivityModel struct {
	ID        int       `column:"id"`
	PostId    int       `column:"post_id"`
	UserId    int       `column:"user_id"`
	IsLiked   bool      `column:"is_liked"`
	CreatedAt time.Time `column:"created_at"`
	UpdatedAt time.Time `column:"updated_at"`
	CreatedBy string    `column:"created_by"`
	UpdatedBy string    `column:"updated_by"`
}
