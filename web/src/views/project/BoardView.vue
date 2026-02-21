<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { issueApi } from '@/api/issue.api'
import type { Issue, CreateIssueRequest } from '@/types/issue.types'
import ViewHeader from '@/components/issues/ViewHeader.vue'
import IssueBoardCard from '@/components/issues/IssueBoardCard.vue'
import CreateIssueModal from '@/components/issues/CreateIssueModal.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import { Circle } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string

const issues = ref<Issue[]>([])
const loading = ref(false)
const showCreateModal = ref(false)

const columns = computed(() => {
  return projectStore.states.map((state) => ({
    state,
    issues: issues.value.filter((i) => i.state_id === state.id),
  }))
})

const defaultStateId = computed(() => {
  return projectStore.states.find((s) => s.is_default)?.id || ''
})

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    await projectStore.fetchStates(slug, projectId)
    const { data } = await issueApi.list(slug, projectId, 1, 200)
    issues.value = data.results
  } finally {
    loading.value = false
  }
})

async function handleCreateIssue(data: { name: string; state_id?: string; priority: string; description_html?: string }) {
  try {
    const req: CreateIssueRequest = {
      name: data.name,
      state_id: data.state_id,
      priority: data.priority as any,
      description_html: data.description_html,
    }
    const { data: newIssue } = await issueApi.create(slug, projectId, req)
    issues.value.unshift(newIssue)
    showCreateModal.value = false
  } catch (e) {
    console.error('Failed to create issue', e)
  }
}

function handleIssueClick(issue: Issue) {
  router.push(`/${slug}/projects/${projectId}/issues/${issue.id}`)
}
</script>

<template>
  <div class="flex h-full flex-col">
    <ViewHeader active-view="board" @create="showCreateModal = true" />

    <!-- Loading -->
    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <PSpinner size="lg" />
    </div>

    <!-- Board -->
    <div v-else class="flex flex-1 gap-4 overflow-x-auto p-4">
      <div
        v-for="column in columns"
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

        <!-- Cards -->
        <div class="space-y-2 rounded-xl bg-custom-background-80 p-2 min-h-[100px]">
          <IssueBoardCard
            v-for="issue in column.issues"
            :key="issue.id"
            :issue="issue"
            :identifier="projectStore.currentProject?.identifier || ''"
            @click="handleIssueClick"
          />

          <div
            v-if="column.issues.length === 0"
            class="rounded-lg border border-dashed border-custom-border-200 p-4 text-center text-xs text-custom-text-300"
          >
            No issues
          </div>
        </div>
      </div>
    </div>

    <!-- Create issue modal -->
    <CreateIssueModal
      v-model:open="showCreateModal"
      :states="projectStore.states"
      :labels="projectStore.labels"
      :default-state-id="defaultStateId"
      @create="handleCreateIssue"
    />
  </div>
</template>
