import type { Label } from './project.types'

export type IssuePriority = 'urgent' | 'high' | 'medium' | 'low' | 'none'

export interface Issue {
  id: string
  project_id: string
  workspace_id: string
  parent_id?: string
  state_id?: string
  name: string
  description_html: string
  description_stripped: string
  priority: IssuePriority
  start_date?: string
  target_date?: string
  sequence_id: number
  sort_order: number
  completed_at?: string
  archived_at?: string
  is_draft: boolean
  created_by?: string
  updated_by?: string
  created_at: string
  updated_at: string
  // Enriched fields
  assignees?: UserSummary[]
  labels?: Label[]
  sub_issues?: Issue[]
}

export interface UserSummary {
  id: string
  email: string
  first_name: string
  last_name: string
  display_name: string
  avatar_url: string
}

export interface IssueComment {
  id: string
  issue_id: string
  comment_html: string
  comment_stripped: string
  actor_id?: string
  first_name: string
  last_name: string
  display_name: string
  avatar_url: string
  created_at: string
  updated_at: string
}

export interface CreateIssueRequest {
  name: string
  description_html?: string
  priority?: IssuePriority
  state_id?: string
  parent_id?: string
  assignee_ids?: string[]
  label_ids?: string[]
}

export interface PaginatedResponse<T> {
  results: T[]
  total_count: number
  next_page?: number
  prev_page?: number
}
