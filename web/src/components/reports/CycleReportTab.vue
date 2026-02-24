<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { reportApi } from '@/api/report.api'
import LineChart from '@/components/ui/charts/LineChart.vue'
import BarChart from '@/components/ui/charts/BarChart.vue'
import DoughnutChart from '@/components/ui/charts/DoughnutChart.vue'
import PSkeleton from '@/components/ui/PSkeleton.vue'
import type { CycleReportData, VelocityPoint, Cycle } from '@/types/report.types'

interface Props {
  slug: string
  projectId: string
}

const props = defineProps<Props>()

const loading = ref(true)
const cycles = ref<Cycle[]>([])
const selectedCycleId = ref<string>('')
const cycleReport = ref<CycleReportData | null>(null)
const velocity = ref<VelocityPoint[]>([])

onMounted(async () => {
  try {
    // Fetch cycles
    const cyclesRes = await reportApi.listCycles(props.slug, props.projectId)
    cycles.value = cyclesRes.data.results || []

    // Fetch velocity
    const velRes = await reportApi.getVelocity(props.slug, props.projectId)
    velocity.value = velRes.data.results || []

    // Select first cycle
    if (cycles.value.length > 0 && cycles.value[0]) {
      selectedCycleId.value = cycles.value[0].id
    }
  } finally {
    loading.value = false
  }
})

watch(selectedCycleId, async (id) => {
  if (!id) return
  try {
    const { data } = await reportApi.getCycleReport(props.slug, props.projectId, id)
    cycleReport.value = data
  } catch {
    cycleReport.value = null
  }
})
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h2 class="text-lg font-semibold text-custom-text-100">Cycle Reports</h2>
      <select
        v-if="cycles.length > 0"
        v-model="selectedCycleId"
        class="rounded-md border border-custom-border-200 bg-custom-background-100 px-3 py-1.5 text-sm text-custom-text-200"
      >
        <option v-for="cycle in cycles" :key="cycle.id" :value="cycle.id">
          {{ cycle.name }}
        </option>
      </select>
    </div>

    <div v-if="loading" class="space-y-4">
      <div class="grid grid-cols-2 gap-4">
        <PSkeleton v-for="i in 2" :key="i" class="h-24 rounded-xl" />
      </div>
      <PSkeleton class="h-[310px] rounded-xl" />
    </div>

    <template v-else>
      <div v-if="cycles.length === 0" class="flex h-40 items-center justify-center text-custom-text-300">
        No cycles found for this project
      </div>

      <template v-else-if="cycleReport">
        <!-- Stat cards -->
        <div class="mb-6 grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
            <p class="text-sm text-custom-text-300">Total Issues</p>
            <p class="mt-1 text-2xl font-bold text-custom-text-100">{{ cycleReport.total_issues }}</p>
          </div>
          <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
            <p class="text-sm text-custom-text-300">Completed</p>
            <p class="mt-1 text-2xl font-bold text-green-500">{{ cycleReport.completed_issues }}</p>
          </div>
          <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-5">
            <p class="text-sm text-custom-text-300">Completion Rate</p>
            <p class="mt-1 text-2xl font-bold text-brand-500">
              {{ cycleReport.total_issues > 0 ? Math.round((cycleReport.completed_issues / cycleReport.total_issues) * 100) : 0 }}%
            </p>
          </div>
        </div>

        <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
          <!-- Burndown chart -->
          <LineChart
            v-if="cycleReport.burndown_points.length > 0"
            title="Burndown Chart"
            :labels="cycleReport.burndown_points.map((p) => p.date.slice(5))"
            :data="cycleReport.burndown_points.map((p) => p.remaining_count)"
            color="#ef4444"
          />

          <!-- Issues by state -->
          <DoughnutChart
            v-if="cycleReport.issues_by_state.length > 0"
            title="Issues by State"
            :labels="cycleReport.issues_by_state.map((s) => s.state_name)"
            :data="cycleReport.issues_by_state.map((s) => s.count)"
            :colors="cycleReport.issues_by_state.map((s) => s.color)"
          />
        </div>
      </template>

      <!-- Velocity chart -->
      <div v-if="velocity.length > 0" class="mt-6">
        <BarChart
          title="Velocity (Completed Issues per Cycle)"
          :labels="velocity.map((v) => v.cycle_name)"
          :data="velocity.map((v) => v.completed_count)"
          :colors="velocity.map(() => '#3b82f6')"
        />
      </div>
    </template>
  </div>
</template>
