package posts

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
)

func (r *postsRepository) GetActivities(ctx context.Context, postId, userId int) (*posts.ActivityModel, error) {
	query := `SELECT id, post_id, user_id, is_liked, created_at, updated_at, created_by, updated_by FROM activities
				WHERE post_id = ? and user_id = ?`
	row := r.DB.QueryRowContext(ctx, query, postId, userId)

	var data posts.ActivityModel
	err := row.Scan(
		&data.ID,
		&data.PostId,
		&data.UserId,
		&data.IsLiked,
		&data.CreatedAt,
		&data.UpdatedAt,
		&data.CreatedBy,
		&data.UpdatedBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetActivities: %w", err)
	}
	return &data, nil
}

func (r *postsRepository) CreateActivities(ctx context.Context, model posts.ActivityModel) error {
	query := `INSERT INTO activities (post_id, user_id, is_liked, created_at, updated_at, created_by, updated_by)
			VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, model.PostId, model.UserId, model.IsLiked, model.CreatedAt, model.UpdatedAt,
		model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return fmt.Errorf("CreateActivities: %w", err)
	}
	return nil
}

func (r *postsRepository) UpdateActivities(ctx context.Context, model posts.ActivityModel) error {
	query := `UPDATE activities SET is_liked = ?, updated_at = ?, updated_by = ? WHERE post_id = ? AND user_id = ?`
	_, err := r.DB.ExecContext(ctx, query, model.IsLiked, model.UpdatedAt, model.UpdatedBy, model.PostId, model.UserId)
	if err != nil {
		return fmt.Errorf("UpdateActivity: %w", err)
	}
	return nil
}
