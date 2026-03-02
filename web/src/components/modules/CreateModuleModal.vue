<script setup lang="ts">
import { ref } from 'vue'
import type { ModuleStatus } from '@/types/module.types'
import PModal from '@/components/ui/PModal.vue'
import PInput from '@/components/ui/PInput.vue'
import PTextarea from '@/components/ui/PTextarea.vue'
import PButton from '@/components/ui/PButton.vue'

interface Props {
  open: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  create: [data: { name: string; description?: string; start_date?: string; target_date?: string; status?: ModuleStatus }]
}>()

const name = ref('')
const description = ref('')
const startDate = ref('')
const targetDate = ref('')
const status = ref<ModuleStatus>('backlog')
const loading = ref(false)

const statusOptions: { value: ModuleStatus; label: string }[] = [
  { value: 'backlog', label: 'Backlog' },
  { value: 'planned', label: 'Planned' },
  { value: 'in-progress', label: 'In Progress' },
  { value: 'paused', label: 'Paused' },
  { value: 'completed', label: 'Completed' },
  { value: 'cancelled', label: 'Cancelled' },
]

function handleSubmit() {
  if (!name.value.trim()) return
  loading.value = true
  emit('create', {
    name: name.value,
    description: description.value || undefined,
    start_date: startDate.value || undefined,
    target_date: targetDate.value || undefined,
    status: status.value,
  })
  name.value = ''
  description.value = ''
  startDate.value = ''
  targetDate.value = ''
  status.value = 'backlog'
  loading.value = false
}
</script>

<template>
  <PModal :open="props.open" @update:open="emit('update:open', $event)" title="Create module" size="md">
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Name</label>
        <PInput v-model="name" placeholder="Module name" autofocus />
      </div>

      <div>
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Description</label>
        <PTextarea v-model="description" placeholder="Describe the module..." :rows="3" />
      </div>

      <div>
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Status</label>
        <select
          v-model="status"
          class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 focus:border-custom-primary-100 focus:outline-none focus:ring-1 focus:ring-custom-primary-100"
        >
          <option v-for="opt in statusOptions" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
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
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Target date</label>
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
      <PButton variant="primary" :loading="loading" @click="handleSubmit">Create module</PButton>
    </template>
  </PModal>
</template>
