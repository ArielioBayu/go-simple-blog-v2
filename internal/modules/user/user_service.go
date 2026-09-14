package user

import (
	"context"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
)

type UserService interface {
	GetUserById(ctx context.Context, id int) (*UserModel, error)
}

type userService struct {
	cfg      *configs.Config
	userRepo UserRepository
}

func NewUserService(cfg *configs.Config, userRepo UserRepository) UserService {
	return &userService{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

func (s *userService) GetUserById(ctx context.Context, id int) (*UserModel, error) {
	return s.userRepo.GetUserById(ctx, id)
}
