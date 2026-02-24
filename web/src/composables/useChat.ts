import { ref } from 'vue'
import { chatApi } from '@/api/chat.api'
import type { ChatMessage } from '@/types/chat.types'

const messages = ref<ChatMessage[]>([])
const loading = ref(false)
const isOpen = ref(false)

export function useChat() {
  function toggle() {
    isOpen.value = !isOpen.value
  }

  async function loadHistory(slug: string, projectId: string) {
    try {
      const { data } = await chatApi.history(slug, projectId)
      messages.value = data.results || []
    } catch {
      messages.value = []
    }
  }

  async function sendMessage(slug: string, projectId: string, content: string) {
    const userMsg: ChatMessage = {
      id: crypto.randomUUID(),
      project_id: projectId,
      workspace_id: '',
      user_id: '',
      role: 'user',
      content,
      created_at: new Date().toISOString(),
    }
    messages.value.push(userMsg)
    loading.value = true

    try {
      const { data } = await chatApi.send(slug, projectId, content)
      messages.value.push(data.message)
    } catch {
      messages.value.push({
        id: crypto.randomUUID(),
        project_id: projectId,
        workspace_id: '',
        user_id: '',
        role: 'assistant',
        content: 'Sorry, I was unable to process your request. Please try again.',
        created_at: new Date().toISOString(),
      })
    } finally {
      loading.value = false
    }
  }

  async function clearHistory(slug: string, projectId: string) {
    try {
      await chatApi.clear(slug, projectId)
      messages.value = []
    } catch {
      // ignore
    }
  }

  return {
    messages,
    loading,
    isOpen,
    toggle,
    loadHistory,
    sendMessage,
    clearHistory,
  }
}
