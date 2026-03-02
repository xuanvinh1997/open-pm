export type ComponentDefaultAssigneeType = 'project_default' | 'component_lead' | 'unassigned'

export interface Component {
  id: string
  project_id: string
  workspace_id: string
  name: string
  description: string
  lead_id?: string
  default_assignee_type: ComponentDefaultAssigneeType
  sort_order: number
  created_at: string
  updated_at: string
  // Enriched lead fields (when using ComponentWithLead)
  lead_email?: string
  lead_first_name?: string
  lead_last_name?: string
  lead_display_name?: string
  lead_avatar_url?: string
}

export interface CreateComponentRequest {
  name: string
  description?: string
  lead_id?: string
  default_assignee_type?: ComponentDefaultAssigneeType
  sort_order?: number
}

export interface UpdateComponentRequest {
  name?: string
  description?: string
  lead_id?: string
  default_assignee_type?: ComponentDefaultAssigneeType
  sort_order?: number
}
