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
import { PRIORITY_CONFIG } from '@/utils/issue-helpers'
import type { IssuePriority } from '@/types/issue.types'
import { Zap, Plus, Calendar, ArrowRight } from 'lucide-vue-next'

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

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    await Promise.all([
      projectStore.fetchStates(slug, projectId),
      epicStore.fetchEpics(slug, projectId),
    ])
  } finally {
    loading.value = false
  }
})

async function handleCreateEpic(data: { name: string; description_html?: string; start_date?: string; target_date?: string }) {
  try {
    await epicStore.createEpic(slug, projectId, data)
    showCreateModal.value = false
    toast.success('Epic created')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to create epic'))
  }
}

function handleEpicClick(epicId: string) {
  router.push(`/${slug}/projects/${projectId}/issues/${epicId}`)
}

function findState(stateId: string | null | undefined) {
  if (!stateId) return undefined
  return projectStore.states.find((s) => s.id === stateId)
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
        description="Create your first epic to organize large features and initiatives."
        :icon="Zap"
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
          v-for="epic in epicStore.epics"
          :key="epic.id"
          class="flex items-center gap-4 rounded-lg border border-custom-border-200 bg-custom-background-100 p-4 hover:bg-custom-background-90 transition-colors cursor-pointer"
          @click="handleEpicClick(epic.id)"
        >
          <Zap class="h-5 w-5 flex-shrink-0" style="color: #F97316" />
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <span class="text-xs text-custom-text-300">{{ projectStore.currentProject?.identifier }}-{{ epic.sequence_id }}</span>
              <h3 class="text-sm font-medium text-custom-text-100 truncate">{{ epic.name }}</h3>
            </div>
            <p v-if="epic.description_stripped" class="mt-0.5 text-xs text-custom-text-300 truncate">{{ epic.description_stripped }}</p>
          </div>
          <div v-if="epic.start_date && epic.target_date" class="flex items-center gap-1.5 text-xs text-custom-text-300 flex-shrink-0">
            <Calendar class="h-3.5 w-3.5" />
            <span>{{ formatDate(epic.start_date) }}</span>
            <ArrowRight class="h-3 w-3" />
            <span>{{ formatDate(epic.target_date) }}</span>
          </div>
          <div v-if="findState(epic.state_id)" class="flex items-center gap-1.5 text-xs text-custom-text-300 flex-shrink-0">
            <span
              class="h-2.5 w-2.5 rounded-full flex-shrink-0"
              :style="{ backgroundColor: findState(epic.state_id)?.color }"
            />
            <span>{{ findState(epic.state_id)?.name }}</span>
          </div>
          <component
            :is="PRIORITY_CONFIG[epic.priority as IssuePriority]?.icon"
            v-if="epic.priority && epic.priority !== 'none'"
            class="h-4 w-4 flex-shrink-0"
            :style="{ color: PRIORITY_CONFIG[epic.priority as IssuePriority]?.color }"
          />
        </div>
      </div>
    </div>

    <CreateEpicModal v-model:open="showCreateModal" @create="handleCreateEpic" />
  </div>
</template>
