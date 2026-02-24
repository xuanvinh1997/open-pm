<script setup lang="ts">
import { computed, watch, nextTick, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useChat } from '@/composables/useChat'
import ChatMessageBubble from './ChatMessage.vue'
import ChatInput from './ChatInput.vue'
import { MessageCircle, X, Trash2 } from 'lucide-vue-next'

const route = useRoute()
const { messages, loading, isOpen, toggle, loadHistory, sendMessage, clearHistory } = useChat()

const slug = computed(() => route.params.workspaceSlug as string)
const projectId = computed(() => route.params.projectId as string)
const isInProject = computed(() => !!slug.value && !!projectId.value)

const messagesContainer = ref<HTMLElement | null>(null)

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

watch(isOpen, (open) => {
  if (open && isInProject.value) {
    loadHistory(slug.value, projectId.value)
  }
})

watch(
  () => messages.value.length,
  () => scrollToBottom(),
)

async function handleSend(content: string) {
  await sendMessage(slug.value, projectId.value, content)
}

function handleClear() {
  clearHistory(slug.value, projectId.value)
}
</script>

<template>
  <div v-if="isInProject" class="fixed bottom-5 right-5 z-50">
    <!-- Chat panel -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 translate-y-4 scale-95"
      enter-to-class="opacity-100 translate-y-0 scale-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0 scale-100"
      leave-to-class="opacity-0 translate-y-4 scale-95"
    >
      <div
        v-if="isOpen"
        class="absolute bottom-16 right-0 flex h-[500px] w-[380px] flex-col overflow-hidden rounded-2xl border border-custom-border-200 bg-custom-background-100 shadow-xl"
      >
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-custom-border-200 px-4 py-3">
          <div class="flex items-center gap-2">
            <MessageCircle class="h-4 w-4 text-brand-500" />
            <h3 class="text-sm font-semibold text-custom-text-100">Project Assistant</h3>
          </div>
          <div class="flex items-center gap-1">
            <button
              class="rounded-md p-1.5 text-custom-text-300 hover:bg-custom-background-80 hover:text-custom-text-200 transition-colors"
              title="Clear history"
              @click="handleClear"
            >
              <Trash2 class="h-3.5 w-3.5" />
            </button>
            <button
              class="rounded-md p-1.5 text-custom-text-300 hover:bg-custom-background-80 hover:text-custom-text-200 transition-colors"
              @click="toggle"
            >
              <X class="h-4 w-4" />
            </button>
          </div>
        </div>

        <!-- Messages -->
        <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-4">
          <div v-if="messages.length === 0" class="flex h-full items-center justify-center">
            <div class="text-center">
              <MessageCircle class="mx-auto h-10 w-10 text-custom-text-300 opacity-50" />
              <p class="mt-3 text-sm text-custom-text-300">
                Ask questions about your project data.
              </p>
              <p class="mt-1 text-xs text-custom-text-300 opacity-70">
                e.g. "What issues are overdue?" or "Summarize sprint progress"
              </p>
            </div>
          </div>
          <ChatMessageBubble v-for="msg in messages" :key="msg.id" :message="msg" />
          <div v-if="loading" class="flex items-center gap-2 text-custom-text-300">
            <div class="flex gap-1">
              <span class="h-2 w-2 animate-bounce rounded-full bg-custom-text-300" style="animation-delay: 0ms" />
              <span class="h-2 w-2 animate-bounce rounded-full bg-custom-text-300" style="animation-delay: 150ms" />
              <span class="h-2 w-2 animate-bounce rounded-full bg-custom-text-300" style="animation-delay: 300ms" />
            </div>
            <span class="text-xs">Thinking...</span>
          </div>
        </div>

        <!-- Input -->
        <ChatInput :disabled="loading" @send="handleSend" />
      </div>
    </Transition>

    <!-- Toggle button -->
    <button
      class="flex h-12 w-12 items-center justify-center rounded-full bg-brand-500 text-white shadow-lg transition-all hover:bg-brand-600 hover:scale-105"
      @click="toggle"
    >
      <MessageCircle v-if="!isOpen" class="h-5 w-5" />
      <X v-else class="h-5 w-5" />
    </button>
  </div>
</template>
