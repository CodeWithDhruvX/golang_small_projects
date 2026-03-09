import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { DocumentService } from '../../services/document.service';
import { AuthService } from '../../services/auth.service';
import { Document, DocumentListResponse } from '../../models/document.model';

@Component({
  selector: 'app-document-list',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './document-list.component.html',
  styleUrls: ['./document-list.component.scss']
})
export class DocumentListComponent implements OnInit {
  documents: Document[] = [];
  isLoading = false;
  errorMessage = '';
  
  // Pagination
  currentPage = 1;
  totalPages = 0;
  totalDocuments = 0;
  pageSize = 20;
  
  // Search and filters
  searchQuery = '';
  filterProcessed: boolean | null = null;
  
  // Sort
  sortBy = 'upload_time';
  sortDirection: 'asc' | 'desc' = 'desc';
  
  // Selected documents for bulk actions
  selectedDocuments: Set<string> = new Set();
  
  @Output() documentSelected = new EventEmitter<Document>();

  constructor(
    private documentService: DocumentService,
    private authService: AuthService,
    private router: Router
  ) {}

  ngOnInit(): void {
    // Check if user is authenticated before loading documents
    if (!this.authService.isAuthenticated) {
      this.router.navigate(['/login']);
      return;
    }
    this.loadDocuments();
  }

  loadDocuments(): void {
    // Check authentication before making API call
    if (!this.authService.isAuthenticated) {
      this.router.navigate(['/login']);
      return;
    }
    
    this.isLoading = true;
    this.errorMessage = '';
    
    const headers = this.authService.getAuthHeaders();
    
    this.documentService.getDocuments(
      headers,
      this.currentPage,
      this.pageSize,
      this.searchQuery || undefined,
      this.filterProcessed !== null ? this.filterProcessed : undefined
    ).subscribe({
      next: (response: DocumentListResponse) => {
        this.documents = response.documents;
        this.totalPages = response.pagination.total_pages;
        this.totalDocuments = response.pagination.total;
        this.isLoading = false;
      },
      error: (error) => {
        this.errorMessage = error.message || 'Failed to load documents';
        this.isLoading = false;
      }
    });
  }

  onSearch(): void {
    this.currentPage = 1;
    this.loadDocuments();
  }

  onFilterChange(): void {
    this.currentPage = 1;
    this.loadDocuments();
  }

  onPageChange(page: number): void {
    if (page >= 1 && page <= this.totalPages) {
      this.currentPage = page;
      this.loadDocuments();
    }
  }

  onSort(column: string): void {
    if (this.sortBy === column) {
      this.sortDirection = this.sortDirection === 'asc' ? 'desc' : 'asc';
    } else {
      this.sortBy = column;
      this.sortDirection = 'desc';
    }
    this.loadDocuments();
  }

  selectDocument(document: Document): void {
    this.documentSelected.emit(document);
  }

  toggleDocumentSelection(documentId: string): void {
    if (this.selectedDocuments.has(documentId)) {
      this.selectedDocuments.delete(documentId);
    } else {
      this.selectedDocuments.add(documentId);
    }
  }

  selectAllDocuments(): void {
    if (this.selectedDocuments.size === this.documents.length) {
      this.selectedDocuments.clear();
    } else {
      this.documents.forEach(doc => this.selectedDocuments.add(doc.id));
    }
  }

  deleteDocument(documentId: string): void {
    if (!confirm('Are you sure you want to delete this document? This action cannot be undone.')) {
      return;
    }

    // Check authentication before making API call
    if (!this.authService.isAuthenticated) {
      this.router.navigate(['/login']);
      return;
    }

    const headers = this.authService.getAuthHeaders();
    
    this.documentService.deleteDocument(documentId, headers).subscribe({
      next: () => {
        this.loadDocuments();
      },
      error: (error) => {
        this.errorMessage = error.message || 'Failed to delete document';
      }
    });
  }

  reindexDocument(documentId: string): void {
    // Check authentication before making API call
    if (!this.authService.isAuthenticated) {
      this.router.navigate(['/login']);
      return;
    }

    const headers = this.authService.getAuthHeaders();
    
    this.documentService.reindexDocument(documentId, headers).subscribe({
      next: () => {
        this.loadDocuments();
      },
      error: (error) => {
        this.errorMessage = error.message || 'Failed to reindex document';
      }
    });
  }

  bulkDelete(): void {
    if (this.selectedDocuments.size === 0) return;
    
    if (!confirm(`Are you sure you want to delete ${this.selectedDocuments.size} document(s)? This action cannot be undone.`)) {
      return;
    }

    // Check authentication before making API call
    if (!this.authService.isAuthenticated) {
      this.router.navigate(['/login']);
      return;
    }

    const headers = this.authService.getAuthHeaders();
    let completed = 0;
    let errors = 0;

    this.selectedDocuments.forEach(documentId => {
      this.documentService.deleteDocument(documentId, headers).subscribe({
        next: () => {
          completed++;
          if (completed + errors === this.selectedDocuments.size) {
            this.selectedDocuments.clear();
            this.loadDocuments();
          }
        },
        error: () => {
          errors++;
          if (completed + errors === this.selectedDocuments.size) {
            this.selectedDocuments.clear();
            this.loadDocuments();
          }
        }
      });
    });
  }

  getSortIcon(column: string): string {
    if (this.sortBy !== column) return '↕️';
    return this.sortDirection === 'asc' ? '↑' : '↓';
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString() + ' ' + new Date(dateString).toLocaleTimeString();
  }

  getProcessingStatusColor(status: string): string {
    return this.documentService.getProcessingStatusColor(status);
  }

  getProcessingStatusIcon(status: string): string {
    return this.documentService.getProcessingStatusIcon(status);
  }

  formatFileSize(bytes: number): string {
    return this.documentService.formatFileSize(bytes);
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

  clearSelection(): void {
    this.selectedDocuments.clear();
  }

  refresh(): void {
    this.loadDocuments();
  }

  min(a: number, b: number): number {
    return a < b ? a : b;
  }
}
