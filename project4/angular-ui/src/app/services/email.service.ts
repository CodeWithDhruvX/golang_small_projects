import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class EmailService {
  constructor(private api: ApiService) {}

  getEmails(page: number = 1, limit: number = 10, startDate?: string, endDate?: string): Observable<any> {
    let params = '';
    if (page) params += `page=${page}&`;
    if (limit) params += `limit=${limit}&`;
    if (startDate) params += `start_date=${startDate}&`;
    if (endDate) params += `end_date=${endDate}&`;
    
    return this.api.get(`/emails?${params.slice(0, -1)}`);
  }

  getEmailById(id: string): Observable<any> {
    return this.api.get(`/emails/${id}`);
  }

  importEmails(): Observable<any> {
    return this.api.post('/emails/import', {});
  }

  syncEmails(startDate?: string, endDate?: string): Observable<any> {
    let params = '';
    if (startDate) params += `start_date=${startDate}&`;
    if (endDate) params += `end_date=${endDate}&`;
    
    const url = params ? `/gmail/sync?${params.slice(0, -1)}` : '/gmail/sync';
    return this.api.post(url, {});
  }

  getGmailStatus(): Observable<any> {
    return this.api.get('/gmail/status');
  }

  connectGmail(): Observable<any> {
    return this.api.get('/auth/gmail');
  }

  uploadEmail(emailData: any): Observable<any> {
    return this.api.post('/emails/upload', emailData);
  }

  classifyEmail(emailId: string, body: string): Observable<any> {
    return this.api.post('/ai/classify-email', { email_id: emailId, body });
  }

  generateReply(emailId: string, body: string, model: string = ''): Observable<any> {
    return this.api.post('/ai/generate-reply', { email_id: emailId, body, model });
  }

  getGPUStatus(): Observable<any> {
    return this.api.getGPUStatus();
  }
}
