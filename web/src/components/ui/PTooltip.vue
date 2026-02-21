<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  content: string
  position?: 'top' | 'bottom' | 'left' | 'right'
  delay?: number
}

const props = withDefaults(defineProps<Props>(), {
  position: 'top',
  delay: 300,
})

const isVisible = ref(false)
let timeout: ReturnType<typeof setTimeout> | null = null

function show() {
  timeout = setTimeout(() => {
    isVisible.value = true
  }, props.delay)
}

function hide() {
  if (timeout) {
    clearTimeout(timeout)
    timeout = null
  }
  isVisible.value = false
}

const positionClasses: Record<string, string> = {
  top: 'bottom-full left-1/2 -translate-x-1/2 mb-2',
  bottom: 'top-full left-1/2 -translate-x-1/2 mt-2',
  left: 'right-full top-1/2 -translate-y-1/2 mr-2',
  right: 'left-full top-1/2 -translate-y-1/2 ml-2',
}
</script>

<template>
  <div class="relative inline-flex" @mouseenter="show" @mouseleave="hide">
    <slot />
    <Transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition duration-75 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="isVisible"
        class="absolute z-50 whitespace-nowrap rounded-md bg-custom-text-100 px-2.5 py-1.5 text-xs text-white shadow-md pointer-events-none"
        :class="positionClasses[props.position]"
      >
        {{ props.content }}
      </div>
    </Transition>
  </div>
</template>
