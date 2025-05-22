// store/comments.js - 评论模块
import { service } from './index'
import type,{ AxiosResponse } from 'axios'
export default {
    namespaced: true,
    state: {
        comments: [],
        loading: false,
        error: null
    },
    getters: {
        comments: state => state.comments,
        loading: state => state.loading,
        error: state => state.error
    },
    mutations: {
        SET_COMMENTS(state, comments) {
            state.comments = comments
        },
        ADD_COMMENT(state, comment) {
            state.comments.unshift(comment)
        },
        UPDATE_COMMENT(state, updatedComment) {
            const index = state.comments.findIndex(comment => comment.id === updatedComment.id)
            if (index !== -1) {
                state.comments.splice(index, 1, updatedComment)
            }
        },
        DELETE_COMMENT(state, commentId) {
            state.comments = state.comments.filter(comment => comment.id !== commentId)
        },
        SET_LOADING(state, loading) {
            state.loading = loading
        },
        SET_ERROR(state, error) {
            state.error = error
        }
    },
    actions: {
        async fetchCommentsByArticle({ commit }, articleId) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                const response = await service.get(`/articles/${articleId}/comments`)
                commit('SET_COMMENTS', response)
                return response
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '获取评论失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        },

        async createComment({ commit }, { articleId, commentData }) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                const response = await service.post(`/articles/${articleId}/comments`, commentData)
                commit('ADD_COMMENT', response)
                return response
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '创建评论失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        },

        async updateComment({ commit }, { commentId, commentData }) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                const response = await service.put(`/comments/${commentId}`, commentData)
                commit('UPDATE_COMMENT', response)
                return response
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '更新评论失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        },

        async deleteComment({ commit }, commentId) {
            commit('SET_LOADING', true)
            commit('SET_ERROR', null)

            try {
                await service.delete(`/comments/${commentId}`)
                commit('DELETE_COMMENT', commentId)
                return true
            } catch (error) {
                commit('SET_ERROR', error.response?.data?.message || '删除评论失败')
                throw error
            } finally {
                commit('SET_LOADING', false)
            }
        }
    }
}