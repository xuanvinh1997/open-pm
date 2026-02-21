<script setup lang="ts">
import { ref, onMounted, toRef } from 'vue'
import { useRoute } from 'vue-router'
import { useProjectStore } from '@/stores/project.store'
import { issueApi } from '@/api/issue.api'
import { useAnalytics } from '@/composables/useAnalytics'
import { useTheme } from '@/composables/useTheme'
import type { Issue } from '@/types/issue.types'
import ViewHeader from '@/components/issues/ViewHeader.vue'
import DoughnutChart from '@/components/ui/charts/DoughnutChart.vue'
import BarChart from '@/components/ui/charts/BarChart.vue'
import LineChart from '@/components/ui/charts/LineChart.vue'
import PSpinner from '@/components/ui/PSpinner.vue'
import { LayoutDashboard, CheckCircle2, Loader2, AlertTriangle } from 'lucide-vue-next'

const route = useRoute()
const projectStore = useProjectStore()
const { theme } = useTheme()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string

const issues = ref<Issue[]>([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    await projectStore.setCurrentProject(slug, projectId)
    await Promise.all([
      projectStore.fetchStates(slug, projectId),
      projectStore.fetchMembers(slug, projectId),
    ])
    const { data } = await issueApi.list(slug, projectId, 1, 200)
    issues.value = data.results
  } finally {
    loading.value = false
  }
})

const { issuesByState, issuesByPriority, completionOverTime, stats } =
  useAnalytics({ issues, states: toRef(projectStore, 'states') })
</script>

<template>
  <div class="flex h-full flex-col">
    <ViewHeader active-view="analytics" />

    <div v-if="loading" class="flex flex-1 items-center justify-center">
      <PSpinner size="lg" />
    </div>

    <div v-else class="flex-1 overflow-y-auto p-6">
      <!-- Stat cards -->
      <div class="mb-6 grid grid-cols-2 gap-4 md:grid-cols-4">
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
          <div class="flex items-center gap-3">
            <div class="rounded-lg bg-brand-500/10 p-2.5">
              <LayoutDashboard class="h-5 w-5 text-brand-600" />
            </div>
            <div>
              <p class="text-2xl font-bold text-custom-text-100">{{ stats.total }}</p>
              <p class="text-sm text-custom-text-300">Total</p>
            </div>
          </div>
        </div>
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
          <div class="flex items-center gap-3">
            <div class="rounded-lg bg-amber-500/10 p-2.5">
              <Loader2 class="h-5 w-5 text-amber-600" />
            </div>
            <div>
              <p class="text-2xl font-bold text-custom-text-100">{{ stats.active }}</p>
              <p class="text-sm text-custom-text-300">Active</p>
            </div>
          </div>
        </div>
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
          <div class="flex items-center gap-3">
            <div class="rounded-lg bg-green-500/10 p-2.5">
              <CheckCircle2 class="h-5 w-5 text-green-600" />
            </div>
            <div>
              <p class="text-2xl font-bold text-custom-text-100">{{ stats.completed }}</p>
              <p class="text-sm text-custom-text-300">Completed</p>
            </div>
          </div>
        </div>
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
          <div class="flex items-center gap-3">
            <div class="rounded-lg bg-red-500/10 p-2.5">
              <AlertTriangle class="h-5 w-5 text-red-600" />
            </div>
            <div>
              <p class="text-2xl font-bold text-custom-text-100">{{ stats.overdue }}</p>
              <p class="text-sm text-custom-text-300">Overdue</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Charts -->
      <div class="grid grid-cols-1 gap-6 md:grid-cols-2" :key="theme">
        <DoughnutChart
          title="Issues by State"
          :labels="issuesByState.map(s => s.name)"
          :data="issuesByState.map(s => s.count)"
          :colors="issuesByState.map(s => s.color)"
        />
        <BarChart
          title="Issues by Priority"
          :labels="issuesByPriority.map(p => p.label)"
          :data="issuesByPriority.map(p => p.count)"
          :colors="issuesByPriority.map(p => p.color)"
        />
        <LineChart
          title="Completions (Last 30 Days)"
          :labels="completionOverTime.map(d => d.date.slice(5))"
          :data="completionOverTime.map(d => d.count)"
          class="md:col-span-2"
        />
      </div>
    </div>
  </div>
</template>
