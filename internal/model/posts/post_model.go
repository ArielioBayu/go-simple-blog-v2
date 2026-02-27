package posts

import "time"

//	entity dari client
type PostRequest struct {
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

type (
	GetAllPostResponse struct {
		Data       []Data     `json:"data"`
		Pagination Pagination `json:"pagination"`
	}

	Data struct {
		ID           int       `json:"id"`
		UserId       int       `json:"user_id"`
		Username     string    `json:"username"`
		PostTitle    string    `json:"post_title"`
		PostContent  string    `json:"post_content"`
		PostHashtags []string  `json:"post_hashtags"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		CreatedBy    string    `json:"created_by"`
		UpdatedBy    string    `json:"updated_by"`
	}

	Pagination struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}
)
