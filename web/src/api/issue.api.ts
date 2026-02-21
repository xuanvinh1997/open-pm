import apiClient from './client'
import type { Issue, IssueComment, CreateIssueRequest, PaginatedResponse } from '@/types/issue.types'

const base = (slug: string, projectId: string) =>
  `/api/v1/workspaces/${slug}/projects/${projectId}`

export const issueApi = {
  list(slug: string, projectId: string, page = 1, perPage = 50) {
    return apiClient.get<PaginatedResponse<Issue>>(
      `${base(slug, projectId)}/issues?page=${page}&per_page=${perPage}`,
    )
  },

  create(slug: string, projectId: string, data: CreateIssueRequest) {
    return apiClient.post<Issue>(`${base(slug, projectId)}/issues`, data)
  },

  get(slug: string, projectId: string, issueId: string) {
    return apiClient.get<Issue>(`${base(slug, projectId)}/issues/${issueId}`)
  },

  update(slug: string, projectId: string, issueId: string, data: Partial<Issue>) {
    return apiClient.put<Issue>(`${base(slug, projectId)}/issues/${issueId}`, data)
  },

  delete(slug: string, projectId: string, issueId: string) {
    return apiClient.delete(`${base(slug, projectId)}/issues/${issueId}`)
  },

  // Comments
  listComments(slug: string, projectId: string, issueId: string) {
    return apiClient.get<{ results: IssueComment[] }>(
      `${base(slug, projectId)}/issues/${issueId}/comments`,
    )
  },

  createComment(slug: string, projectId: string, issueId: string, data: { comment_html: string }) {
    return apiClient.post<IssueComment>(
      `${base(slug, projectId)}/issues/${issueId}/comments`,
      data,
    )
  },
}
