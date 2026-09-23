package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository interface {
	GetUser(ctx context.Context, email, username string) (*UserModel, error)
	GetUserById(ctx context.Context, id int) (*UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
	CreateUser(ctx context.Context, model *UserModel) error
	GetProfile(ctx context.Context, id int) (*ProfileResponse, error)
	UpdateProfile(ctx context.Context, id int, req UpdateProfileRequest) (*ProfileResponse, error)
	CheckUsernameExistsExcludeSelf(ctx context.Context, id int, username string) (bool, error)
	UpdateUserVerification(ctx context.Context, userID int, isVerified bool) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUser(ctx context.Context, email, username string) (*UserModel, error) {
	query := `SELECT id, email, username, COALESCE(bio, ''), COALESCE(avatar_url, ''), COALESCE(banner_url, ''), is_verified, created_at, updated_at, created_by, updated_by FROM users 
			WHERE email = ? OR username = ?`
	row := r.db.QueryRowContext(ctx, query, email, username)

	var response UserModel
	err := row.Scan(
		&response.ID,
		&response.Email,
		&response.Username,
		&response.Bio,
		&response.AvatarURL,
		&response.BannerURL,
		&response.IsVerified,
		&response.CreatedAt,
		&response.UpdatedAt,
		&response.CreatedBy,
		&response.UpdatedBy,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repository get user: %w", err)
	}

	return &response, nil
}

func (r *userRepository) GetUserById(ctx context.Context, id int) (*UserModel, error) {
	query := `SELECT id, email, username, COALESCE(bio, ''), COALESCE(avatar_url, ''), COALESCE(banner_url, ''), is_verified, created_at, updated_at, created_by, updated_by FROM users
				WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var response UserModel
	err := row.Scan(
		&response.ID,
		&response.Email,
		&response.Username,
		&response.Bio,
		&response.AvatarURL,
		&response.BannerURL,
		&response.IsVerified,
		&response.CreatedAt,
		&response.UpdatedAt,
		&response.CreatedBy,
		&response.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repository get user by id: %w", err)
	}

	return &response, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*UserModel, error) {
	query := `SELECT id, email, password, username, COALESCE(bio, ''), COALESCE(avatar_url, ''), COALESCE(banner_url, ''), is_verified, created_at, updated_at, created_by, updated_by 
				FROM users WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email)

	var response UserModel
	err := row.Scan(
		&response.ID,
		&response.Email,
		&response.Password,
		&response.Username,
		&response.Bio,
		&response.AvatarURL,
		&response.BannerURL,
		&response.IsVerified,
		&response.CreatedAt,
		&response.UpdatedAt,
		&response.CreatedBy,
		&response.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repository get user by email: %w", err)
	}
	return &response, nil
}

func (r *userRepository) CreateUser(ctx context.Context, model *UserModel) error {
	query := `INSERT INTO users (email, password, username, bio, avatar_url, banner_url, is_verified, created_at, updated_at, created_by, updated_by)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, model.Email, model.Password, model.Username, model.Bio, model.AvatarURL, model.BannerURL, model.IsVerified,
		model.CreatedAt, model.UpdatedAt, model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return fmt.Errorf("repository create user: %w", err)
	}

	id, err := res.LastInsertId()
	if err == nil {
		model.ID = id
	}

	return nil
}

func (r *userRepository) UpdateUserVerification(ctx context.Context, userID int, isVerified bool) error {
	query := `UPDATE users SET is_verified = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, isVerified, userID)
	if err != nil {
		return fmt.Errorf("repository UpdateUserVerification: %w", err)
	}
	return nil
}

func (r *userRepository) GetProfile(ctx context.Context, id int) (*ProfileResponse, error) {
	query := `SELECT 
				u.id, 
				u.username, 
				u.email, 
				COALESCE(u.bio, ''), 
				COALESCE(u.avatar_url, ''), 
				COALESCE(u.banner_url, ''), 
				COALESCE(u.is_verified, false),
				u.created_at,
				(SELECT COUNT(p.id) FROM posts p WHERE p.user_id = u.id) AS stories_count,
				(SELECT COUNT(a.id) FROM activities a JOIN posts p ON a.post_id = p.id WHERE p.user_id = u.id AND a.is_liked = true) AS likes_count
			  FROM users u 
			  WHERE u.id = ?`

	var profile ProfileResponse
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&profile.ID,
		&profile.Username,
		&profile.Email,
		&profile.Bio,
		&profile.AvatarURL,
		&profile.BannerURL,
		&profile.IsVerified,
		&profile.CreatedAt,
		&profile.Stats.StoriesCount,
		&profile.Stats.LikesCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repository GetProfile: %w", err)
	}

	return &profile, nil
}

func (r *userRepository) CheckUsernameExistsExcludeSelf(ctx context.Context, id int, username string) (bool, error) {
	query := `SELECT 1 FROM users WHERE username = ? AND id != ? LIMIT 1`
	var exists int
	err := r.db.QueryRowContext(ctx, query, username, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("repository CheckUsernameExistsExcludeSelf: %w", err)
	}
	return true, nil
}

func (r *userRepository) UpdateProfile(ctx context.Context, id int, req UpdateProfileRequest) (*ProfileResponse, error) {
	query := `UPDATE users 
	          SET username = COALESCE(NULLIF(?, ''), username),
	              bio = CASE WHEN ? != '' THEN ? ELSE bio END, 
	              avatar_url = CASE WHEN ? != '' THEN ? ELSE avatar_url END, 
	              banner_url = CASE WHEN ? != '' THEN ? ELSE banner_url END, 
	              updated_at = NOW() 
	          WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, req.Username, req.Bio, req.Bio, req.AvatarURL, req.AvatarURL, req.BannerURL, req.BannerURL, id)
	if err != nil {
		return nil, fmt.Errorf("repository UpdateProfile: %w", err)
	}

	return r.GetProfile(ctx, id)
}
