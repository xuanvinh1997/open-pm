import apiClient from './client'
import type { Module, CreateModuleRequest, UpdateModuleRequest, ModuleDetail } from '@/types/module.types'

const base = (slug: string, projectId: string) =>
  `/api/v1/workspaces/${slug}/projects/${projectId}`

export const moduleApi = {
  list(slug: string, projectId: string) {
    return apiClient.get<{ results: Module[] }>(`${base(slug, projectId)}/modules`)
  },

  create(slug: string, projectId: string, data: CreateModuleRequest) {
    return apiClient.post<Module>(`${base(slug, projectId)}/modules`, data)
  },

  get(slug: string, projectId: string, moduleId: string) {
    return apiClient.get<ModuleDetail>(`${base(slug, projectId)}/modules/${moduleId}`)
  },

  update(slug: string, projectId: string, moduleId: string, data: UpdateModuleRequest) {
    return apiClient.put<Module>(`${base(slug, projectId)}/modules/${moduleId}`, data)
  },

  delete(slug: string, projectId: string, moduleId: string) {
    return apiClient.delete(`${base(slug, projectId)}/modules/${moduleId}`)
  },

  addIssue(slug: string, projectId: string, moduleId: string, issueId: string) {
    return apiClient.post(`${base(slug, projectId)}/modules/${moduleId}/issues`, { issue_id: issueId })
  },

  removeIssue(slug: string, projectId: string, moduleId: string, issueId: string) {
    return apiClient.delete(`${base(slug, projectId)}/modules/${moduleId}/issues/${issueId}`)
  },
}
