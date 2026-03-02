<script setup lang="ts">
import { computed } from 'vue'
import type { Issue, IssuePriority, IssueType, WorkLog, CreateWorkLogRequest, UpdateWorkLogRequest, EstimateSystem } from '@/types/issue.types'
import type { State, Label, ProjectMember } from '@/types/project.types'
import type { Sprint } from '@/types/sprint.types'
import { useAuthStore } from '@/stores/auth.store'
import StateSelector from './StateSelector.vue'
import PrioritySelector from './PrioritySelector.vue'
import TypeSelector from './TypeSelector.vue'
import EstimateSelector from './EstimateSelector.vue'
import AssigneeSelector from './AssigneeSelector.vue'
import LabelSelector from './LabelSelector.vue'
import PButton from '@/components/ui/PButton.vue'
import { formatDate } from '@/utils/helpers'
import { Calendar, Clock, Plus, Pencil, Trash2, Repeat, Zap, X } from 'lucide-vue-next'

interface Props {
  issue: Issue
  states: State[]
  labels: Label[]
  members: ProjectMember[]
  sprints: Sprint[]
  epics: Issue[]
  workLogs: WorkLog[]
  totalMinutes: number
  estimateSystem: EstimateSystem | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:state': [stateId: string]
  'update:priority': [priority: IssuePriority]
  'update:type': [issueType: IssueType]
  'update:estimate_point': [value: number | undefined]
  'update:assignees': [assigneeIds: string[]]
  'update:labels': [labelIds: string[]]
  'update:start_date': [date: string]
  'update:target_date': [date: string]
  'update:sprint': [sprintId: string]
  'remove:sprint': [sprintId: string]
  'update:parent': [parentId: string | null]
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

const assigneeIds = computed(() =>
  (props.issue.assignees || []).map((a) => a.id)
)

const labelIds = computed(() =>
  (props.issue.labels || []).map((l) => l.id)
)

const membersAsUserSummary = computed(() =>
  props.members.map((m) => ({
    id: m.user_id,
    email: m.email || '',
    first_name: m.first_name,
    last_name: m.last_name,
    display_name: m.display_name,
    avatar_url: m.avatar_url,
  }))
)

const isEpic = computed(() => props.issue.issue_type === 'epic')

const parentEpic = computed(() => {
  if (!props.issue.parent_id) return null
  return props.epics.find((e) => e.id === props.issue.parent_id) || null
})

function formatDuration(minutes: number): string {
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  if (h > 0 && m > 0) return `${h}h ${m}m`
  if (h > 0) return `${h}h`
  return `${m}m`
}

function formatLogDate(start: string, end: string): string {
  const s = new Date(start)
  const e = new Date(end)
  const fmt = (d: Date) => d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  if (start.substring(0, 10) === end.substring(0, 10)) return fmt(s)
  return `${fmt(s)} – ${fmt(e)}`
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

    <!-- Estimate -->
    <div v-if="props.estimateSystem">
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Estimate</label>
      <EstimateSelector
        :model-value="props.issue.estimate_point"
        :estimate-system="props.estimateSystem"
        @update:model-value="emit('update:estimate_point', $event)"
      />
    </div>

    <!-- Assignees -->
    <div>
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Assignees</label>
      <AssigneeSelector
        :model-value="assigneeIds"
        :members="membersAsUserSummary"
        @update:model-value="emit('update:assignees', $event)"
      />
    </div>

    <!-- Labels -->
    <div>
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Labels</label>
      <LabelSelector
        :model-value="labelIds"
        :labels="props.labels"
        @update:model-value="emit('update:labels', $event)"
      />
    </div>

    <!-- Epic (Parent) — only for non-epic issues -->
    <div v-if="!isEpic && props.epics.length > 0">
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Epic</label>
      <div v-if="parentEpic" class="flex items-center gap-2 rounded-md border border-custom-border-200 px-2.5 py-1.5 text-xs">
        <Zap class="h-3 w-3 text-orange-500 flex-shrink-0" />
        <span class="flex-1 text-custom-text-200 truncate">{{ parentEpic.name }}</span>
        <button
          @click="emit('update:parent', null)"
          class="rounded p-0.5 text-custom-text-300 hover:bg-red-50 hover:text-red-500 transition-colors"
          title="Remove from epic"
        >
          <X class="h-3 w-3" />
        </button>
      </div>
      <select
        v-if="!parentEpic"
        class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-2.5 py-1.5 text-xs text-custom-text-300"
        @change="(e) => { const target = e.target as HTMLSelectElement; if (target.value) { emit('update:parent', target.value); target.value = '' } }"
      >
        <option value="">Select epic...</option>
        <option
          v-for="epic in props.epics"
          :key="epic.id"
          :value="epic.id"
        >
          {{ epic.name }}
        </option>
      </select>
      <select
        v-else
        class="mt-1.5 w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-2.5 py-1.5 text-xs text-custom-text-300"
        @change="(e) => { const target = e.target as HTMLSelectElement; if (target.value) { emit('update:parent', target.value); target.value = '' } }"
      >
        <option value="">Change epic...</option>
        <option
          v-for="epic in props.epics.filter(e => e.id !== parentEpic?.id)"
          :key="epic.id"
          :value="epic.id"
        >
          {{ epic.name }}
        </option>
      </select>
    </div>

    <!-- Sprint — only for non-epic issues (epics span multiple sprints) -->
    <div v-if="!isEpic && props.sprints.length > 0">
      <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Sprint</label>
      <div class="space-y-1">
        <div
          v-for="sprint in (props.issue.sprints || [])"
          :key="sprint.id"
          class="flex items-center gap-2 rounded-md border border-custom-border-200 px-2.5 py-1.5 text-xs"
        >
          <Repeat class="h-3 w-3 text-custom-text-300 flex-shrink-0" />
          <span class="flex-1 text-custom-text-200 truncate">{{ sprint.name }}</span>
          <span
            class="rounded-full px-1.5 py-0.5 text-2xs font-medium"
            :class="sprint.status === 'active' ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400' : sprint.status === 'completed' ? 'bg-gray-100 text-gray-600 dark:bg-gray-800/50 dark:text-gray-400' : 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'"
          >
            {{ sprint.status }}
          </span>
          <button
            @click="emit('remove:sprint', sprint.id)"
            class="rounded p-0.5 text-custom-text-300 hover:bg-red-50 hover:text-red-500 transition-colors"
            title="Remove from sprint"
          >
            <Trash2 class="h-3 w-3" />
          </button>
        </div>
      </div>
      <select
        class="mt-1.5 w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-2.5 py-1.5 text-xs text-custom-text-300"
        @change="(e) => { const target = e.target as HTMLSelectElement; if (target.value) { emit('update:sprint', target.value); target.value = '' } }"
      >
        <option value="">+ Add to sprint</option>
        <option
          v-for="s in props.sprints.filter(sp => !(props.issue.sprints || []).some(is => is.id === sp.id))"
          :key="s.id"
          :value="s.id"
        >
          {{ s.name }}{{ s.status === 'active' ? ' (Active)' : '' }}
        </option>
      </select>
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
              <span class="text-custom-text-300">{{ formatLogDate(log.start_date, log.end_date) }}</span>
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
