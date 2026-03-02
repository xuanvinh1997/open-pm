export type SprintStatus = 'planned' | 'active' | 'completed'

export interface Sprint {
  id: string
  project_id: string
  workspace_id: string
  name: string
  description: string
  status: SprintStatus
  goal: string
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
  goal?: string
  start_date?: string
  end_date?: string
}

export interface UpdateSprintRequest {
  name?: string
  description?: string
  goal?: string
  start_date?: string
  end_date?: string
  status?: SprintStatus
  sort_order?: number
}

export interface CompleteSprintRequest {
  move_to_sprint_id?: string
}

export interface SprintDetail {
  sprint: Sprint
  issues: import('@/types/issue.types').Issue[]
  total_issues: number
}
