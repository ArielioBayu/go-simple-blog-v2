package posts

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
)

func (s *postsService) CreatePost(ctx context.Context, userId int, request posts.PostRequest) error {
	poshHashtags := strings.Join(request.PostHashtags, ",")

	time := time.Now()
	model := posts.PostModel{
		UserId:       userId,
		PostTitle:    request.PostTitle,
		PostContent:  request.PostContent,
		PostHashtags: poshHashtags,
		CreatedAt:    utils.JsonTime(time),
		UpdatedAt:    utils.JsonTime(time),
		CreatedBy:    strconv.Itoa(userId),
		UpdatedBy:    strconv.Itoa(userId),
	}

	err := s.postsRepo.CreatePost(ctx, model)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *postsService) GetAllPost(ctx context.Context, pageSize, pageIndex, userID int) (posts.GetAllPostResponse, error) {
	limit := pageSize
	offset := pageSize * (pageIndex - 1)

	response, err := s.postsRepo.GetAllPost(ctx, limit, offset, userID)
	if err != nil {
		return response, fmt.Errorf("Service GetAllPost: %w", err)
	}

	return response, nil
}

func (s *postsService) GetPostById(ctx context.Context, id, userID int) (*posts.GetPostResponse, error) {
	data, err := s.postsRepo.GetPostById(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("Service GetPostById: %w", err)
	}

	counts, err := s.postsRepo.CountLikedByPostID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Service CountLikedByPostID: %w", err)
	}

	comments, err := s.postsRepo.GetCommentById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Service GetCommentById: %w", err)
	}

	result := &posts.GetPostResponse{
		DetailPost: posts.Data{
			ID:           data.ID,
			UserId:       data.UserId,
			Username:     data.Username,
			PostTitle:    data.PostTitle,
			PostContent:  data.PostContent,
			PostHashtags: data.PostHashtags,
			IsLiked:      data.IsLiked,
			CreatedAt:    data.CreatedAt,
			UpdatedAt:    data.UpdatedAt,
			// CreatedBy:    data.UpdatedBy,
			// UpdatedBy:    data.UpdatedBy,
		},
		LikedCount: counts,
		Comments:   comments,
	}

	return result, nil
}
