<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { issueApi } from '@/api/issue.api'
import type { IssueFilterParams } from '@/api/issue.api'
import type { Issue, CreateIssueRequest } from '@/types/issue.types'
import type { State } from '@/types/project.types'
import { calculateSortOrder } from '@/utils/sort-order'
import draggable from 'vuedraggable'
import ViewHeader from '@/components/issues/ViewHeader.vue'
import IssueBoardCard from '@/components/issues/IssueBoardCard.vue'
import FilterBar from '@/components/issues/FilterBar.vue'
import CreateIssueModal from '@/components/issues/CreateIssueModal.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import { useToast } from '@/composables/useToast'
import { extractErrorMessage } from '@/utils/api-error'
import { Circle } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string

const toast = useToast()

const issues = ref<Issue[]>([])
const loading = ref(false)
const showCreateModal = ref(false)
const currentFilters = ref<IssueFilterParams>({})

interface ColumnData {
  state: State
  issues: Issue[]
}

const columnData = ref<ColumnData[]>([])

function buildColumns() {
  columnData.value = projectStore.states.map((state) => ({
    state,
    issues: issues.value
      .filter((i) => i.state_id === state.id)
      .sort((a, b) => a.sort_order - b.sort_order),
  }))
}

const defaultStateId = computed(() => {
  return projectStore.states.find((s) => s.is_default)?.id || ''
})

async function fetchIssues(filters?: IssueFilterParams) {
  const { data } = await issueApi.list(slug, projectId, 1, 200, filters)
  issues.value = data.results
  buildColumns()
}

async function handleFilterChange(filters: Record<string, string>) {
  currentFilters.value = filters
  loading.value = true
  try {
    await fetchIssues(filters)
  } finally {
    loading.value = false
  }
}

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
    const q = route.query
    const initialFilters: IssueFilterParams = {}
    if (q.search) initialFilters.search = q.search as string
    if (q.priority) initialFilters.priority = q.priority as string
    if (q.type) initialFilters.type = q.type as string
    if (q.state) initialFilters.state = q.state as string
    if (q.assignee) initialFilters.assignee = q.assignee as string
    if (q.label) initialFilters.label = q.label as string
    if (q.sort_by) initialFilters.sort_by = q.sort_by as string
    if (q.sort_order) initialFilters.sort_order = q.sort_order as string
    currentFilters.value = initialFilters
    await fetchIssues(initialFilters)
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
    buildColumns()
    showCreateModal.value = false
    toast.success('Issue created')
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to create issue'))
  }
}

function handleIssueClick(issue: Issue) {
  router.push(`/${slug}/projects/${projectId}/issues/${issue.id}`)
}

async function handleDragChange(
  column: ColumnData,
  event: { added?: { element: Issue; newIndex: number }; moved?: { element: Issue; newIndex: number } },
) {
  const added = event.added
  const moved = event.moved

  if (!added && !moved) return

  const issue = (added || moved)!.element
  const newIndex = (added || moved)!.newIndex
  const columnIssues = column.issues

  const prev = newIndex > 0 ? columnIssues[newIndex - 1].sort_order : null
  const next = newIndex < columnIssues.length - 1 ? columnIssues[newIndex + 1].sort_order : null
  const newSortOrder = calculateSortOrder(prev, next)

  const updates: Partial<Issue> = { sort_order: newSortOrder }
  if (added) {
    updates.state_id = column.state.id
  }

  // Optimistic local update
  issue.sort_order = newSortOrder
  if (added) {
    issue.state_id = column.state.id
  }

  // Update canonical issues array
  const idx = issues.value.findIndex((i) => i.id === issue.id)
  if (idx !== -1) {
    Object.assign(issues.value[idx], updates)
  }

  try {
    await issueApi.update(slug, projectId, issue.id, updates)
  } catch (e) {
    toast.error(extractErrorMessage(e, 'Failed to update issue position'))
    buildColumns()
  }
}
</script>

<template>
  <div class="flex h-full flex-col">
    <ViewHeader active-view="board" @create="showCreateModal = true" />

    <!-- Filters -->
    <div class="border-b border-custom-border-200 px-4 py-2">
      <FilterBar
        :states="projectStore.states"
        :labels="projectStore.labels"
        :members="projectStore.members"
        @filter-change="handleFilterChange"
      />
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <PSpinner size="lg" />
    </div>

    <!-- Board -->
    <div v-else class="flex flex-1 gap-4 overflow-x-auto p-4">
      <div
        v-for="column in columnData"
        :key="column.state.id"
        class="w-[300px] flex-shrink-0"
      >
        <!-- Column header -->
        <div class="mb-3 flex items-center gap-2 px-2">
          <Circle class="h-3 w-3 flex-shrink-0" :style="{ color: column.state.color }" fill="currentColor" />
          <span class="text-sm font-medium text-custom-text-100">{{ column.state.name }}</span>
          <span class="ml-auto rounded-full bg-custom-background-80 px-2 py-0.5 text-2xs text-custom-text-300">
            {{ column.issues.length }}
          </span>
        </div>

        <!-- Draggable cards -->
        <draggable
          v-model="column.issues"
          group="board-columns"
          item-key="id"
          :animation="200"
          ghost-class="board-card-ghost"
          drag-class="board-card-drag"
          class="space-y-2 rounded-xl bg-custom-background-80 p-2 min-h-[100px]"
          @change="(e: any) => handleDragChange(column, e)"
        >
          <template #item="{ element }">
            <IssueBoardCard
              :issue="element"
              :identifier="projectStore.currentProject?.identifier || ''"
              @click="handleIssueClick"
            />
          </template>
        </draggable>
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
