package service

import (
	"context"
	"log"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/activity"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/auth"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/comment"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/follow"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/notification"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/post"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/saves"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/upload"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/user"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/repository"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/storage"
)

type Services struct {
	Auth         auth.AuthService
	User         user.UserService
	Post         post.PostService
	Comment      comment.CommentService
	Activity     activity.ActivityService
	Upload       upload.UploadService
	Saves        saves.SavesService
	Follow       follow.FollowService
	Notification notification.NotificationService
}

func InitServices(repos *repository.Repositories, cfg *configs.Config) *Services {
	notifSvc := notification.NewNotificationService(cfg, repos.Notification)

	var r2Storage storage.R2Storage
	r2Cfg := cfg.GetR2Config()
	if r2Cfg.AccountID != "" && r2Cfg.AccessKeyID != "" && r2Cfg.SecretAccessKey != "" && r2Cfg.BucketName != "" {
		client, err := storage.NewR2Storage(context.Background(), r2Cfg)
		if err != nil {
			log.Printf("[InitServices] warning: failed to initialize Cloudflare R2 storage: %v", err)
		} else {
			r2Storage = client
			log.Printf("[InitServices] Cloudflare R2 initialized successfully for bucket: %s", r2Cfg.BucketName)
		}
	} else {
		log.Println("[InitServices] Cloudflare R2 credentials not configured, falling back to local storage")
	}

	return &Services{
		Auth:         auth.NewAuthService(cfg, repos.Auth, repos.User),
		User:         user.NewUserService(cfg, repos.User),
		Post:         post.NewPostService(cfg, repos.Post, repos.Comment, repos.Activity, repos.Upload, repos.Follow),
		Comment:      comment.NewCommentService(cfg, repos.Comment, repos.Post, notifSvc),
		Activity:     activity.NewActivityService(cfg, repos.Activity, notifSvc),
		Upload:       upload.NewUploadService(cfg, repos.Upload, r2Storage),
		Saves:        saves.NewSavesService(cfg, repos.Saves),
		Follow:       follow.NewFollowService(cfg, repos.Follow, notifSvc),
		Notification: notifSvc,
	}
}

