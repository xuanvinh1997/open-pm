export type EpicStatus = 'backlog' | 'planned' | 'in-progress' | 'paused' | 'completed' | 'cancelled'

export interface Epic {
  id: string
  project_id: string
  workspace_id: string
  name: string
  description: string
  description_html: string
  start_date?: string
  target_date?: string
  status: EpicStatus
  lead_id?: string
  sort_order: number
  archived_at?: string
  created_by?: string
  created_at: string
  updated_at: string
}

export interface CreateEpicRequest {
  name: string
  description?: string
  start_date?: string
  target_date?: string
  status?: EpicStatus
  lead_id?: string
}

export interface UpdateEpicRequest {
  name?: string
  description?: string
  start_date?: string
  target_date?: string
  status?: EpicStatus
  lead_id?: string
  sort_order?: number
}

export interface EpicDetail {
  epic: Epic
  issues: import('@/types/issue.types').Issue[]
  total_issues: number
}
