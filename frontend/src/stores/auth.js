import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI } from '../api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(null)

  const isLoggedIn = computed(() => !!token.value)

  async function login(credentials) {
    const res = await authAPI.login(credentials)
    token.value = res.token
    localStorage.setItem('token', res.token)
    user.value = res.user
    return res
  }

  async function register(data) {
    const res = await authAPI.register(data)
    token.value = res.token
    localStorage.setItem('token', res.token)
    user.value = res.user
    return res
  }

  async function fetchUser() {
    if (!token.value) return null
    try {
      user.value = await authAPI.me()
      return user.value
    } catch {
      logout()
      return null
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  return { token, user, isLoggedIn, login, register, fetchUser, logout }
})