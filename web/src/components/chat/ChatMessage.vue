<script setup lang="ts">
import type { ChatMessage } from '@/types/chat.types'
import { Bot, User } from 'lucide-vue-next'

interface Props {
  message: ChatMessage
}

defineProps<Props>()
</script>

<template>
  <div class="flex gap-3" :class="message.role === 'user' ? 'flex-row-reverse' : ''">
    <div
      class="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-full"
      :class="message.role === 'user' ? 'bg-brand-500' : 'bg-custom-background-80'"
    >
      <User v-if="message.role === 'user'" class="h-3.5 w-3.5 text-white" />
      <Bot v-else class="h-3.5 w-3.5 text-custom-text-200" />
    </div>
    <div
      class="max-w-[80%] rounded-xl px-3.5 py-2.5 text-sm leading-relaxed"
      :class="
        message.role === 'user'
          ? 'bg-brand-500 text-white'
          : 'bg-custom-background-80 text-custom-text-200'
      "
    >
      <div class="whitespace-pre-wrap">{{ message.content }}</div>
    </div>
  </div>
</template>
