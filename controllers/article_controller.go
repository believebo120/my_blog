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

type ArticleController struct {
	articleService *services.ArticleService
}

func NewArticleController() *ArticleController {
	return &ArticleController{
		articleService: &services.ArticleService{},
	}
}

// GetArticles 获取所有文章
func (c *ArticleController) GetArticles(w http.ResponseWriter, r *http.Request) {
	// 检查服务是否初始化
	if c.articleService == nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "服务未初始化")
		return
	}

	articles, err := c.articleService.GetAllArticles()
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "获取文章列表失败")
		return
	}

	utils.SendResponse(w, http.StatusOK, "成功", articles)
}

// GetArticle 获取单个文章
func (c *ArticleController) GetArticle(w http.ResponseWriter, r *http.Request) {
	if c.articleService == nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "服务未初始化")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的文章ID")
		return
	}

	article, err := c.articleService.GetArticleByID(id)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "获取文章失败")
		return
	}

	if article == nil {
		utils.SendErrorResponse(w, http.StatusNotFound, "文章不存在")
		return
	}

	utils.SendResponse(w, http.StatusOK, "成功", article)
}

// CreateArticle 创建文章
func (c *ArticleController) CreateArticle(w http.ResponseWriter, r *http.Request) {
	if c.articleService == nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "服务未初始化")
		return
	}

	var req struct {
		Title        string  `json:"title"`
		Content      string  `json:"content"`
		Author       string  `json:"author"`
		ImagePath    *string `json:"image_path,omitempty"`
		CategoryName string  `json:"category_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	if req.Title == "" || req.Content == "" || req.Author == "" || req.CategoryName == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "标题、内容、作者和分类名不能为空")
		return
	}

	article := &models.Article{
		Title:     req.Title,
		Content:   req.Content,
		Author:    req.Author,
		ImagePath: req.ImagePath,
	}

	id, err := c.articleService.CreateArticle(article, req.CategoryName)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "创建文章失败: "+err.Error())
		return
	}

	utils.SendResponse(w, http.StatusCreated, "文章创建成功", map[string]interface{}{"id": id})
}

// UpdateArticle 更新文章
func (c *ArticleController) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	if c.articleService == nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "服务未初始化")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的文章ID")
		return
	}

	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	if err := c.articleService.UpdateArticle(id, &article); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "更新文章失败: "+err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, "文章更新成功", nil)
}

// DeleteArticle 删除文章
func (c *ArticleController) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	if c.articleService == nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "服务未初始化")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的文章ID")
		return
	}

	if err := c.articleService.DeleteArticle(id); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "删除文章失败: "+err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, "文章删除成功", nil)
}

// GetArticlesByCategory 获取分类下的所有文章
func (c *ArticleController) GetArticlesByCategory(w http.ResponseWriter, r *http.Request) {
	if c.articleService == nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "服务未初始化")
		return
	}

	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的分类ID")
		return
	}

	articles, err := c.articleService.GetArticlesByCategory(categoryID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "获取文章列表失败: "+err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, "成功", articles)
}
