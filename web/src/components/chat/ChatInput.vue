<script setup lang="ts">
import { ref } from 'vue'
import { Send } from 'lucide-vue-next'

interface Props {
  disabled?: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  send: [message: string]
}>()

const text = ref('')

function handleSend() {
  const trimmed = text.value.trim()
  if (!trimmed) return
  emit('send', trimmed)
  text.value = ''
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}
</script>

<template>
  <div class="flex items-end gap-2 border-t border-custom-border-200 p-3">
    <textarea
      v-model="text"
      :disabled="disabled"
      placeholder="Ask about your project..."
      rows="1"
      class="flex-1 resize-none rounded-lg border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-200 placeholder:text-custom-text-300 focus:border-brand-500 focus:outline-none disabled:opacity-50"
      @keydown="handleKeydown"
    />
    <button
      :disabled="disabled || !text.trim()"
      class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg bg-brand-500 text-white transition-colors hover:bg-brand-600 disabled:opacity-50 disabled:cursor-not-allowed"
      @click="handleSend"
    >
      <Send class="h-4 w-4" />
    </button>
  </div>
</template>
