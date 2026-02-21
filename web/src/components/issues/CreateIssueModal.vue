<script setup lang="ts">
import { ref } from 'vue'
import type { IssuePriority, IssueType } from '@/types/issue.types'
import type { State, Label } from '@/types/project.types'
import PModal from '@/components/ui/PModal.vue'
import PInput from '@/components/ui/PInput.vue'
import PTextarea from '@/components/ui/PTextarea.vue'
import PButton from '@/components/ui/PButton.vue'
import StateSelector from './StateSelector.vue'
import PrioritySelector from './PrioritySelector.vue'
import TypeSelector from './TypeSelector.vue'

interface Props {
  open: boolean
  states: State[]
  labels: Label[]
  defaultStateId?: string
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
    start_date?: string
    target_date?: string
  }]
}>()

const name = ref('')
const description = ref('')
const stateId = ref(props.defaultStateId || '')
const priority = ref<IssuePriority>('none')
const issueType = ref<IssueType>('task')
const startDate = ref('')
const targetDate = ref('')
const loading = ref(false)

function handleSubmit() {
  if (!name.value.trim()) return
  loading.value = true
  emit('create', {
    name: name.value,
    state_id: stateId.value || undefined,
    priority: priority.value,
    issue_type: issueType.value,
    description_html: description.value || undefined,
    start_date: startDate.value || undefined,
    target_date: targetDate.value || undefined,
  })
  // Reset
  name.value = ''
  description.value = ''
  priority.value = 'none'
  issueType.value = 'task'
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
        <PTextarea v-model="description" placeholder="Add a description..." :rows="3" />
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
