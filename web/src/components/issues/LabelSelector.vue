<script setup lang="ts">
import type { Label } from '@/types/project.types'
import PDropdown from '@/components/ui/PDropdown.vue'
import { ChevronDown, Check } from 'lucide-vue-next'

interface Props {
  modelValue?: string[]
  labels: Label[]
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => [],
})

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
}>()

function toggleLabel(labelId: string) {
  const current = [...props.modelValue]
  const index = current.indexOf(labelId)
  if (index >= 0) {
    current.splice(index, 1)
  } else {
    current.push(labelId)
  }
  emit('update:modelValue', current)
}

function isSelected(labelId: string) {
  return props.modelValue.includes(labelId)
}
</script>

<template>
  <PDropdown align="left" width="14rem">
    <template #trigger>
      <button
        type="button"
        class="flex items-center gap-2 rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 hover:bg-custom-background-90 transition-colors w-full"
      >
        <span v-if="props.modelValue.length > 0" class="text-sm">{{ props.modelValue.length }} label{{ props.modelValue.length > 1 ? 's' : '' }}</span>
        <span v-else class="text-custom-text-300">Labels</span>
        <ChevronDown class="h-3.5 w-3.5 text-custom-text-300 flex-shrink-0 ml-auto" />
      </button>
    </template>
    <template #default>
      <button
        v-for="label in props.labels"
        :key="label.id"
        class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-custom-background-80 transition-colors text-custom-text-200"
        @click="toggleLabel(label.id)"
      >
        <span class="h-3 w-3 rounded-full flex-shrink-0" :style="{ backgroundColor: label.color }" />
        <span class="flex-1 truncate">{{ label.name }}</span>
        <Check v-if="isSelected(label.id)" class="h-3.5 w-3.5 text-brand-600 flex-shrink-0" />
      </button>
      <div v-if="props.labels.length === 0" class="px-3 py-4 text-center text-xs text-custom-text-300">
        No labels
      </div>
    </template>
  </PDropdown>
</template>
