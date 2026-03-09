import { Component, OnInit, OnDestroy, ViewChild, ElementRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { ChatService } from '../../services/chat.service';
import { AuthService } from '../../services/auth.service';
import { 
  ChatSession, 
  CreateSessionRequest, 
  ChatMessage, 
  ChatRequest, 
  ChatStreamEvent,
  ChatHistoryResponse 
} from '../../models/chat.model';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-chat',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.scss']
})
export class ChatComponent implements OnInit, OnDestroy {
  @ViewChild('messagesContainer') private messagesContainer!: ElementRef;
  
  // Chat sessions
  sessions: ChatSession[] = [];
  currentSession: ChatSession | null = null;
  
  // Messages
  messages: ChatMessage[] = [];
  newMessage = '';
  isTyping = false;
  streamingMessage = '';
  isStreaming = false;
  
  // UI State
  isLoadingSessions = false;
  isLoadingMessages = false;
  showNewSessionModal = false;
  
  // New session form
  newSessionName = '';
  selectedModel = 'llama3.1:8b';
  availableModels = [
    { name: 'llama3.1:8b', displayName: 'Llama 3.1 8B' },
    { name: 'phi3:latest', displayName: 'Phi-3' },
    { name: 'qwen2.5-coder:3b', displayName: 'Qwen2.5-Coder 3B' }
  ];
  
  private subscriptions: Subscription[] = [];

  constructor(
    private chatService: ChatService,
    private authService: AuthService,
    private router: Router
  ) {}

  ngOnInit(): void {
    // Check authentication before loading sessions
    if (!this.authService.isAuthenticated) {
      this.router.navigate(['/login']);
      return;
    }
    this.loadSessions();
  }

  ngOnDestroy(): void {
    this.subscriptions.forEach(sub => sub.unsubscribe());
  }

  loadSessions(): void {
    // Check authentication before making API call
    if (!this.authService.isAuthenticated) {
      this.router.navigate(['/login']);
      return;
    }
    
    this.isLoadingSessions = true;
    const headers = this.authService.getAuthHeaders();
    
    this.chatService.getChatSessions(headers).subscribe({
      next: (response) => {
        this.sessions = response.sessions || [];
        this.isLoadingSessions = false;
        
        // Select first session if available
        if (this.sessions.length > 0 && !this.currentSession) {
          this.selectSession(this.sessions[0]);
        }
      },
      error: (error) => {
        console.error('Failed to load sessions:', error);
        this.sessions = [];
        this.isLoadingSessions = false;
      }
    });
  }

  selectSession(session: ChatSession): void {
    this.currentSession = session;
    this.loadMessages();
  }

  loadMessages(): void {
    if (!this.currentSession) return;
    
    // Check authentication before making API call
    if (!this.authService.isAuthenticated) {
      this.router.navigate(['/login']);
      return;
    }
    
    this.isLoadingMessages = true;
    const headers = this.authService.getAuthHeaders();
    
    this.chatService.getChatHistory(this.currentSession.session_id, 1, 50, headers).subscribe({
      next: (response: ChatHistoryResponse) => {
        this.messages = response?.messages ?? [];
        this.isLoadingMessages = false;
        this.scrollToBottom();
      },
      error: (error) => {
        console.error('Failed to load messages:', error);
        this.isLoadingMessages = false;
      }
    });
  }

  createNewSession(): void {
    if (!this.newSessionName.trim()) {
      return;
    }

    // Check authentication before making API call
    if (!this.authService.isAuthenticated) {
      this.router.navigate(['/login']);
      return;
    }

    const request: CreateSessionRequest = {
      session_name: this.newSessionName.trim(),
      model: this.selectedModel
    };

    const headers = this.authService.getAuthHeaders();
    
    this.chatService.createChatSession(request, headers).subscribe({
      next: (response) => {
        const newSession: ChatSession = {
          session_id: response.session_id,
          session_name: response.session_name,
          model: response.model,
          created_at: response.created_at,
          last_activity: response.last_activity
        };
        
        this.sessions.unshift(newSession);
        this.selectSession(newSession);
        this.showNewSessionModal = false;
        this.newSessionName = '';
      },
      error: (error) => {
        console.error('Failed to create session:', error);
      }
    });
  }

  deleteSession(sessionId: string): void {
    if (!confirm('Are you sure you want to delete this chat session? This action cannot be undone.')) {
      return;
    }

    // Check authentication before making API call
    if (!this.authService.isAuthenticated) {
      this.router.navigate(['/login']);
      return;
    }

    const headers = this.authService.getAuthHeaders();
    
    this.chatService.deleteChatSession(sessionId, headers).subscribe({
      next: () => {
        this.sessions = this.sessions.filter(s => s.session_id !== sessionId);
        
        if (this.currentSession?.session_id === sessionId) {
          this.currentSession = null;
          this.messages = [];
          
          if (this.sessions.length > 0) {
            this.selectSession(this.sessions[0]);
          }
        }
      },
      error: (error) => {
        console.error('Failed to delete session:', error);
      }
    });
  }

