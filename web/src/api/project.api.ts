import apiClient from './client'
import type { Project, CreateProjectRequest, State, Label } from '@/types/project.types'

const base = (slug: string) => `/api/v1/workspaces/${slug}`

export const projectApi = {
  list(slug: string) {
    return apiClient.get<{ results: Project[] }>(`${base(slug)}/projects`)
  },

  create(slug: string, data: CreateProjectRequest) {
    return apiClient.post<Project>(`${base(slug)}/projects`, data)
  },

  get(slug: string, projectId: string) {
    return apiClient.get<Project>(`${base(slug)}/projects/${projectId}`)
  },

  update(slug: string, projectId: string, data: Partial<Project>) {
    return apiClient.put<Project>(`${base(slug)}/projects/${projectId}`, data)
  },

  delete(slug: string, projectId: string) {
    return apiClient.delete(`${base(slug)}/projects/${projectId}`)
  },

  // States
  listStates(slug: string, projectId: string) {
    return apiClient.get<{ results: State[] }>(`${base(slug)}/projects/${projectId}/states`)
  },

  createState(slug: string, projectId: string, data: Partial<State>) {
    return apiClient.post<State>(`${base(slug)}/projects/${projectId}/states`, data)
  },

  // Labels
  listLabels(slug: string, projectId: string) {
    return apiClient.get<{ results: Label[] }>(`${base(slug)}/projects/${projectId}/labels`)
  },

  createLabel(slug: string, projectId: string, data: Partial<Label>) {
    return apiClient.post<Label>(`${base(slug)}/projects/${projectId}/labels`, data)
  },
}
