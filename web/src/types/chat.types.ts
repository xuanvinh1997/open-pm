export interface ChatMessage {
  id: string
  project_id: string
  workspace_id: string
  user_id: string
  role: 'user' | 'assistant'
  content: string
  metadata?: Record<string, unknown>
  created_at: string
}

export interface SendChatRequest {
  message: string
}

export interface ChatResponse {
  message: ChatMessage
}
