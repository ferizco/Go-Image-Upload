// database/database.go

package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-image-upload/models"
)

func InitDatabase() *gorm.DB {
	// Initialize database connection
	dsn := "user=go_nico password=12345 dbname=goimageupload port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Image{})

	return db
}
