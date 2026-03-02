export interface Version {
  id: string
  project_id: string
  workspace_id: string
  name: string
  description: string
  start_date?: string
  release_date?: string
  released: boolean
  released_at?: string
  archived_at?: string
  sort_order: number
  created_by?: string
  created_at: string
  updated_at: string
}

export interface CreateVersionRequest {
  name: string
  description?: string
  start_date?: string
  release_date?: string
}

export interface UpdateVersionRequest {
  name?: string
  description?: string
  start_date?: string
  release_date?: string
  released?: boolean
  sort_order?: number
}
