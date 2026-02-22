<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useNotificationStore } from '@/stores/notification.store'
import { formatRelativeDate } from '@/utils/helpers'
import PDropdown from '@/components/ui/PDropdown.vue'
import PAvatar from '@/components/ui/PAvatar.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import { Bell, CheckCheck } from 'lucide-vue-next'
import type { Notification } from '@/types/notification.types'

const route = useRoute()
const router = useRouter()
const notificationStore = useNotificationStore()
const dropdownRef = ref<InstanceType<typeof PDropdown>>()

const slug = computed(() => route.params.workspaceSlug as string)

onMounted(() => {
  if (slug.value) {
    notificationStore.startPolling(slug.value)
  }
})

onUnmounted(() => {
  notificationStore.stopPolling()
})

async function handleOpen() {
  await notificationStore.fetchNotifications(slug.value)
}

async function handleClick(notification: Notification) {
  if (!notification.read_at) {
    await notificationStore.markRead(slug.value, notification.id)
  }
  dropdownRef.value?.close()
  if (
    (notification.entity_type === 'issue' || notification.entity_type === 'issue_comment') &&
    notification.project_id &&
    notification.entity_id
  ) {
    router.push(
      `/${slug.value}/projects/${notification.project_id}/issues/${notification.entity_id}`,
    )
  }
}

function handleMarkAllRead() {
  notificationStore.markAllRead(slug.value)
}

function senderName(n: Notification): string {
  if (n.sender_display_name) return n.sender_display_name
  if (n.sender_first_name) return `${n.sender_first_name} ${n.sender_last_name || ''}`.trim()
  return 'Someone'
}
</script>

<template>
  <PDropdown ref="dropdownRef" align="right" width="22rem">
    <template #trigger>
      <button
        class="relative rounded p-1 text-custom-text-300 hover:bg-custom-background-90 hover:text-custom-text-200 transition-colors"
        title="Notifications"
        @click="handleOpen"
      >
        <Bell class="h-4 w-4" />
        <span
          v-if="notificationStore.unreadCount > 0"
          class="absolute -right-0.5 -top-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-red-500 text-[9px] font-bold text-white"
        >
          {{ notificationStore.unreadCount > 9 ? '9+' : notificationStore.unreadCount }}
        </span>
      </button>
    </template>
    <template #default="{ close }">
      <div class="flex items-center justify-between border-b border-custom-border-200 px-3 py-2">
        <span class="text-sm font-semibold text-custom-text-100">Notifications</span>
        <button
          v-if="notificationStore.unreadCount > 0"
          class="flex items-center gap-1 rounded px-1.5 py-0.5 text-xs text-custom-text-300 hover:bg-custom-background-80 hover:text-custom-text-200 transition-colors"
          @click="handleMarkAllRead"
        >
          <CheckCheck class="h-3 w-3" />
          Mark all read
        </button>
      </div>

      <div class="max-h-80 overflow-y-auto">
        <div v-if="notificationStore.loading" class="flex items-center justify-center py-8">
          <PSpinner />
        </div>
        <div
          v-else-if="notificationStore.notifications.length === 0"
          class="py-8 text-center text-sm text-custom-text-300"
        >
          No notifications yet
        </div>
        <button
          v-else
          v-for="n in notificationStore.notifications"
          :key="n.id"
          class="flex w-full items-start gap-2.5 px-3 py-2.5 text-left hover:bg-custom-background-80 transition-colors"
          :class="!n.read_at ? 'bg-custom-background-90/50' : ''"
          @click="handleClick(n)"
        >
          <PAvatar
            :name="senderName(n)"
            :src="n.sender_avatar_url"
            size="sm"
            class="mt-0.5 flex-shrink-0"
          />
          <div class="min-w-0 flex-1">
            <div class="flex items-start gap-1.5">
              <p
                class="flex-1 text-xs leading-relaxed"
                :class="!n.read_at ? 'text-custom-text-100 font-medium' : 'text-custom-text-200'"
              >
                {{ n.title }}
              </p>
              <span
                v-if="!n.read_at"
                class="mt-1.5 h-1.5 w-1.5 flex-shrink-0 rounded-full bg-brand-500"
              />
            </div>
            <p v-if="n.message" class="mt-0.5 truncate text-2xs text-custom-text-300">
              {{ n.message }}
            </p>
            <p class="mt-0.5 text-2xs text-custom-text-300">
              {{ formatRelativeDate(n.created_at) }}
            </p>
          </div>
        </button>
      </div>
    </template>
  </PDropdown>
</template>
