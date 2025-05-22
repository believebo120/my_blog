import axios from 'axios';

// 分类管理
export const getCategories = () => axios.get('/categories');
export const getCategoryById = (id) => axios.get(`/categories/${id}`);
export const createCategory = (data) => axios.post('/categories', data);
export const updateCategory = (id, data) => axios.put(`/categories/${id}`, data);
export const deleteCategory = (id) => axios.delete(`/categories/${id}`);

// 分类下的文章
export const getArticlesByCategory = (id) => axios.get(`/categories/${id}/articles`);