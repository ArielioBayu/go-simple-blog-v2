package memberships

import (
	"context"
	"database/sql"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
)

type MembershipRepository interface {
	GetUser(ctx context.Context, email, username string) (*memberships.UserModel, error)
	GetUserById(ctx context.Context, id int) (*memberships.UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (*memberships.UserModel, error)
	GetRefreshToken(ctx context.Context, userId int, now time.Time) (*memberships.RefreshTokenModel, error)
	CreateUser(ctx context.Context, model memberships.UserModel) error
	InsertRefreshToken(ctx context.Context, model memberships.RefreshTokenModel) error
}

type repository struct {
	DB *sql.DB
}

func NewMembershipsRepository(db *sql.DB) MembershipRepository {
	return &repository{db}
}
