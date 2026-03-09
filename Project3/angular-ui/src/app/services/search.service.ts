import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { 
  SearchRequest, 
  SearchResponse, 
  SearchResult,
  HybridSearchRequest, 
  HybridSearchResponse 
} from '../models/search.model';

@Injectable({
  providedIn: 'root'
})
export class SearchService {
  private readonly API_URL = '/api/v1';

  constructor(private http: HttpClient) {}

  semanticSearch(request: SearchRequest, headers: HttpHeaders): Observable<SearchResponse> {
    return this.http.post<SearchResponse>(`${this.API_URL}/search`, request, { headers });
  }

  hybridSearch(request: HybridSearchRequest, headers: HttpHeaders): Observable<HybridSearchResponse> {
    return this.http.post<HybridSearchResponse>(`${this.API_URL}/search/hybrid`, request, { headers });
  }

  formatScore(score: number): string {
    return (score * 100).toFixed(1) + '%';
  }

  getScoreColor(score: number): string {
    if (score >= 0.9) return 'text-green-600';
    if (score >= 0.7) return 'text-blue-600';
    if (score >= 0.5) return 'text-yellow-600';
    return 'text-red-600';
  }

  highlightSearchTerms(content: string, query: string): string {
    if (!query || !content) return content;
    
    const terms = query.split(/\s+/).filter(term => term.length > 2);
    let highlightedContent = content;
    
    terms.forEach(term => {
      const regex = new RegExp(`(${term})`, 'gi');
      highlightedContent = highlightedContent.replace(regex, '<mark class="bg-yellow-200">$1</mark>');
    });
    
    return highlightedContent;
  }

  truncateContent(content: string, maxLength: number = 200): string {
    if (content.length <= maxLength) return content;
    return content.substring(0, maxLength) + '...';
  }

  extractSnippet(content: string, query: string, contextLength: number = 100): string {
    if (!query || !content) return this.truncateContent(content, contextLength * 2);
    
    const terms = query.split(/\s+/).filter(term => term.length > 2);
    let bestMatch = -1;
    let bestScore = 0;
    
    terms.forEach(term => {
      const regex = new RegExp(term, 'gi');
      const matches = Array.from(content.matchAll(regex));
      
      matches.forEach(match => {
        if (match.index !== undefined) {
          const score = term.length;
          if (score > bestScore) {
            bestScore = score;
            bestMatch = match.index;
          }
        }
      });
    });
    
    if (bestMatch === -1) {
      return this.truncateContent(content, contextLength * 2);
    }
    
    const start = Math.max(0, bestMatch - contextLength);
    const end = Math.min(content.length, bestMatch + terms[0].length + contextLength);
    
    let snippet = content.substring(start, end);
    
    if (start > 0) snippet = '...' + snippet;
    if (end < content.length) snippet = snippet + '...';
    
    return snippet;
  }
}
