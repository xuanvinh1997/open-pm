export interface Page {
  id: string
  project_id?: string
  workspace_id: string
  name: string
  description_html: string
  description_json?: Record<string, unknown>
  description_stripped: string
  color: string
  is_locked: boolean
  archived_at?: string
  owned_by: string
  parent_id?: string
  created_at: string
  updated_at: string
}

export interface CreatePageRequest {
  name: string
  description_html?: string
  description_json?: Record<string, unknown>
  parent_id?: string
  color?: string
}

export interface UpdatePageRequest {
  name?: string
  description_html?: string
  description_json?: Record<string, unknown>
  description_stripped?: string
  color?: string
  is_locked?: boolean
  parent_id?: string
}
