<script setup lang="ts">
import type { State } from '@/types/project.types'
import PDropdown from '@/components/ui/PDropdown.vue'
import { Circle, ChevronDown } from 'lucide-vue-next'

interface Props {
  modelValue?: string
  states: State[]
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: 'Select state',
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

function selectedState() {
  return props.states.find((s) => s.id === props.modelValue)
}
</script>

<template>
  <PDropdown align="left" width="14rem">
    <template #trigger>
      <button
        type="button"
        class="flex items-center gap-2 rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 hover:bg-custom-background-90 transition-colors w-full"
      >
        <template v-if="selectedState()">
          <Circle class="h-3 w-3 flex-shrink-0" :style="{ color: selectedState()!.color }" fill="currentColor" />
          <span class="flex-1 text-left truncate">{{ selectedState()!.name }}</span>
        </template>
        <span v-else class="flex-1 text-left text-custom-text-300">{{ props.placeholder }}</span>
        <ChevronDown class="h-3.5 w-3.5 text-custom-text-300 flex-shrink-0" />
      </button>
    </template>
    <template #default="{ close }">
      <button
        v-for="state in props.states"
        :key="state.id"
        class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-custom-background-80 transition-colors"
        :class="state.id === props.modelValue ? 'bg-custom-background-80 text-custom-text-100' : 'text-custom-text-200'"
        @click="emit('update:modelValue', state.id); close()"
      >
        <Circle class="h-3 w-3 flex-shrink-0" :style="{ color: state.color }" fill="currentColor" />
        <span>{{ state.name }}</span>
      </button>
    </template>
  </PDropdown>
</template>
