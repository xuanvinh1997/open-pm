<script setup lang="ts">
import { ref, watch } from 'vue'
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
  create: [data: { name: string; identifier: string; description?: string }]
}>()

const name = ref('')
const identifier = ref('')
const description = ref('')
const loading = ref(false)

watch(name, (val) => {
  identifier.value = val
    .toUpperCase()
    .replace(/[^A-Z0-9]/g, '')
    .slice(0, 5)
})

function handleSubmit() {
  if (!name.value.trim() || !identifier.value.trim()) return
  loading.value = true
  emit('create', {
    name: name.value,
    identifier: identifier.value,
    description: description.value || undefined,
  })
  name.value = ''
  identifier.value = ''
  description.value = ''
  loading.value = false
}
</script>

<template>
  <PModal :open="props.open" @update:open="emit('update:open', $event)" title="Create project" size="md">
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <div>
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Name</label>
        <PInput v-model="name" placeholder="Project name" autofocus />
      </div>

      <div>
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Identifier</label>
        <PInput v-model="identifier" placeholder="PROJ" />
        <p class="mt-1 text-xs text-custom-text-300">Used as prefix for issues (e.g. PROJ-123)</p>
      </div>

      <div>
        <label class="mb-1.5 block text-sm font-medium text-custom-text-200">Description</label>
        <PTextarea v-model="description" placeholder="What's this project about?" :rows="2" />
      </div>
    </form>

    <template #footer>
      <PButton variant="secondary" @click="emit('update:open', false)">Cancel</PButton>
      <PButton variant="primary" :loading="loading" @click="handleSubmit">Create project</PButton>
    </template>
  </PModal>
</template>
