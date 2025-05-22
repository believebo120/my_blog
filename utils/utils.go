package utils

import (
	"encoding/json"
	"my_blog/models"
	"net/http"
)

// SendResponse 发送统一格式的JSON响应
func SendResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.APIResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

// SendErrorResponse 发送错误响应
func SendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	SendResponse(w, statusCode, message, nil)
}

// Constants 定义常量
const (
	RoleAdmin = 1
	RoleUser  = 2
	RoleGuest = 3
)
