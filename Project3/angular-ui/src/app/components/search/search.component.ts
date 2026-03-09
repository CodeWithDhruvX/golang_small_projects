import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SearchService } from '../../services/search.service';
import { AuthService } from '../../services/auth.service';
import { SearchRequest, SearchResponse, SearchResult, HybridSearchRequest, HybridSearchResponse, HybridSearchResult } from '../../models/search.model';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss']
})
export class SearchComponent {
  searchQuery = '';
  searchType: 'semantic' | 'hybrid' = 'semantic';
  searchResults: SearchResult[] | HybridSearchResult[] = [];
  isLoading = false;
  errorMessage = '';
  
  // Search options
  searchLimit = 10;
  searchThreshold = 0.7;
  semanticWeight = 0.7;
  keywordWeight = 0.3;
  
  // UI state
  hasSearched = false;

  constructor(
    private searchService: SearchService,
    private authService: AuthService
  ) {}

  onSearch(): void {
    if (!this.searchQuery.trim()) {
      this.errorMessage = 'Please enter a search query';
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';
    this.hasSearched = true;

    const headers = this.authService.getAuthHeaders();

    if (this.searchType === 'semantic') {
      const request: SearchRequest = {
        query: this.searchQuery.trim(),
        limit: this.searchLimit,
        threshold: this.searchThreshold
      };

      this.searchService.semanticSearch(request, headers).subscribe({
        next: (response: SearchResponse) => {
          this.searchResults = response.results;
          this.isLoading = false;
        },
        error: (error) => {
          this.errorMessage = error.message || 'Search failed. Please try again.';
          this.isLoading = false;
        }
      });
    } else {
      const request: HybridSearchRequest = {
        query: this.searchQuery.trim(),
        limit: this.searchLimit,
        semantic_weight: this.semanticWeight,
        keyword_weight: this.keywordWeight
      };

      this.searchService.hybridSearch(request, headers).subscribe({
        next: (response: HybridSearchResponse) => {
          this.searchResults = response.results;
          this.isLoading = false;
        },
        error: (error) => {
          this.errorMessage = error.message || 'Search failed. Please try again.';
          this.isLoading = false;
        }
      });
    }
  }

  onKeyPress(event: KeyboardEvent): void {
    if (event.key === 'Enter') {
      this.onSearch();
    }
  }

  clearSearch(): void {
    this.searchQuery = '';
    this.searchResults = [];
    this.hasSearched = false;
    this.errorMessage = '';
  }

  getScore(result: SearchResult | HybridSearchResult): number {
    if ('similarity' in result && !('keyword_score' in result)) {
      return (result as SearchResult).similarity;
    } else if ('combined_score' in result) {
      return (result as HybridSearchResult).combined_score;
    }
    return 0;
  }

  getKeywordScore(result: SearchResult | HybridSearchResult): number {
    if ('keyword_score' in result) {
      return (result as HybridSearchResult).keyword_score || 0;
    }
    return 0;
  }

  getSimilarityScore(result: SearchResult | HybridSearchResult): number {
    if ('similarity' in result) {
      return (result as HybridSearchResult).similarity || 0;
    }
    return 0;
  }

  formatScore(score: number): string {
    return this.searchService.formatScore(score);
  }

  getScoreColor(score: number): string {
    return this.searchService.getScoreColor(score);
  }

  highlightSearchTerms(content: string): string {
    return this.searchService.highlightSearchTerms(content, this.searchQuery);
  }

  truncateContent(content: string): string {
    return this.searchService.truncateContent(content);
  }

  extractSnippet(content: string): string {
    return this.searchService.extractSnippet(content, this.searchQuery);
  }

  formatFileSize(bytes: number): string {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }

  getFileIcon(contentType: string): string {
    if (contentType.includes('pdf')) return '📄';
    if (contentType.includes('word')) return '📝';
    if (contentType.includes('excel')) return '📊';
    if (contentType.includes('powerpoint')) return '📈';
    if (contentType.includes('text')) return '📃';
    if (contentType.includes('markdown')) return '📋';
    if (contentType.includes('html')) return '🌐';
    return '📄';
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString();
  }

  onSearchTypeChange(): void {
    this.clearSearch();
  }

  trackByResultId(index: number, result: SearchResult | HybridSearchResult): string {
    return result.chunk.id;
  }
}
