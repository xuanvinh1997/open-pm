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
import { Repeat, Calendar, ArrowRight, LayoutList } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const sprintStore = useSprintStore()
const toast = useToast()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string
const sprintId = route.params.sprintId as string

const loading = ref(false)

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
</script>

<template>
  <div class="flex h-full flex-col">
    <!-- Header -->
    <div class="border-b border-custom-border-200 px-4 py-3">
      <PBreadcrumb :items="breadcrumbs" />
      <div v-if="sprintStore.currentSprint" class="mt-2 flex items-center gap-4">
        <Repeat class="h-5 w-5 text-custom-text-300" />
        <h1 class="text-lg font-semibold text-custom-text-100">{{ sprintStore.currentSprint.name }}</h1>
        <div v-if="sprintStore.currentSprint.start_date && sprintStore.currentSprint.end_date" class="flex items-center gap-1.5 text-xs text-custom-text-300">
          <Calendar class="h-3.5 w-3.5" />
          <span>{{ formatDate(sprintStore.currentSprint.start_date) }}</span>
          <ArrowRight class="h-3 w-3" />
          <span>{{ formatDate(sprintStore.currentSprint.end_date) }}</span>
        </div>
        <span class="rounded-full bg-custom-background-80 px-2 py-0.5 text-2xs text-custom-text-300">
          {{ sprintStore.totalIssues }} issue{{ sprintStore.totalIssues !== 1 ? 's' : '' }}
        </span>
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
  </div>
</template>
