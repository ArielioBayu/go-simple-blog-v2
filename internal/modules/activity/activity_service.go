package activity

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type ActivityService interface {
	InsertUpdateActivities(ctx context.Context, postId, userId int, request ActivityRequest) error
	CountLikes(ctx context.Context, postId int) (int, error)
}

type activityService struct {
	cfg          *configs.Config
	activityRepo ActivityRepository
}

func NewActivityService(cfg *configs.Config, activityRepo ActivityRepository) ActivityService {
	return &activityService{
		cfg:          cfg,
		activityRepo: activityRepo,
	}
}

func (s *activityService) InsertUpdateActivities(ctx context.Context, postId, userId int, request ActivityRequest) error {
	now := time.Now()
	model := ActivityModel{
		PostId:    postId,
		UserId:    userId,
		IsLiked:   request.IsLiked,
		CreatedAt: utils.JsonTime(now),
		UpdatedAt: utils.JsonTime(now),
		CreatedBy: strconv.Itoa(userId),
		UpdatedBy: strconv.Itoa(userId),
	}

	err := s.activityRepo.UpsertActivities(ctx, model)
	if err != nil {
		return fmt.Errorf("service UpsertActivities: %w", err)
	}

	return nil
}

func (s *activityService) CountLikes(ctx context.Context, postId int) (int, error) {
	return s.activityRepo.CountLikedByPostID(ctx, postId)
}
