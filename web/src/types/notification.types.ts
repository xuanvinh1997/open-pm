export interface Notification {
  id: string
  workspace_id: string
  project_id?: string
  title: string
  message: string
  data?: Record<string, unknown>
  entity_type?: string
  entity_id?: string
  sender_id?: string
  receiver_id: string
  read_at?: string
  snoozed_till?: string
  archived_at?: string
  created_at: string
  sender_first_name: string
  sender_last_name: string
  sender_display_name: string
  sender_avatar_url: string
}

export interface NotificationListResponse {
  results: Notification[]
  unread_count: number
}

export interface UnreadCountResponse {
  unread_count: number
}
