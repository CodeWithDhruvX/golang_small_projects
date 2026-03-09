import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { 
  Document, 
  DocumentUploadResponse, 
  DocumentListResponse, 
  DocumentChunk, 
  DocumentChunksResponse 
} from '../models/document.model';

@Injectable({
  providedIn: 'root'
})
export class DocumentService {
  private readonly API_URL = '/api/v1';

  constructor(private http: HttpClient) {}

  uploadDocument(file: File, headers: HttpHeaders): Observable<DocumentUploadResponse> {
    const formData = new FormData();
    formData.append('file', file);
    
    return this.http.post<DocumentUploadResponse>(`${this.API_URL}/documents/upload`, formData, { headers });
  }

  getDocuments(
    headers: HttpHeaders,
    page: number = 1, 
    limit: number = 20, 
    search?: string, 
    processed?: boolean
  ): Observable<DocumentListResponse> {
    let params = `?page=${page}&limit=${limit}`;
    if (search) params += `&search=${encodeURIComponent(search)}`;
    if (processed !== undefined) params += `&processed=${processed}`;
    
    return this.http.get<DocumentListResponse>(`${this.API_URL}/documents${params}`, { headers }).pipe(
      map(res => {
        res.documents = res.documents.map(d => ({
          ...d,
          processing_status: d.processing_status || (d.processed ? 'completed' : 'pending'),
          chunk_count: d.chunk_count ?? 0
        }));
        return res;
      })
    );
  }

  getDocument(documentId: string, headers: HttpHeaders): Observable<Document> {
    return this.http.get<Document>(`${this.API_URL}/documents/${documentId}`, { headers }).pipe(
      map(d => ({
        ...d,
        processing_status: d.processing_status || (d.processed ? 'completed' : 'pending'),
        chunk_count: d.chunk_count ?? 0
      }))
    );
  }

  deleteDocument(documentId: string, headers: HttpHeaders): Observable<any> {
    return this.http.delete(`${this.API_URL}/documents/${documentId}`, { headers });
  }

  reindexDocument(documentId: string, headers: HttpHeaders): Observable<any> {
    return this.http.post(`${this.API_URL}/documents/${documentId}/reindex`, {}, { headers });
  }

  getDocumentChunks(
    headers: HttpHeaders,
    documentId: string, 
    page: number = 1, 
    limit: number = 50
  ): Observable<DocumentChunksResponse> {
    const params = `?page=${page}&limit=${limit}`;
    return this.http.get<DocumentChunksResponse>(`${this.API_URL}/documents/${documentId}/chunks${params}`, { headers });
  }

  formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }

  getProcessingStatusColor(status: string): string {
    switch (status) {
      case 'completed': return 'text-green-600';
      case 'processing': return 'text-blue-600';
      case 'failed': return 'text-red-600';
      case 'pending': return 'text-yellow-600';
      default: return 'text-gray-600';
    }
  }

  getProcessingStatusIcon(status: string): string {
    switch (status) {
      case 'completed': return '✅';
      case 'processing': return '⚡';
      case 'failed': return '❌';
      case 'pending': return '⏳';
      default: return '❓';
    }
  }
}
