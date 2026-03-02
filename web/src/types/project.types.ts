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
  sprint_view: boolean
  epic_view: boolean
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

export interface ProjectMember {
  id: string
  project_id: string
  user_id: string
  role: number
  is_active: boolean
  email?: string
  first_name: string
  last_name: string
  display_name: string
  avatar_url: string
  created_at: string
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
