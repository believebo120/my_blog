package models

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	ID        int       `json:"id"`
	ArticleID int       `json:"article_id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreateAt  time.Time `json:"create_at"`
}

// Article 文章模型
type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreateAt  time.Time `json:"create_at"`
	ImagePath *string   `json:"image_path,omitempty"`
	Category  Category  `json:"category"`
	Views     int       `json:"views"`
}

// UserArticleCount 用户文章统计
type UserArticleCount struct {
	Username     string `json:"username"`
	ArticleCount int    `json:"article_count"`
}

// Category 分类模型
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// User 用户模型
type User struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"` // 密码字段，不输出到JSON
	ImageData       string `json:"image_data"`
	RoleID          int    `json:"role_id"`
	Status          int    `json:"status"`
	BackgroundImage string `json:"background_image"`
}

// APIResponse API响应格式
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
