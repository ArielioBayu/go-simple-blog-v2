package saves

import (
	"context"
	"fmt"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
)

type SavesService interface {
	SavePost(ctx context.Context, postId, userID int, request SaveRequest) error
	GetSavedPosts(ctx context.Context, userID, page, limit int) (*GetSavedPostsResponse, error)
	GetSavedPostIDs(ctx context.Context, userID int) ([]int, error)
}

type savesService struct {
	cfg       *configs.Config
	savesRepo SavesRepository
}

func NewSavesService(cfg *configs.Config, savesRepo SavesRepository) SavesService {
	return &savesService{
		cfg:       cfg,
		savesRepo: savesRepo,
	}
}

func (s *savesService) SavePost(ctx context.Context, postId, userID int, request SaveRequest) error {
	exists, err := s.savesRepo.CheckPostExists(ctx, postId)
	if err != nil {
		return fmt.Errorf("service SavePost: %w", err)
	}
	if !exists {
		return constants.ErrPostNotFound
	}

	if request.IsSaved {
		if err := s.savesRepo.SavePost(ctx, userID, postId); err != nil {
			return fmt.Errorf("service SavePost save: %w", err)
		}
	} else {
		if err := s.savesRepo.UnsavePost(ctx, userID, postId); err != nil {
			return fmt.Errorf("service SavePost unsave: %w", err)
		}
	}

	return nil
}

func (s *savesService) GetSavedPosts(ctx context.Context, userID, page, limit int) (*GetSavedPostsResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	items, totalData, err := s.savesRepo.GetSavedPosts(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("service GetSavedPosts: %w", err)
	}

	totalPage := 0
	if totalData > 0 {
		totalPage = (totalData + limit - 1) / limit
	}

	return &GetSavedPostsResponse{
		Data: items,
		Pagination: SavedPagination{
			Page:      page,
			Limit:     limit,
			TotalPage: totalPage,
			TotalData: totalData,
		},
	}, nil
}

func (s *savesService) GetSavedPostIDs(ctx context.Context, userID int) ([]int, error) {
	ids, err := s.savesRepo.GetSavedPostIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service GetSavedPostIDs: %w", err)
	}
	return ids, nil
}

