import apiClient from './client'
import type { Workspace, WorkspaceMember, CreateWorkspaceRequest } from '@/types/workspace.types'

export const workspaceApi = {
  list() {
    return apiClient.get<{ results: Workspace[] }>('/api/v1/workspaces')
  },

  create(data: CreateWorkspaceRequest) {
    return apiClient.post<Workspace>('/api/v1/workspaces', data)
  },

  get(slug: string) {
    return apiClient.get<Workspace>(`/api/v1/workspaces/${slug}`)
  },

  update(slug: string, data: Partial<Workspace>) {
    return apiClient.put<Workspace>(`/api/v1/workspaces/${slug}`, data)
  },

  delete(slug: string) {
    return apiClient.delete(`/api/v1/workspaces/${slug}`)
  },

  listMembers(slug: string) {
    return apiClient.get<{ results: WorkspaceMember[] }>(`/api/v1/workspaces/${slug}/members`)
  },

  invite(slug: string, data: { email: string; role: number; message?: string }) {
    return apiClient.post(`/api/v1/workspaces/${slug}/invites`, data)
  },

  acceptInvite(token: string) {
    return apiClient.post(`/api/v1/invites/${token}/accept`)
  },
}
