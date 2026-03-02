import { defineStore } from 'pinia'
import { ref } from 'vue'
import { epicApi } from '@/api/epic.api'
import type { Epic, CreateEpicRequest, UpdateEpicRequest } from '@/types/epic.types'
import type { Issue } from '@/types/issue.types'

export const useEpicStore = defineStore('epic', () => {
  const epics = ref<Epic[]>([])
  const currentEpic = ref<Epic | null>(null)
  const epicIssues = ref<Issue[]>([])
  const totalIssues = ref(0)
  const loading = ref(false)

  async function fetchEpics(slug: string, projectId: string) {
    loading.value = true
    try {
      const { data } = await epicApi.list(slug, projectId)
      epics.value = data.results
    } finally {
      loading.value = false
    }
  }

  async function fetchEpic(slug: string, projectId: string, epicId: string) {
    const { data } = await epicApi.get(slug, projectId, epicId)
    currentEpic.value = data.epic
    epicIssues.value = data.issues
    totalIssues.value = data.total_issues
    return data
  }

  async function createEpic(slug: string, projectId: string, data: CreateEpicRequest) {
    const { data: mod } = await epicApi.create(slug, projectId, data)
    epics.value.push(mod)
    return mod
  }

  async function updateEpic(slug: string, projectId: string, epicId: string, data: UpdateEpicRequest) {
    const { data: updated } = await epicApi.update(slug, projectId, epicId, data)
    const idx = epics.value.findIndex((m) => m.id === epicId)
    if (idx >= 0) epics.value[idx] = updated
    if (currentEpic.value?.id === epicId) currentEpic.value = updated
    return updated
  }

  async function deleteEpic(slug: string, projectId: string, epicId: string) {
    await epicApi.delete(slug, projectId, epicId)
    epics.value = epics.value.filter((m) => m.id !== epicId)
    if (currentEpic.value?.id === epicId) currentEpic.value = null
  }

  async function addIssueToEpic(slug: string, projectId: string, epicId: string, issueId: string) {
    await epicApi.addIssue(slug, projectId, epicId, issueId)
  }

  async function removeIssueFromEpic(slug: string, projectId: string, epicId: string, issueId: string) {
    await epicApi.removeIssue(slug, projectId, epicId, issueId)
    epicIssues.value = epicIssues.value.filter((i) => i.id !== issueId)
    totalIssues.value = Math.max(0, totalIssues.value - 1)
  }

  return {
    epics,
    currentEpic,
    epicIssues,
    totalIssues,
    loading,
    fetchEpics,
    fetchEpic,
    createEpic,
    updateEpic,
    deleteEpic,
    addIssueToEpic,
    removeIssueFromEpic,
  }
})
