package repository

import (
	"database/sql"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/activity"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/auth"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/comment"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/post"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/saves"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/upload"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/user"
)

type Repositories struct {
	Auth     auth.AuthRepository
	User     user.UserRepository
	Post     post.PostRepository
	Comment  comment.CommentRepository
	Activity activity.ActivityRepository
	Upload   upload.UploadRepository
	Saves    saves.SavesRepository
}

func InitRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Auth:     auth.NewAuthRepository(db),
		User:     user.NewUserRepository(db),
		Post:     post.NewPostRepository(db),
		Comment:  comment.NewCommentRepository(db),
		Activity: activity.NewActivityRepository(db),
		Upload:   upload.NewUploadRepository(db),
		Saves:    saves.NewSavesRepository(db),
	}
}
