import apiClient from './client'
import type { ChatMessage, ChatResponse } from '@/types/chat.types'

const base = (slug: string, projectId: string) =>
  `/api/v1/workspaces/${slug}/projects/${projectId}`

export const chatApi = {
  send(slug: string, projectId: string, message: string) {
    return apiClient.post<ChatResponse>(`${base(slug, projectId)}/chat`, { message })
  },

  history(slug: string, projectId: string, limit = 50) {
    return apiClient.get<{ results: ChatMessage[] }>(`${base(slug, projectId)}/chat?limit=${limit}`)
  },

  clear(slug: string, projectId: string) {
    return apiClient.delete(`${base(slug, projectId)}/chat`)
  },
}
