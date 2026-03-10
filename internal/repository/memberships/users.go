package memberships

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
)

func (r *repository) GetUser(ctx context.Context, email, username string) (*memberships.UserModel, error) {
	query := `SELECT id, email, username, created_at, updated_at, created_by, updated_by FROM users 
			WHERE email = ? OR username = ?`
	row := r.DB.QueryRowContext(ctx, query, email, username)

	var response memberships.UserModel
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

func (r *repository) GetUserById(ctx context.Context, id int) (*memberships.UserModel, error) {
	query := `SELECT id, email, username, created_at, updated_at, created_by, updated_by FROM users
				WHERE id = ?`
	row := r.DB.QueryRowContext(ctx, query, id)

	var response memberships.UserModel
	err := row.Scan(
		response.ID,
		response.Email,
		response.Username,
		response.CreatedAt,
		response.UpdatedAt,
		response.CreatedBy,
		response.UpdatedBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("repository get user by id: %w", err)
	}

	return &response, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*memberships.UserModel, error) {
	query := `SELECT id, email, password, created_at, updated_at, created_by, updated_by, username 
				FROM users WHERE email = ?`
	row := r.DB.QueryRowContext(ctx, query, email)

	var response memberships.UserModel
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

func (r *repository) CreateUser(ctx context.Context, model memberships.UserModel) error {
	query := `INSERT INTO users (id, email, password, username, created_at, updated_at, created_by, updated_by)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, model.ID, model.Email, model.Password, model.Username, model.CreatedAt, model.UpdatedAt,
		model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return fmt.Errorf("repository create user: %w", err)
	}

	return nil
}
