package posts

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
	repo "github.com/ArielioBayu/go-simple-blog-v2/internal/repository/posts"
)

type PostsService interface {
	CreatePost(ctx context.Context, userId int, request posts.PostsRequest) error
}

type postsService struct {
	cfg       *configs.Config
	postsRepo repo.PostsRepository
}

func NewPostsService(cfg *configs.Config, postRepo repo.PostsRepository) PostsService {
	return &postsService{
		cfg:       cfg,
		postsRepo: postRepo,
	}
}

func (s *postsService) CreatePost(ctx context.Context, userId int, request posts.PostsRequest) error {
	poshHashtags := strings.Join(request.PostHashtags, ",")

	time := time.Now()
	model := posts.PostModel{
		UserId:       userId,
		PostTitle:    request.PostTitle,
		PostContent:  request.PostContent,
		PostHashtags: poshHashtags,
		CreatedAt:    time,
		UpdatedAt:    time,
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
