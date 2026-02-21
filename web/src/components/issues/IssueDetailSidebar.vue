<script setup lang="ts">
import type { Issue, IssuePriority } from '@/types/issue.types'
import type { State, Label } from '@/types/project.types'
import StateSelector from './StateSelector.vue'
import PrioritySelector from './PrioritySelector.vue'
import { formatDate } from '@/utils/helpers'
import { Calendar } from 'lucide-vue-next'

interface Props {
  issue: Issue
  states: State[]
  labels: Label[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:state': [stateId: string]
  'update:priority': [priority: IssuePriority]
}>()
</script>

<template>
  <div class="space-y-4">
    <h3 class="text-sm font-semibold text-custom-text-100">Details</h3>

    <!-- State -->
    <div>
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">State</label>
      <StateSelector
        :model-value="props.issue.state_id || ''"
        :states="props.states"
        @update:model-value="emit('update:state', $event)"
      />
    </div>

    <!-- Priority -->
    <div>
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Priority</label>
      <PrioritySelector
        :model-value="props.issue.priority"
        @update:model-value="emit('update:priority', $event)"
      />
    </div>

    <!-- Assignees -->
    <div v-if="props.issue.assignees && props.issue.assignees.length > 0">
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Assignees</label>
      <div class="flex flex-wrap gap-1">
        <span
          v-for="assignee in props.issue.assignees"
          :key="assignee.id"
          class="inline-flex items-center gap-1 rounded-full bg-custom-background-80 px-2 py-1 text-xs text-custom-text-200"
        >
          {{ assignee.first_name }} {{ assignee.last_name }}
        </span>
      </div>
    </div>

    <!-- Labels -->
    <div v-if="props.issue.labels && props.issue.labels.length > 0">
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Labels</label>
      <div class="flex flex-wrap gap-1">
        <span
          v-for="label in props.issue.labels"
          :key="label.id"
          class="inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-xs"
          :style="{ backgroundColor: label.color + '1A', color: label.color }"
        >
          <span class="h-2 w-2 rounded-full" :style="{ backgroundColor: label.color }" />
          {{ label.name }}
        </span>
      </div>
    </div>

    <!-- Dates -->
    <div v-if="props.issue.start_date || props.issue.target_date">
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Dates</label>
      <div class="space-y-1 text-sm text-custom-text-200">
        <div v-if="props.issue.start_date" class="flex items-center gap-2">
          <Calendar class="h-3.5 w-3.5 text-custom-text-300" />
          <span>Start: {{ formatDate(props.issue.start_date) }}</span>
        </div>
        <div v-if="props.issue.target_date" class="flex items-center gap-2">
          <Calendar class="h-3.5 w-3.5 text-custom-text-300" />
          <span>Due: {{ formatDate(props.issue.target_date) }}</span>
        </div>
      </div>
    </div>

    <!-- Created -->
    <div>
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Created</label>
      <p class="text-sm text-custom-text-200">{{ formatDate(props.issue.created_at) }}</p>
    </div>
  </div>
</template>
