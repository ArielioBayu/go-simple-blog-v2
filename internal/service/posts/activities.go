package posts

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

func (s *postsService) InsertUpdateActivities(ctx context.Context, postId, userId int, request posts.ActivityRequest) error {
	userActivity, err := s.postsRepo.GetActivities(ctx, postId, userId)
	if err != nil {
		return fmt.Errorf("Service GetActivities: %w", err)
	}

	time := time.Now()
	model := posts.ActivityModel{
		PostId:    postId,
		UserId:    userId,
		IsLiked:   request.IsLiked,
		CreatedAt: utils.JsonTime(time),
		UpdatedAt: utils.JsonTime(time),
		CreatedBy: strconv.Itoa(userId),
		UpdatedBy: strconv.Itoa(userId),
	}

	if userActivity == nil {
		err := s.postsRepo.CreateActivities(ctx, model)
		if err != nil {
			return fmt.Errorf("service CreateActivities: %w", err)
		}
	} else {
		err := s.postsRepo.UpdateActivities(ctx, model)
		if err != nil {
			return fmt.Errorf("service UpdateActivities: %w", err)
		}
	}

	return nil
}
