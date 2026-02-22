<script setup lang="ts">
import type { EstimateSystem } from '@/types/issue.types'
import PDropdown from '@/components/ui/PDropdown.vue'
import { ChevronDown, Gauge } from 'lucide-vue-next'

interface Props {
  modelValue?: number
  estimateSystem: EstimateSystem | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: number | undefined]
}>()

function getLabel(value?: number): string {
  if (value == null || !props.estimateSystem) return 'No estimate'
  const opt = props.estimateSystem.estimates.find((e) => String(e.key) === String(value))
  return opt?.value ?? String(value)
}
</script>

<template>
  <PDropdown align="left" width="12rem">
    <template #trigger>
      <button
        type="button"
        class="flex items-center gap-2 rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 hover:bg-custom-background-90 transition-colors w-full"
      >
        <Gauge class="h-3.5 w-3.5 flex-shrink-0 text-custom-text-300" />
        <span class="flex-1 text-left">{{ getLabel(props.modelValue) }}</span>
        <ChevronDown class="h-3.5 w-3.5 text-custom-text-300 flex-shrink-0" />
      </button>
    </template>
    <template #default="{ close }">
      <button
        class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-custom-background-80 transition-colors"
        :class="props.modelValue == null ? 'bg-custom-background-80 text-custom-text-100' : 'text-custom-text-200'"
        @click="emit('update:modelValue', undefined); close()"
      >
        <span>No estimate</span>
      </button>
      <button
        v-for="opt in (props.estimateSystem?.estimates ?? [])"
        :key="opt.key"
        class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-custom-background-80 transition-colors"
        :class="String(props.modelValue) === String(opt.key) ? 'bg-custom-background-80 text-custom-text-100' : 'text-custom-text-200'"
        @click="emit('update:modelValue', Number(opt.key)); close()"
      >
        <span>{{ opt.value }}</span>
      </button>
    </template>
  </PDropdown>
</template>
