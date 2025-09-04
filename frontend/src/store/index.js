// src/store/index.js

import { createStore } from 'vuex';
import api from '../api/axiosInstance'; // 导入 Axios 实例

export default createStore({
    state: {
        isAuthenticated: false,
        user: null
    },
    mutations: {
        setAuthentication(state, status) {
            state.isAuthenticated = status;
        },
        setUser(state, user) {
            state.user = user;
        }
    },
    actions: {
        async login({ commit }, { username, password }) {
            try {
                const response = await api.post('/auth/login', { username, password });
                if (response.data.token) {
                    // 存储 JWT 和用户信息
                    localStorage.setItem('token', response.data.token);
                    localStorage.setItem('username', response.data.username);
                    localStorage.setItem('role', response.data.role);
                    
                    commit('setAuthentication', true);
                    commit('setUser', {
                        username: response.data.username,
                        role: response.data.role
                    });
                    return true;
                }
            } catch (error) {
                console.error('Login failed:', error);
                return false;
            }
        },
        async logout({ commit }) {
            try {
                // 清除所有认证信息
                localStorage.removeItem('token');
                localStorage.removeItem('username');
                localStorage.removeItem('role');
                commit('setAuthentication', false);
                commit('setUser', null);
            } catch (error) {
                console.error('Logout failed:', error);
            }
        },
        async checkAuth({ commit }) {
            try {
                const response = await api.get('/auth/check');
                if (response.data.authenticated) {
                    commit('setAuthentication', true);
                    commit('setUser', {
                        username: response.data.username,
                        role: response.data.role
                    });
                } else {
                    commit('setAuthentication', false);
                    commit('setUser', null);
                }
            } catch (error) {
                commit('setAuthentication', false);
                commit('setUser', null);
                // 清除过期token
                localStorage.removeItem('token');
                localStorage.removeItem('username');
                localStorage.removeItem('role');
            }
        },
        
        // 从本地存储恢复用户状态
        restoreAuth({ commit }) {
            const token = localStorage.getItem('token');
            const username = localStorage.getItem('username');
            const role = localStorage.getItem('role');
            
            if (token && username && role) {
                commit('setAuthentication', true);
                commit('setUser', { username, role });
            }
        }
    }
});
