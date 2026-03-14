import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class EmailService {
  constructor(private api: ApiService) {}

  getEmails(): Observable<any> {
    return this.api.get('/emails');
  }

  getEmailById(id: string): Observable<any> {
    return this.api.get(`/emails/${id}`);
  }

  importEmails(): Observable<any> {
    return this.api.post('/emails/import', {});
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
