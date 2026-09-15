package service

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/activity"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/auth"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/comment"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/post"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/upload"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/user"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/repository"
)

type Services struct {
	Auth     auth.AuthService
	User     user.UserService
	Post     post.PostService
	Comment  comment.CommentService
	Activity activity.ActivityService
	Upload   upload.UploadService
}

func InitServices(repos *repository.Repositories, cfg *configs.Config) *Services {
	return &Services{
		Auth:     auth.NewAuthService(cfg, repos.Auth, repos.User),
		User:     user.NewUserService(cfg, repos.User),
		Post:     post.NewPostService(cfg, repos.Post, repos.Comment, repos.Activity),
		Comment:  comment.NewCommentService(cfg, repos.Comment, repos.Post),
		Activity: activity.NewActivityService(cfg, repos.Activity),
		Upload:   upload.NewUploadService(cfg, repos.Upload),
	}
}
