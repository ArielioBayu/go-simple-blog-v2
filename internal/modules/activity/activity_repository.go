package activity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ActivityRepository interface {
	CountLikedByPostID(ctx context.Context, postId int) (int, error)
	GetActivities(ctx context.Context, postId, userId int) (*ActivityModel, error)
	UpsertActivities(ctx context.Context, model ActivityModel) error
	GetPostAuthorID(ctx context.Context, postId int) (int64, error)
	CheckPostExists(ctx context.Context, postId int) (bool, error)
}

type activityRepository struct {
	db *sql.DB
}

func NewActivityRepository(db *sql.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) GetPostAuthorID(ctx context.Context, postId int) (int64, error) {
	query := `SELECT user_id FROM posts WHERE id = ? LIMIT 1`
	var authorID int64
	err := r.db.QueryRowContext(ctx, query, postId).Scan(&authorID)
	if err != nil {
		return 0, fmt.Errorf("repository GetPostAuthorID: %w", err)
	}
	return authorID, nil
}

func (r *activityRepository) CountLikedByPostID(ctx context.Context, postId int) (int, error) {
	query := `SELECT COUNT(id) FROM activities WHERE post_id = ? AND is_liked = true`
	row := r.db.QueryRowContext(ctx, query, postId)

	var response int
	err := row.Scan(&response)
	if err != nil {
		return response, fmt.Errorf("repository CountLikedByPostID: %w", err)
	}

	return response, nil
}

func (r *activityRepository) GetActivities(ctx context.Context, postId, userId int) (*ActivityModel, error) {
	query := `SELECT id, post_id, user_id, is_liked, created_at, updated_at, created_by, updated_by FROM activities
				WHERE post_id = ? and user_id = ?`
	row := r.db.QueryRowContext(ctx, query, postId, userId)

	var data ActivityModel
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
		return nil, fmt.Errorf("repository GetActivities: %w", err)
	}
	return &data, nil
}

func (r *activityRepository) UpsertActivities(ctx context.Context, model ActivityModel) error {
	query := `INSERT INTO activities (post_id, user_id, is_liked, created_at, updated_at, created_by, updated_by)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				is_liked = VALUES(is_liked),
				updated_at = VALUES(updated_at),
				updated_by = VALUES(updated_by)`
	_, err := r.db.ExecContext(ctx, query, model.PostId, model.UserId, model.IsLiked, model.CreatedAt, model.UpdatedAt,
		model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return fmt.Errorf("repository UpsertActivities: %w", err)
	}
	return nil
}

func (r *activityRepository) CheckPostExists(ctx context.Context, postId int) (bool, error) {
	query := `SELECT 1 FROM posts WHERE id = ? LIMIT 1`
	var exists int
	err := r.db.QueryRowContext(ctx, query, postId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("repository CheckPostExists: %w", err)
	}
	return true, nil
}
