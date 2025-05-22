// store/index.js - 主入口
import { createStore } from 'vuex'
import axios from 'axios'
import router from '@/router'

// 创建独立的axios实例
export const service = axios.create({
  baseURL: 'http://localhost:8080',
  timeout: 5000
})

// 请求拦截器
service.interceptors.request.use(
    config => {
      const token = localStorage.getItem('token')
      if (token) {
        config.headers['Authorization'] = `Bearer ${token}`
      }
      return config
    },
    error => {
      console.error('请求错误:', error)
      return Promise.reject(error)
    }
)

// 响应拦截器
service.interceptors.response.use(
    response => {
        if(response.status >= 200&&response.status < 300){
            return response;
        }
      return response.data
    },
    error => {
      console.error('响应错误:', error)

      if (error.response && error.response.status === 401) {
        localStorage.removeItem('token')
        router.push('/login')
      }

      return Promise.reject(error.response?.data?.message || '服务器错误')
    }
)

// 导入所有模块
import user from './user'
import blogs from './blogs'
import comments from './comments'
import categories from './categories'

export default createStore({
  modules: {
    user,
    blogs,
    comments,
    categories
  },
  strict: process.env.NODE_ENV !== 'production'
})