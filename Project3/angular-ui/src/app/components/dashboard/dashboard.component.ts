import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { DocumentService } from '../../services/document.service';
import { ChatService } from '../../services/chat.service';
import { UserProfile } from '../../models/auth.model';
import { Document } from '../../models/document.model';
import { ChatSession } from '../../models/chat.model';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  userProfile: UserProfile | null = null;
  recentDocuments: Document[] = [];
  recentSessions: ChatSession[] = [];
  isLoading = true;
  errorMessage = '';

  constructor(
    private authService: AuthService,
    private documentService: DocumentService,
    private chatService: ChatService,
    private router: Router
  ) {}

  ngOnInit(): void {
    // Check if user is authenticated before loading data
    this.authService.isAuthenticated$.subscribe(isAuthenticated => {
      if (isAuthenticated) {
        this.loadDashboardData();
      } else {
        // Redirect to login if not authenticated
        this.router.navigate(['/login']);
      }
    });
  }

  loadDashboardData(): void {
    this.isLoading = true;
    this.errorMessage = '';

    const headers = this.authService.getAuthHeaders();
    const token = this.authService.getStoredToken();
    
    console.log('Auth headers:', headers);
    console.log('Is authenticated:', this.authService.isAuthenticated);
    console.log('Stored token:', token);

    // Check if token exists
    if (!token) {
      console.error('No token found, redirecting to login');
      this.router.navigate(['/login']);
      return;
    }

    // Load user profile
    this.authService.getUserProfile().subscribe({
      next: (profile) => {
        console.log('User profile loaded:', profile);
        this.userProfile = profile;
      },
      error: (error) => {
        console.error('Failed to load user profile:', error);
        if (error.status === 401) {
          console.log('User profile failed - token invalid, redirecting to login');
          this.router.navigate(['/login']);
        }
      }
    });

    // Load recent documents
    this.documentService.getDocuments(headers, 1, 5).subscribe({
      next: (response) => {
        console.log('Documents loaded:', response);
        this.recentDocuments = response?.documents ?? [];
      },
      error: (error) => {
        console.error('Failed to load recent documents:', error);
        console.error('Error details:', {
          status: error.status,
          statusText: error.statusText,
          url: error.url,
          message: error.message,
          headers: error.headers
        });
        if (error.status === 401) {
          console.log('Documents failed - token invalid, redirecting to login');
          this.router.navigate(['/login']);
        }
      }
    });

    // Load recent chat sessions
    this.chatService.getChatSessions(headers).subscribe({
      next: (response) => {
        console.log('Chat sessions loaded:', response);
        const sessions = response?.sessions ?? [];
        this.recentSessions = sessions.slice(0, 5);
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Failed to load chat sessions:', error);
        console.error('Error details:', {
          status: error.status,
          statusText: error.statusText,
          url: error.url,
          message: error.message,
          headers: error.headers
        });
        if (error.status === 401) {
          console.log('Chat sessions failed - token invalid, redirecting to login');
          this.router.navigate(['/login']);
        }
        this.isLoading = false;
      }
    });
  }

  navigateToDocuments(): void {
    this.router.navigate(['/documents']);
  }

  navigateToUpload(): void {
    this.router.navigate(['/upload']);
  }

  navigateToChat(): void {
    this.router.navigate(['/chat']);
  }

  navigateToSearch(): void {
    this.router.navigate(['/search']);
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString();
  }

  formatFileSize(bytes: number): string {
    return this.documentService.formatFileSize(bytes);
  }

  getProcessingStatusColor(status: string): string {
    return this.documentService.getProcessingStatusColor(status);
  }

  getProcessingStatusIcon(status: string): string {
    return this.documentService.getProcessingStatusIcon(status);
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

  getModelDisplayName(model: string): string {
    return this.chatService.getModelDisplayName(model);
  }
}
