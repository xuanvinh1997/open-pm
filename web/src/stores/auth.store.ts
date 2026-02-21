import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth.api'
import type { User, LoginRequest, SignupRequest } from '@/types/auth.types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(localStorage.getItem('access_token'))
  const refreshToken = ref<string | null>(localStorage.getItem('refresh_token'))
  const loading = ref(false)

  const isAuthenticated = computed(() => !!accessToken.value)

  async function login(credentials: LoginRequest) {
    loading.value = true
    try {
      const { data } = await authApi.login(credentials)
      setTokens(data.access_token, data.refresh_token)
      user.value = data.user
    } finally {
      loading.value = false
    }
  }

  async function signup(data: SignupRequest) {
    loading.value = true
    try {
      const { data: resp } = await authApi.signup(data)
      setTokens(resp.access_token, resp.refresh_token)
      user.value = resp.user
    } finally {
      loading.value = false
    }
  }

  async function fetchUser() {
    if (!accessToken.value) return
    try {
      const { data } = await authApi.getUser()
      user.value = data
    } catch {
      logout()
    }
  }

  function setTokens(access: string, refresh: string) {
    accessToken.value = access
    refreshToken.value = refresh
    localStorage.setItem('access_token', access)
    localStorage.setItem('refresh_token', refresh)
  }

  function logout() {
    authApi.logout().catch(() => {})
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  return {
    user,
    accessToken,
    refreshToken,
    loading,
    isAuthenticated,
    login,
    signup,
    fetchUser,
    logout,
  }
})
