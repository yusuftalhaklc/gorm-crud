package database

import (
	"gorm-crud/config"
	"gorm-crud/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := config.Config("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	err = DB.AutoMigrate(&models.User{}, &models.Post{})
	if err != nil {
		log.Fatal(err)
	}
}