  sendMessage(): void {
    if (!this.newMessage.trim() || !this.currentSession || this.isStreaming) {
      return;
    }

    // Check authentication before making API call
    if (!this.authService.isAuthenticated) {
      this.router.navigate(['/login']);
      return;
    }

    if (!Array.isArray(this.messages)) {
      this.messages = [];
    }

    const messageContent = this.newMessage.trim();
    this.newMessage = '';
    this.isTyping = false;

    // Add user message immediately
    const userMessage: ChatMessage = {
      id: Date.now().toString(),
      session_id: this.currentSession.session_id,
      message_type: 'user',
      content: messageContent,
      created_at: new Date().toISOString()
    };
    
    this.messages.push(userMessage);
    this.scrollToBottom();

    // Start streaming response
    this.streamResponse(messageContent);
  }

  private streamResponse(message: string): void {
    if (!this.currentSession) return;

    this.isStreaming = true;
    this.streamingMessage = '';

    const request: ChatRequest = {
      session_id: this.currentSession.session_id,
      message: message,
      model: this.currentSession.model,
      stream: true
    };

    const headers = this.authService.getAuthHeaders();
    
    this.chatService.sendMessageStream(request, headers).subscribe({
      next: (event: ChatStreamEvent) => {
        if (event.type === 'start') {
          // Streaming started
        } else if (event.type === 'chunk') {
          this.streamingMessage += event.content || '';
          this.scrollToBottom();
        } else if (event.type === 'end') {
          // Streaming ended, add final message
          const assistantMessage: ChatMessage = {
            id: event.message_id || Date.now().toString(),
            session_id: this.currentSession!.session_id,
            message_type: 'assistant',
            content: this.streamingMessage,
            created_at: new Date().toISOString(),
            citations: event.citations,
            tokens_used: event.tokens_used
          };
          
          this.messages.push(assistantMessage);
          this.streamingMessage = '';
          this.isStreaming = false;
          this.scrollToBottom();
        } else if (event.type === 'error') {
          // Server signaled an error over SSE
          this.isStreaming = false;
          const errorMessage: ChatMessage = {
            id: Date.now().toString(),
            session_id: this.currentSession!.session_id,
            message_type: 'assistant',
            content: 'Sorry, I encountered an error while processing your message. Please try again.',
            created_at: new Date().toISOString()
          };
          this.messages.push(errorMessage);
          this.streamingMessage = '';
          this.scrollToBottom();
        }
      },
      error: (error) => {
        // Silent: fallback is handled inside service; present a friendly message if both fail
        this.isStreaming = false;
        this.streamingMessage = '';
        
        // Add error message
        const errorMessage: ChatMessage = {
          id: Date.now().toString(),
          session_id: this.currentSession!.session_id,
          message_type: 'assistant',
          content: 'Sorry, I encountered an error while processing your message. Please try again.',
          created_at: new Date().toISOString()
        };
        
        this.messages.push(errorMessage);
        this.scrollToBottom();
      },
      complete: () => {
        // Ensure UI never stays stuck if stream completes without 'end'
        if (this.isStreaming) {
          if (this.streamingMessage) {
            const assistantMessage: ChatMessage = {
              id: Date.now().toString(),
              session_id: this.currentSession!.session_id,
              message_type: 'assistant',
              content: this.streamingMessage,
              created_at: new Date().toISOString()
            };
            this.messages.push(assistantMessage);
          }
          this.streamingMessage = '';
          this.isStreaming = false;
          this.scrollToBottom();
        }
      }
    });
  }

  onTyping(): void {
    this.isTyping = this.newMessage.trim().length > 0;
  }

  scrollToBottom(): void {
    setTimeout(() => {
      if (this.messagesContainer) {
        this.messagesContainer.nativeElement.scrollTop = this.messagesContainer.nativeElement.scrollHeight;
      }
    }, 100);
  }

  formatTimestamp(timestamp: string): string {
    return this.chatService.formatTimestamp(timestamp);
  }

  getModelDisplayName(model: string): string {
    return this.chatService.getModelDisplayName(model);
  }

  openNewSessionModal(): void {
    this.showNewSessionModal = true;
    this.newSessionName = `Chat ${(this.sessions?.length || 0) + 1}`;
  }

  closeNewSessionModal(): void {
    this.showNewSessionModal = false;
    this.newSessionName = '';
  }

  onKeyPress(event: KeyboardEvent): void {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      this.sendMessage();
    }
  }
}
