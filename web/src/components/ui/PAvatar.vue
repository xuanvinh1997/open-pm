<script setup lang="ts">
import { computed } from 'vue'
import { getInitials, hashStringToColor } from '@/utils/helpers'

interface Props {
  src?: string
  name?: string
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
}

const props = withDefaults(defineProps<Props>(), {
  src: '',
  name: '',
  size: 'md',
})

const sizeMap: Record<string, string> = {
  xs: 'h-5 w-5 text-2xs',
  sm: 'h-6 w-6 text-2xs',
  md: 'h-7 w-7 text-xs',
  lg: 'h-8 w-8 text-xs',
  xl: 'h-10 w-10 text-sm',
}

const initials = computed(() => getInitials(props.name || '?'))
const bgColor = computed(() => hashStringToColor(props.name || 'default'))
</script>

<template>
  <div
    class="inline-flex items-center justify-center rounded-full font-medium flex-shrink-0"
    :class="sizeMap[props.size]"
  >
    <img
      v-if="props.src"
      :src="props.src"
      :alt="props.name"
      class="h-full w-full rounded-full object-cover"
    />
    <span
      v-else
      class="flex h-full w-full items-center justify-center rounded-full text-white"
      :style="{ backgroundColor: bgColor }"
    >
      {{ initials }}
    </span>
  </div>
</template>
