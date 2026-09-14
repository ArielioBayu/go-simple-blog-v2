package user

import (
	"context"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	GetUser(ctx context.Context, email, username string) (*UserModel, error)
	GetUserById(ctx context.Context, id int) (*UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (*UserModel, error)
	CreateUser(ctx context.Context, model UserModel) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUser(ctx context.Context, email, username string) (*UserModel, error) {
	query := `SELECT id, email, username, created_at, updated_at, created_by, updated_by FROM users 
			WHERE email = ? OR username = ?`
	row := r.db.QueryRowContext(ctx, query, email, username)

	var response UserModel
	err := row.Scan(
		&response.ID,
		&response.Email,
		&response.Username,
		&response.CreatedAt,
		&response.UpdatedAt,
		&response.CreatedBy,
		&response.UpdatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("repository get user: %w", err)
	}

	return &response, nil
}

func (r *userRepository) GetUserById(ctx context.Context, id int) (*UserModel, error) {
	query := `SELECT id, email, username, created_at, updated_at, created_by, updated_by FROM users
				WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var response UserModel
	err := row.Scan(
		&response.ID,
		&response.Email,
		&response.Username,
		&response.CreatedAt,
		&response.UpdatedAt,
		&response.CreatedBy,
		&response.UpdatedBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("repository get user by id: %w", err)
	}

	return &response, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*UserModel, error) {
	query := `SELECT id, email, password, created_at, updated_at, created_by, updated_by, username 
				FROM users WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email)

	var response UserModel
	err := row.Scan(
		&response.ID,
		&response.Email,
		&response.Password,
		&response.CreatedAt,
		&response.UpdatedAt,
		&response.CreatedBy,
		&response.UpdatedBy,
		&response.Username,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("repository get user by email: %w", err)
	}
	return &response, nil
}

func (r *userRepository) CreateUser(ctx context.Context, model UserModel) error {
	query := `INSERT INTO users (email, password, username, created_at, updated_at, created_by, updated_by)
	VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, model.Email, model.Password, model.Username, model.CreatedAt, model.UpdatedAt,
		model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return fmt.Errorf("repository create user: %w", err)
	}

	return nil
}
