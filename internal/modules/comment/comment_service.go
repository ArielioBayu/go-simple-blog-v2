package comment

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type PostChecker interface {
	CheckPostExists(ctx context.Context, postId int) (bool, error)
}

type CommentService interface {
	CreateComment(ctx context.Context, postId, userId int, request CommentRequest) error
}

type commentService struct {
	cfg         *configs.Config
	commentRepo CommentRepository
	postChecker PostChecker
}

func NewCommentService(cfg *configs.Config, commentRepo CommentRepository, postChecker PostChecker) CommentService {
	return &commentService{
		cfg:         cfg,
		commentRepo: commentRepo,
		postChecker: postChecker,
	}
}

func (s *commentService) CreateComment(ctx context.Context, postId, userId int, request CommentRequest) error {
	if s.postChecker != nil {
		exists, err := s.postChecker.CheckPostExists(ctx, postId)
		if err != nil {
			return err
		}
		if !exists {
			return constants.ErrPostNotFound
		}
	}

	now := time.Now()
	model := CommentModel{
		PostId:         postId,
		UserId:         userId,
		CommentContent: request.CommentContent,
		CreatedAt:      utils.JsonTime(now),
		UpdatedAt:      utils.JsonTime(now),
		CreatedBy:      strconv.Itoa(userId),
		UpdatedBy:      strconv.Itoa(userId),
	}

	err := s.commentRepo.CreateComment(ctx, model)
	if err != nil {
		return fmt.Errorf("service CreateComment: %w", err)
	}

	return nil
}
