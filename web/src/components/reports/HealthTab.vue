<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { reportApi } from '@/api/report.api'
import type { ProjectHealthReport } from '@/types/report.types'
import PSkeleton from '@/components/ui/PSkeleton.vue'
import { exportToCSV } from '@/utils/csv-export'
import ReportExportButton from './ReportExportButton.vue'
import { AlertTriangle, Ban, HelpCircle } from 'lucide-vue-next'

interface Props {
  slug: string
  projectId: string
}

const props = defineProps<Props>()
const loading = ref(true)
const report = ref<ProjectHealthReport | null>(null)

onMounted(async () => {
  try {
    const { data } = await reportApi.getHealth(props.slug, props.projectId)
    report.value = data
  } finally {
    loading.value = false
  }
})

function handleExportCSV() {
  if (!report.value) return
  const headers = ['Category', 'Issue Name', 'Sequence ID', 'Priority', 'Target Date']
  const rows: (string | number)[][] = []
  report.value.overdue_issues.forEach((i) =>
    rows.push(['Overdue', i.name, i.sequence_id, i.priority, i.target_date || '']),
  )
  report.value.blocked_issues.forEach((i) =>
    rows.push(['Blocked', i.name, i.sequence_id, i.priority, i.target_date || '']),
  )
  report.value.unestimated_issues.forEach((i) =>
    rows.push(['Unestimated', i.name, i.sequence_id, i.priority, i.target_date || '']),
  )
  exportToCSV('project-health.csv', headers, rows)
}
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h2 class="text-lg font-semibold text-custom-text-100">Project Health</h2>
      <ReportExportButton
        report-type="health"
        :slug="slug"
        :project-id="projectId"
        @export-csv="handleExportCSV"
      />
    </div>

    <div v-if="loading" class="space-y-4">
      <div class="grid grid-cols-3 gap-4">
        <PSkeleton v-for="i in 3" :key="i" class="h-24 rounded-xl" />
      </div>
      <PSkeleton class="h-64 rounded-xl" />
    </div>

    <template v-else-if="report">
      <!-- Summary cards -->
      <div class="mb-6 grid grid-cols-1 gap-4 sm:grid-cols-3">
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
          <div class="flex items-center gap-2">
            <AlertTriangle class="h-4 w-4 text-red-500" />
            <p class="text-sm text-custom-text-300">Overdue</p>
          </div>
          <p class="mt-1 text-2xl font-bold" :class="report.overdue_count > 0 ? 'text-red-500' : 'text-green-500'">
            {{ report.overdue_count }}
          </p>
        </div>
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
          <div class="flex items-center gap-2">
            <Ban class="h-4 w-4 text-orange-500" />
            <p class="text-sm text-custom-text-300">Blocked</p>
          </div>
          <p class="mt-1 text-2xl font-bold" :class="report.blocked_count > 0 ? 'text-orange-500' : 'text-green-500'">
            {{ report.blocked_count }}
          </p>
        </div>
        <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
          <div class="flex items-center gap-2">
            <HelpCircle class="h-4 w-4 text-yellow-500" />
            <p class="text-sm text-custom-text-300">Unestimated</p>
          </div>
          <p class="mt-1 text-2xl font-bold" :class="report.unestimated_count > 0 ? 'text-yellow-500' : 'text-green-500'">
            {{ report.unestimated_count }}
          </p>
        </div>
      </div>

      <!-- Overdue issues table -->
      <div v-if="report.overdue_issues.length > 0" class="mb-4 rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
        <h3 class="mb-3 flex items-center gap-2 text-sm font-semibold text-custom-text-100">
          <AlertTriangle class="h-4 w-4 text-red-500" />
          Overdue Issues ({{ report.overdue_count }})
        </h3>
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-custom-border-200">
              <th class="pb-2 text-left font-medium text-custom-text-300">Issue</th>
              <th class="pb-2 text-left font-medium text-custom-text-300">Priority</th>
              <th class="pb-2 text-left font-medium text-custom-text-300">Due Date</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="issue in report.overdue_issues" :key="issue.id" class="border-b border-custom-border-100">
              <td class="py-2 text-custom-text-200">
                <span class="font-mono text-xs text-custom-text-300">#{{ issue.sequence_id }}</span>
                {{ issue.name }}
              </td>
              <td class="py-2">
                <span class="rounded px-1.5 py-0.5 text-xs font-medium capitalize" :class="{
                  'bg-red-100 text-red-700': issue.priority === 'urgent',
                  'bg-orange-100 text-orange-700': issue.priority === 'high',
                  'bg-yellow-100 text-yellow-700': issue.priority === 'medium',
                  'bg-green-100 text-green-700': issue.priority === 'low',
                  'bg-gray-100 text-gray-600': issue.priority === 'none',
                }">{{ issue.priority }}</span>
              </td>
              <td class="py-2 text-red-500">{{ issue.target_date?.slice(0, 10) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Blocked issues table -->
      <div v-if="report.blocked_issues.length > 0" class="mb-4 rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
        <h3 class="mb-3 flex items-center gap-2 text-sm font-semibold text-custom-text-100">
          <Ban class="h-4 w-4 text-orange-500" />
          Blocked Issues ({{ report.blocked_count }})
        </h3>
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-custom-border-200">
              <th class="pb-2 text-left font-medium text-custom-text-300">Issue</th>
              <th class="pb-2 text-left font-medium text-custom-text-300">Priority</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="issue in report.blocked_issues" :key="issue.id" class="border-b border-custom-border-100">
              <td class="py-2 text-custom-text-200">
                <span class="font-mono text-xs text-custom-text-300">#{{ issue.sequence_id }}</span>
                {{ issue.name }}
              </td>
              <td class="py-2">
                <span class="rounded px-1.5 py-0.5 text-xs font-medium capitalize" :class="{
                  'bg-red-100 text-red-700': issue.priority === 'urgent',
                  'bg-orange-100 text-orange-700': issue.priority === 'high',
                  'bg-yellow-100 text-yellow-700': issue.priority === 'medium',
                  'bg-green-100 text-green-700': issue.priority === 'low',
                  'bg-gray-100 text-gray-600': issue.priority === 'none',
                }">{{ issue.priority }}</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Unestimated issues table -->
      <div v-if="report.unestimated_issues.length > 0" class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
        <h3 class="mb-3 flex items-center gap-2 text-sm font-semibold text-custom-text-100">
          <HelpCircle class="h-4 w-4 text-yellow-500" />
          Unestimated Issues ({{ report.unestimated_count }})
        </h3>
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-custom-border-200">
              <th class="pb-2 text-left font-medium text-custom-text-300">Issue</th>
              <th class="pb-2 text-left font-medium text-custom-text-300">Priority</th>
              <th class="pb-2 text-left font-medium text-custom-text-300">Type</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="issue in report.unestimated_issues" :key="issue.id" class="border-b border-custom-border-100">
              <td class="py-2 text-custom-text-200">
                <span class="font-mono text-xs text-custom-text-300">#{{ issue.sequence_id }}</span>
                {{ issue.name }}
              </td>
              <td class="py-2">
                <span class="rounded px-1.5 py-0.5 text-xs font-medium capitalize" :class="{
                  'bg-red-100 text-red-700': issue.priority === 'urgent',
                  'bg-orange-100 text-orange-700': issue.priority === 'high',
                  'bg-yellow-100 text-yellow-700': issue.priority === 'medium',
                  'bg-green-100 text-green-700': issue.priority === 'low',
                  'bg-gray-100 text-gray-600': issue.priority === 'none',
                }">{{ issue.priority }}</span>
              </td>
              <td class="py-2 text-custom-text-300 capitalize">{{ issue.issue_type }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="report.overdue_count === 0 && report.blocked_count === 0 && report.unestimated_count === 0"
        class="flex h-40 items-center justify-center text-green-500 font-medium">
        All clear! No health issues found.
      </div>
    </template>
  </div>
</template>
