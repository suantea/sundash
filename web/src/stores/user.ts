import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '../api'
import type { User } from '../types'
import router from '../router'

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(localStorage.getItem('sundash-token'))
  const user = ref<User | null>(null)

  const isLoggedIn = computed(() => !!token.value)

  async function login(username: string, password: string) {
    const res = await api.post('auth/login', { username, password })
    token.value = res.data.token
    user.value = res.data.user
    localStorage.setItem('sundash-token', res.data.token)
    api.defaults.headers.common['Authorization'] = `Bearer ${res.data.token}`
    router.push('/')
  }

  async function register(username: string, password: string, displayName?: string) {
    const res = await api.post('auth/register', { username, password, display_name: displayName })
    // If pending approval, no token is returned
    if (res.data.status === 'pending') {
      return res.data
    }
    token.value = res.data.token
    user.value = res.data.user
    localStorage.setItem('sundash-token', res.data.token)
    api.defaults.headers.common['Authorization'] = `Bearer ${res.data.token}`
    router.push('/')
    return res.data
  }

  async function fetchProfile() {
    if (!token.value) return
    api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
    try {
      const res = await api.get('profile')
      user.value = res.data
    } catch {
      logout()
    }
  }

  // Direct state fill used by the home page bootstrap request.
  function setUser(u: User | null) {
    user.value = u
  }

  // Fill token + user (used by the first-time setup flow after POST /auth/setup).
  function setToken(t: string, u: User) {
    token.value = t
    user.value = u
    localStorage.setItem('sundash-token', t)
    api.defaults.headers.common['Authorization'] = `Bearer ${t}`
    router.push('/')
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('sundash-token')
    delete api.defaults.headers.common['Authorization']
    router.push('/login')
  }

  // Initialize token header
  if (token.value) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  return { token, user, isLoggedIn, login, register, fetchProfile, setUser, setToken, logout }
})
