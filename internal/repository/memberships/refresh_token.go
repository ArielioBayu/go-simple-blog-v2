package memberships

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
)

func (r *repository) InsertRefreshToken(ctx context.Context, model memberships.RefreshTokenModel) error {
	query := `INSERT INTO refresh_tokens (user_id, refresh_token, expired_at, created_at, updated_at, created_by, updated_by) 
				VALUES(?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query,
		model.UserId, model.RefreshToken, model.ExpiredAt, model.CreatedAt, model.UpdatedAt, model.CreatedBy, model.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("repository insert refresh token: %w", err)
	}
	return nil
}

func (r *repository) GetRefreshToken(ctx context.Context, userId int, now time.Time) (*memberships.RefreshTokenModel, error) {
	query := `SELECT id, user_id, refresh_token, expired_at, created_at, updated_at FROM refresh_tokens 
				WHERE user_id = ? AND expired_at >= ?`
	row := r.DB.QueryRowContext(ctx, query, userId, now)

	var model memberships.RefreshTokenModel
	err := row.Scan(
		&model.ID,
		&model.UserId,
		&model.RefreshToken,
		&model.ExpiredAt,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("repository get refresh token: %w", err)
	}

	return &model, nil
}

func (r *repository) GetIdRefreshToken(ctx context.Context, request memberships.RefreshTokenRequest) (*memberships.RefreshTokenModel, error) {
	query := `SELECT id, user_id, expired_at FROM refresh_tokens WHERE refresh_token = ?`
	row := r.DB.QueryRowContext(ctx, query, request.Token)

	var model memberships.RefreshTokenModel
	err := row.Scan(
		&model.ID,
		&model.UserId,
		&model.ExpiredAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repository get id refresh token: %w", err)
	}

	return &model, nil
}
