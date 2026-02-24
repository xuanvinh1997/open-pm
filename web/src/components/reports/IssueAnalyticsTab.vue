<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { reportApi } from '@/api/report.api'
import DoughnutChart from '@/components/ui/charts/DoughnutChart.vue'
import BarChart from '@/components/ui/charts/BarChart.vue'
import ReportExportButton from './ReportExportButton.vue'
import type { IssueAnalyticsReport } from '@/types/report.types'
import { exportToCSV } from '@/utils/csv-export'
import PSkeleton from '@/components/ui/PSkeleton.vue'

interface Props {
  slug: string
  projectId: string
}

const props = defineProps<Props>()
const loading = ref(true)
const report = ref<IssueAnalyticsReport | null>(null)

onMounted(async () => {
  try {
    const { data } = await reportApi.getIssueAnalytics(props.slug, props.projectId)
    report.value = data
  } finally {
    loading.value = false
  }
})

const stateLabels = computed(() => report.value?.by_state.map((s) => s.state_name) ?? [])
const stateData = computed(() => report.value?.by_state.map((s) => s.count) ?? [])
const stateColors = computed(() => report.value?.by_state.map((s) => s.color) ?? [])

const priorityLabels = computed(() => report.value?.by_priority.map((p) => p.priority) ?? [])
const priorityData = computed(() => report.value?.by_priority.map((p) => p.count) ?? [])
const priorityColors = computed(() =>
  report.value?.by_priority.map((p) => {
    const colorMap: Record<string, string> = {
      urgent: '#ef4444',
      high: '#f97316',
      medium: '#eab308',
      low: '#22c55e',
      none: '#94a3b8',
    }
    return colorMap[p.priority] || '#94a3b8'
  }) ?? [],
)

const typeLabels = computed(() => report.value?.by_type.map((t) => t.issue_type) ?? [])
const typeData = computed(() => report.value?.by_type.map((t) => t.count) ?? [])
const typeColors = computed(() => ['#3b82f6', '#ef4444', '#8b5cf6', '#f59e0b'])

const assigneeLabels = computed(() => report.value?.by_assignee.map((a) => a.display_name || 'Unnamed') ?? [])
const assigneeData = computed(() => report.value?.by_assignee.map((a) => a.count) ?? [])
const assigneeColors = computed((): string[] => {
  const palette = ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#22c55e', '#06b6d4']
  return report.value?.by_assignee.map((_, i) => palette[i % palette.length]!) ?? []
})

function handleExportCSV() {
  if (!report.value) return
  const headers = ['Category', 'Item', 'Count']
  const rows: (string | number)[][] = []
  report.value.by_state.forEach((s) => rows.push(['State', s.state_name, s.count]))
  report.value.by_priority.forEach((p) => rows.push(['Priority', p.priority, p.count]))
  report.value.by_type.forEach((t) => rows.push(['Type', t.issue_type, t.count]))
  report.value.by_assignee.forEach((a) => rows.push(['Assignee', a.display_name, a.count]))
  exportToCSV('issue-analytics.csv', headers, rows)
}
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h2 class="text-lg font-semibold text-custom-text-100">Issue Analytics</h2>
      <ReportExportButton
        report-type="issue_analytics"
        :slug="slug"
        :project-id="projectId"
        @export-csv="handleExportCSV"
      />
    </div>

    <div v-if="loading" class="space-y-4">
      <div class="grid grid-cols-3 gap-4">
        <PSkeleton v-for="i in 3" :key="i" class="h-24 rounded-xl" />
      </div>
      <div class="grid grid-cols-2 gap-4">
        <PSkeleton v-for="i in 4" :key="i" class="h-[310px] rounded-xl" />
      </div>
    </div>

    <template v-else-if="report">
      <!-- Stat cards -->
      <div class="mb-6 grid grid-cols-1 gap-4 sm:grid-cols-3">
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
          <p class="text-sm text-custom-text-300">Total Issues</p>
          <p class="mt-1 text-2xl font-bold text-custom-text-100">{{ report.total }}</p>
        </div>
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
          <p class="text-sm text-custom-text-300">Completed</p>
          <p class="mt-1 text-2xl font-bold text-green-500">{{ report.completed }}</p>
        </div>
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
          <p class="text-sm text-custom-text-300">Overdue</p>
          <p class="mt-1 text-2xl font-bold text-red-500">{{ report.overdue }}</p>
        </div>
      </div>

      <!-- Charts -->
      <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <DoughnutChart title="Issues by State" :labels="stateLabels" :data="stateData" :colors="stateColors" />
        <BarChart title="Issues by Priority" :labels="priorityLabels" :data="priorityData" :colors="priorityColors" />
        <BarChart title="Issues by Type" :labels="typeLabels" :data="typeData" :colors="typeColors" />
        <BarChart title="Issues by Assignee" :labels="assigneeLabels" :data="assigneeData" :colors="assigneeColors" />
      </div>
    </template>
  </div>
</template>
