package main

import (
	"log"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/handlers/memberships"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/handlers/posts"

	"github.com/ArielioBayu/go-simple-blog-v2/pkg/internalsql"

	_ "github.com/go-sql-driver/mysql"

	membershipsRepo "github.com/ArielioBayu/go-simple-blog-v2/internal/repository/memberships"
	postsRepo "github.com/ArielioBayu/go-simple-blog-v2/internal/repository/posts"
	membershipsSrv "github.com/ArielioBayu/go-simple-blog-v2/internal/service/memberships"
	postsSrv "github.com/ArielioBayu/go-simple-blog-v2/internal/service/posts"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	var cfg *configs.Config

	err := configs.Init(
		configs.WithConfigFolder(
			[]string{"./internal/configs/"},
		),
		configs.WithConfigFile("config"),
		configs.WithConfigType("yaml"),
	)
	if err != nil {
		log.Fatal("Gagal Inisiasi Config : ", err)
	}

	cfg = configs.Get()

	db, err := internalsql.Connect(cfg.Database.DbSourceName)
	if err != nil {
		log.Fatal("Gagal Inisiasi Database", err)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//	Repository
	membershipsRepo := membershipsRepo.NewMembershipsRepository(db)
	postsRepo := postsRepo.NewPostsRepository(db)

	//	Service
	membershipsService := membershipsSrv.NewMembershipsService(cfg, membershipsRepo)
	postsService := postsSrv.NewPostsService(cfg, postsRepo)

	//	Handler
	membershipsHandler := memberships.NewHandler(r, membershipsService)
	membershipsHandler.RegisterRoute()

	postsHandler := posts.NewHandler(r, postsService)
	postsHandler.RegisterRoute()

	r.Run(cfg.Service.Port)
}
