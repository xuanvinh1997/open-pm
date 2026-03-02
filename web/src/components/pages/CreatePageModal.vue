<script setup lang="ts">
import { ref } from 'vue'
import type { Page } from '@/types/page.types'
import PModal from '@/components/ui/PModal.vue'
import PInput from '@/components/ui/PInput.vue'
import PButton from '@/components/ui/PButton.vue'

interface Props {
  open: boolean
  pages?: Page[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  create: [data: { name: string; parent_id?: string; color?: string }]
}>()

const name = ref('')
const parentId = ref('')
const loading = ref(false)

function handleSubmit() {
  if (!name.value.trim()) return
  loading.value = true
  emit('create', {
    name: name.value,
    parent_id: parentId.value || undefined,
  })
  name.value = ''
  parentId.value = ''
  loading.value = false
}
</script>

<template>
  <PModal :open="props.open" @update:open="emit('update:open', $event)" title="Create page" size="md">
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Title</label>
        <PInput v-model="name" placeholder="Page title" autofocus />
      </div>

      <div v-if="props.pages && props.pages.length > 0">
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Parent page</label>
        <select
          v-model="parentId"
          class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-2 text-sm text-custom-text-100 focus:border-custom-primary-100 focus:outline-none focus:ring-1 focus:ring-custom-primary-100"
        >
          <option value="">None (top level)</option>
          <option v-for="page in props.pages" :key="page.id" :value="page.id">
            {{ page.name }}
          </option>
        </select>
      </div>
    </form>

    <template #footer>
      <PButton variant="secondary" @click="emit('update:open', false)">Cancel</PButton>
      <PButton variant="primary" :loading="loading" @click="handleSubmit">Create page</PButton>
    </template>
  </PModal>
</template>
