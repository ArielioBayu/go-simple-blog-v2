package posts

import (
	"context"
	"database/sql"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
)

type PostsRepository interface {
	CountLikedByPostID(ctx context.Context, postId int) (int, error)
	CreatePost(ctx context.Context, model posts.PostModel) error
	CreateComment(ctx context.Context, model posts.CommentModel) error
	CreateActivities(ctx context.Context, model posts.ActivityModel) error
	GetActivities(ctx context.Context, postId, userId int) (*posts.ActivityModel, error)
	GetAllPost(ctx context.Context, limit, offset, userID int) (posts.GetAllPostResponse, error)
	GetPostById(ctx context.Context, id, userID int) (*posts.Data, error)
	GetCommentById(ctx context.Context, postId int) ([]posts.GetComment, error)
	UpdateActivities(ctx context.Context, model posts.ActivityModel) error
}

type postsRepository struct {
	DB *sql.DB
}

func NewPostsRepository(db *sql.DB) *postsRepository {
	return &postsRepository{
		DB: db,
	}
}
