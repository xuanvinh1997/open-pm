export interface Workspace {
  id: string
  name: string
  slug: string
  logo_url: string
  owner_id: string
  organization_size?: string
  created_at: string
  updated_at: string
}

export interface WorkspaceMember {
  id: string
  workspace_id: string
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

export interface WorkspaceInvite {
  id: string
  workspace_id: string
  email: string
  role: number
  accepted: boolean
  message?: string
  responded_at?: string
  created_at: string
}

export interface CreateWorkspaceRequest {
  name: string
  slug: string
  organization_size?: string
}
