package post

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/comment"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type PostRequest struct {
	PostTitle    string   `json:"post_title"`
	PostContent  string   `json:"post_content"`
	PostHashtags []string `json:"post_hashtags"`
}

type PostModel struct {
	ID           int            `column:"id"`
	UserId       int            `column:"user_id"`
	PostTitle    string         `column:"post_title"`
	PostContent  string         `column:"post_content"`
	PostHashtags string         `column:"post_hashtags"`
	CreatedAt    utils.JsonTime `column:"created_at"`
	UpdatedAt    utils.JsonTime `column:"updated_at"`
	CreatedBy    string         `column:"created_by"`
	UpdatedBy    string         `column:"updated_by"`
}

type (
	GetAllPostResponse struct {
		Data       []Data     `json:"data"`
		Pagination Pagination `json:"pagination"`
	}

	Data struct {
		ID           int            `json:"id" column:"id"`
		UserId       int            `json:"user_id" column:"user_id"`
		Username     string         `json:"username" column:"username"`
		PostTitle    string         `json:"post_title" column:"post_title"`
		PostContent  string         `json:"post_content" column:"post_content"`
		PostHashtags []string       `json:"post_hashtags" column:"post_hashtags"`
		IsLiked      bool           `json:"is_liked" column:"is_liked"`
		CreatedAt    utils.JsonTime `json:"created_at" column:"created_at"`
		UpdatedAt    utils.JsonTime `json:"updated_at" column:"updated_at"`
	}

	Pagination struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}
)

type GetPostResponse struct {
	DetailPost Data                 `json:"detail_post"`
	LikedCount int                  `json:"liked_count"`
	Comments   []comment.GetComment `json:"comments"`
}
