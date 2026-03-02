export interface Sprint {
  id: string
  project_id: string
  workspace_id: string
  name: string
  description: string
  start_date?: string
  end_date?: string
  owned_by: string
  sort_order: number
  progress_snapshot?: Record<string, unknown>
  archived_at?: string
  created_at: string
  updated_at: string
}

export interface CreateSprintRequest {
  name: string
  description?: string
  start_date?: string
  end_date?: string
}

export interface UpdateSprintRequest {
  name?: string
  description?: string
  start_date?: string
  end_date?: string
  sort_order?: number
}

export interface SprintDetail {
  sprint: Sprint
  issues: import('@/types/issue.types').Issue[]
  total_issues: number
}
