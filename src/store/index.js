// src/store/index.js
import { createStore } from 'vuex'

export default createStore({
    state: {
        isAuthenticated: JSON.parse(localStorage.getItem('isAuthenticated')) || false
    },
    mutations: {
        setAuthentication(state, status) {
            state.isAuthenticated = status
            localStorage.setItem('isAuthenticated', JSON.stringify(status))
        }
    },
    actions: {
        login({ commit }) {
            commit('setAuthentication', true)
        },
        logout({ commit }) {
            commit('setAuthentication', false)
        }
    }
})
