package posts

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

func (s *postsService) CreateComment(ctx context.Context, postId, userId int, request posts.CommentRequest) error {
	// Verifikasi keberadaan post terlebih dahulu
	post, err := s.postsRepo.GetPostById(ctx, postId, userId)
	if err != nil {
		return err
	}
	if post == nil {
		return constants.ErrPostNotFound
	}

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

	err = s.postsRepo.CreateComment(ctx, model)
	if err != nil {
		return fmt.Errorf("Service CreateComment: %w", err)
	}

	return nil
}
