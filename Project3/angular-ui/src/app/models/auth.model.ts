export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user_id: string;
  username: string;
  expires_in: number;
}

export interface RefreshTokenResponse {
  token: string;
  expires_in: number;
}

export interface UserProfile {
  user_id: string;
  username: string;
  created_at: string;
  last_login: string;
  statistics: {
    documents_uploaded: number;
    chat_sessions: number;
    total_messages: number;
    tokens_used: number;
  };
}
