<script setup lang="ts">
import type { IssueType } from '@/types/issue.types'
import { ISSUE_TYPE_CONFIG } from '@/utils/issue-helpers'
import PDropdown from '@/components/ui/PDropdown.vue'
import { ChevronDown } from 'lucide-vue-next'

interface Props {
  modelValue?: IssueType
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: 'task',
})

const emit = defineEmits<{
  'update:modelValue': [value: IssueType]
}>()

const issueTypes: IssueType[] = ['story', 'bug', 'task', 'epic']
</script>

<template>
  <PDropdown align="left" width="12rem">
    <template #trigger>
      <button
        type="button"
        class="flex items-center gap-2 rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 hover:bg-custom-background-90 transition-colors w-full"
      >
        <component
          :is="ISSUE_TYPE_CONFIG[props.modelValue].icon"
          class="h-3.5 w-3.5 flex-shrink-0"
          :style="{ color: ISSUE_TYPE_CONFIG[props.modelValue].color }"
        />
        <span class="flex-1 text-left">{{ ISSUE_TYPE_CONFIG[props.modelValue].label }}</span>
        <ChevronDown class="h-3.5 w-3.5 text-custom-text-300 flex-shrink-0" />
      </button>
    </template>
    <template #default="{ close }">
      <button
        v-for="t in issueTypes"
        :key="t"
        class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-custom-background-80 transition-colors"
        :class="t === props.modelValue ? 'bg-custom-background-80 text-custom-text-100' : 'text-custom-text-200'"
        @click="emit('update:modelValue', t); close()"
      >
        <component
          :is="ISSUE_TYPE_CONFIG[t].icon"
          class="h-3.5 w-3.5 flex-shrink-0"
          :style="{ color: ISSUE_TYPE_CONFIG[t].color }"
        />
        <span>{{ ISSUE_TYPE_CONFIG[t].label }}</span>
      </button>
    </template>
  </PDropdown>
</template>
