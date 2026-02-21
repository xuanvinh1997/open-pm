<script setup lang="ts">
import { watch } from 'vue'
import { X } from 'lucide-vue-next'

interface Props {
  open: boolean
  size?: 'sm' | 'md' | 'lg' | 'xl'
  title?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  title: '',
})

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const sizeMap: Record<string, string> = {
  sm: 'max-w-sm',
  md: 'max-w-lg',
  lg: 'max-w-2xl',
  xl: 'max-w-4xl',
}

function close() {
  emit('update:open', false)
}

function handleBackdropClick(e: MouseEvent) {
  if (e.target === e.currentTarget) {
    close()
  }
}

watch(() => props.open, (val) => {
  if (val) {
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="props.open"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
        @click="handleBackdropClick"
      >
        <div
          class="w-full rounded-xl border border-custom-border-200 bg-custom-background-100 shadow-custom-lg overflow-hidden"
          :class="sizeMap[props.size]"
        >
          <!-- Header -->
          <div v-if="props.title || $slots.header" class="flex items-center justify-between border-b border-custom-border-200 px-6 py-4">
            <slot name="header">
              <h3 class="text-lg font-semibold text-custom-text-100">{{ props.title }}</h3>
            </slot>
            <button
              @click="close"
              class="rounded-md p-1 text-custom-text-300 hover:bg-custom-background-80 hover:text-custom-text-100 transition-colors"
            >
              <X class="h-5 w-5" />
            </button>
          </div>

          <!-- Body -->
          <div class="px-6 py-4">
            <slot />
          </div>

          <!-- Footer -->
          <div v-if="$slots.footer" class="flex items-center justify-end gap-3 border-t border-custom-border-200 px-6 py-4">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
