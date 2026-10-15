package user

import (
	"context"
	"math"
	"strings"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type UserService interface {
	GetUserById(ctx context.Context, id int) (*UserModel, error)
	GetProfile(ctx context.Context, id int) (*ProfileResponse, error)
	GetProfileByID(ctx context.Context, id int) (*ProfileResponse, error)
	UpdateProfile(ctx context.Context, id int, req UpdateProfileRequest) (*ProfileResponse, error)
	UpdatePrivacy(ctx context.Context, id int, req UpdatePrivacyRequest) error
	SearchUsers(ctx context.Context, viewerID int, query string, page, limit int) (*SearchUserResponse, error)
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

func (s *userService) GetProfileByID(ctx context.Context, id int) (*ProfileResponse, error) {
	return s.userRepo.GetProfile(ctx, id)
}

func (s *userService) UpdateProfile(ctx context.Context, id int, req UpdateProfileRequest) (*ProfileResponse, error) {
	// Jika tidak ada field yang dikirim untuk diupdate, return profile saat ini tanpa query UPDATE
	if req.Username == nil && req.Bio == nil && req.AvatarURL == nil && req.BannerURL == nil {
		return s.userRepo.GetProfile(ctx, id)
	}

	if req.Username != nil {
		trimmed := strings.TrimSpace(*req.Username)
		req.Username = &trimmed
		if *req.Username != "" {
			isTaken, err := s.userRepo.CheckUsernameExistsExcludeSelf(ctx, id, *req.Username)
			if err != nil {
				return nil, err
			}
			if isTaken {
				return nil, constants.ErrUsernameOrEmailAlreadyExists
			}
		}
	}

	if req.Bio != nil {
		trimmed := strings.TrimSpace(*req.Bio)
		req.Bio = &trimmed
	}

	if req.AvatarURL != nil {
		trimmed := strings.TrimSpace(*req.AvatarURL)
		req.AvatarURL = &trimmed
	}

	if req.BannerURL != nil {
		trimmed := strings.TrimSpace(*req.BannerURL)
		req.BannerURL = &trimmed

		// Jika user menghapus banner (string kosong), bersihkan juga file fisik lama di disk jika merupakan upload lokal
		if *req.BannerURL == "" {
			currentProfile, _ := s.userRepo.GetProfile(ctx, id)
			if currentProfile != nil && currentProfile.BannerURL != "" && strings.HasPrefix(currentProfile.BannerURL, "uploads/") {
				_ = utils.RemoveUploadFile(currentProfile.BannerURL)
			}
		}
	}

	return s.userRepo.UpdateProfile(ctx, id, req)
}

func (s *userService) UpdatePrivacy(ctx context.Context, id int, req UpdatePrivacyRequest) error {
	if req.IsPrivate == nil {
		return constants.ErrDataNotFound
	}
	return s.userRepo.UpdatePrivacy(ctx, id, *req.IsPrivate)
}

func (s *userService) SearchUsers(ctx context.Context, viewerID int, query string, page, limit int) (*SearchUserResponse, error) {
	query = strings.TrimSpace(query)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	if query == "" {
		return &SearchUserResponse{
			Data: []SearchUserItem{},
			Pagination: SearchUserPagination{
				Page:      page,
				Limit:     limit,
				TotalPage: 0,
				TotalData: 0,
			},
		}, nil
	}

	offset := (page - 1) * limit
	items, total, err := s.userRepo.SearchUsers(ctx, viewerID, query, limit, offset)
	if err != nil {
		return nil, err
	}

	totalPage := 0
	if total > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &SearchUserResponse{
		Data: items,
		Pagination: SearchUserPagination{
			Page:      page,
			Limit:     limit,
			TotalPage: totalPage,
			TotalData: total,
		},
	}, nil
}

