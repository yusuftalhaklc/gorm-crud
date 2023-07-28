package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID       int    `json:"user_id"`
	Content      string `json:"content"`
	CommentCount int    `json:"comment_count"`
	LikeCount    int    `json:"like_count"`
}
type PostLike struct {
	gorm.Model
	PostID   int    `json:"post_id"`
	Username string `json:"username"`
}
type Comment struct {
	gorm.Model
	PostId   int    `json:"post_id"`
	Username string `json:"username"`
	Comment  string `json:"Comment"`
}
