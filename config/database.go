package config

import (
	"log"

	"github.com/ArielioBayu/go-simple-blog-v2/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB initializes database connection
func ConnectDB(config *Config) error {
	dsn := config.GetDSN()

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	log.Println("Database connected successfully")

	// Auto migrate models
	if err := DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Category{},
		&models.Tag{},
		&models.Post{},
		&models.Comment{},
	); err != nil {
		return err
	}

	log.Println("Database migration completed")

	// Seed default roles if not exists
	if err := seedDefaultRoles(); err != nil {
		return err
	}

	return nil
}

// seedDefaultRoles creates default roles
func seedDefaultRoles() error {
	roles := []models.Role{
		{Name: "admin", Description: "Administrator with full access"},
		{Name: "author", Description: "Can create and manage own posts"},
		{Name: "reader", Description: "Can read and comment on posts"},
	}

	for _, role := range roles {
		var existingRole models.Role
		if err := DB.Where("name = ?", role.Name).First(&existingRole).Error; err == gorm.ErrRecordNotFound {
			if err := DB.Create(&role).Error; err != nil {
				return err
			}
			log.Printf("Created default role: %s\n", role.Name)
		}
	}

	return nil
}
