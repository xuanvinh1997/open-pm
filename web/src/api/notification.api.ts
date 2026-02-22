import apiClient from './client'
import type { NotificationListResponse, UnreadCountResponse } from '@/types/notification.types'

const base = (slug: string) => `/api/v1/workspaces/${slug}/notifications`

export const notificationApi = {
  list(slug: string, page = 1, perPage = 20) {
    return apiClient.get<NotificationListResponse>(
      `${base(slug)}?page=${page}&per_page=${perPage}`,
    )
  },

  unreadCount(slug: string) {
    return apiClient.get<UnreadCountResponse>(`${base(slug)}/unread-count`)
  },

  markRead(slug: string, notificationId: string) {
    return apiClient.post(`${base(slug)}/${notificationId}/read`)
  },

  markAllRead(slug: string) {
    return apiClient.post(`${base(slug)}/read-all`)
  },
}
