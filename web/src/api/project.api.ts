import apiClient from './client'
import type { Project, CreateProjectRequest, State, Label, ProjectMember } from '@/types/project.types'
import type { EstimateSystem, EstimateOption } from '@/types/issue.types'

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

  // Members
  listMembers(slug: string, projectId: string) {
    return apiClient.get<{ results: ProjectMember[] }>(`${base(slug)}/projects/${projectId}/members`)
  },

  addMember(slug: string, projectId: string, data: { user_id: string; role: number }) {
    return apiClient.post<ProjectMember>(`${base(slug)}/projects/${projectId}/members`, data)
  },

  updateMemberRole(slug: string, projectId: string, memberId: string, role: number) {
    return apiClient.put(`${base(slug)}/projects/${projectId}/members/${memberId}`, { role })
  },

  removeMember(slug: string, projectId: string, memberId: string) {
    return apiClient.delete(`${base(slug)}/projects/${projectId}/members/${memberId}`)
  },

  // Estimates
  getEstimateSystem(slug: string, projectId: string) {
    return apiClient.get<EstimateSystem>(`${base(slug)}/projects/${projectId}/estimates`)
  },

  updateEstimateSystem(slug: string, projectId: string, data: { type: string; estimates: EstimateOption[] }) {
    return apiClient.put<EstimateSystem>(`${base(slug)}/projects/${projectId}/estimates`, data)
  },
}
