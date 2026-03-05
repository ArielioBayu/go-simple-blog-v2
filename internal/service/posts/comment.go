package posts

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

func (s *postsService) CreateComment(ctx context.Context, postId, userId int, request posts.CommentRequest) error {
	time := time.Now()
	model := posts.CommentModel{
		PostId:         postId,
		UserId:         userId,
		CommentContent: request.CommentContent,
		CreatedAt:      utils.JsonTime(time),
		UpdatedAt:      utils.JsonTime(time),
		CreatedBy:      strconv.Itoa(userId),
		UpdatedBy:      strconv.Itoa(userId),
	}

	err := s.postsRepo.CreateComment(ctx, model)
	if err != nil {
		return fmt.Errorf("Service CreateComment: %w", err)
	}

	return nil
}
