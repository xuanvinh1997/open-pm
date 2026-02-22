<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { issueApi } from '@/api/issue.api'
import type { Issue, CreateIssueRequest } from '@/types/issue.types'
import ViewHeader from '@/components/issues/ViewHeader.vue'
import IssueListItem from '@/components/issues/IssueListItem.vue'
import CreateIssueModal from '@/components/issues/CreateIssueModal.vue'
import PEmptyState from '@/components/ui/PEmptyState.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import PButton from '@/components/ui/PButton.vue'
import { useToast } from '@/composables/useToast'
import { extractErrorMessage } from '@/utils/api-error'
import { LayoutList, Plus, ChevronRight, Circle } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string

const toast = useToast()

const issues = ref<Issue[]>([])
const totalCount = ref(0)
const loading = ref(false)
const showCreateModal = ref(false)
const collapsedGroups = reactive(new Set<string>())

const groupedIssues = computed(() => {
  const groups: Record<string, { state: typeof projectStore.states[0]; issues: Issue[] }> = {}
  for (const state of projectStore.states) {
    groups[state.id] = { state, issues: [] }
  }
  groups['none'] = {
    state: { id: 'none', name: 'No State', color: '#94A3B8', group: 'backlog' } as any,
    issues: [],
  }
  for (const issue of issues.value) {
    const key = issue.state_id || 'none'
    if (groups[key]) {
      groups[key].issues.push(issue)
    } else {
      groups['none'].issues.push(issue)
    }
  }
  return Object.values(groups).filter((g) => g.issues.length > 0 || g.state.id !== 'none')
})

const defaultStateId = computed(() => {
  return projectStore.states.find((s) => s.is_default)?.id || ''
})

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    await Promise.all([
      projectStore.fetchStates(slug, projectId),
      projectStore.fetchLabels(slug, projectId),
      projectStore.fetchMembers(slug, projectId),
      projectStore.fetchEstimateSystem(slug, projectId),
    ])
    const { data } = await issueApi.list(slug, projectId)
    issues.value = data.results
    totalCount.value = data.total_count
  } finally {
    loading.value = false
  }
})

async function handleCreateIssue(data: { name: string; state_id?: string; priority: string; issue_type?: string; description_html?: string; start_date?: string; target_date?: string; estimate_point?: number; assignee_ids?: string[]; label_ids?: string[] }) {
  try {
    const req: CreateIssueRequest = {
      name: data.name,
      state_id: data.state_id,
      priority: data.priority as any,
      issue_type: data.issue_type as any,
      description_html: data.description_html,
      start_date: data.start_date,
      target_date: data.target_date,
      estimate_point: data.estimate_point,
      assignee_ids: data.assignee_ids,
      label_ids: data.label_ids,
    }
    const { data: newIssue } = await issueApi.create(slug, projectId, req)
    issues.value.unshift(newIssue)
    totalCount.value++
    showCreateModal.value = false
    toast.success('Issue created')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to create issue'))
  }
}

function handleIssueClick(issue: Issue) {
  router.push(`/${slug}/projects/${projectId}/issues/${issue.id}`)
}

function toggleGroup(stateId: string) {
  if (collapsedGroups.has(stateId)) {
    collapsedGroups.delete(stateId)
  } else {
    collapsedGroups.add(stateId)
  }
}

function findState(stateId: string | null | undefined) {
  if (!stateId) return undefined
  return projectStore.states.find((s) => s.id === stateId)
}
</script>

<template>
  <div class="flex h-full flex-col">
    <ViewHeader active-view="list" @create="showCreateModal = true" />

    <!-- Loading -->
    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <PSpinner size="lg" />
    </div>

    <!-- Empty state -->
    <div v-else-if="issues.length === 0" class="flex-1">
      <PEmptyState
        title="No issues"
        description="Create your first issue to start tracking work."
        :icon="LayoutList"
      >
        <PButton variant="primary" @click="showCreateModal = true">
          <Plus class="h-4 w-4" />
          New issue
        </PButton>
      </PEmptyState>
    </div>

    <!-- Issues grouped by state -->
    <div v-else class="flex-1 overflow-y-auto">
      <div v-for="group in groupedIssues" :key="group.state.id">
        <button
          class="flex w-full items-center gap-2 px-4 py-2 hover:bg-custom-background-80 transition-colors"
          @click="toggleGroup(group.state.id)"
        >
          <ChevronRight
            class="h-3.5 w-3.5 text-custom-text-300 transition-transform duration-100"
            :class="{ 'rotate-90': !collapsedGroups.has(group.state.id) }"
          />
          <Circle
            class="h-3 w-3 flex-shrink-0"
            :style="{ color: group.state.color }"
            fill="currentColor"
          />
          <span class="text-sm font-medium text-custom-text-100">{{ group.state.name }}</span>
          <span class="rounded-full bg-custom-background-80 px-2 py-0.5 text-2xs text-custom-text-300">
            {{ group.issues.length }}
          </span>
        </button>
        <div v-if="!collapsedGroups.has(group.state.id)">
          <IssueListItem
            v-for="issue in group.issues"
            :key="issue.id"
            :issue="issue"
            :identifier="projectStore.currentProject?.identifier || ''"
            :state="findState(issue.state_id)"
            @click="handleIssueClick"
          />
        </div>
      </div>
    </div>

    <!-- Create issue modal -->
    <CreateIssueModal
      v-model:open="showCreateModal"
      :states="projectStore.states"
      :labels="projectStore.labels"
      :members="projectStore.members"
      :default-state-id="defaultStateId"
      :estimate-system="projectStore.estimateSystem"
      @create="handleCreateIssue"
    />
  </div>
</template>
