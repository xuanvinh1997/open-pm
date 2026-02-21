<script setup lang="ts">
import { ref } from 'vue'
import type { State } from '@/types/project.types'
import { ChevronRight, Circle } from 'lucide-vue-next'

interface Props {
  state: State
  count: number
}

const props = defineProps<Props>()
const isExpanded = ref(true)

defineExpose({ isExpanded })
</script>

<template>
  <button
    class="flex w-full items-center gap-2 px-4 py-2 hover:bg-custom-background-80 transition-colors"
    @click="isExpanded = !isExpanded"
  >
    <ChevronRight
      class="h-3.5 w-3.5 text-custom-text-300 transition-transform duration-100"
      :class="{ 'rotate-90': isExpanded }"
    />
    <Circle
      class="h-3 w-3 flex-shrink-0"
      :style="{ color: props.state.color }"
      fill="currentColor"
    />
    <span class="text-sm font-medium text-custom-text-100">{{ props.state.name }}</span>
    <span class="rounded-full bg-custom-background-80 px-2 py-0.5 text-2xs text-custom-text-300">
      {{ props.count }}
    </span>
  </button>
</template>
