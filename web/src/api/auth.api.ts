import apiClient from './client'
import type { LoginRequest, SignupRequest, TokenResponse, User } from '@/types/auth.types'

export const authApi = {
  signup(data: SignupRequest) {
    return apiClient.post<TokenResponse>('/auth/signup', data)
  },

  login(data: LoginRequest) {
    return apiClient.post<TokenResponse>('/auth/token?grant_type=password', data)
  },

  refresh(refreshToken: string) {
    return apiClient.post<TokenResponse>('/auth/token?grant_type=refresh_token', {
      refresh_token: refreshToken,
    })
  },

  getUser() {
    return apiClient.get<User>('/auth/user')
  },

  updateUser(data: Partial<User>) {
    return apiClient.put<User>('/auth/user', data)
  },

  logout() {
    return apiClient.post('/auth/logout')
  },

  recover(email: string) {
    return apiClient.post('/auth/recover', { email })
  },
}
