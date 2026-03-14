import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Observable, BehaviorSubject, tap } from 'rxjs';

export interface GmailStatus {
  connected: boolean;
  email?: string;
  last_sync_at?: string;
  created_at?: string;
  message?: string;
}

export interface GmailAuthResponse {
  auth_url: string;
  state: string;
  message: string;
}

export interface SyncResponse {
  message: string;
  processed_count: number;
  total_fetched: number;
}

export interface SendEmailRequest {
  to: string;
  subject: string;
  body: string;
}

@Injectable({
  providedIn: 'root'
})
export class GmailService {
  private gmailStatusSubject = new BehaviorSubject<GmailStatus>({
    connected: false,
    message: 'Not connected'
  });
  
  public gmailStatus$ = this.gmailStatusSubject.asObservable();

  constructor(private api: ApiService) {}

  // Get Gmail integration status
  getStatus(): Observable<GmailStatus> {
    return this.api.get<GmailStatus>('/gmail/status').pipe(
      tap(status => {
        this.gmailStatusSubject.next(status);
      })
    );
  }

  // Initiate Gmail OAuth2 flow
  initiateAuth(): Observable<GmailAuthResponse> {
    return this.api.get<GmailAuthResponse>('/auth/gmail');
  }

  // Sync emails from Gmail
  syncEmails(maxResults: number = 50): Observable<SyncResponse> {
    return this.api.post<SyncResponse>('/gmail/sync', { max_results: maxResults });
  }

  // Send email via Gmail
  sendEmail(emailData: SendEmailRequest): Observable<any> {
    return this.api.post<any>('/gmail/send', emailData);
  }

  // Disconnect Gmail account
  disconnect(): Observable<any> {
    return this.api.delete<any>('/gmail/disconnect').pipe(
      tap(() => {
        this.gmailStatusSubject.next({
          connected: false,
          message: 'Gmail account disconnected'
        });
      })
    );
  }

  // Handle OAuth2 callback (called from callback component)
  handleCallback(): Observable<any> {
    // This will be called after OAuth2 redirect
    // The callback is handled by the backend, we just need to update status
    return this.getStatus();
  }

  // Refresh Gmail status
  refreshStatus(): void {
    this.getStatus().subscribe();
  }

  // Get current status value
  getCurrentStatus(): GmailStatus {
    return this.gmailStatusSubject.value;
  }
}
