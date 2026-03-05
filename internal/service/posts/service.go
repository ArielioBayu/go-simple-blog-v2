package posts

import (
	"context"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
	repo "github.com/ArielioBayu/go-simple-blog-v2/internal/repository/posts"
)

type PostsService interface {
	CreatePost(ctx context.Context, userId int, request posts.PostRequest) error
	CreateComment(ctx context.Context, postId, userId int, request posts.CommentRequest) error
	GetAllPost(ctx context.Context, pageSize, pageIndex int) (posts.GetAllPostResponse, error)
	GetPostById(ctx context.Context, id int) (*posts.GetPostResponse, error)
	InsertUpdateActivities(ctx context.Context, postId, userId int, request posts.ActivityRequest) error
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
