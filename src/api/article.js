// src/api/article.js
import axios from 'axios';

// 获取文章列表
export const getArticles = (params = {}) => {
    return axios.get('/articles', { params });
};

// 获取单篇文章
export const getArticleById = (id) => {
    return axios.get(`/articles/${id}`);
};

// 创建文章
export const createArticle = (data) => {
    return axios.post('/articles', data);
};

// 更新文章
export const updateArticle = (id, data) => {
    return axios.put(`/articles/${id}`, data);
};

// 删除文章
export const deleteArticle = (id) => {
    return axios.delete(`/articles/${id}`);
};

// 获取文章评论
export const getCommentsByArticleId = (id) => {
    return axios.get(`/articles/${id}/comments`);
};

// 创建评论
export const createComment = (articleId, data) => {
    return axios.post(`/articles/${articleId}/comments`, data);
};