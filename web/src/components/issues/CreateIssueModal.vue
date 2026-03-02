<script setup lang="ts">
import { ref, computed } from 'vue'
import type { IssuePriority, IssueType, EstimateSystem, UserSummary } from '@/types/issue.types'
import type { State, Label, ProjectMember } from '@/types/project.types'
import PModal from '@/components/ui/PModal.vue'
import PInput from '@/components/ui/PInput.vue'
import PButton from '@/components/ui/PButton.vue'
import { RichTextEditor } from '@/components/editor'
import StateSelector from './StateSelector.vue'
import PrioritySelector from './PrioritySelector.vue'
import TypeSelector from './TypeSelector.vue'
import EstimateSelector from './EstimateSelector.vue'
import AssigneeSelector from './AssigneeSelector.vue'
import LabelSelector from './LabelSelector.vue'

interface Props {
  open: boolean
  states: State[]
  labels: Label[]
  members?: ProjectMember[]
  defaultStateId?: string
  estimateSystem?: EstimateSystem | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  create: [data: {
    name: string
    state_id?: string
    priority: IssuePriority
    issue_type: IssueType
    description_html?: string
    description_json?: Record<string, unknown>
    description_stripped?: string
    start_date?: string
    target_date?: string
    estimate_point?: number
    assignee_ids?: string[]
    label_ids?: string[]
  }]
}>()

const name = ref('')
const descriptionHtml = ref('')
const descriptionJson = ref<Record<string, unknown>>({})
const descriptionStripped = ref('')
const stateId = ref(props.defaultStateId || '')
const priority = ref<IssuePriority>('none')
const issueType = ref<IssueType>('task')
const estimatePoint = ref<number | undefined>(undefined)
const assigneeIds = ref<string[]>([])
const labelIds = ref<string[]>([])
const startDate = ref('')
const targetDate = ref('')
const loading = ref(false)

const membersAsUserSummary = computed<UserSummary[]>(() =>
  (props.members || []).map((m) => ({
    id: m.user_id,
    email: m.email || '',
    first_name: m.first_name,
    last_name: m.last_name,
    display_name: m.display_name,
    avatar_url: m.avatar_url,
  }))
)

function handleSubmit() {
  if (!name.value.trim()) return
  loading.value = true
  emit('create', {
    name: name.value,
    state_id: stateId.value || undefined,
    priority: priority.value,
    issue_type: issueType.value,
    description_html: descriptionHtml.value || undefined,
    description_json: Object.keys(descriptionJson.value).length > 0 ? descriptionJson.value : undefined,
    description_stripped: descriptionStripped.value || undefined,
    start_date: startDate.value || undefined,
    target_date: targetDate.value || undefined,
    estimate_point: estimatePoint.value,
    assignee_ids: assigneeIds.value.length > 0 ? assigneeIds.value : undefined,
    label_ids: labelIds.value.length > 0 ? labelIds.value : undefined,
  })
  // Reset
  name.value = ''
  descriptionHtml.value = ''
  descriptionJson.value = {}
  descriptionStripped.value = ''
  priority.value = 'none'
  issueType.value = 'task'
  estimatePoint.value = undefined
  assigneeIds.value = []
  labelIds.value = []
  startDate.value = ''
  targetDate.value = ''
  loading.value = false
}
</script>

<template>
  <PModal :open="props.open" @update:open="emit('update:open', $event)" title="Create issue" size="md">
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Title</label>
        <PInput v-model="name" placeholder="Issue title" autofocus />
      </div>

      <div>
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Description</label>
        <RichTextEditor
          v-model="descriptionHtml"
          toolbar="compact"
          placeholder="Add a description..."
          min-height="80px"
          max-height="200px"
          @update:json="(v) => descriptionJson = v"
          @update:stripped="(v) => descriptionStripped = v"
        />
      </div>

      <div class="grid grid-cols-3 gap-3">
        <div>
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Type</label>
          <TypeSelector v-model="issueType" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">State</label>
          <StateSelector v-model="stateId" :states="props.states" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Priority</label>
          <PrioritySelector v-model="priority" />
        </div>
      </div>

      <div class="grid grid-cols-2 gap-3">
        <div v-if="props.estimateSystem">
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Estimate</label>
          <EstimateSelector v-model="estimatePoint" :estimate-system="props.estimateSystem" />
        </div>
        <div v-if="props.members && props.members.length > 0">
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Assignees</label>
          <AssigneeSelector v-model="assigneeIds" :members="membersAsUserSummary" />
        </div>
      </div>

      <div v-if="props.labels && props.labels.length > 0">
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Labels</label>
        <LabelSelector v-model="labelIds" :labels="props.labels" />
      </div>

      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Start date</label>
          <input
            v-model="startDate"
            type="date"
            class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 focus:border-custom-primary-100 focus:outline-none focus:ring-1 focus:ring-custom-primary-100"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Due date</label>
          <input
            v-model="targetDate"
            type="date"
            class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 focus:border-custom-primary-100 focus:outline-none focus:ring-1 focus:ring-custom-primary-100"
          />
        </div>
      </div>
    </form>

    <template #footer>
      <PButton variant="secondary" @click="emit('update:open', false)">Cancel</PButton>
      <PButton variant="primary" :loading="loading" @click="handleSubmit">Create issue</PButton>
    </template>
  </PModal>
</template>
