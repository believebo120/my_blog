// store/user.js - 用户模块
import { service } from './index'
import type,{ AxiosResponse } from 'axios'
import router from "@/router";
export default {
    namespaced: true,
    state: {
        user: {
            isLoggedIn: false,
            userInfo: null
        },
        loading: false,
        error: null
    },
    getters: {
        isLoggedIn: state => state.user.isLoggedIn,
        userInfo: state => state.user.userInfo,
        loading: state => state.loading,
        error: state => state.error
    },
    mutations: {
        SET_LOGIN_STATE(state, isLoggedIn) {
            state.user.isLoggedIn = isLoggedIn
        },
        SET_USER_INFO(state, userInfo) {
            state.user.userInfo = userInfo
        },
        SET_LOADING(state, loading) {
            state.loading = loading
        },
        SET_ERROR(state, error) {
            state.error = error
        }
    },
    actions: {
        async login({ commit }, credentials) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                const response = await service.post('/login', credentials)
                localStorage.setItem('token', response.data.token)
                commit('SET_LOGIN_STATE', true)
                commit('SET_USER_INFO', response.data.user)
                await router.push('/home')
                return response
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '登录失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        },

        async register({ commit }, formData) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                const response = await service.post('/register', formData)
                localStorage.setItem('token', response.data.token)
                commit('SET_LOGIN_STATE', true)
                commit('SET_USER_INFO', response.data.user)
                return response
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '注册失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        },

        async logout({ commit }) {
            try {
                localStorage.removeItem('token')
                commit('SET_LOGIN_STATE', false)
                commit('SET_USER_INFO', null)
                return true
            } catch (error) {
                commit('SET_ERROR', '登出失败')
                throw error
            }
        },

        async initAuth({ commit }) {
            const token = localStorage.getItem('token')
            if (!token) return false

            try {
                service.defaults.headers.common['Authorization'] = `Bearer ${token}`
                const response = await service.get('/users/me')
                commit('SET_LOGIN_STATE', true)
                commit('SET_USER_INFO', response)
                return true
            } catch (error) {
                localStorage.removeItem('token')
                commit('SET_LOGIN_STATE', false)
                commit('SET_USER_INFO', null)
                return false
            }
        },

        async updateUser({ commit }, userData) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                const response = await service.put('/users/me', userData)
                commit('SET_USER_INFO', response)
                return response
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '更新失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        }
    }
}