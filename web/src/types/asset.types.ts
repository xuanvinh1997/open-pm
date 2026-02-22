export interface FileAsset {
  id: string
  workspace_id: string
  entity_type: string
  entity_id: string
  file_name: string
  file_size: number
  content_type: string
  storage_key: string
  uploaded_by: string
  download_url?: string
  created_at: string
}
