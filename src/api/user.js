import axios from 'axios';

// 用户认证
export const register = (data) => axios.post('/register', data);
export const login = (data) => axios.post('/login', data);

// 用户信息管理
export const getUserInfo = () => axios.get('/users/me');
export const updateUserInfo = (userId, data) => axios.put(`/users/${userId}`, data);

// 用户统计
export const getArticleCount = () => axios.get('/users/me/articles/count');

// 管理员功能
export const getAllUsers = () => axios.get('/users');
export const deleteUser = (userId) => axios.delete(`/users/${userId}`);
export const updateUserRole = (userId, role) => axios.put(`/users/${userId}/role`, { role });
export const changePassword = (userId, params) => {
    return axios.post(`/api/change-password/${userId}`, params);
};