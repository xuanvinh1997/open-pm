import { defineStore } from 'pinia'
import { ref } from 'vue'
import { issueApi } from '@/api/issue.api'
import type { Issue, CreateIssueRequest } from '@/types/issue.types'

export const useEpicStore = defineStore('epic', () => {
  const epics = ref<Issue[]>([])
  const currentEpic = ref<Issue | null>(null)
  const epicIssues = ref<Issue[]>([])
  const totalIssues = ref(0)
  const loading = ref(false)

  async function fetchEpics(slug: string, projectId: string) {
    loading.value = true
    try {
      const { data } = await issueApi.list(slug, projectId, 1, 200, { type: 'epic' })
      epics.value = data.results
    } finally {
      loading.value = false
    }
  }

  async function fetchEpic(slug: string, projectId: string, epicId: string) {
    const { data: epic } = await issueApi.get(slug, projectId, epicId)
    currentEpic.value = epic
    // Sub-issues are the children of this epic issue (linked via parent_id)
    epicIssues.value = epic.sub_issues || []
    totalIssues.value = epicIssues.value.length
    return epic
  }

  async function createEpic(
    slug: string,
    projectId: string,
    data: { name: string; description_html?: string; start_date?: string; target_date?: string },
  ) {
    const req: CreateIssueRequest = {
      name: data.name,
      description_html: data.description_html,
      issue_type: 'epic',
      start_date: data.start_date,
      target_date: data.target_date,
    }
    const { data: newIssue } = await issueApi.create(slug, projectId, req)
    epics.value.push(newIssue)
    return newIssue
  }

  async function deleteEpic(slug: string, projectId: string, epicId: string) {
    await issueApi.delete(slug, projectId, epicId)
    epics.value = epics.value.filter((e) => e.id !== epicId)
    if (currentEpic.value?.id === epicId) currentEpic.value = null
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
    deleteEpic,
  }
})
