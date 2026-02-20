package posts

import (
	"context"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
)

func (r *postsRepository) CreateComment(ctx context.Context, model posts.CommentModel) error {
	query := `INSERT INTO comments (post_id, user_id, comment_content, created_at, updated_at, created_by, updated_by)
				VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, model.PostId, model.UserId, model.CommentContent, model.CreatedAt, model.UpdatedAt,
		model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return err
	}

	return nil
}
