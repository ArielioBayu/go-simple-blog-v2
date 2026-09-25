package handlers

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/activity"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/auth"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/comment"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/post"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/saves"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/upload"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/user"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/service"
)

type Handlers struct {
	Auth     *auth.AuthHandler
	User     *user.UserHandler
	Post     *post.PostHandler
	Comment  *comment.CommentHandler
	Activity *activity.ActivityHandler
	Upload   *upload.UploadHandler
	Saves    *saves.SavesHandler
}

func InitHandlers(services *service.Services, cfg *configs.Config) *Handlers {
	return &Handlers{
		Auth:     auth.NewAuthHandler(services.Auth, cfg),
		User:     user.NewUserHandler(services.User, cfg),
		Post:     post.NewPostHandler(services.Post, cfg),
		Comment:  comment.NewCommentHandler(services.Comment, cfg),
		Activity: activity.NewActivityHandler(services.Activity, cfg),
		Upload:   upload.NewUploadHandler(services.Upload, cfg),
		Saves:    saves.NewSavesHandler(services.Saves, cfg),
	}
}
