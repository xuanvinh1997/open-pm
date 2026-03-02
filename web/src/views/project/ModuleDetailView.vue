<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { useModuleStore } from '@/stores/module.store'
import IssueListItem from '@/components/issues/IssueListItem.vue'
import PBreadcrumb from '@/components/ui/PBreadcrumb.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import PEmptyState from '@/components/ui/PEmptyState.vue'
import PBadge from '@/components/ui/PBadge.vue'
import { useToast } from '@/composables/useToast'
import { extractErrorMessage } from '@/utils/api-error'
import { formatDate } from '@/utils/helpers'
import { Layers, Calendar, ArrowRight, LayoutList } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const moduleStore = useModuleStore()
const toast = useToast()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string
const moduleId = route.params.moduleId as string

const loading = ref(false)

const breadcrumbs = computed(() => [
  { label: projectStore.currentProject?.name || 'Project', to: `/${slug}/projects` },
  { label: 'Modules', to: `/${slug}/projects/${projectId}/modules` },
  { label: moduleStore.currentModule?.name || 'Module' },
])

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    await Promise.all([
      projectStore.fetchStates(slug, projectId),
      moduleStore.fetchModule(slug, projectId, moduleId),
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
</script>

<template>
  <div class="flex h-full flex-col">
    <!-- Header -->
    <div class="border-b border-custom-border-200 px-4 py-3">
      <PBreadcrumb :items="breadcrumbs" />
      <div v-if="moduleStore.currentModule" class="mt-2 flex items-center gap-4">
        <Layers class="h-5 w-5 text-custom-text-300" />
        <h1 class="text-lg font-semibold text-custom-text-100">{{ moduleStore.currentModule.name }}</h1>
        <PBadge>{{ moduleStore.currentModule.status }}</PBadge>
        <div v-if="moduleStore.currentModule.start_date && moduleStore.currentModule.target_date" class="flex items-center gap-1.5 text-xs text-custom-text-300">
          <Calendar class="h-3.5 w-3.5" />
          <span>{{ formatDate(moduleStore.currentModule.start_date) }}</span>
          <ArrowRight class="h-3 w-3" />
          <span>{{ formatDate(moduleStore.currentModule.target_date) }}</span>
        </div>
        <span class="rounded-full bg-custom-background-80 px-2 py-0.5 text-2xs text-custom-text-300">
          {{ moduleStore.totalIssues }} issue{{ moduleStore.totalIssues !== 1 ? 's' : '' }}
        </span>
      </div>
      <p v-if="moduleStore.currentModule?.description" class="mt-1 ml-9 text-sm text-custom-text-300">
        {{ moduleStore.currentModule.description }}
      </p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <PSpinner size="lg" />
    </div>

    <!-- Empty state -->
    <div v-else-if="moduleStore.moduleIssues.length === 0" class="flex-1">
      <PEmptyState
        title="No issues in this module"
        description="Add issues to this module from the issue detail page."
        :icon="LayoutList"
      />
    </div>

    <!-- Issues list -->
    <div v-else class="flex-1 overflow-y-auto">
      <IssueListItem
        v-for="issue in moduleStore.moduleIssues"
        :key="issue.id"
        :issue="issue"
        :identifier="projectStore.currentProject?.identifier || ''"
        :state="findState(issue.state_id)"
        @click="handleIssueClick"
      />
    </div>
  </div>
</template>
