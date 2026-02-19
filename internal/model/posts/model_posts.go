package posts

import "time"

//	entity dari client
type PostsRequest struct {
	PostTitle    string   `json:"post_title"`
	PostContent  string   `json:"post_content"`
	PostHashtags []string `json:"post_hashtags"`
}

//	entity ke databasse
type PostModel struct {
	ID           int       `column:"id"`
	UserId       int       `column:"user_id"`
	PostTitle    string    `column:"post_title"`
	PostContent  string    `column:"post_content"`
	PostHashtags string    `column:"post_hashtags"`
	CreatedAt    time.Time `column:"created_at"`
	UpdatedAt    time.Time `column:"updated_at"`
	CreatedBy    string    `column:"created_by"`
	UpdatedBy    string    `column:"updated_by"`
}
