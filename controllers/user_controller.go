package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"my_blog/middleware"
	"my_blog/models"
	"my_blog/services"
	"my_blog/utils"
	"net/http"
	"strconv"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: &services.UserService{},
	}
}

// Register 用户注册（支持文件上传）
func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的表单数据")
		return
	}

	var user models.User
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.Email = r.FormValue("email")

	if user.Username == "" || user.Password == "" || user.Email == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "用户名、密码和邮箱不能为空")
		return
	}

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "头像为必填项")
		return
	}
	defer file.Close()

	if err := c.userService.Register(&user, fileHeader); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "注册失败")
		return
	}

	utils.SendResponse(w, http.StatusCreated, "注册成功", nil)
}

// Login 用户登录
func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	user, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "生成token失败")
		return
	}

	utils.SendResponse(w, http.StatusOK, "登录成功", map[string]string{"token": token})
}

// GetCurrentUser 获取当前用户信息
func (c *UserController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "未登录")
		return
	}

	user, err := c.userService.GetUserByID(userID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "获取用户信息失败")
		return
	}

	if user == nil {
		utils.SendErrorResponse(w, http.StatusNotFound, "用户不存在")
		return
	}

	utils.SendResponse(w, http.StatusOK, "成功", user)
}

// UpdateUser 更新用户信息
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的用户ID")
		return
	}

	// 解析表单数据（支持文件上传）
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		// 如果解析失败，尝试作为JSON解析（不包含文件的情况）
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
			return
		}

		if err := c.userService.UpdateUser(id, &user); err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "更新用户信息失败: "+err.Error())
			return
		}

		utils.SendResponse(w, http.StatusOK, "用户信息更新成功", nil)
		return
	}

	// 处理包含文件上传的情况
	var user models.User
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.Email = r.FormValue("email")
	user.Status = parseStatus(r.FormValue("status"))

	// 处理文件上传（如果有）
	file, fileHeader, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		// 生成唯一文件名（假设 services 包中有此函数，返回 string）
		uniqueName := services.GenerateUniqueFileName(fileHeader.Filename)
		user.ImageData = "/uploads/" + uniqueName
		// 保存文件，SaveAvatar 仅返回 error
		if err := c.userService.SaveAvatar(file, uniqueName); err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "保存头像失败")
			return
		}
	}

	if err := c.userService.UpdateUser(id, &user); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "更新用户信息失败: "+err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, "用户信息更新成功", nil)
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的用户ID")
		return
	}

	if err := c.userService.DeleteUser(id); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "删除用户失败")
		return
	}

	utils.SendResponse(w, http.StatusOK, "用户删除成功", nil)
}

// GetAllUsers 获取所有用户
func (c *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.userService.GetAllUsers()
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	utils.SendResponse(w, http.StatusOK, "成功", users)
}

// UpdateUserRole 更新用户角色
func (c *UserController) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req struct {
		RoleID int `json:"role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	if err := c.userService.UpdateUserRole(id, req.RoleID); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "更新用户角色失败")
		return
	}

	utils.SendResponse(w, http.StatusOK, "用户角色更新成功", nil)
}

// 辅助函数：解析状态值
func parseStatus(statusStr string) int {
	if statusStr == "" {
		return 0 // 默认值
	}
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		return 0 // 解析失败返回默认值
	}
	return status
}

// UpdateUserBackgroundImage 更新用户背景图
func (c *UserController) UpdateUserBackgroundImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的用户ID")
		return
	}

	// 解析表单数据
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的表单数据")
		return
	}

	file, fileHeader, err := r.FormFile("background_image")
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "背景图为必填项")
		return
	}
	defer file.Close()

	// 生成唯一文件名
	uniqueName := services.GenerateUniqueFileName(fileHeader.Filename)
	backgroundImagePath := "/uploads/backgrounds/" + uniqueName

	// 保存文件
	if err := c.userService.SaveBackgroundImage(file, uniqueName); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "保存背景图失败")
		return
	}

	// 更新用户背景图
	if err := c.userService.UpdateUserBackgroundImage(id, backgroundImagePath); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "更新用户背景图失败")
		return
	}

	utils.SendResponse(w, http.StatusOK, "用户背景图更新成功", nil)
}
