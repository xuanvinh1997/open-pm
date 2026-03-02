import apiClient from './client'
import type { Cycle, CreateCycleRequest, UpdateCycleRequest, CycleDetail } from '@/types/cycle.types'

const base = (slug: string, projectId: string) =>
  `/api/v1/workspaces/${slug}/projects/${projectId}`

export const cycleApi = {
  list(slug: string, projectId: string) {
    return apiClient.get<{ results: Cycle[] }>(`${base(slug, projectId)}/cycles`)
  },

  create(slug: string, projectId: string, data: CreateCycleRequest) {
    return apiClient.post<Cycle>(`${base(slug, projectId)}/cycles`, data)
  },

  get(slug: string, projectId: string, cycleId: string) {
    return apiClient.get<CycleDetail>(`${base(slug, projectId)}/cycles/${cycleId}`)
  },

  update(slug: string, projectId: string, cycleId: string, data: UpdateCycleRequest) {
    return apiClient.put<Cycle>(`${base(slug, projectId)}/cycles/${cycleId}`, data)
  },

  delete(slug: string, projectId: string, cycleId: string) {
    return apiClient.delete(`${base(slug, projectId)}/cycles/${cycleId}`)
  },

  addIssue(slug: string, projectId: string, cycleId: string, issueId: string) {
    return apiClient.post(`${base(slug, projectId)}/cycles/${cycleId}/issues`, { issue_id: issueId })
  },

  removeIssue(slug: string, projectId: string, cycleId: string, issueId: string) {
    return apiClient.delete(`${base(slug, projectId)}/cycles/${cycleId}/issues/${issueId}`)
  },
}
