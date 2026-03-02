import apiClient from './client'
import type { Issue } from '@/types/issue.types'
import type { Page } from '@/types/page.types'

export interface SearchResponse {
  issues: Issue[]
  pages: Page[]
}

export const searchApi = {
  search(slug: string, query: string) {
    return apiClient.get<SearchResponse>(
      `/api/v1/workspaces/${slug}/search?q=${encodeURIComponent(query)}`,
    )
  },
}
