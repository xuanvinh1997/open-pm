<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { useEpicStore } from '@/stores/epic.store'
import CreateEpicModal from '@/components/epics/CreateEpicModal.vue'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PEmptyState from '@/components/ui/PEmptyState.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import PButton from '@/components/ui/PButton.vue'
import { useToast } from '@/composables/useToast'
import { extractErrorMessage } from '@/utils/api-error'
import { formatDate } from '@/utils/helpers'
import { Layers, Plus, Calendar, ArrowRight } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const epicStore = useEpicStore()
const toast = useToast()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string

const loading = ref(false)
const showCreateModal = ref(false)

const breadcrumbs = computed(() => [
  { label: projectStore.currentProject?.name || 'Project', to: `/${slug}/projects` },
  { label: 'Epics' },
])

const MODULE_STATUS_CONFIG: Record<string, { label: string; color: string }> = {
  'backlog': { label: 'Backlog', color: 'bg-gray-100 text-gray-800 dark:bg-gray-800/50 dark:text-gray-400' },
  'planned': { label: 'Planned', color: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400' },
  'in-progress': { label: 'In Progress', color: 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400' },
  'paused': { label: 'Paused', color: 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400' },
  'completed': { label: 'Completed', color: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400' },
  'cancelled': { label: 'Cancelled', color: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400' },
}

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    await epicStore.fetchEpics(slug, projectId)
  } finally {
    loading.value = false
  }
})

async function handleCreateEpic(data: { name: string; description?: string; start_date?: string; target_date?: string; status?: string }) {
  try {
    await epicStore.createEpic(slug, projectId, data as any)
    showCreateModal.value = false
    toast.success('Epic created')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to create epic'))
  }
}

function handleEpicClick(epicId: string) {
  router.push(`/${slug}/projects/${projectId}/epics/${epicId}`)
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
        New epic
      </PButton>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <PSpinner size="lg" />
    </div>

    <!-- Empty state -->
    <div v-else-if="epicStore.epics.length === 0" class="flex-1">
      <PEmptyState
        title="No epics"
        description="Create your first epic to organize features."
        :icon="Layers"
      >
        <PButton variant="primary" @click="showCreateModal = true">
          <Plus class="h-4 w-4" />
          New epic
        </PButton>
      </PEmptyState>
    </div>

    <!-- Epics list -->
    <div v-else class="flex-1 overflow-y-auto p-4">
      <div class="space-y-2">
        <div
          v-for="mod in epicStore.epics"
          :key="mod.id"
          class="flex items-center gap-4 rounded-lg border border-custom-border-200 bg-custom-background-100 p-4 hover:bg-custom-background-90 transition-colors cursor-pointer"
          @click="handleEpicClick(mod.id)"
        >
          <Layers class="h-5 w-5 text-custom-text-300 flex-shrink-0" />
          <div class="flex-1 min-w-0">
            <h3 class="text-sm font-medium text-custom-text-100 truncate">{{ mod.name }}</h3>
            <p v-if="mod.description" class="mt-0.5 text-xs text-custom-text-300 truncate">{{ mod.description }}</p>
          </div>
          <div v-if="mod.start_date && mod.target_date" class="flex items-center gap-1.5 text-xs text-custom-text-300 flex-shrink-0">
            <Calendar class="h-3.5 w-3.5" />
            <span>{{ formatDate(mod.start_date) }}</span>
            <ArrowRight class="h-3 w-3" />
            <span>{{ formatDate(mod.target_date) }}</span>
          </div>
          <span
            class="rounded-full px-2 py-0.5 text-2xs font-medium flex-shrink-0"
            :class="MODULE_STATUS_CONFIG[mod.status]?.color || MODULE_STATUS_CONFIG['backlog'].color"
          >
            {{ MODULE_STATUS_CONFIG[mod.status]?.label || mod.status }}
          </span>
        </div>
      </div>
    </div>

    <CreateEpicModal v-model:open="showCreateModal" @create="handleCreateEpic" />
  </div>
</template>
