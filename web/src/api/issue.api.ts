import apiClient from './client'
import type { Issue, IssueComment, CreateIssueRequest, PaginatedResponse, WorkLog, CreateWorkLogRequest, UpdateWorkLogRequest, WorkLogListResponse, IssueRelation, RelationType, IssueLink } from '@/types/issue.types'

const base = (slug: string, projectId: string) =>
  `/api/v1/workspaces/${slug}/projects/${projectId}`

export interface IssueFilterParams {
  page?: number
  per_page?: number
  priority?: string
  type?: string
  state?: string
  assignee?: string
  label?: string
  search?: string
  sort_by?: string
  sort_order?: string
}

export const issueApi = {
  list(slug: string, projectId: string, page = 1, perPage = 50, filters?: IssueFilterParams) {
    const params = new URLSearchParams()
    params.set('page', String(page))
    params.set('per_page', String(perPage))
    if (filters?.priority) params.set('priority', filters.priority)
    if (filters?.type) params.set('type', filters.type)
    if (filters?.state) params.set('state', filters.state)
    if (filters?.assignee) params.set('assignee', filters.assignee)
    if (filters?.label) params.set('label', filters.label)
    if (filters?.search) params.set('search', filters.search)
    if (filters?.sort_by) params.set('sort_by', filters.sort_by)
    if (filters?.sort_order) params.set('sort_order', filters.sort_order)
    return apiClient.get<PaginatedResponse<Issue>>(
      `${base(slug, projectId)}/issues?${params.toString()}`,
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

  createComment(slug: string, projectId: string, issueId: string, data: { comment_html: string; comment_json?: Record<string, unknown>; comment_stripped?: string }) {
    return apiClient.post<IssueComment>(
      `${base(slug, projectId)}/issues/${issueId}/comments`,
      data,
    )
  },

  updateComment(slug: string, projectId: string, issueId: string, commentId: string, data: { comment_html: string; comment_json?: Record<string, unknown>; comment_stripped?: string }) {
    return apiClient.put<IssueComment>(
      `${base(slug, projectId)}/issues/${issueId}/comments/${commentId}`,
      data,
    )
  },

  deleteComment(slug: string, projectId: string, issueId: string, commentId: string) {
    return apiClient.delete(
      `${base(slug, projectId)}/issues/${issueId}/comments/${commentId}`,
    )
  },

  // Work Logs
  listWorkLogs(slug: string, projectId: string, issueId: string) {
    return apiClient.get<WorkLogListResponse>(
      `${base(slug, projectId)}/issues/${issueId}/work-logs`,
    )
  },

  createWorkLog(slug: string, projectId: string, issueId: string, data: CreateWorkLogRequest) {
    return apiClient.post<WorkLog>(
      `${base(slug, projectId)}/issues/${issueId}/work-logs`,
      data,
    )
  },

  updateWorkLog(slug: string, projectId: string, issueId: string, workLogId: string, data: UpdateWorkLogRequest) {
    return apiClient.put<WorkLog>(
      `${base(slug, projectId)}/issues/${issueId}/work-logs/${workLogId}`,
      data,
    )
  },

  deleteWorkLog(slug: string, projectId: string, issueId: string, workLogId: string) {
    return apiClient.delete(
      `${base(slug, projectId)}/issues/${issueId}/work-logs/${workLogId}`,
    )
  },

  // Relations
  listRelations(slug: string, projectId: string, issueId: string) {
    return apiClient.get<{ results: IssueRelation[] }>(
      `${base(slug, projectId)}/issues/${issueId}/relations`,
    )
  },

  addRelation(slug: string, projectId: string, issueId: string, data: { related_issue_id: string; relation_type: RelationType }) {
    return apiClient.post<IssueRelation>(
      `${base(slug, projectId)}/issues/${issueId}/relations`,
      data,
    )
  },

  removeRelation(slug: string, projectId: string, issueId: string, relationId: string) {
    return apiClient.delete(
      `${base(slug, projectId)}/issues/${issueId}/relations/${relationId}`,
    )
  },

  // Links
  listLinks(slug: string, projectId: string, issueId: string) {
    return apiClient.get<{ results: IssueLink[] }>(
      `${base(slug, projectId)}/issues/${issueId}/links`,
    )
  },

  createLink(slug: string, projectId: string, issueId: string, data: { title: string; url: string }) {
    return apiClient.post<IssueLink>(
      `${base(slug, projectId)}/issues/${issueId}/links`,
      data,
    )
  },

  updateLink(slug: string, projectId: string, issueId: string, linkId: string, data: { title?: string; url?: string }) {
    return apiClient.put<IssueLink>(
      `${base(slug, projectId)}/issues/${issueId}/links/${linkId}`,
      data,
    )
  },

  deleteLink(slug: string, projectId: string, issueId: string, linkId: string) {
    return apiClient.delete(
      `${base(slug, projectId)}/issues/${issueId}/links/${linkId}`,
    )
  },
}
