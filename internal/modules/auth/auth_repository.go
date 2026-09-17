package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type AuthRepository interface {
	InsertRefreshToken(ctx context.Context, model RefreshTokenModel) error
	GetRefreshToken(ctx context.Context, userId int, now time.Time) (*RefreshTokenModel, error)
	GetIdRefreshToken(ctx context.Context, request RefreshTokenRequest) (*RefreshTokenModel, error)
	DeleteExpiredRefreshTokens(ctx context.Context, userId int, now time.Time) error
	DeleteRefreshTokenByUserId(ctx context.Context, userId int) error
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) InsertRefreshToken(ctx context.Context, model RefreshTokenModel) error {
	query := `INSERT INTO refresh_tokens (user_id, refresh_token, expired_at, created_at, updated_at, created_by, updated_by) 
				VALUES(?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		model.UserId, model.RefreshToken, model.ExpiredAt, model.CreatedAt, model.UpdatedAt, model.CreatedBy, model.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("repository insert refresh token: %w", err)
	}
	return nil
}

func (r *authRepository) GetRefreshToken(ctx context.Context, userId int, now time.Time) (*RefreshTokenModel, error) {
	query := `SELECT id, user_id, refresh_token, expired_at, created_at, updated_at FROM refresh_tokens 
				WHERE user_id = ? AND expired_at >= ?`
	row := r.db.QueryRowContext(ctx, query, userId, now)

	var model RefreshTokenModel
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

func (r *authRepository) GetIdRefreshToken(ctx context.Context, request RefreshTokenRequest) (*RefreshTokenModel, error) {
	query := `SELECT id, user_id, expired_at FROM refresh_tokens WHERE refresh_token = ?`
	row := r.db.QueryRowContext(ctx, query, request.Token)

	var model RefreshTokenModel
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

func (r *authRepository) DeleteExpiredRefreshTokens(ctx context.Context, userId int, now time.Time) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = ? AND expired_at < ?`
	_, err := r.db.ExecContext(ctx, query, userId, now)
	if err != nil {
		return fmt.Errorf("repository delete expired refresh tokens: %w", err)
	}
	return nil
}

func (r *authRepository) DeleteRefreshTokenByUserId(ctx context.Context, userId int) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("repository delete refresh token by user id: %w", err)
	}
	return nil
}
