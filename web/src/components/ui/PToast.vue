<script setup lang="ts">
import { useToast } from '@/composables/useToast'
import { CheckCircle, AlertCircle, AlertTriangle, Info, X } from 'lucide-vue-next'

const { toasts, removeToast } = useToast()

const config = {
  success: { icon: CheckCircle, bg: 'bg-green-50 dark:bg-green-900/20', border: 'border-green-200 dark:border-green-800', text: 'text-green-800 dark:text-green-200', iconColor: 'text-green-500' },
  error: { icon: AlertCircle, bg: 'bg-red-50 dark:bg-red-900/20', border: 'border-red-200 dark:border-red-800', text: 'text-red-800 dark:text-red-200', iconColor: 'text-red-500' },
  warning: { icon: AlertTriangle, bg: 'bg-yellow-50 dark:bg-yellow-900/20', border: 'border-yellow-200 dark:border-yellow-800', text: 'text-yellow-800 dark:text-yellow-200', iconColor: 'text-yellow-500' },
  info: { icon: Info, bg: 'bg-blue-50 dark:bg-blue-900/20', border: 'border-blue-200 dark:border-blue-800', text: 'text-blue-800 dark:text-blue-200', iconColor: 'text-blue-500' },
}
</script>

<template>
  <Teleport to="body">
    <div class="fixed bottom-4 right-4 z-[9999] flex flex-col gap-2 max-w-sm">
      <TransitionGroup
        enter-active-class="transition-all duration-300 ease-out"
        leave-active-class="transition-all duration-200 ease-in"
        enter-from-class="translate-x-full opacity-0"
        enter-to-class="translate-x-0 opacity-100"
        leave-from-class="translate-x-0 opacity-100"
        leave-to-class="translate-x-full opacity-0"
      >
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="flex items-start gap-2.5 rounded-lg border px-4 py-3 shadow-lg"
          :class="[config[toast.type].bg, config[toast.type].border]"
        >
          <component
            :is="config[toast.type].icon"
            class="h-4 w-4 mt-0.5 flex-shrink-0"
            :class="config[toast.type].iconColor"
          />
          <p class="flex-1 text-sm font-medium" :class="config[toast.type].text">
            {{ toast.message }}
          </p>
          <button
            class="flex-shrink-0 rounded p-0.5 opacity-60 hover:opacity-100 transition-opacity"
            :class="config[toast.type].text"
            @click="removeToast(toast.id)"
          >
            <X class="h-3.5 w-3.5" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>
