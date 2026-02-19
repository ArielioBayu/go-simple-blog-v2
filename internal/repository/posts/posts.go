package posts

import (
	"context"
	"database/sql"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
)

type PostsRepository interface {
	CreatePost(ctx context.Context, model posts.PostModel) error
}

type postsRepository struct {
	DB *sql.DB
}

func NewPostsRepository(db *sql.DB) *postsRepository {
	return &postsRepository{
		DB: db,
	}
}

func (r *postsRepository) CreatePost(ctx context.Context, model posts.PostModel) error {
	query := `INSERT INTO posts (user_id, post_title, post_content, post_hashtags, created_at, updated_at,
				created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, model.UserId, model.PostTitle, model.PostContent, model.PostHashtags, model.CreatedAt, model.UpdatedAt,
		model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return err
	}

	return nil
}
