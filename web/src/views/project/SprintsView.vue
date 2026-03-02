<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { useSprintStore } from '@/stores/sprint.store'
import CreateSprintModal from '@/components/sprints/CreateSprintModal.vue'
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
const sprintStore = useSprintStore()
const toast = useToast()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string

const loading = ref(false)
const showCreateModal = ref(false)

const breadcrumbs = computed(() => [
  { label: projectStore.currentProject?.name || 'Project', to: `/${slug}/projects` },
  { label: 'Sprints' },
])

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    await sprintStore.fetchSprints(slug, projectId)
  } finally {
    loading.value = false
  }
})

async function handleCreateSprint(data: { name: string; description?: string; start_date?: string; end_date?: string }) {
  try {
    await sprintStore.createSprint(slug, projectId, data)
    showCreateModal.value = false
    toast.success('Sprint created')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to create sprint'))
  }
}

function handleSprintClick(sprintId: string) {
  router.push(`/${slug}/projects/${projectId}/sprints/${sprintId}`)
}

function getStatusLabel(status: string) {
  switch (status) {
    case 'active': return 'Active'
    case 'completed': return 'Completed'
    default: return 'Planned'
  }
}

function getStatusColor(status: string) {
  switch (status) {
    case 'active': return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
    case 'completed': return 'bg-gray-100 text-gray-800 dark:bg-gray-800/50 dark:text-gray-400'
    default: return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'
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
        New sprint
      </PButton>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <PSpinner size="lg" />
    </div>

    <!-- Empty state -->
    <div v-else-if="sprintStore.sprints.length === 0" class="flex-1">
      <PEmptyState
        title="No sprints"
        description="Create your first sprint to start planning sprints."
        :icon="Repeat"
      >
        <PButton variant="primary" @click="showCreateModal = true">
          <Plus class="h-4 w-4" />
          New sprint
        </PButton>
      </PEmptyState>
    </div>

    <!-- Sprints list -->
    <div v-else class="flex-1 overflow-y-auto p-4">
      <div class="space-y-2">
        <div
          v-for="sprint in sprintStore.sprints"
          :key="sprint.id"
          class="flex items-center gap-4 rounded-lg border border-custom-border-200 bg-custom-background-100 p-4 hover:bg-custom-background-90 transition-colors cursor-pointer"
          @click="handleSprintClick(sprint.id)"
        >
          <Repeat class="h-5 w-5 text-custom-text-300 flex-shrink-0" />
          <div class="flex-1 min-w-0">
            <h3 class="text-sm font-medium text-custom-text-100 truncate">{{ sprint.name }}</h3>
            <p v-if="sprint.description" class="mt-0.5 text-xs text-custom-text-300 truncate">{{ sprint.description }}</p>
          </div>
          <div v-if="sprint.start_date && sprint.end_date" class="flex items-center gap-1.5 text-xs text-custom-text-300 flex-shrink-0">
            <Calendar class="h-3.5 w-3.5" />
            <span>{{ formatDate(sprint.start_date) }}</span>
            <ArrowRight class="h-3 w-3" />
            <span>{{ formatDate(sprint.end_date) }}</span>
          </div>
          <span
            class="rounded-full px-2 py-0.5 text-2xs font-medium flex-shrink-0"
            :class="getStatusColor(sprint.status)"
          >
            {{ getStatusLabel(sprint.status) }}
          </span>
        </div>
      </div>
    </div>

    <CreateSprintModal v-model:open="showCreateModal" @create="handleCreateSprint" />
  </div>
</template>
