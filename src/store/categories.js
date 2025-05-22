// store/categories.js - 分类模块
import { service } from './index'
import type,{ AxiosResponse } from 'axios'
export default {
    namespaced: true,
    state: {
        categories: [],
        loading: false,
        error: null
    },
    getters: {
        categories: state => state.categories,
        loading: state => state.loading,
        error: state => state.error
    },
    mutations: {
        SET_CATEGORIES(state, categories) {
            state.categories = categories
        },
        ADD_CATEGORY(state, category) {
            state.categories.push(category)
        },
        UPDATE_CATEGORY(state, updatedCategory) {
            const index = state.categories.findIndex(cat => cat.id === updatedCategory.id)
            if (index !== -1) {
                state.categories.splice(index, 1, updatedCategory)
            }
        },
        DELETE_CATEGORY(state, categoryId) {
            state.categories = state.categories.filter(cat => cat.id !== categoryId)
        },
        SET_LOADING(state, loading) {
            state.loading = loading
        },
        SET_ERROR(state, error) {
            state.error = error
        }
    },
    actions: {
        async fetchCategories({ commit }) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                const response = await service.get('/categories')
                commit('SET_CATEGORIES', response)
                return response
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '获取分类失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        },

        async createCategory({ commit }, categoryData) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                const response = await service.post('/categories', categoryData)
                commit('ADD_CATEGORY', response)
                return response
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '创建分类失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        },

        async updateCategory({ commit }, { categoryId, categoryData }) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                const response = await service.put(`/categories/${categoryId}`, categoryData)
                commit('UPDATE_CATEGORY', response)
                return response
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '更新分类失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        },

        async deleteCategory({ commit }, categoryId) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                await service.delete(`/categories/${categoryId}`)
                commit('DELETE_CATEGORY', categoryId)
                return true
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '删除分类失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        }
    }
}