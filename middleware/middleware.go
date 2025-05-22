package middleware

import (
	"context"
	"errors"
	"fmt"
	"my_blog/config"
	"net/http"
	"strings"
	"time" // 添加 time 包导入

	"github.com/dgrijalva/jwt-go"
)

const secretKey = "your-secret-key" // 生产环境应使用更安全的方式存储

// Claims JWT声明结构
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// CorsMiddleware 处理跨域请求
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin,Accept")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware 验证用户身份
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "未提供认证信息", http.StatusUnauthorized)
			return
		}

		userID, err := validateToken(token)
		if err != nil {
			http.Error(w, "无效的认证信息", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleMiddleware 验证用户角色权限
func RoleMiddleware(requiredRole int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value("userID").(int)
			if !ok {
				http.Error(w, "无法获取用户信息", http.StatusInternalServerError)
				return
			}

			var roleID int
			err := config.DB.QueryRow("SELECT role_id FROM users WHERE id = ?", userID).Scan(&roleID)
			if err != nil {
				http.Error(w, "无法获取用户角色", http.StatusInternalServerError)
				return
			}

			if roleID < requiredRole {
				http.Error(w, "权限不足", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GenerateToken 生成JWT Token
func GenerateToken(userID int, username string, expiresIn time.Duration) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresIn).Unix(),
			Issuer:    "your-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
/ VerifyToken 验证JWT Token
func VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
