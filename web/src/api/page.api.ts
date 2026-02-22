import apiClient from './client'
import type { Page, CreatePageRequest, UpdatePageRequest } from '@/types/page.types'

const base = (slug: string, projectId: string) =>
  `/api/v1/workspaces/${slug}/projects/${projectId}`

export const pageApi = {
  list(slug: string, projectId: string) {
    return apiClient.get<{ results: Page[] }>(`${base(slug, projectId)}/pages`)
  },

  create(slug: string, projectId: string, data: CreatePageRequest) {
    return apiClient.post<Page>(`${base(slug, projectId)}/pages`, data)
  },

  get(slug: string, projectId: string, pageId: string) {
    return apiClient.get<Page>(`${base(slug, projectId)}/pages/${pageId}`)
  },

  update(slug: string, projectId: string, pageId: string, data: UpdatePageRequest) {
    return apiClient.put<Page>(`${base(slug, projectId)}/pages/${pageId}`, data)
  },

  delete(slug: string, projectId: string, pageId: string) {
    return apiClient.delete(`${base(slug, projectId)}/pages/${pageId}`)
  },
}
