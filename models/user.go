package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName string `json:"full_name"`
	Gender   string `json:"gender"`
	Username string `json:"username"`
	Password string `json:"password"`
	Verified bool   `json:"verified"`
}
