package user

import (
	"context"
	"strings"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
)

type UserService interface {
	GetUserById(ctx context.Context, id int) (*UserModel, error)
	GetProfile(ctx context.Context, id int) (*ProfileResponse, error)
	UpdateProfile(ctx context.Context, id int, req UpdateProfileRequest) (*ProfileResponse, error)
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

func (s *userService) GetProfile(ctx context.Context, id int) (*ProfileResponse, error) {
	return s.userRepo.GetProfile(ctx, id)
}

func (s *userService) UpdateProfile(ctx context.Context, id int, req UpdateProfileRequest) (*ProfileResponse, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.Bio = strings.TrimSpace(req.Bio)
	req.AvatarURL = strings.TrimSpace(req.AvatarURL)
	req.BannerURL = strings.TrimSpace(req.BannerURL)

	// Jika tidak ada field yang dikirim untuk diupdate, return profile saat ini tanpa query UPDATE
	if req.Username == "" && req.Bio == "" && req.AvatarURL == "" && req.BannerURL == "" {
		return s.userRepo.GetProfile(ctx, id)
	}

	if req.Username != "" {
		isTaken, err := s.userRepo.CheckUsernameExistsExcludeSelf(ctx, id, req.Username)
		if err != nil {
			return nil, err
		}
		if isTaken {
			return nil, constants.ErrUsernameOrEmailAlreadyExists
		}
	}

	return s.userRepo.UpdateProfile(ctx, id, req)
}
