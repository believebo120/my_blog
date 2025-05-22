package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB        *sql.DB
	JWTSecret = "your-secret-key-here"
)

// InitDB 初始化数据库连接
func InitDB() {
	var err error
	DB, err = sql.Open("mysql", "root:zrb123456@tcp(127.0.0.1:3306)/test?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	createTables()
}

// createTables 创建数据库表
func createTables() {
	// 创建分类表
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS categories (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL UNIQUE,
		description TEXT
	)`)
	if err != nil {
		log.Fatal("创建categories表失败:", err)
	}

	// 创建文章表
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS articles (
		id INT PRIMARY KEY AUTO_INCREMENT,
		author VARCHAR(200) NOT NULL,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		create_at DATETIME NOT NULL,
		image_path VARCHAR(255) DEFAULT NULL,
		views INT NOT NULL,
		category_id INT,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
	)`)
	if err != nil {
		log.Fatal("创建articles表失败:", err)
	}

	// 创建用户表
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT PRIMARY KEY AUTO_INCREMENT,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		image_data VARCHAR(255) DEFAULT '/uploads/default_avatar.jpg', -- 默认头像
		background_image VARCHAR(255) DEFAULT '/uploads/default_bg.jpg', 
		role_id INT DEFAULT 2 -- 默认普通用户
	)`)
	if err != nil {
		log.Fatal("创建users表失败:", err)
	}

	// 创建评论表
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS comments (
		id INT PRIMARY KEY AUTO_INCREMENT,
		article_id INT NOT NULL,
		content TEXT NOT NULL,
		author VARCHAR(200) NOT NULL,
		create_at DATETIME NOT NULL,
		FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatal("创建comments表失败:", err)
	}

	// 创建角色表
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS roles (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL UNIQUE,
		description TEXT
	)`)
	if err != nil {
		log.Fatal("创建roles表失败:", err)
	}

	// 插入默认角色
	_, err = DB.Exec(`INSERT INTO roles (name, description) 
		VALUES ('admin', '管理员'), ('user', '普通用户'), ('guest', '访客') 
		ON DUPLICATE KEY UPDATE name=name`)
	if err != nil {
		log.Fatal("插入默认角色失败:", err)
	}
}
