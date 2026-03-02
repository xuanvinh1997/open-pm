export type ModuleStatus = 'backlog' | 'planned' | 'in-progress' | 'paused' | 'completed' | 'cancelled'

export interface Module {
  id: string
  project_id: string
  workspace_id: string
  name: string
  description: string
  description_html: string
  start_date?: string
  target_date?: string
  status: ModuleStatus
  lead_id?: string
  sort_order: number
  archived_at?: string
  created_by?: string
  created_at: string
  updated_at: string
}

export interface CreateModuleRequest {
  name: string
  description?: string
  start_date?: string
  target_date?: string
  status?: ModuleStatus
  lead_id?: string
}

export interface UpdateModuleRequest {
  name?: string
  description?: string
  start_date?: string
  target_date?: string
  status?: ModuleStatus
  lead_id?: string
  sort_order?: number
}

export interface ModuleDetail {
  module: Module
  issues: import('@/types/issue.types').Issue[]
  total_issues: number
}
