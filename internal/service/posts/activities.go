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
	now := time.Now()
	model := posts.ActivityModel{
		PostId:    postId,
		UserId:    userId,
		IsLiked:   request.IsLiked,
		CreatedAt: utils.JsonTime(now),
		UpdatedAt: utils.JsonTime(now),
		CreatedBy: strconv.Itoa(userId),
		UpdatedBy: strconv.Itoa(userId),
	}

	err := s.postsRepo.UpsertActivities(ctx, model)
	if err != nil {
		return fmt.Errorf("service UpsertActivities: %w", err)
	}

	return nil
}
