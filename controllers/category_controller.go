package controllers

import (
	"encoding/json"
	"my_blog/models"
	"my_blog/services"
	"my_blog/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController() *CategoryController {
	return &CategoryController{
		categoryService: &services.CategoryService{},
	}
}

// GetCategories 获取所有分类
func (c *CategoryController) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := c.categoryService.GetAllCategories()
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "获取分类列表失败")
		return
	}

	utils.SendResponse(w, http.StatusOK, "成功", categories)
}

// GetCategory 获取单个分类
func (c *CategoryController) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的分类ID")
		return
	}

	category, err := c.categoryService.GetCategoryByID(id)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "获取分类失败")
		return
	}

	if category == nil {
		utils.SendErrorResponse(w, http.StatusNotFound, "分类不存在")
		return
	}

	utils.SendResponse(w, http.StatusOK, "成功", category)
}

// CreateCategory 创建分类
func (c *CategoryController) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	if category.Name == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "分类名称不能为空")
		return
	}

	id, err := c.categoryService.CreateCategory(&category)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "创建分类失败")
		return
	}

	utils.SendResponse(w, http.StatusCreated, "分类创建成功", map[string]interface{}{"id": id})
}

// UpdateCategory 更新分类
func (c *CategoryController) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的分类ID")
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	if err := c.categoryService.UpdateCategory(id, &category); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "更新分类失败")
		return
	}

	utils.SendResponse(w, http.StatusOK, "分类更新成功", nil)
}

// DeleteCategory 删除分类
func (c *CategoryController) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的分类ID")
		return
	}

	if err := c.categoryService.DeleteCategory(id); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "删除分类失败")
		return
	}

	utils.SendResponse(w, http.StatusOK, "分类删除成功", nil)
}
