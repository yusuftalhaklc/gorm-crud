package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Content  string `json:"content"`
	Username string `json:"username"`
	Tags     string `json:"tags"`
}
