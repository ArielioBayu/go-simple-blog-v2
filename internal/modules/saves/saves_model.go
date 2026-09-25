package saves

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/post"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type SaveRequest struct {
	IsSaved bool `json:"is_saved"`
}

type SaveResponse struct {
	PostID  int  `json:"post_id"`
	IsSaved bool `json:"is_saved"`
}

type SavedPostItem struct {
	ID           int              `json:"id"`
	UserID       int              `json:"user_id"`
	Username     string           `json:"username"`
	AvatarURL    string           `json:"avatar_url"`
	PostTitle    string           `json:"post_title"`
	PostContent  string           `json:"post_content"`
	PostHashtags []string         `json:"post_hashtags"`
	IsLiked      bool             `json:"is_liked"`
	IsSaved      bool             `json:"is_saved"`
	UploadID     *int64           `json:"upload_id,omitempty"`
	FilePath     string           `json:"file_path"`
	Filepath     string           `json:"filepath"`
	FileType     string           `json:"file_type"`
	FileSize     int64            `json:"file_size"`
	Media        []post.PostMedia `json:"media"`
	CreatedAt    utils.JsonTime   `json:"created_at"`
	UpdatedAt    utils.JsonTime   `json:"updated_at"`
	SavedAt      utils.JsonTime   `json:"saved_at"`
}

type SavedPagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}

type GetSavedPostsResponse struct {
	Data       []SavedPostItem `json:"data"`
	Pagination SavedPagination `json:"pagination"`
}

type UserSavedPostModel struct {
	ID        int64          `column:"id"`
	UserID    int64          `column:"user_id"`
	PostID    int64          `column:"post_id"`
	CreatedAt utils.JsonTime `column:"created_at"`
}
