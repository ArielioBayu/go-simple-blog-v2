package activity

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/notification"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type ActivityService interface {
	InsertUpdateActivities(ctx context.Context, postId, userId int, request ActivityRequest) error
	CountLikes(ctx context.Context, postId int) (int, error)
}

type activityService struct {
	cfg          *configs.Config
	activityRepo ActivityRepository
	notifSvc     notification.NotificationService
}

func NewActivityService(cfg *configs.Config, activityRepo ActivityRepository, notifSvc ...notification.NotificationService) ActivityService {
	var svc notification.NotificationService
	if len(notifSvc) > 0 {
		svc = notifSvc[0]
	}
	return &activityService{
		cfg:          cfg,
		activityRepo: activityRepo,
		notifSvc:     svc,
	}
}

func (s *activityService) InsertUpdateActivities(ctx context.Context, postId, userId int, request ActivityRequest) error {
	exists, err := s.activityRepo.CheckPostExists(ctx, postId)
	if err != nil {
		return fmt.Errorf("service InsertUpdateActivities check: %w", err)
	}
	if !exists {
		return constants.ErrPostNotFound
	}

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

	err = s.activityRepo.UpsertActivities(ctx, model)
	if err != nil {
		if strings.Contains(err.Error(), "fk_activities_post_id") || strings.Contains(err.Error(), "1452") {
			return constants.ErrPostNotFound
		}
		return fmt.Errorf("service UpsertActivities: %w", err)
	}

	if s.notifSvc != nil {
		if request.IsLiked {
			authorID, err := s.activityRepo.GetPostAuthorID(ctx, postId)
			if err == nil && authorID > 0 {
				_ = s.notifSvc.CreateNotification(ctx, notification.CreateNotificationInput{
					UserID:     authorID,
					ActorID:    int64(userId),
					Type:       notification.NotificationTypeLike,
					EntityID:   int64(postId),
					EntityType: notification.EntityTypePost,
					Message:    constants.MsgLikePost,
				})
			}
		} else {
			_ = s.notifSvc.DeleteNotification(ctx, int64(userId), notification.NotificationTypeLike, int64(postId))
		}
	}

	return nil
}

func (s *activityService) CountLikes(ctx context.Context, postId int) (int, error) {
	exists, err := s.activityRepo.CheckPostExists(ctx, postId)
	if err != nil {
		return 0, fmt.Errorf("service CountLikes check: %w", err)
	}
	if !exists {
		return 0, constants.ErrPostNotFound
	}
	return s.activityRepo.CountLikedByPostID(ctx, postId)
}
