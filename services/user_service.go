package services

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"my_blog/config"
	"my_blog/models"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

// UserService 用户服务
type UserService struct{}

// GenerateUniqueFileName 生成唯一文件名（时间戳+随机数）
func GenerateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().UnixNano()
	random := time.Now().Nanosecond() % 1000000
	return fmt.Sprintf("%d_%06d%s", timestamp, random, ext)
}

// SaveAvatar 保存文件到本地服务器
func (s *UserService) SaveAvatar(file io.Reader, filename string) error {
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0750); err != nil {
		return err
	}
	filePath := filepath.Join(uploadDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	return err
}

// Register 注册用户（存储路径）
func (s *UserService) Register(user *models.User, fileHeader *multipart.FileHeader) error {
	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码哈希失败")
	}
	user.Password = string(hashedPassword) // 更新为哈希后的密码
	if fileHeader != nil {
		// 验证文件类型和大小
		ext := filepath.Ext(fileHeader.Filename)
		if !allowedExtensions[strings.ToLower(ext)] || fileHeader.Size > 2*1024*1024 {
			return errors.New("不支持的文件类型或文件过大")
		}

		// 生成唯一文件名
		uniqueName := GenerateUniqueFileName(fileHeader.Filename)
		user.ImageData = "/uploads/" + uniqueName // 存储相对路径

		// 打开文件并保存
		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()
		if err := s.SaveAvatar(file, uniqueName); err != nil {
			return err
		}
	}

	// 插入数据库（仅存储路径）
	result, err := config.DB.Exec(`
		INSERT INTO users (username, password, email, image_data) 
		VALUES (?, ?, ?, ?)
	`, user.Username, user.Password, user.Email, user.ImageData)
	if err != nil {
		return err
	}
	_ = result // 若不需要 sql.Result，可忽略
	return nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*models.User, error) {
	var user models.User
	var hashedPassword string

	err := config.DB.QueryRow(`
		SELECT id, username, password, email
		FROM users WHERE username = ?
	`, username).Scan(
		&user.ID,
		&user.Username,
		&hashedPassword,
		&user.Email,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	// 智能验证：判断密码是否为哈希格式
	if strings.HasPrefix(hashedPassword, "$2a$") {
		// 是哈希密码，使用 bcrypt 验证
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			return nil, err // 密码不匹配
		}
		// 验证成功，检查是否需要升级（可选）
		if needUpgradeHash(hashedPassword) {
			go s.upgradePassword(user.ID, password) // 异步升级密码
		}
	} else {
		// 是明文密码，直接对比（仅针对老用户）
		if hashedPassword != password {
			return nil, errors.New("password mismatch")
		}
		// 明文验证成功，立即升级为哈希
		go s.upgradePassword(user.ID, password) // 异步升级密码
	}

	return &user, nil
}

// 检查哈希是否需要升级（例如，成本因子变化时）
func needUpgradeHash(hash string) bool {
	cost, err := bcrypt.Cost([]byte(hash))
	return err != nil || cost < bcrypt.DefaultCost
}
func (s *UserService) upgradePassword(userId int, password string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password for user %d: %v", userId, err)
		return
	}
	_, err = config.DB.Exec("UPDATE users SET password = ? WHERE id = ?", hashedPassword, userId)
	if err != nil {
		log.Printf("Failed to upgrade password for user %d: %v", userId, err)
	}
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := config.DB.QueryRow(`
		SELECT id, username, email, image_data
		FROM users WHERE id = ?
	`, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.ImageData,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id int, user *models.User) error {
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	_, err := config.DB.Exec(`
		UPDATE users
		SET username = ?, email = ?, password = ?, image_data = ?
		WHERE id = ?
	`,
		user.Username,
		user.Email,
		user.Password,
		user.ImageData,
		id,
	)
	return err
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int) error {
	result, err := config.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// UpdateUserRole 更新用户角色
func (s *UserService) UpdateUserRole(userID, roleID int) error {
	_, err := config.DB.Exec("UPDATE users SET role_id = ? WHERE id = ?", roleID, userID)
	return err
}

// GetAllUsers 获取所有用户
func (s *UserService) GetAllUsers() ([]models.User, error) {
	rows, err := config.DB.Query("SELECT id, username, email, image_data FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.ImageData,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// SaveBackgroundImage 保存背景图到本地服务器
func (s *UserService) SaveBackgroundImage(file io.Reader, filename string) error {
	uploadDir := "uploads/backgrounds"
	if err := os.MkdirAll(uploadDir, 0750); err != nil {
		return err
	}
	filePath := filepath.Join(uploadDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	return err
}

// UpdateUserBackgroundImage 更新用户背景图
func (s *UserService) UpdateUserBackgroundImage(id int, backgroundImage string) error {
	_, err := config.DB.Exec(`
        UPDATE users
        SET background_image = ?
        WHERE id = ?
    `, backgroundImage, id)
	return err
}
