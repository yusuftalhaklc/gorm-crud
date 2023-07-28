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
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get DB object: " + err.Error())
	}
	sqlDB.SetMaxIdleConns(10)  // Set maximum number of idle connections in the pool
	sqlDB.SetMaxOpenConns(100) // Set maximum number of open connections in the pool

	DB = db
	err = DB.AutoMigrate(&models.User{}, &models.Post{}, &models.PostLike{}, &models.Comment{})
	if err != nil {
		log.Fatal(err)
	}

}
