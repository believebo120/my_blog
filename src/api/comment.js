import axios from 'axios';

// 评论列表
export const getCommentsByArticleId = (articleId) => axios.get(`/articles/${articleId}/comments`);

// 评论操作
export const createComment = (articleId, data) => axios.post(`/articles/${articleId}/comments`, data);
export const updateComment = (commentId, data) => axios.put(`/comments/${commentId}`, data);
export const deleteComment = (commentId) => axios.delete(`/comments/${commentId}`);