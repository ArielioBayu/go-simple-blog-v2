package main

import (
	"go-simple-blog-v2/internal/configs"
	"go-simple-blog-v2/internal/handlers/memberships"
	"go-simple-blog-v2/pkg/internalsql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	membershipsRepo "go-simple-blog-v2/internal/repository/memberships"

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
	log.Println("config", cfg)

	db, err := internalsql.Connect(cfg.Database.DbSourceName)
	if err != nil {
		log.Fatal("Gagal Inisiasi Database", err)
	}

	_ = membershipsRepo.NewRepository(db)

	membershipsHandler := memberships.NewHandler(r)
	membershipsHandler.RegisterRoute()

	r.Run(cfg.Service.Port)

}
