<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { useCycleStore } from '@/stores/cycle.store'
import CreateCycleModal from '@/components/cycles/CreateCycleModal.vue'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PEmptyState from '@/components/ui/PEmptyState.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import PButton from '@/components/ui/PButton.vue'
import { useToast } from '@/composables/useToast'
import { extractErrorMessage } from '@/utils/api-error'
import { formatDate } from '@/utils/helpers'
import { Repeat, Plus, Calendar, ArrowRight } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const cycleStore = useCycleStore()
const toast = useToast()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string

const loading = ref(false)
const showCreateModal = ref(false)

const breadcrumbs = computed(() => [
  { label: projectStore.currentProject?.name || 'Project', to: `/${slug}/projects` },
  { label: 'Cycles' },
])

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    await cycleStore.fetchCycles(slug, projectId)
  } finally {
    loading.value = false
  }
})

async function handleCreateCycle(data: { name: string; description?: string; start_date?: string; end_date?: string }) {
  try {
    await cycleStore.createCycle(slug, projectId, data)
    showCreateModal.value = false
    toast.success('Cycle created')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to create cycle'))
  }
}

function handleCycleClick(cycleId: string) {
  router.push(`/${slug}/projects/${projectId}/cycles/${cycleId}`)
}

function getCycleStatus(cycle: typeof cycleStore.cycles[0]) {
  if (!cycle.start_date || !cycle.end_date) return 'Draft'
  const now = new Date()
  const start = new Date(cycle.start_date)
  const end = new Date(cycle.end_date)
  if (now < start) return 'Upcoming'
  if (now > end) return 'Completed'
  return 'Active'
}

function getStatusColor(status: string) {
  switch (status) {
    case 'Active': return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
    case 'Upcoming': return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
    case 'Completed': return 'bg-gray-100 text-gray-800 dark:bg-gray-800/50 dark:text-gray-400'
    default: return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400'
  }
}
</script>

<template>
  <div class="flex h-full flex-col">
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-custom-border-200 px-4 py-3">
      <div class="flex items-center gap-3">
        <PBreadcrumb :items="breadcrumbs" />
      </div>
      <PButton variant="primary" size="sm" @click="showCreateModal = true">
        <Plus class="h-4 w-4" />
        New cycle
      </PButton>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <PSpinner size="lg" />
    </div>

    <!-- Empty state -->
    <div v-else-if="cycleStore.cycles.length === 0" class="flex-1">
      <PEmptyState
        title="No cycles"
        description="Create your first cycle to start planning sprints."
        :icon="Repeat"
      >
        <PButton variant="primary" @click="showCreateModal = true">
          <Plus class="h-4 w-4" />
          New cycle
        </PButton>
      </PEmptyState>
    </div>

    <!-- Cycles list -->
    <div v-else class="flex-1 overflow-y-auto p-4">
      <div class="space-y-2">
        <div
          v-for="cycle in cycleStore.cycles"
          :key="cycle.id"
          class="flex items-center gap-4 rounded-lg border border-custom-border-200 bg-custom-background-100 p-4 hover:bg-custom-background-90 transition-colors cursor-pointer"
          @click="handleCycleClick(cycle.id)"
        >
          <Repeat class="h-5 w-5 text-custom-text-300 flex-shrink-0" />
          <div class="flex-1 min-w-0">
            <h3 class="text-sm font-medium text-custom-text-100 truncate">{{ cycle.name }}</h3>
            <p v-if="cycle.description" class="mt-0.5 text-xs text-custom-text-300 truncate">{{ cycle.description }}</p>
          </div>
          <div v-if="cycle.start_date && cycle.end_date" class="flex items-center gap-1.5 text-xs text-custom-text-300 flex-shrink-0">
            <Calendar class="h-3.5 w-3.5" />
            <span>{{ formatDate(cycle.start_date) }}</span>
            <ArrowRight class="h-3 w-3" />
            <span>{{ formatDate(cycle.end_date) }}</span>
          </div>
          <span
            class="rounded-full px-2 py-0.5 text-2xs font-medium flex-shrink-0"
            :class="getStatusColor(getCycleStatus(cycle))"
          >
            {{ getCycleStatus(cycle) }}
          </span>
        </div>
      </div>
    </div>

    <CreateCycleModal v-model:open="showCreateModal" @create="handleCreateCycle" />
  </div>
</template>
