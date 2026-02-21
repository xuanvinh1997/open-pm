<script setup lang="ts">
import { computed } from 'vue'
import type { Issue, IssuePriority, IssueType, WorkLog, CreateWorkLogRequest, UpdateWorkLogRequest } from '@/types/issue.types'
import type { State, Label } from '@/types/project.types'
import { useAuthStore } from '@/stores/auth.store'
import StateSelector from './StateSelector.vue'
import PrioritySelector from './PrioritySelector.vue'
import TypeSelector from './TypeSelector.vue'
import PButton from '@/components/ui/PButton.vue'
import { formatDate } from '@/utils/helpers'
import { Calendar, Clock, Plus, Pencil, Trash2 } from 'lucide-vue-next'

interface Props {
  issue: Issue
  states: State[]
  labels: Label[]
  workLogs: WorkLog[]
  totalMinutes: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:state': [stateId: string]
  'update:priority': [priority: IssuePriority]
  'update:type': [issueType: IssueType]
  'update:start_date': [date: string]
  'update:target_date': [date: string]
  'log-work': []
  'edit-work-log': [workLog: WorkLog]
  'delete-work-log': [id: string]
}>()

const authStore = useAuthStore()

const startDateValue = computed(() => {
  if (!props.issue.start_date) return ''
  return props.issue.start_date.substring(0, 10)
})

const targetDateValue = computed(() => {
  if (!props.issue.target_date) return ''
  return props.issue.target_date.substring(0, 10)
})

function formatDuration(minutes: number): string {
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  if (h > 0 && m > 0) return `${h}h ${m}m`
  if (h > 0) return `${h}h`
  return `${m}m`
}

function formatLogDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

function isOwnLog(log: WorkLog): boolean {
  return authStore.user?.id === log.logged_by
}
</script>

<template>
  <div class="space-y-4">
    <h3 class="text-sm font-semibold text-custom-text-100">Details</h3>

    <!-- Type -->
    <div>
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Type</label>
      <TypeSelector
        :model-value="props.issue.issue_type || 'task'"
        @update:model-value="emit('update:type', $event)"
      />
    </div>

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
    <div>
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Start date</label>
      <div class="flex items-center gap-2">
        <Calendar class="h-3.5 w-3.5 text-custom-text-300 flex-shrink-0" />
        <input
          type="date"
          :value="startDateValue"
          @change="emit('update:start_date', ($event.target as HTMLInputElement).value)"
          class="flex-1 rounded-md border border-custom-border-200 bg-custom-background-100 px-2 py-1 text-sm text-custom-text-200 focus:border-custom-primary-100 focus:outline-none"
        />
      </div>
    </div>

    <div>
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Due date</label>
      <div class="flex items-center gap-2">
        <Calendar class="h-3.5 w-3.5 text-custom-text-300 flex-shrink-0" />
        <input
          type="date"
          :value="targetDateValue"
          @change="emit('update:target_date', ($event.target as HTMLInputElement).value)"
          class="flex-1 rounded-md border border-custom-border-200 bg-custom-background-100 px-2 py-1 text-sm text-custom-text-200 focus:border-custom-primary-100 focus:outline-none"
        />
      </div>
    </div>

    <!-- Time Tracking -->
    <div>
      <div class="mb-1.5 flex items-center justify-between">
        <label class="text-xs font-medium text-custom-text-300">Time tracking</label>
        <button
          @click="emit('log-work')"
          class="flex items-center gap-1 rounded-md px-1.5 py-0.5 text-xs text-custom-primary-100 hover:bg-custom-background-80 transition-colors"
        >
          <Plus class="h-3 w-3" />
          Log work
        </button>
      </div>

      <!-- Total time -->
      <div class="flex items-center gap-2 mb-2">
        <Clock class="h-3.5 w-3.5 text-custom-text-300 flex-shrink-0" />
        <span class="text-sm text-custom-text-200">
          {{ props.totalMinutes > 0 ? formatDuration(props.totalMinutes) : 'No time logged' }}
        </span>
      </div>

      <!-- Work log entries -->
      <div v-if="props.workLogs.length > 0" class="space-y-1.5">
        <div
          v-for="log in props.workLogs"
          :key="log.id"
          class="group flex items-start gap-2 rounded-md border border-custom-border-200 px-2.5 py-2 text-xs"
        >
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-1.5">
              <span class="font-medium text-custom-text-100">{{ formatDuration(log.duration_minutes) }}</span>
              <span class="text-custom-text-300">&middot;</span>
              <span class="text-custom-text-300">{{ formatLogDate(log.logged_at) }}</span>
            </div>
            <p v-if="log.description" class="mt-0.5 text-custom-text-300 truncate">{{ log.description }}</p>
            <p class="mt-0.5 text-custom-text-300">{{ log.display_name || `${log.first_name} ${log.last_name}` }}</p>
          </div>
          <div v-if="isOwnLog(log)" class="flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0">
            <button
              @click="emit('edit-work-log', log)"
              class="rounded p-1 text-custom-text-300 hover:bg-custom-background-80 hover:text-custom-text-100 transition-colors"
            >
              <Pencil class="h-3 w-3" />
            </button>
            <button
              @click="emit('delete-work-log', log.id)"
              class="rounded p-1 text-custom-text-300 hover:bg-red-50 hover:text-red-500 transition-colors"
            >
              <Trash2 class="h-3 w-3" />
            </button>
          </div>
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
