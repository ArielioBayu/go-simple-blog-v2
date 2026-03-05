package posts

import (
	"context"
	"fmt"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
)

func (r *postsRepository) GetCommentById(ctx context.Context, postId int) ([]posts.GetComment, error) {
	query := `SELECT c.id, c.user_id, u.username, c.comment_content
				FROM comments as c
				JOIN users as u ON c.user_id = u.id
				WHERE c.post_id = ?`
	rows, err := r.DB.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, fmt.Errorf("Repository GetCommentById: %w", err)
	}
	defer rows.Close()

	response := make([]posts.GetComment, 0)
	for rows.Next() {
		var model posts.GetComment
		var username string

		err := rows.Scan(
			&model.ID,
			&model.UserId,
			&username,
			&model.CommentContent,
		)
		if err != nil {
			return nil, fmt.Errorf("Repository GetCommentById: %w", err)
		}

		response = append(response, posts.GetComment{
			ID:             model.ID,
			UserId:         model.UserId,
			Username:       username,
			CommentContent: model.CommentContent,
		})
	}
	return response, nil
}

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
