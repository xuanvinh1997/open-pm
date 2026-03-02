<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { useSprintStore } from '@/stores/sprint.store'
import { issueApi } from '@/api/issue.api'
import type { Issue } from '@/types/issue.types'
import IssueListItem from '@/components/issues/IssueListItem.vue'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PEmptyState from '@/components/ui/PEmptyState.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import PButton from '@/components/ui/PButton.vue'
import { useToast } from '@/composables/useToast'
import { extractErrorMessage } from '@/utils/api-error'
import { ListTodo, Repeat, ChevronRight, ChevronDown } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const sprintStore = useSprintStore()
const toast = useToast()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string

const loading = ref(false)
const backlogIssues = ref<Issue[]>([])
const totalCount = ref(0)
const expandedSprints = ref<Set<string>>(new Set())

const breadcrumbs = computed(() => [
  { label: projectStore.currentProject?.name || 'Project', to: `/${slug}/projects` },
  { label: 'Backlog' },
])

const plannedSprints = computed(() =>
  sprintStore.sprints.filter((s) => s.status !== 'completed')
)

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    await Promise.all([
      projectStore.fetchStates(slug, projectId),
      sprintStore.fetchSprints(slug, projectId),
      fetchBacklog(),
    ])
  } finally {
    loading.value = false
  }
})

async function fetchBacklog() {
  const { data } = await issueApi.listBacklog(slug, projectId)
  backlogIssues.value = data.results
  totalCount.value = data.total_count
}

function findState(stateId: string | null | undefined) {
  if (!stateId) return undefined
  return projectStore.states.find((s) => s.id === stateId)
}

function handleIssueClick(issue: { id: string }) {
  router.push(`/${slug}/projects/${projectId}/issues/${issue.id}`)
}

function toggleSprint(sprintId: string) {
  if (expandedSprints.value.has(sprintId)) {
    expandedSprints.value.delete(sprintId)
  } else {
    expandedSprints.value.add(sprintId)
  }
}

async function addToSprint(issueId: string, sprintId: string) {
  try {
    await sprintStore.addIssueToSprint(slug, projectId, sprintId, issueId)
    backlogIssues.value = backlogIssues.value.filter((i) => i.id !== issueId)
    totalCount.value = Math.max(0, totalCount.value - 1)
    toast.success('Issue added to sprint')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to add issue to sprint'))
  }
}

function handleSprintClick(sprintId: string) {
  router.push(`/${slug}/projects/${projectId}/sprints/${sprintId}`)
}
</script>

<template>
  <div class="flex h-full flex-col">
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-custom-border-200 px-4 py-3">
      <PBreadcrumb :items="breadcrumbs" />
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <PSpinner size="lg" />
    </div>

    <div v-else class="flex flex-1 overflow-hidden">
      <!-- Main area -->
      <div class="flex-1 overflow-y-auto">
        <!-- Sprint containers -->
        <div v-if="plannedSprints.length > 0" class="border-b border-custom-border-200">
          <div v-for="sprint in plannedSprints" :key="sprint.id" class="border-b border-custom-border-100 last:border-b-0">
            <button
              class="flex w-full items-center gap-2 px-4 py-2.5 text-sm font-medium text-custom-text-100 hover:bg-custom-background-90 transition-colors"
              @click="toggleSprint(sprint.id)"
            >
              <component
                :is="expandedSprints.has(sprint.id) ? ChevronDown : ChevronRight"
                class="h-3.5 w-3.5 text-custom-text-300"
              />
              <Repeat class="h-4 w-4 text-custom-text-300" />
              <span>{{ sprint.name }}</span>
              <span
                class="ml-1 rounded-full px-1.5 py-0.5 text-2xs font-medium"
                :class="sprint.status === 'active' ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400' : 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400'"
              >
                {{ sprint.status === 'active' ? 'Active' : 'Planned' }}
              </span>
              <span
                class="ml-auto text-xs text-custom-text-300 cursor-pointer hover:text-brand-500"
                @click.stop="handleSprintClick(sprint.id)"
              >
                View
              </span>
            </button>
          </div>
        </div>

        <!-- Backlog section -->
        <div class="px-4 py-3">
          <div class="mb-3 flex items-center gap-2">
            <ListTodo class="h-4 w-4 text-custom-text-300" />
            <h2 class="text-sm font-semibold text-custom-text-100">Backlog</h2>
            <span class="rounded-full bg-custom-background-80 px-2 py-0.5 text-2xs text-custom-text-300">
              {{ totalCount }}
            </span>
          </div>

          <div v-if="backlogIssues.length === 0">
            <PEmptyState
              title="Backlog is empty"
              description="All issues are assigned to active sprints."
              :icon="ListTodo"
            />
          </div>

          <div v-else class="space-y-0.5">
            <div
              v-for="issue in backlogIssues"
              :key="issue.id"
              class="group flex items-center"
            >
              <div class="flex-1 min-w-0">
                <IssueListItem
                  :issue="issue"
                  :identifier="projectStore.currentProject?.identifier || ''"
                  :state="findState(issue.state_id)"
                  @click="handleIssueClick"
                />
              </div>
              <!-- Add to sprint dropdown -->
              <div v-if="plannedSprints.length > 0" class="flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity pr-2">
                <select
                  class="rounded border border-custom-border-200 bg-custom-background-100 px-1.5 py-0.5 text-2xs text-custom-text-300"
                  @change="(e) => { const target = e.target as HTMLSelectElement; if (target.value) { addToSprint(issue.id, target.value); target.value = '' } }"
                >
                  <option value="">+ Sprint</option>
                  <option v-for="s in plannedSprints" :key="s.id" :value="s.id">{{ s.name }}</option>
                </select>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
