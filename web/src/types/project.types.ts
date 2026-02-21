export interface Project {
  id: string
  workspace_id: string
  name: string
  description: string
  identifier: string
  network: number
  emoji?: string
  cover_image_url?: string
  default_assignee_id?: string
  project_lead_id?: string
  cycle_view: boolean
  module_view: boolean
  page_view: boolean
  inbox_view: boolean
  sort_order: number
  created_by?: string
  created_at: string
  updated_at: string
}

export interface CreateProjectRequest {
  name: string
  identifier: string
  description?: string
}

export interface State {
  id: string
  project_id: string
  workspace_id: string
  name: string
  description: string
  color: string
  group: 'backlog' | 'unstarted' | 'started' | 'completed' | 'cancelled' | 'triage'
  sequence: number
  is_default: boolean
}

export interface Label {
  id: string
  project_id?: string
  workspace_id: string
  parent_id?: string
  name: string
  description: string
  color: string
  sort_order: number
}
