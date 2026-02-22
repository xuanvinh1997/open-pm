<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import type { WorkLog, CreateWorkLogRequest, UpdateWorkLogRequest } from '@/types/issue.types'
import PModal from '@/components/ui/PModal.vue'
import PButton from '@/components/ui/PButton.vue'
import PInput from '@/components/ui/PInput.vue'
import PTextarea from '@/components/ui/PTextarea.vue'

interface Props {
  open: boolean
  workLog?: WorkLog | null
}

const props = withDefaults(defineProps<Props>(), {
  workLog: null,
})

const emit = defineEmits<{
  'update:open': [value: boolean]
  create: [data: CreateWorkLogRequest]
  update: [id: string, data: UpdateWorkLogRequest]
}>()

const hours = ref(0)
const minutes = ref(0)
const description = ref('')
const startDate = ref('')
const endDate = ref('')
const submitting = ref(false)

const isEdit = computed(() => !!props.workLog)

const title = computed(() => (isEdit.value ? 'Edit Work Log' : 'Log Work'))

const isValid = computed(() => {
  const totalMinutes = hours.value * 60 + minutes.value
  return totalMinutes > 0 && startDate.value !== '' && endDate.value !== '' && endDate.value >= startDate.value
})

watch(
  () => props.open,
  (val) => {
    if (val) {
      const today = new Date().toISOString().substring(0, 10)
      if (props.workLog) {
        const totalMinutes = props.workLog.duration_minutes
        hours.value = Math.floor(totalMinutes / 60)
        minutes.value = totalMinutes % 60
        description.value = props.workLog.description || ''
        startDate.value = props.workLog.start_date.substring(0, 10)
        endDate.value = props.workLog.end_date.substring(0, 10)
      } else {
        hours.value = 0
        minutes.value = 0
        description.value = ''
        startDate.value = today
        endDate.value = today
      }
    }
  },
)

// When start date changes, ensure end date is not before it
watch(startDate, (val) => {
  if (endDate.value && val > endDate.value) {
    endDate.value = val
  }
})

function handleSubmit() {
  if (!isValid.value || submitting.value) return
  submitting.value = true

  const totalMinutes = hours.value * 60 + minutes.value

  if (isEdit.value && props.workLog) {
    emit('update', props.workLog.id, {
      duration_minutes: totalMinutes,
      description: description.value,
      start_date: startDate.value,
      end_date: endDate.value,
    })
  } else {
    emit('create', {
      duration_minutes: totalMinutes,
      description: description.value || undefined,
      start_date: startDate.value,
      end_date: endDate.value,
    })
  }

  submitting.value = false
  emit('update:open', false)
}
</script>

<template>
  <PModal :open="props.open" :title="title" size="sm" @update:open="emit('update:open', $event)">
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <!-- Duration -->
      <div>
        <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Duration</label>
        <div class="flex items-center gap-2">
          <div class="flex items-center gap-1">
            <PInput
              v-model.number="hours"
              type="number"
              :min="0"
              :max="999"
              size="sm"
              class="w-20"
            />
            <span class="text-xs text-custom-text-300">h</span>
          </div>
          <div class="flex items-center gap-1">
            <PInput
              v-model.number="minutes"
              type="number"
              :min="0"
              :max="59"
              size="sm"
              class="w-20"
            />
            <span class="text-xs text-custom-text-300">m</span>
          </div>
        </div>
      </div>

      <!-- Date Range -->
      <div class="flex items-end gap-3">
        <div class="flex-1">
          <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Start date</label>
          <input
            v-model="startDate"
            type="date"
            class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-1.5 text-sm text-custom-text-200 focus:border-custom-primary-100 focus:outline-none"
          />
        </div>
        <span class="pb-2 text-xs text-custom-text-300">to</span>
        <div class="flex-1">
          <label class="mb-1.5 block text-xs font-medium text-custom-text-300">End date</label>
          <input
            v-model="endDate"
            type="date"
            :min="startDate"
            class="w-full rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-1.5 text-sm text-custom-text-200 focus:border-custom-primary-100 focus:outline-none"
          />
        </div>
      </div>

      <!-- Description -->
      <div>
        <label class="mb-1.5 block text-xs font-medium text-custom-text-300">Description (optional)</label>
        <PTextarea
          v-model="description"
          placeholder="What did you work on?"
          :rows="3"
        />
      </div>
    </form>

    <template #footer>
      <PButton variant="secondary" size="sm" @click="emit('update:open', false)">Cancel</PButton>
      <PButton size="sm" :disabled="!isValid || submitting" @click="handleSubmit">
        {{ isEdit ? 'Update' : 'Log Work' }}
      </PButton>
    </template>
  </PModal>
</template>
