// store/blogs.js - 博客文章模块
import { service } from './index'
import type,{ AxiosResponse } from 'axios'
export default {
    namespaced: true,
    state: {
        blogs: [],
        currentBlog: null,
        loading: false,
        error: null
    },
    getters: {
        blogs: state => state.blogs,
        currentBlog: state => state.currentBlog,
        loading: state => state.loading,
        error: state => state.error
    },
    mutations: {
        SET_BLOGS(state, blogs) {
            state.blogs = blogs
        },
        SET_CURRENT_BLOG(state, blog) {
            state.currentBlog = blog
        },
        ADD_BLOG(state, blog) {
            state.blogs.unshift(blog)
        },
        UPDATE_BLOG(state, updatedBlog) {
            const index = state.blogs.findIndex(blog => blog.id === updatedBlog.id)
            if (index !== -1) {
                state.blogs.splice(index, 1, updatedBlog)
            }
        },
        DELETE_BLOG(state, blogId) {
            state.blogs = state.blogs.filter(blog => blog.id !== blogId)
        },
        SET_LOADING(state, loading) {
            state.loading = loading
        },
        SET_ERROR(state, error) {
            state.error = error
        }
    },
    actions: {
        async fetchBlogs({ commit }) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                const response = await service.get('/articles')
                commit('SET_BLOGS', response)
                return response
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '获取文章失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        },

        async fetchBlog({ commit }, blogId) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                const response = await service.get(`/articles/${blogId}`)
                commit('SET_CURRENT_BLOG', response)
                return response
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '获取文章失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        },

        async createBlog({ commit }, blogData) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                const response = await service.post('/articles', blogData)
                commit('ADD_BLOG', response)
                return response
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '创建文章失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        },

        async updateBlog({ commit }, { blogId, blogData }) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                const response = await service.put(`/articles/${blogId}`, blogData)
                commit('UPDATE_BLOG', response)
                return response
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '更新文章失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        },

        async deleteBlog({ commit }, blogId) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                await service.delete(`/articles/${blogId}`)
                commit('DELETE_BLOG', blogId)
                return true
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '删除文章失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        }
    }
}