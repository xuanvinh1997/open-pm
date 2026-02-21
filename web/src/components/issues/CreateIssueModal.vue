<script setup lang="ts">
import { ref } from 'vue'
import type { IssuePriority } from '@/types/issue.types'
import type { State, Label } from '@/types/project.types'
import PModal from '@/components/ui/PModal.vue'
import PInput from '@/components/ui/PInput.vue'
import PTextarea from '@/components/ui/PTextarea.vue'
import PButton from '@/components/ui/PButton.vue'
import StateSelector from './StateSelector.vue'
import PrioritySelector from './PrioritySelector.vue'

interface Props {
  open: boolean
  states: State[]
  labels: Label[]
  defaultStateId?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  create: [data: { name: string; state_id?: string; priority: IssuePriority; description_html?: string }]
}>()

const name = ref('')
const description = ref('')
const stateId = ref(props.defaultStateId || '')
const priority = ref<IssuePriority>('none')
const loading = ref(false)

function handleSubmit() {
  if (!name.value.trim()) return
  loading.value = true
  emit('create', {
    name: name.value,
    state_id: stateId.value || undefined,
    priority: priority.value,
    description_html: description.value || undefined,
  })
  // Reset
  name.value = ''
  description.value = ''
  priority.value = 'none'
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

      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">State</label>
          <StateSelector v-model="stateId" :states="props.states" />
        </div>
        <div>
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Priority</label>
          <PrioritySelector v-model="priority" />
        </div>
      </div>
    </form>

    <template #footer>
      <PButton variant="secondary" @click="emit('update:open', false)">Cancel</PButton>
      <PButton variant="primary" :loading="loading" @click="handleSubmit">Create issue</PButton>
    </template>
  </PModal>
</template>
