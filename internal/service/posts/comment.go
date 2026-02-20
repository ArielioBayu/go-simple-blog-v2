package posts

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
)

func (s *postsService) CreateComment(ctx context.Context, postId, userId int, request posts.CommentRequest) error {
	time := time.Now()
	model := posts.CommentModel{
		PostId:         postId,
		UserId:         userId,
		CommentContent: request.CommentContent,
		CreatedAt:      time,
		UpdatedAt:      time,
		CreatedBy:      strconv.Itoa(userId),
		UpdatedBy:      strconv.Itoa(userId),
	}

	err := s.postsRepo.CreateComment(ctx, model)
	if err != nil {
		log.Println("error :", err)
		return err
	}

	return nil
}
