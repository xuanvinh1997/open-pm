<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { reportApi } from '@/api/report.api'
import BarChart from '@/components/ui/charts/BarChart.vue'
import LineChart from '@/components/ui/charts/LineChart.vue'
import ReportExportButton from './ReportExportButton.vue'
import type { WorkLogReport } from '@/types/report.types'
import { exportToCSV, formatMinutes } from '@/utils/csv-export'
import PSkeleton from '@/components/ui/PSkeleton.vue'
import dayjs from 'dayjs'

interface Props {
  slug: string
  projectId: string
}

const props = defineProps<Props>()
const loading = ref(true)
const report = ref<WorkLogReport | null>(null)

const startDate = ref(dayjs().subtract(30, 'day').format('YYYY-MM-DD'))
const endDate = ref(dayjs().format('YYYY-MM-DD'))

async function fetchReport() {
  loading.value = true
  try {
    const { data } = await reportApi.getWorkLogs(props.slug, props.projectId, startDate.value, endDate.value)
    report.value = data
  } finally {
    loading.value = false
  }
}

onMounted(fetchReport)

function handleExportCSV() {
  if (!report.value) return
  const headers = ['Member', 'Total Minutes', 'Total Time', 'Log Count']
  const rows = report.value.by_member.map((m) => [
    m.display_name,
    m.total_minutes,
    formatMinutes(m.total_minutes),
    m.log_count,
  ])
  exportToCSV('work-logs.csv', headers, rows)
}
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h2 class="text-lg font-semibold text-custom-text-100">Work Log Summary</h2>
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-2">
          <input
            v-model="startDate"
            type="date"
            class="rounded-md border border-custom-border-200 bg-custom-background-100 px-2.5 py-1.5 text-sm text-custom-text-200"
            @change="fetchReport"
          />
          <span class="text-sm text-custom-text-300">to</span>
          <input
            v-model="endDate"
            type="date"
            class="rounded-md border border-custom-border-200 bg-custom-background-100 px-2.5 py-1.5 text-sm text-custom-text-200"
            @change="fetchReport"
          />
        </div>
        <ReportExportButton
          report-type="work_logs"
          :slug="slug"
          :project-id="projectId"
          @export-csv="handleExportCSV"
        />
      </div>
    </div>

    <div v-if="loading" class="space-y-4">
      <PSkeleton class="h-24 rounded-xl" />
      <div class="grid grid-cols-2 gap-4">
        <PSkeleton v-for="i in 2" :key="i" class="h-[310px] rounded-xl" />
      </div>
    </div>

    <template v-else-if="report">
      <!-- Summary card -->
      <div class="mb-6 rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
        <p class="text-sm text-custom-text-300">Total Time Logged</p>
        <p class="mt-1 text-2xl font-bold text-custom-text-100">{{ formatMinutes(report.total_minutes) }}</p>
        <p class="mt-0.5 text-xs text-custom-text-300">
          {{ report.start_date }} to {{ report.end_date }}
        </p>
      </div>

      <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <!-- By member chart -->
        <BarChart
          title="Time by Member"
          :labels="report.by_member.map((m) => m.display_name || 'Unnamed')"
          :data="report.by_member.map((m) => Math.round(m.total_minutes / 60 * 10) / 10)"
          :colors="report.by_member.map((_, i) => ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#22c55e', '#06b6d4'][i % 6]!)"
        />

        <!-- By date chart -->
        <LineChart
          title="Daily Work Log (hours)"
          :labels="report.by_date.map((d) => d.date.slice(5))"
          :data="report.by_date.map((d) => Math.round(d.total_minutes / 60 * 10) / 10)"
          color="#3b82f6"
        />
      </div>

      <!-- By issue table -->
      <div class="mt-4 rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
        <h3 class="mb-4 text-sm font-semibold text-custom-text-100">Time by Issue</h3>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-custom-border-200">
                <th class="pb-2 text-left font-medium text-custom-text-300">Issue</th>
                <th class="pb-2 text-right font-medium text-custom-text-300">Time</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="issue in report.by_issue" :key="issue.issue_id" class="border-b border-custom-border-100">
                <td class="py-2 text-custom-text-200">
                  <span class="font-mono text-xs text-custom-text-300">#{{ issue.sequence_id }}</span>
                  {{ issue.issue_name }}
                </td>
                <td class="py-2 text-right text-custom-text-200">{{ formatMinutes(issue.total_minutes) }}</td>
              </tr>
              <tr v-if="report.by_issue.length === 0">
                <td colspan="2" class="py-4 text-center text-custom-text-300">No work logs in this period</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>
