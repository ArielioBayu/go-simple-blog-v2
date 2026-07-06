package memberships

import (
	"context"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
	repo "github.com/ArielioBayu/go-simple-blog-v2/internal/repository/memberships"
)

type MembershipsService interface {
	GetUser(ctx context.Context, request memberships.SignUpRequest) (*memberships.UserModel, error)
	GetIdRefreshToken(ctx context.Context, request memberships.RefreshTokenRequest) (*memberships.RefreshTokenModel, error)
	SignUp(ctx context.Context, request memberships.SignUpRequest) error
	SignIn(ctx context.Context, request memberships.SignInRequest) (string, string, error)
	ValidateRefreshToken(ctx context.Context, userId int, request memberships.RefreshTokenRequest) (string, error)
}

type membershipsService struct {
	cfg             *configs.Config
	membershipsRepo repo.MembershipRepository
}

func NewMembershipsService(cfg *configs.Config, membershipsRepo repo.MembershipRepository) MembershipsService {
	return &membershipsService{cfg, membershipsRepo}
}
