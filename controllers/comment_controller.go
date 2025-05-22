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

type CommentController struct {
	commentService *services.CommentService
}

func NewCommentController() *CommentController {
	return &CommentController{
		commentService: &services.CommentService{},
	}
}

// GetCommentsByArticle 获取文章的所有评论
func (c *CommentController) GetCommentsByArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的文章ID")
		return
	}

	comments, err := c.commentService.GetCommentsByArticle(articleID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "获取评论列表失败")
		return
	}

	utils.SendResponse(w, http.StatusOK, "成功", comments)
}

// CreateComment 创建评论
func (c *CommentController) CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	articleID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的文章ID")
		return
	}

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	comment.ArticleID = articleID

	if comment.Content == "" || comment.Author == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "评论内容和作者不能为空")
		return
	}

	id, err := c.commentService.CreateComment(&comment)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "创建评论失败")
		return
	}

	utils.SendResponse(w, http.StatusCreated, "评论创建成功", map[string]interface{}{"id": id})
}

// UpdateComment 更新评论
func (c *CommentController) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的评论ID")
		return
	}

	var req struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的请求数据")
		return
	}

	if req.Content == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "评论内容不能为空")
		return
	}

	if err := c.commentService.UpdateComment(id, req.Content); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "更新评论失败")
		return
	}

	utils.SendResponse(w, http.StatusOK, "评论更新成功", nil)
}

// DeleteComment 删除评论
func (c *CommentController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "无效的评论ID")
		return
	}

	if err := c.commentService.DeleteComment(id); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "删除评论失败")
		return
	}

	utils.SendResponse(w, http.StatusOK, "评论删除成功", nil)
}
