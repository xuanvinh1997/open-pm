<script setup lang="ts">
import type { IssuePriority } from '@/types/issue.types'
import { PRIORITY_CONFIG } from '@/utils/issue-helpers'
import PDropdown from '@/components/ui/PDropdown.vue'
import { ChevronDown } from 'lucide-vue-next'

interface Props {
  modelValue?: IssuePriority
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: 'none',
})

const emit = defineEmits<{
  'update:modelValue': [value: IssuePriority]
}>()

const priorities: IssuePriority[] = ['urgent', 'high', 'medium', 'low', 'none']
</script>

<template>
  <PDropdown align="left" width="12rem">
    <template #trigger>
      <button
        type="button"
        class="flex items-center gap-2 rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 hover:bg-custom-background-90 transition-colors w-full"
      >
        <component
          :is="PRIORITY_CONFIG[props.modelValue].icon"
          class="h-3.5 w-3.5 flex-shrink-0"
          :style="{ color: PRIORITY_CONFIG[props.modelValue].color }"
        />
        <span class="flex-1 text-left">{{ PRIORITY_CONFIG[props.modelValue].label }}</span>
        <ChevronDown class="h-3.5 w-3.5 text-custom-text-300 flex-shrink-0" />
      </button>
    </template>
    <template #default="{ close }">
      <button
        v-for="p in priorities"
        :key="p"
        class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-custom-background-80 transition-colors"
        :class="p === props.modelValue ? 'bg-custom-background-80 text-custom-text-100' : 'text-custom-text-200'"
        @click="emit('update:modelValue', p); close()"
      >
        <component
          :is="PRIORITY_CONFIG[p].icon"
          class="h-3.5 w-3.5 flex-shrink-0"
          :style="{ color: PRIORITY_CONFIG[p].color }"
        />
        <span>{{ PRIORITY_CONFIG[p].label }}</span>
      </button>
    </template>
  </PDropdown>
</template>
