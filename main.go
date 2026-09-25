package main

import (
	"log"
	"os"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/handlers"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/repository"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/router"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/service"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/internalsql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
)

func main() {
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

	db, err := internalsql.Connect(cfg.Database.GetDSN())
	if err != nil {
		log.Fatal("Gagal Inisiasi Database: ", err)
	}
	defer db.Close()

	// Jalankan migration secara otomatis saat startup
	internalsql.RunMigration(db, "./scripts/migrations")

	// Buat direktori upload jika belum ada
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		log.Fatalf("Gagal membuat direktori uploads: %v", err)
	}

	// Inisialisasi Modular Application Containers (Registries)
	repos := repository.InitRepositories(db)
	services := service.InitServices(repos, cfg)
	handlersList := handlers.InitHandlers(services, cfg)

	// Setup Gin Router
	r := gin.Default()
	r.Use(middleware.CorsMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Static route untuk file upload
	r.Static("/uploads", "./uploads")

	// Setup Routes
	router.SetupRoutes(r, handlersList)

	log.Printf("Server running on port %s", cfg.Service.Port)
	if err := r.Run(cfg.Service.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
