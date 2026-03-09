import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { 
  ChatSession, 
  CreateSessionRequest, 
  CreateSessionResponse, 
  ChatMessage, 
  ChatRequest, 
  ChatResponse, 
  ChatHistoryResponse,
  ChatStreamEvent
} from '../models/chat.model';

@Injectable({
  providedIn: 'root'
})
export class ChatService {
  private readonly API_URL = '/api/v1';
  private readonly STREAM_BASE =
    typeof window !== 'undefined' && window.location && window.location.port === '4300'
      ? 'http://localhost:8080'
      : '';

  constructor(private http: HttpClient) {}

  createChatSession(request: CreateSessionRequest, headers: HttpHeaders): Observable<CreateSessionResponse> {
    return this.http.post<CreateSessionResponse>(`${this.API_URL}/chat/sessions`, request, { headers });
  }

  getChatSessions(headers: HttpHeaders): Observable<{ sessions: ChatSession[] }> {
    return this.http.get<{ sessions: ChatSession[] }>(`${this.API_URL}/chat/sessions`, { headers });
  }

  deleteChatSession(sessionId: string, headers: HttpHeaders): Observable<any> {
    return this.http.delete(`${this.API_URL}/chat/sessions/${sessionId}`, { headers });
  }

  sendMessage(request: ChatRequest, headers: HttpHeaders): Observable<ChatResponse> {
    return this.http.post<ChatResponse>(`${this.API_URL}/chat`, request, { headers });
  }

  sendMessageStream(request: ChatRequest, headers: HttpHeaders): Observable<ChatStreamEvent> {
    return new Observable<ChatStreamEvent>(observer => {
      const streamHeaders = headers.set('Accept', 'text/event-stream');
      
      // Use fetch with streaming for better compatibility
      const url = `${this.STREAM_BASE}${this.API_URL}/chat/stream`;
      fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': headers.get('Authorization') || '',
          'Accept': 'text/event-stream'
        },
        body: JSON.stringify({ ...request, stream: true })
      })
      .then(response => {
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const reader = response.body?.getReader();
        if (!reader) {
          throw new Error('Response body is not readable');
        }
        
        const decoder = new TextDecoder();
        let buffer = '';
        
        function processText(text: string) {
          buffer += text;
          const lines = buffer.split('\n');
          buffer = lines.pop() || '';
          
          lines.forEach(line => {
            const trimmed = line.trim();
            if (trimmed.length === 0) return;

            // Process SSE 'data:' lines; tolerate optional space and CRLF endings
            if (trimmed.startsWith('data:')) {
              const data = trimmed.replace(/^data:\s?/, '').trim();
              if (data === '[DONE]') {
                observer.complete();
                return;
              }
              
              try {
                let event = JSON.parse(data) as ChatStreamEvent;
                // Normalize potential server error payloads without explicit type
                if (!event.type && (event as any).error) {
                  event = { ...(event as any), type: 'error' } as ChatStreamEvent;
                }
                observer.next(event);
                
                if (event.type === 'end') {
                  observer.complete();
                }
              } catch (error) {
                console.error('Error parsing SSE data:', error);
              }
            }
          });
        }
        
        function read() {
          reader?.read().then(({ done, value }) => {
            if (done) {
              processText(buffer);
              observer.complete();
              return;
            }
            
            const text = decoder.decode(value, { stream: true });
            processText(text);
            read();
          }).catch(error => {
            observer.error(error);
          });
        }
        
        read();
      })
      .catch(error => {
        // Fallback to non-streaming request if streaming fails (e.g., proxy aborts)
        try {
          observer.next({ type: 'start' } as ChatStreamEvent);
          this.http.post<ChatResponse>(`${this.API_URL}/chat`, request, { headers })
            .subscribe({
              next: (resp) => {
                observer.next({
                  type: 'chunk',
                  content: resp.content,
                  session_id: resp.session_id
                } as ChatStreamEvent);
                observer.next({
                  type: 'end',
                  session_id: resp.session_id,
                  message_id: resp.message_id,
                  citations: resp.citations,
                  tokens_used: resp.tokens_used
                } as ChatStreamEvent);
                observer.complete();
              },
              error: (e) => {
                observer.error(e);
              }
            });
        } catch (e) {
          observer.error(error);
        }
      });
      
      return () => {
        // Cleanup if needed
      };
    });
  }

  getChatHistory(sessionId: string, page: number = 1, limit: number = 50, headers: HttpHeaders): Observable<ChatHistoryResponse> {
    const params = `?page=${page}&limit=${limit}`;
    return this.http.get<ChatHistoryResponse>(`${this.API_URL}/chat/sessions/${sessionId}/messages${params}`, { headers });
  }

  formatTimestamp(timestamp: string): string {
    const date = new Date(timestamp);
    return date.toLocaleString();
  }

  formatResponseTime(responseTimeMs: number): string {
    if (responseTimeMs < 1000) {
      return `${responseTimeMs}ms`;
    } else if (responseTimeMs < 60000) {
      return `${(responseTimeMs / 1000).toFixed(1)}s`;
    } else {
      return `${(responseTimeMs / 60000).toFixed(1)}m`;
    }
  }

  getModelDisplayName(model: string): string {
    const modelNames: { [key: string]: string } = {
      'llama3.1:8b': 'Llama 3.1 8B',
      'phi3:latest': 'Phi-3',
      'qwen2.5-coder:3b': 'Qwen2.5-Coder 3B'
    };
    return modelNames[model] || model;
  }
}
