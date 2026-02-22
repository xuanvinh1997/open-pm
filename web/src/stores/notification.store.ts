import { defineStore } from 'pinia'
import { ref } from 'vue'
import { notificationApi } from '@/api/notification.api'
import type { Notification } from '@/types/notification.types'

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([])
  const unreadCount = ref(0)
  const loading = ref(false)

  let pollInterval: ReturnType<typeof setInterval> | null = null

  async function fetchNotifications(slug: string) {
    loading.value = true
    try {
      const { data } = await notificationApi.list(slug)
      notifications.value = data.results
      unreadCount.value = data.unread_count
    } finally {
      loading.value = false
    }
  }

  async function fetchUnreadCount(slug: string) {
    try {
      const { data } = await notificationApi.unreadCount(slug)
      unreadCount.value = data.unread_count
    } catch {
      // Silently fail -- polling should not disrupt UX
    }
  }

  async function markRead(slug: string, notificationId: string) {
    await notificationApi.markRead(slug, notificationId)
    const n = notifications.value.find((n) => n.id === notificationId)
    if (n) n.read_at = new Date().toISOString()
    unreadCount.value = Math.max(0, unreadCount.value - 1)
  }

  async function markAllRead(slug: string) {
    await notificationApi.markAllRead(slug)
    notifications.value.forEach((n) => {
      if (!n.read_at) n.read_at = new Date().toISOString()
    })
    unreadCount.value = 0
  }

  function startPolling(slug: string, intervalMs = 30000) {
    stopPolling()
    fetchUnreadCount(slug)
    pollInterval = setInterval(() => fetchUnreadCount(slug), intervalMs)
  }

  function stopPolling() {
    if (pollInterval) {
      clearInterval(pollInterval)
      pollInterval = null
    }
  }

  return {
    notifications,
    unreadCount,
    loading,
    fetchNotifications,
    fetchUnreadCount,
    markRead,
    markAllRead,
    startPolling,
    stopPolling,
  }
})
