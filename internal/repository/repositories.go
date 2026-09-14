package repository

import (
	"database/sql"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/activity"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/auth"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/comment"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/post"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/user"
)

type Repositories struct {
	Auth     auth.AuthRepository
	User     user.UserRepository
	Post     post.PostRepository
	Comment  comment.CommentRepository
	Activity activity.ActivityRepository
}

func InitRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Auth:     auth.NewAuthRepository(db),
		User:     user.NewUserRepository(db),
		Post:     post.NewPostRepository(db),
		Comment:  comment.NewCommentRepository(db),
		Activity: activity.NewActivityRepository(db),
	}
}
