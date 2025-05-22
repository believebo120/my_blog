import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import HomeLoggedIn from '../views/HomeLoggedIn.vue'
import store from '../store'
import { ElMessage } from 'element-plus'
import { refreshToken } from '@/api/auth' // 假设存在刷新token的接口

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView
  },
  {
    path: '/home',
    name: 'HomeLoggedIn',
    component: HomeLoggedIn,
    meta: { requiresAuth: true }
  },
  {
    path: '/blog',
    name: 'blog',
    component: () => import('../views/BlogList.vue')
  },
  {
    path: '/blog/:id',
    name: 'blogDetail',
    component: () => import('../views/BlogDetail.vue')
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('../views/Register.vue')
  },
  {
    path: '/personal',
    name: 'personal',
    component: () => import('../views/Personal.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/add-post',
    name: 'addPost',
    component: () => import('../views/AddPost.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/edit-post/:id',
    name: 'editPost',
    component: () => import('../views/EditPost.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

// 全局前置守卫
router.beforeEach(async (to, from, next) => {
  // 检查路由是否需要身份验证
  if (to.meta.requiresAuth) {
    // 获取token和登录状态
    const token = localStorage.getItem('token')
    const isLoggedIn = store.getters['user/isLoggedIn']

    // 如果没有token，直接跳转登录页
    if (!token) {
      return next({ path: '/login', query: { redirect: to.fullPath } })
    }

    // 如果有token但状态显示未登录，尝试刷新token
    if (token && !isLoggedIn) {
      try {
        // 调用刷新token的接口
        const response = await refreshToken()
        // 更新store中的用户信息
        store.commit('user/SET_USER', response.data.user)
        store.commit('user/SET_TOKEN', response.data.token)
        next()
      } catch (error) {
        // 刷新失败，清除token并重定向到登录页
        localStorage.removeItem('token')
        ElMessage.error('登录状态已过期，请重新登录')
        next({ path: '/login', query: { redirect: to.fullPath } })
      }
      return
    }

    // 已登录，允许访问
    next()
  } else {
    // 不需要身份验证，直接通过
    next()
  }
})

// 添加axios响应拦截器处理401错误
import axios from 'axios'
axios.interceptors.response.use(
    response => response,
    error => {
      if (error.response && error.response.status === 401) {
        // 清除token和用户信息
        store.commit('user/CLEAR_USER')
        localStorage.removeItem('token')

        // 显示错误提示
        ElMessage.error('登录状态已过期，请重新登录')

        // 获取当前路由
        const currentRoute = router.currentRoute.value

        // 如果不是在登录页，才跳转到登录页
        if (currentRoute.name !== 'login') {
          router.push({
            path: '/login',
            query: { redirect: currentRoute.fullPath }
          })
        }
      }
      return Promise.reject(error)
    }
)

export default router