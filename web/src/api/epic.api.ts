import apiClient from './client'
import type { Epic, CreateEpicRequest, UpdateEpicRequest, EpicDetail } from '@/types/epic.types'

const base = (slug: string, projectId: string) =>
  `/api/v1/workspaces/${slug}/projects/${projectId}`

export const epicApi = {
  list(slug: string, projectId: string) {
    return apiClient.get<{ results: Epic[] }>(`${base(slug, projectId)}/epics`)
  },

  create(slug: string, projectId: string, data: CreateEpicRequest) {
    return apiClient.post<Epic>(`${base(slug, projectId)}/epics`, data)
  },

  get(slug: string, projectId: string, epicId: string) {
    return apiClient.get<EpicDetail>(`${base(slug, projectId)}/epics/${epicId}`)
  },

  update(slug: string, projectId: string, epicId: string, data: UpdateEpicRequest) {
    return apiClient.put<Epic>(`${base(slug, projectId)}/epics/${epicId}`, data)
  },

  delete(slug: string, projectId: string, epicId: string) {
    return apiClient.delete(`${base(slug, projectId)}/epics/${epicId}`)
  },

  addIssue(slug: string, projectId: string, epicId: string, issueId: string) {
    return apiClient.post(`${base(slug, projectId)}/epics/${epicId}/issues`, { issue_id: issueId })
  },

  removeIssue(slug: string, projectId: string, epicId: string, issueId: string) {
    return apiClient.delete(`${base(slug, projectId)}/epics/${epicId}/issues/${issueId}`)
  },
}
