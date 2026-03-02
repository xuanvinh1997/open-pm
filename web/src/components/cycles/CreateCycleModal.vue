<script setup lang="ts">
import { ref } from 'vue'
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
  create: [data: { name: string; description?: string; start_date?: string; end_date?: string }]
}>()

const name = ref('')
const description = ref('')
const startDate = ref('')
const endDate = ref('')
const loading = ref(false)

function handleSubmit() {
  if (!name.value.trim()) return
  loading.value = true
  emit('create', {
    name: name.value,
    description: description.value || undefined,
    start_date: startDate.value || undefined,
    end_date: endDate.value || undefined,
  })
  name.value = ''
  description.value = ''
  startDate.value = ''
  endDate.value = ''
  loading.value = false
}
</script>

<template>
  <PModal :open="props.open" @update:open="emit('update:open', $event)" title="Create cycle" size="md">
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Name</label>
        <PInput v-model="name" placeholder="Cycle name (e.g. Sprint 1)" autofocus />
      </div>

      <div>
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Description</label>
        <PTextarea v-model="description" placeholder="What's the goal of this cycle?" :rows="3" />
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
          <label class="mb-1.5 block text-sm font-medium text-custom-text-200">End date</label>
          <input
            v-model="endDate"
            type="date"
            class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 focus:border-custom-primary-100 focus:outline-none focus:ring-1 focus:ring-custom-primary-100"
          />
        </div>
      </div>
    </form>

    <template #footer>
      <PButton variant="secondary" @click="emit('update:open', false)">Cancel</PButton>
      <PButton variant="primary" :loading="loading" @click="handleSubmit">Create cycle</PButton>
    </template>
  </PModal>
</template>
