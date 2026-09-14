package post

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/activity"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/comment"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

type PostService interface {
	CreatePost(ctx context.Context, userId int, request PostRequest) error
	GetAllPost(ctx context.Context, pageSize, pageIndex, userID int) (GetAllPostResponse, error)
	GetPostById(ctx context.Context, id, userID int) (*GetPostResponse, error)
}

type postService struct {
	cfg          *configs.Config
	postRepo     PostRepository
	commentRepo  comment.CommentRepository
	activityRepo activity.ActivityRepository
}

func NewPostService(
	cfg *configs.Config,
	postRepo PostRepository,
	commentRepo comment.CommentRepository,
	activityRepo activity.ActivityRepository,
) PostService {
	return &postService{
		cfg:          cfg,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
		activityRepo: activityRepo,
	}
}

func (s *postService) CreatePost(ctx context.Context, userId int, request PostRequest) error {
	postHashtags := strings.Join(request.PostHashtags, ",")

	now := time.Now()
	model := PostModel{
		UserId:       userId,
		PostTitle:    request.PostTitle,
		PostContent:  request.PostContent,
		PostHashtags: postHashtags,
		CreatedAt:    utils.JsonTime(now),
		UpdatedAt:    utils.JsonTime(now),
		CreatedBy:    strconv.Itoa(userId),
		UpdatedBy:    strconv.Itoa(userId),
	}

	err := s.postRepo.CreatePost(ctx, model)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *postService) GetAllPost(ctx context.Context, pageSize, pageIndex, userID int) (GetAllPostResponse, error) {
	limit := pageSize
	offset := pageSize * (pageIndex - 1)

	response, err := s.postRepo.GetAllPost(ctx, limit, offset, userID)
	if err != nil {
		return response, fmt.Errorf("service GetAllPost: %w", err)
	}

	return response, nil
}

func (s *postService) GetPostById(ctx context.Context, id, userID int) (*GetPostResponse, error) {
	data, err := s.postRepo.GetPostById(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("service GetPostById: %w", err)
	}

	var counts int
	if s.activityRepo != nil {
		counts, err = s.activityRepo.CountLikedByPostID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("service CountLikedByPostID: %w", err)
		}
	}

	var comments []comment.GetComment
	if s.commentRepo != nil {
		comments, err = s.commentRepo.GetCommentById(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("service GetCommentById: %w", err)
		}
	}

	result := &GetPostResponse{
		DetailPost: Data{
			ID:           data.ID,
			UserId:       data.UserId,
			Username:     data.Username,
			PostTitle:    data.PostTitle,
			PostContent:  data.PostContent,
			PostHashtags: data.PostHashtags,
			IsLiked:      data.IsLiked,
			CreatedAt:    data.CreatedAt,
			UpdatedAt:    data.UpdatedAt,
		},
		LikedCount: counts,
		Comments:   comments,
	}

	return result, nil
}
