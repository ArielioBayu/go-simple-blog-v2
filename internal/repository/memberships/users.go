package memberships

import (
	"context"
	"database/sql"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
)

type MembershipRepository interface {
	GetUser(ctx context.Context, email, username string) (*memberships.UserModel, error)
	CreateUser(ctx context.Context, model memberships.UserModel) error
}

type repository struct {
	DB *sql.DB
}

func NewMembershipsRepository(db *sql.DB) MembershipRepository {
	return &repository{db}
}

func (r *repository) GetUser(ctx context.Context, email, username string) (*memberships.UserModel, error) {
	query := `SELECT id, email, username, created_at, updated_at, created_by, updated_by FROM users WHERE email = ? AND username = ?`
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
		return nil, err
	}

	return &response, nil
}

func (r *repository) CreateUser(ctx context.Context, model memberships.UserModel) error {
	query := `INSERT INTO users (id, email, password, username, created_at, updated_at, created_by, updated_by)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, model.ID, model.Email, model.Password, model.Username, model.CreatedAt, model.UpdatedAt,
		model.CreatedBy, model.UpdatedBy)
	if err != nil {
		return err
	}

	return nil
}
