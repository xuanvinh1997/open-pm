import apiClient from './client'
import type { Sprint, CreateSprintRequest, UpdateSprintRequest, CompleteSprintRequest, SprintDetail } from '@/types/sprint.types'

const base = (slug: string, projectId: string) =>
  `/api/v1/workspaces/${slug}/projects/${projectId}`

export const sprintApi = {
  list(slug: string, projectId: string) {
    return apiClient.get<{ results: Sprint[] }>(`${base(slug, projectId)}/sprints`)
  },

  create(slug: string, projectId: string, data: CreateSprintRequest) {
    return apiClient.post<Sprint>(`${base(slug, projectId)}/sprints`, data)
  },

  get(slug: string, projectId: string, sprintId: string) {
    return apiClient.get<SprintDetail>(`${base(slug, projectId)}/sprints/${sprintId}`)
  },

  update(slug: string, projectId: string, sprintId: string, data: UpdateSprintRequest) {
    return apiClient.put<Sprint>(`${base(slug, projectId)}/sprints/${sprintId}`, data)
  },

  delete(slug: string, projectId: string, sprintId: string) {
    return apiClient.delete(`${base(slug, projectId)}/sprints/${sprintId}`)
  },

  complete(slug: string, projectId: string, sprintId: string, data?: CompleteSprintRequest) {
    return apiClient.post<Sprint>(`${base(slug, projectId)}/sprints/${sprintId}/complete`, data || {})
  },

  addIssue(slug: string, projectId: string, sprintId: string, issueId: string) {
    return apiClient.post(`${base(slug, projectId)}/sprints/${sprintId}/issues`, { issue_id: issueId })
  },

  removeIssue(slug: string, projectId: string, sprintId: string, issueId: string) {
    return apiClient.delete(`${base(slug, projectId)}/sprints/${sprintId}/issues/${issueId}`)
  },
}
