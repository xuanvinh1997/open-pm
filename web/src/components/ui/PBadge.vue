<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  color?: string
  variant?: 'solid' | 'soft' | 'outline'
  size?: 'xs' | 'sm'
}

const props = withDefaults(defineProps<Props>(), {
  color: '#A3A3A3',
  variant: 'soft',
  size: 'xs',
})

const sizeClasses: Record<string, string> = {
  xs: 'px-1.5 py-0.5 text-2xs',
  sm: 'px-2 py-0.5 text-xs',
}

const style = computed(() => {
  switch (props.variant) {
    case 'solid':
      return { backgroundColor: props.color, color: '#FFFFFF' }
    case 'soft':
      return { backgroundColor: props.color + '1A', color: props.color }
    case 'outline':
      return { borderColor: props.color, color: props.color }
    default:
      return {}
  }
})
</script>

<template>
  <span
    class="inline-flex items-center gap-1 rounded-full font-medium whitespace-nowrap"
    :class="[
      sizeClasses[props.size],
      props.variant === 'outline' ? 'border' : '',
    ]"
    :style="style"
  >
    <slot />
  </span>
</template>
