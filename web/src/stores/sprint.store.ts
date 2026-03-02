import { defineStore } from 'pinia'
import { ref } from 'vue'
import { sprintApi } from '@/api/sprint.api'
import type { Sprint, CreateSprintRequest, UpdateSprintRequest, SprintDetail } from '@/types/sprint.types'
import type { Issue } from '@/types/issue.types'

export const useSprintStore = defineStore('sprint', () => {
  const sprints = ref<Sprint[]>([])
  const currentSprint = ref<Sprint | null>(null)
  const sprintIssues = ref<Issue[]>([])
  const totalIssues = ref(0)
  const loading = ref(false)

  async function fetchSprints(slug: string, projectId: string) {
    loading.value = true
    try {
      const { data } = await sprintApi.list(slug, projectId)
      sprints.value = data.results
    } finally {
      loading.value = false
    }
  }

  async function fetchSprint(slug: string, projectId: string, sprintId: string) {
    const { data } = await sprintApi.get(slug, projectId, sprintId)
    currentSprint.value = data.sprint
    sprintIssues.value = data.issues
    totalIssues.value = data.total_issues
    return data
  }

  async function createSprint(slug: string, projectId: string, data: CreateSprintRequest) {
    const { data: sprint } = await sprintApi.create(slug, projectId, data)
    sprints.value.push(sprint)
    return sprint
  }

  async function updateSprint(slug: string, projectId: string, sprintId: string, data: UpdateSprintRequest) {
    const { data: updated } = await sprintApi.update(slug, projectId, sprintId, data)
    const idx = sprints.value.findIndex((c) => c.id === sprintId)
    if (idx >= 0) sprints.value[idx] = updated
    if (currentSprint.value?.id === sprintId) currentSprint.value = updated
    return updated
  }

  async function deleteSprint(slug: string, projectId: string, sprintId: string) {
    await sprintApi.delete(slug, projectId, sprintId)
    sprints.value = sprints.value.filter((c) => c.id !== sprintId)
    if (currentSprint.value?.id === sprintId) currentSprint.value = null
  }

  async function addIssueToSprint(slug: string, projectId: string, sprintId: string, issueId: string) {
    await sprintApi.addIssue(slug, projectId, sprintId, issueId)
  }

  async function removeIssueFromSprint(slug: string, projectId: string, sprintId: string, issueId: string) {
    await sprintApi.removeIssue(slug, projectId, sprintId, issueId)
    sprintIssues.value = sprintIssues.value.filter((i) => i.id !== issueId)
    totalIssues.value = Math.max(0, totalIssues.value - 1)
  }

  return {
    sprints,
    currentSprint,
    sprintIssues,
    totalIssues,
    loading,
    fetchSprints,
    fetchSprint,
    createSprint,
    updateSprint,
    deleteSprint,
    addIssueToSprint,
    removeIssueFromSprint,
  }
})
