export interface ChatSession {
  session_id: string;
  session_name: string;
  model: string;
  created_at: string;
  last_activity: string;
  message_count?: number;
}

export interface CreateSessionRequest {
  session_name: string;
  model: string;
}

export interface CreateSessionResponse {
  session_id: string;
  session_name: string;
  model: string;
  created_at: string;
  last_activity: string;
}

export interface ChatMessage {
  id: string;
  session_id: string;
  message_type: 'user' | 'assistant';
  content: string;
  created_at: string;
  citations?: Citation[];
  model_used?: string;
  tokens_used?: number;
  response_time_ms?: number;
}

export interface Citation {
  document_id: string;
  chunk_id: string;
  content: string;
  page_number?: number;
}

export interface ChatRequest {
  session_id: string;
  message: string;
  model: string;
  stream?: boolean;
}

export interface ChatResponse {
  session_id: string;
  message_id: string;
  content: string;
  citations: Citation[];
  model_used: string;
  tokens_used: number;
  response_time_ms: number;
  created_at: string;
}

export interface ChatHistoryResponse {
  messages: ChatMessage[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export interface ChatStreamEvent {
  type: 'start' | 'chunk' | 'end' | 'error';
  content?: string;
  session_id?: string;
  message_id?: string;
  citations?: Citation[];
  tokens_used?: number;
  error?: string;
}
