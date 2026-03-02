<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { useSprintStore } from '@/stores/sprint.store'
import IssueListItem from '@/components/issues/IssueListItem.vue'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import PEmptyState from '@/components/ui/PEmptyState.vue'
import { useToast } from '@/composables/useToast'
import { extractErrorMessage } from '@/utils/api-error'
import { formatDate } from '@/utils/helpers'
import PButton from '@/components/ui/PButton.vue'
import CompleteSprintModal from '@/components/sprints/CompleteSprintModal.vue'
import { Repeat, Calendar, ArrowRight, LayoutList, CheckCircle, Play } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const sprintStore = useSprintStore()
const toast = useToast()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string
const sprintId = route.params.sprintId as string

const loading = ref(false)
const showCompleteModal = ref(false)

const breadcrumbs = computed(() => [
  { label: projectStore.currentProject?.name || 'Project', to: `/${slug}/projects` },
  { label: 'Sprints', to: `/${slug}/projects/${projectId}/sprints` },
  { label: sprintStore.currentSprint?.name || 'Sprint' },
])

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    await Promise.all([
      projectStore.fetchStates(slug, projectId),
      sprintStore.fetchSprint(slug, projectId, sprintId),
      sprintStore.fetchSprints(slug, projectId),
    ])
  } finally {
    loading.value = false
  }
})

function findState(stateId: string | null | undefined) {
  if (!stateId) return undefined
  return projectStore.states.find((s) => s.id === stateId)
}

function handleIssueClick(issue: { id: string }) {
  router.push(`/${slug}/projects/${projectId}/issues/${issue.id}`)
}

async function handleRemoveIssue(issueId: string) {
  try {
    await sprintStore.removeIssueFromSprint(slug, projectId, sprintId, issueId)
    toast.success('Issue removed from sprint')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to remove issue'))
  }
}

async function handleStartSprint() {
  try {
    await sprintStore.updateSprint(slug, projectId, sprintId, { status: 'active' })
    toast.success('Sprint started')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to start sprint'))
  }
}

function handleCompleteSprint() {
  showCompleteModal.value = true
}

async function handleConfirmComplete(moveToSprintId?: string) {
  try {
    const req = moveToSprintId ? { move_to_sprint_id: moveToSprintId } : undefined
    await sprintStore.completeSprint(slug, projectId, sprintId, req)
    showCompleteModal.value = false
    toast.success('Sprint completed')
    // Re-fetch to get updated data
    await sprintStore.fetchSprint(slug, projectId, sprintId)
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to complete sprint'))
  }
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
    <div class="border-b border-custom-border-200 px-4 py-3">
      <PBreadcrumb :items="breadcrumbs" />
      <div v-if="sprintStore.currentSprint" class="mt-2 flex items-center gap-4">
        <Repeat class="h-5 w-5 text-custom-text-300" />
        <h1 class="text-lg font-semibold text-custom-text-100">{{ sprintStore.currentSprint.name }}</h1>
        <span
          class="rounded-full px-2 py-0.5 text-2xs font-medium"
          :class="getStatusColor(sprintStore.currentSprint.status)"
        >
          {{ getStatusLabel(sprintStore.currentSprint.status) }}
        </span>
        <div v-if="sprintStore.currentSprint.start_date && sprintStore.currentSprint.end_date" class="flex items-center gap-1.5 text-xs text-custom-text-300">
          <Calendar class="h-3.5 w-3.5" />
          <span>{{ formatDate(sprintStore.currentSprint.start_date) }}</span>
          <ArrowRight class="h-3 w-3" />
          <span>{{ formatDate(sprintStore.currentSprint.end_date) }}</span>
        </div>
        <span class="rounded-full bg-custom-background-80 px-2 py-0.5 text-2xs text-custom-text-300">
          {{ sprintStore.totalIssues }} issue{{ sprintStore.totalIssues !== 1 ? 's' : '' }}
        </span>
        <div class="ml-auto flex items-center gap-2">
          <PButton
            v-if="sprintStore.currentSprint.status === 'planned'"
            variant="primary"
            size="sm"
            @click="handleStartSprint"
          >
            <Play class="h-3.5 w-3.5" />
            Start Sprint
          </PButton>
          <PButton
            v-if="sprintStore.currentSprint.status === 'active'"
            variant="primary"
            size="sm"
            @click="handleCompleteSprint"
          >
            <CheckCircle class="h-3.5 w-3.5" />
            Complete Sprint
          </PButton>
        </div>
      </div>
      <p v-if="sprintStore.currentSprint?.description" class="mt-1 ml-9 text-sm text-custom-text-300">
        {{ sprintStore.currentSprint.description }}
      </p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <PSpinner size="lg" />
    </div>

    <!-- Empty state -->
    <div v-else-if="sprintStore.sprintIssues.length === 0" class="flex-1">
      <PEmptyState
        title="No issues in this sprint"
        description="Add issues to this sprint from the issue detail page."
        :icon="LayoutList"
      />
    </div>

    <!-- Issues list -->
    <div v-else class="flex-1 overflow-y-auto">
      <IssueListItem
        v-for="issue in sprintStore.sprintIssues"
        :key="issue.id"
        :issue="issue"
        :identifier="projectStore.currentProject?.identifier || ''"
        :state="findState(issue.state_id)"
        @click="handleIssueClick"
      />
    </div>

    <!-- Complete Sprint Modal -->
    <CompleteSprintModal
      v-if="sprintStore.currentSprint"
      v-model:open="showCompleteModal"
      :sprint="sprintStore.currentSprint"
      :issues="sprintStore.sprintIssues"
      :other-sprints="sprintStore.sprints"
      @complete="handleConfirmComplete"
    />
  </div>
</template>
