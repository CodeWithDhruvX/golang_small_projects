import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private apiUrl = 'http://localhost:8082/api/v1';

  constructor(private http: HttpClient) {}

  private getHeaders(): HttpHeaders {
    const token = localStorage.getItem('token');
    return new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': token ? `Bearer ${token}` : ''
    });
  }

  // Add public method for endpoints that don't require auth
  getPublic<T>(path: string): Observable<T> {
    return this.http.get<T>(`${this.apiUrl}${path}`, { 
      headers: new HttpHeaders({ 'Content-Type': 'application/json' })
    });
  }

  get<T>(path: string): Observable<T> {
    return this.http.get<T>(`${this.apiUrl}${path}`, { headers: this.getHeaders() });
  }

  getGPUStatus(): Observable<any> {
    return this.http.get(`${this.apiUrl}/ai/gpu-status`, { headers: this.getHeaders() });
  }

  post<T>(path: string, body: any): Observable<T> {
    return this.http.post<T>(`${this.apiUrl}${path}`, body, { headers: this.getHeaders() });
  }

  put<T>(path: string, body: any): Observable<T> {
    return this.http.put<T>(`${this.apiUrl}${path}`, body, { headers: this.getHeaders() });
  }

  delete<T>(path: string): Observable<T> {
    return this.http.delete<T>(`${this.apiUrl}${path}`, { headers: this.getHeaders() });
  }
}
