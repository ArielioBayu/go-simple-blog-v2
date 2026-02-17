package main

import (
	"log"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/handlers/memberships"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/internalsql"

	_ "github.com/go-sql-driver/mysql"

	membershipsRepo "github.com/ArielioBayu/go-simple-blog-v2/internal/repository/memberships"
	membershipsSrv "github.com/ArielioBayu/go-simple-blog-v2/internal/service/memberships"

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

	membershipsRepo := membershipsRepo.NewMembershipsRepository(db)
	membershipsService := membershipsSrv.NewMembershipsService(cfg, membershipsRepo)

	membershipsHandler := memberships.NewHandler(r, membershipsService)
	membershipsHandler.RegisterRoute()

	r.Run(cfg.Service.Port)

}
