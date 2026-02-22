import { defineStore } from 'pinia'
import { ref } from 'vue'
import { cycleApi } from '@/api/cycle.api'
import type { Cycle, CreateCycleRequest, UpdateCycleRequest, CycleDetail } from '@/types/cycle.types'
import type { Issue } from '@/types/issue.types'

export const useCycleStore = defineStore('cycle', () => {
  const cycles = ref<Cycle[]>([])
  const currentCycle = ref<Cycle | null>(null)
  const cycleIssues = ref<Issue[]>([])
  const totalIssues = ref(0)
  const loading = ref(false)

  async function fetchCycles(slug: string, projectId: string) {
    loading.value = true
    try {
      const { data } = await cycleApi.list(slug, projectId)
      cycles.value = data.results
    } finally {
      loading.value = false
    }
  }

  async function fetchCycle(slug: string, projectId: string, cycleId: string) {
    const { data } = await cycleApi.get(slug, projectId, cycleId)
    currentCycle.value = data.cycle
    cycleIssues.value = data.issues
    totalIssues.value = data.total_issues
    return data
  }

  async function createCycle(slug: string, projectId: string, data: CreateCycleRequest) {
    const { data: cycle } = await cycleApi.create(slug, projectId, data)
    cycles.value.push(cycle)
    return cycle
  }

  async function updateCycle(slug: string, projectId: string, cycleId: string, data: UpdateCycleRequest) {
    const { data: updated } = await cycleApi.update(slug, projectId, cycleId, data)
    const idx = cycles.value.findIndex((c) => c.id === cycleId)
    if (idx >= 0) cycles.value[idx] = updated
    if (currentCycle.value?.id === cycleId) currentCycle.value = updated
    return updated
  }

  async function deleteCycle(slug: string, projectId: string, cycleId: string) {
    await cycleApi.delete(slug, projectId, cycleId)
    cycles.value = cycles.value.filter((c) => c.id !== cycleId)
    if (currentCycle.value?.id === cycleId) currentCycle.value = null
  }

  async function addIssueToCycle(slug: string, projectId: string, cycleId: string, issueId: string) {
    await cycleApi.addIssue(slug, projectId, cycleId, issueId)
  }

  async function removeIssueFromCycle(slug: string, projectId: string, cycleId: string, issueId: string) {
    await cycleApi.removeIssue(slug, projectId, cycleId, issueId)
    cycleIssues.value = cycleIssues.value.filter((i) => i.id !== issueId)
    totalIssues.value = Math.max(0, totalIssues.value - 1)
  }

  return {
    cycles,
    currentCycle,
    cycleIssues,
    totalIssues,
    loading,
    fetchCycles,
    fetchCycle,
    createCycle,
    updateCycle,
    deleteCycle,
    addIssueToCycle,
    removeIssueFromCycle,
  }
})
