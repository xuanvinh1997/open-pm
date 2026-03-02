<script setup lang="ts">
import { computed, ref } from 'vue'
import PModal from '@/components/ui/PModal.vue'
import PButton from '@/components/ui/PButton.vue'
import type { Issue } from '@/types/issue.types'
import type { Sprint } from '@/types/sprint.types'
import { AlertTriangle } from 'lucide-vue-next'

interface Props {
  open: boolean
  sprint: Sprint
  issues: Issue[]
  otherSprints: Sprint[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  complete: [moveToSprintId?: string]
}>()

const moveToSprintId = ref('')
const loading = ref(false)

const incompleteIssues = computed(() =>
  props.issues.filter((i) => !i.completed_at)
)

const completedIssues = computed(() =>
  props.issues.filter((i) => !!i.completed_at)
)

const availableSprints = computed(() =>
  props.otherSprints.filter((s) => s.id !== props.sprint.id && s.status !== 'completed')
)

function handleComplete() {
  loading.value = true
  emit('complete', moveToSprintId.value || undefined)
}

function handleClose() {
  moveToSprintId.value = ''
  loading.value = false
  emit('update:open', false)
}
</script>

<template>
  <PModal :open="props.open" @update:open="handleClose" title="Complete sprint" size="md">
    <div class="space-y-4">
      <!-- Summary -->
      <div class="rounded-md border border-custom-border-200 p-3">
        <p class="text-sm font-medium text-custom-text-100 mb-2">{{ props.sprint.name }}</p>
        <div class="flex gap-4 text-xs text-custom-text-300">
          <span>{{ completedIssues.length }} completed</span>
          <span>{{ incompleteIssues.length }} incomplete</span>
          <span>{{ props.issues.length }} total</span>
        </div>
      </div>

      <!-- Warning for incomplete issues -->
      <div v-if="incompleteIssues.length > 0" class="rounded-md border border-amber-200 bg-amber-50 p-3 dark:border-amber-800 dark:bg-amber-900/20">
        <div class="flex items-start gap-2">
          <AlertTriangle class="h-4 w-4 text-amber-600 dark:text-amber-400 flex-shrink-0 mt-0.5" />
          <div>
            <p class="text-sm font-medium text-amber-800 dark:text-amber-300">
              {{ incompleteIssues.length }} issue{{ incompleteIssues.length !== 1 ? 's are' : ' is' }} not completed
            </p>
            <p class="mt-1 text-xs text-amber-700 dark:text-amber-400">
              Choose where to move incomplete issues, or leave them unassigned.
            </p>
          </div>
        </div>

        <!-- Move to sprint selector -->
        <div class="mt-3">
          <label class="mb-1.5 block text-xs font-medium text-amber-800 dark:text-amber-300">
            Move incomplete issues to
          </label>
          <select
            v-model="moveToSprintId"
            class="w-full rounded-md border border-amber-300 bg-white px-2.5 py-1.5 text-sm text-custom-text-100 dark:border-amber-700 dark:bg-custom-background-100 focus:border-custom-primary-100 focus:outline-none"
          >
            <option value="">Don't move (remove from sprint)</option>
            <option
              v-for="s in availableSprints"
              :key="s.id"
              :value="s.id"
            >
              {{ s.name }}{{ s.status === 'active' ? ' (Active)' : '' }}
            </option>
          </select>
        </div>

        <!-- Incomplete issues list -->
        <div class="mt-3 max-h-32 overflow-y-auto space-y-1">
          <div
            v-for="issue in incompleteIssues"
            :key="issue.id"
            class="flex items-center gap-2 text-xs text-amber-800 dark:text-amber-400"
          >
            <span class="w-1.5 h-1.5 rounded-full bg-amber-500 flex-shrink-0"></span>
            <span class="truncate">{{ issue.name }}</span>
          </div>
        </div>
      </div>

      <!-- All done message -->
      <div v-else class="rounded-md border border-green-200 bg-green-50 p-3 dark:border-green-800 dark:bg-green-900/20">
        <p class="text-sm text-green-800 dark:text-green-300">
          All issues in this sprint are completed. You're good to go!
        </p>
      </div>
    </div>

    <template #footer>
      <PButton variant="secondary" @click="handleClose">Cancel</PButton>
      <PButton variant="primary" :loading="loading" @click="handleComplete">
        Complete Sprint
      </PButton>
    </template>
  </PModal>
</template>
