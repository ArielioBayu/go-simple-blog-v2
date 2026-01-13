package main

import (
	"go-simple-blog-v2/internal/configs"
	"log"

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

	r.Run(cfg.Service.Port)

}
