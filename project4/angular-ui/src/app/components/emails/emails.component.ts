import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { EmailService } from '../../services/email.service';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-emails',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="min-h-screen bg-gray-50">
      <!-- Header -->
      <header class="bg-white shadow">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16">
            <div class="flex items-center">
              <h1 class="text-2xl font-bold text-gray-900">Email Inbox</h1>
              <!-- GPU Status Indicator -->
              <div class="ml-4 flex items-center space-x-2">
                <span class="text-sm text-gray-600">AI Processing:</span>
                <span *ngIf="gpuStatus" class="px-2 py-1 text-xs font-medium rounded-full"
                      [ngClass]="gpuStatus.enabled ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'">
                  <ng-container *ngIf="gpuStatus.enabled; else cpuMode">
                    🚀 GPU ({{ gpuStatus.gpu_type }})
                  </ng-container>
                  <ng-template #cpuMode>
                    💻 CPU
                  </ng-template>
                </span>
              </div>
            </div>
            <div class="flex items-center space-x-4">
              <button class="btn btn-primary" (click)="importEmails()">Import Emails</button>
              <button class="btn btn-info" (click)="syncEmails()">Sync Emails</button>
              <button class="btn btn-secondary" (click)="goToDashboard()">Dashboard</button>
            </div>
          </div>
        </div>
      </header>

      <!-- Main content -->
      <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <!-- Filters and Date Range -->
        <div class="bg-white rounded-lg shadow p-4 mb-6">
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 items-end">
            <!-- Search -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Search</label>
              <input type="text" placeholder="Search emails..." class="form-input" (input)="onSearch($event)">
            </div>
            
            <!-- Date Range -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Start Date</label>
              <input type="date" [(ngModel)]="startDate" class="form-input" (change)="onDateRangeChange()">
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">End Date</label>
              <input type="date" [(ngModel)]="endDate" class="form-input" (change)="onDateRangeChange()">
            </div>
            
            <!-- AI Model -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">AI Model</label>
              <select [(ngModel)]="selectedModel" class="form-input">
                <option *ngFor="let model of availableModels" [value]="model.id">
                  {{ model.name }}
                </option>
              </select>
            </div>
          </div>
          
          <!-- Quick Date Range Buttons -->
          <div class="flex flex-wrap gap-2 mt-4">
            <button class="btn btn-sm btn-secondary" (click)="setQuickDateRange(7)">Last 7 Days</button>
            <button class="btn btn-sm btn-secondary" (click)="setQuickDateRange(30)">Last 30 Days</button>
            <button class="btn btn-sm btn-secondary" (click)="setQuickDateRange(90)">Last 90 Days</button>
            <button class="btn btn-sm btn-secondary" (click)="clearDateRange()">All Time</button>
          </div>
        </div>

        <!-- Email List Header with Controls -->
        <div class="bg-white rounded-lg shadow overflow-hidden">
          <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-900">Emails</h2>
            <div class="flex items-center space-x-4">
              <!-- Page Size Selector -->
              <div class="flex items-center space-x-2">
                <label class="text-sm text-gray-600">Show:</label>
                <select [(ngModel)]="pageSize" (change)="onPageSizeChange()" class="form-input text-sm py-1">
                  <option [value]="10">10</option>
                  <option [value]="25">25</option>
                  <option [value]="50">50</option>
                  <option [value]="100">100</option>
                </select>
                <span class="text-sm text-gray-600">per page</span>
              </div>
            </div>
          </div>
          <div class="divide-y divide-gray-200" *ngIf="emails.length > 0; else noEmails">
            <!-- Email Item -->
            <div *ngFor="let email of emails" class="p-6 hover:bg-gray-50">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <div class="flex items-center space-x-3 mb-2">
                    <span *ngIf="email.is_recruiter" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                      Recruiter
                    </span>
                    <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full" 
                          [ngClass]="email.status === 'processed' ? 'bg-blue-100 text-blue-800' : 'bg-yellow-100 text-yellow-800'">
                      {{ email.status | titlecase }}
                    </span>
                    <span class="text-sm text-gray-500">{{ email.created_at | date:'short' }}</span>
                  </div>
                  <h3 class="text-lg font-medium text-gray-900 mb-1">
                    {{ email.subject }}
                  </h3>
                  <p class="text-sm text-gray-600 mb-2">
                    From: {{ email.sender_name || email.sender_email || 'Unknown' }}
                  </p>
                  <p class="text-sm text-gray-700 line-clamp-2">
                    {{ email.body }}
                  </p>
                </div>
                <div class="ml-4 flex flex-col space-y-2">
                  <button class="btn btn-primary btn-sm" 
                          (click)="generateReply(email)" 
                          [disabled]="isGenerating || !hasValidBody(email) || email.generating"
                          [title]="hasValidBody(email) ? '' : 'Email body is empty - cannot generate reply'">
                    <ng-container *ngIf="email.generating; else normalButton">
                      <span class="spinner-border spinner-border-sm me-1" role="status"></span>
                      Generating...
                    </ng-container>
                    <ng-template #normalButton>
                      {{ email.has_reply ? 'View Reply' : 'Generate Reply' }}
                    </ng-template>
                  </button>
                  <button *ngIf="!email.is_recruiter && email.status === 'pending'" 
                          class="btn btn-secondary btn-sm" 
                          (click)="recallEmail(email)">
                    Recall
                  </button>
                  <button *ngIf="!email.is_recruiter && email.status === 'pending'" 
                          class="btn btn-secondary btn-sm" 
                          (click)="classifyEmail(email)">
                    Classify AI
                  </button>
                </div>
              </div>
              <!-- AI Reply Display -->
              <div *ngIf="email.showReply && email.ai_reply" class="mt-4 p-4 bg-purple-50 rounded-lg border border-purple-100">
                <div class="flex justify-between items-start mb-2">
                  <h4 class="text-sm font-bold text-purple-900">AI Draft Reply:</h4>
                  <div class="flex items-center space-x-2">
                    <span class="text-xs text-gray-600">Generated in {{ email.response_time || 'N/A' }}ms</span>
                    <span *ngIf="email.gpu_enabled" class="px-2 py-1 text-xs font-medium rounded-full bg-green-100 text-green-800">
                      🚀 GPU ({{ email.gpu_type }})
                    </span>
                    <span *ngIf="!email.gpu_enabled" class="px-2 py-1 text-xs font-medium rounded-full bg-yellow-100 text-yellow-800">
                      💻 CPU
                    </span>
                  </div>
                </div>
                <p class="text-sm text-gray-800 whitespace-pre-wrap">{{ email.ai_reply }}</p>
              </div>
            </div>
          </div>
          <ng-template #noEmails>
            <div class="p-10 text-center text-gray-500">
              No emails found for the selected date range. Try importing some or adjusting the date range.
            </div>
          </ng-template>
        </div>

        <!-- Pagination -->
        <div class="bg-white rounded-lg shadow p-4 mt-6" *ngIf="totalPages > 1">
          <div class="flex items-center justify-between">
            <div class="text-sm text-gray-700">
              Showing page {{ currentPage }} of {{ totalPages }} ({{ totalCount }} emails)
            </div>
            <div class="flex space-x-2">
              <button class="btn btn-secondary btn-sm" 
                      (click)="previousPage()" 
                      [disabled]="currentPage <= 1">
                Previous
              </button>
              <span class="px-3 py-1 text-sm bg-gray-100 rounded-md">
                {{ currentPage }}
              </span>
              <button class="btn btn-secondary btn-sm" 
                      (click)="nextPage()" 
                      [disabled]="currentPage >= totalPages">
                Next
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  `,
  styles: [`
    .line-clamp-2 {
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
    
    button:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
    
    button:disabled:hover {
      opacity: 0.5;
    }
    
    .spinner-border {
      display: inline-block;
      width: 1rem;
      height: 1rem;
      vertical-align: text-bottom;
      border: 0.125em solid currentColor;
      border-right-color: transparent;
      border-radius: 50%;
      animation: spinner-border .75s linear infinite;
    }
    
    .spinner-border-sm {
      width: 0.75rem;
      height: 0.75rem;
    }
    
    @keyframes spinner-border {
      to { transform: rotate(360deg); }
    }
    
    .me-1 {
      margin-right: 0.25rem;
    }
    
    .form-input {
      display: block;
      width: 100%;
      padding: 0.375rem 0.75rem;
      font-size: 0.875rem;
      line-height: 1.5;
      color: #495057;
      background-color: #fff;
      background-clip: padding-box;
      border: 1px solid #ced4da;
      border-radius: 0.375rem;
      transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
    }
    
    .form-input:focus {
      color: #495057;
      background-color: #fff;
      border-color: #80bdff;
      outline: 0;
      box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
    }
    
    .btn {
      display: inline-block;
      font-weight: 400;
      text-align: center;
      white-space: nowrap;
      vertical-align: middle;
      user-select: none;
      border: 1px solid transparent;
      padding: 0.375rem 0.75rem;
      font-size: 0.875rem;
      line-height: 1.5;
      border-radius: 0.375rem;
      transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out, border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
      cursor: pointer;
    }
    
    .btn-primary {
      color: #fff;
      background-color: #007bff;
      border-color: #007bff;
    }
    
    .btn-primary:hover {
      color: #fff;
      background-color: #0069d9;
      border-color: #0062cc;
    }
    
    .btn-secondary {
      color: #fff;
      background-color: #6c757d;
      border-color: #6c757d;
    }
    
    .btn-secondary:hover {
      color: #fff;
      background-color: #5a6268;
      border-color: #545b62;
    }
    
    .btn-info {
      color: #fff;
      background-color: #17a2b8;
      border-color: #17a2b8;
    }
    
    .btn-info:hover {
      color: #fff;
      background-color: #138496;
      border-color: #117a8b;
    }
    
    .btn-sm {
      padding: 0.25rem 0.5rem;
      font-size: 0.875rem;
      line-height: 1.5;
      border-radius: 0.2rem;
    }
  `]
})
export class EmailsComponent implements OnInit {
  emails: any[] = [];
  isGenerating = false;
  gpuStatus: any = null;
  
  selectedModel = 'phi3:latest';
  availableModels = [
    { id: 'llama3.1:8b', name: 'Llama 3.1 8B' },
    { id: 'phi3:latest', name: 'Phi-3' },
    { id: 'qwen2.5-coder:3b', name: 'Qwen 2.5 Coder 3B' }
  ];

  // Pagination
  currentPage = 1;
  pageSize = 25; // Increased from 10 to show more emails
  totalCount = 0;
  totalPages = 0;

  // Date range
  startDate: string = '';
  endDate: string = '';

  constructor(
    private router: Router,
    private emailService: EmailService
  ) {}

  ngOnInit(): void {
    this.setDefaultDateRange();
    this.loadEmails();
    this.loadGPUStatus();
  }

  setDefaultDateRange(): void {
    const today = new Date();
    const oneWeekAgo = new Date(today);
    oneWeekAgo.setDate(today.getDate() - 7);
    
    this.endDate = today.toISOString().split('T')[0];
    this.startDate = oneWeekAgo.toISOString().split('T')[0];
  }

  loadEmails(): void {
    this.emailService.getEmails(this.currentPage, this.pageSize, this.startDate, this.endDate).subscribe({
      next: (res: any) => {
        console.log('Loaded emails:', res.emails);
        this.emails = res.emails || [];
        
        // Update pagination info
        this.currentPage = res.page || 1;
        this.pageSize = res.limit || 10;
        
        // Calculate total count and pages (this would ideally come from the API)
        // For now, we'll estimate based on current page data
        if (this.emails.length < this.pageSize) {
          this.totalCount = (this.currentPage - 1) * this.pageSize + this.emails.length;
        } else {
          this.totalCount = this.currentPage * this.pageSize; // Estimate
        }
        this.totalPages = Math.ceil(this.totalCount / this.pageSize);
        
        // Log each email's body content
        this.emails.forEach((email, index) => {
          console.log(`Email ${index + 1}:`, {
            id: email.id,
            subject: email.subject,
            hasBody: !!email.body,
            bodyLength: email.body?.length || 0,
            bodyPreview: email.body?.substring(0, 100) + '...'
          });
        });
      },
      error: (err) => console.error('Error loading emails', err)
    });
  }

  loadGPUStatus(): void {
    this.emailService.getGPUStatus().subscribe({
      next: (status: any) => {
        this.gpuStatus = status;
        console.log('GPU Status:', status);
      },
      error: (err) => console.error('Error loading GPU status', err)
    });
  }

  onDateRangeChange(): void {
    this.currentPage = 1; // Reset to first page when date range changes
    this.loadEmails();
  }

  onPageSizeChange(): void {
    this.currentPage = 1; // Reset to first page when page size changes
    this.loadEmails();
  }

  setQuickDateRange(days: number): void {
    const today = new Date();
    const pastDate = new Date(today);
    pastDate.setDate(today.getDate() - days);
    
    this.endDate = today.toISOString().split('T')[0];
    this.startDate = pastDate.toISOString().split('T')[0];
    
    this.onDateRangeChange();
  }

  clearDateRange(): void {
    this.startDate = '';
    this.endDate = '';
    this.onDateRangeChange();
  }

  previousPage(): void {
    if (this.currentPage > 1) {
      this.currentPage--;
      this.loadEmails();
    }
  }

  nextPage(): void {
    if (this.currentPage < this.totalPages) {
      this.currentPage++;
      this.loadEmails();
    }
  }

  importEmails(): void {
    this.emailService.importEmails().subscribe({
      next: () => this.loadEmails(),
      error: (err) => alert('Import failed: ' + (err.error?.error || 'Unknown error'))
    });
  }

  syncEmails(): void {
    console.log('Syncing emails...');
    
    // Check if user is authenticated
    const token = localStorage.getItem('token');
    if (!token) {
      console.error('No authentication token found');
      
      // For testing purposes, create a test user and login
      if (confirm('No authentication found. Would you like to create a test user and login for testing?')) {
        this.createTestUserAndLogin();
        return;
      }
      
      alert('Please login first to sync emails');
      return;
    }
    
    console.log('Authentication token found:', token.substring(0, 20) + '...');
    
    // First check Gmail connection status
    console.log('Checking Gmail connection status...');
    this.emailService.getGmailStatus().subscribe({
      next: (status) => {
        console.log('Gmail status:', status);
        
        if (!status.connected) {
          // Gmail not connected, guide user through OAuth2 setup
          if (confirm('Gmail account not connected. Would you like to connect your Gmail account to sync emails?')) {
            this.connectGmailAccount();
            return;
          } else {
            alert('Gmail sync requires account connection. You can connect your Gmail account to enable email syncing.');
            return;
          }
        }
        
        // Gmail is connected, proceed with sync
        console.log('Gmail connected, proceeding with sync...');
        this.performGmailSync();
      },
      error: (err) => {
        console.error('Failed to check Gmail status:', err);
        // If status check fails, try sync anyway (might be a different error)
        this.performGmailSync();
      }
    });
  }

  connectGmailAccount(): void {
    console.log('Initiating Gmail OAuth2 connection...');
    
    this.emailService.connectGmail().subscribe({
      next: (response) => {
        console.log('Gmail OAuth URL received:', response);
        
        if (response.auth_url) {
          // Open OAuth URL in new window
          const popup = window.open(response.auth_url, 'gmail_oauth', 'width=500,height=600,scrollbars=yes,resizable=yes');
          
          // Poll for popup closure (simplified approach)
          const checkClosed = setInterval(() => {
            if (popup?.closed) {
              clearInterval(checkClosed);
              console.log('OAuth popup closed, checking connection status...');
              
              // Check if connection was successful
              setTimeout(() => {
                this.emailService.getGmailStatus().subscribe({
                  next: (status) => {
                    if (status.connected) {
                      alert('Gmail account connected successfully! You can now sync emails.');
                      this.loadEmails(); // Refresh to show updated status
                    } else {
                      alert('Gmail connection was not completed. Please try again.');
                    }
                  },
                  error: (err) => {
                    console.error('Error checking connection status:', err);
                    alert('There was an issue connecting your Gmail account. Please try again.');
                  }
                });
              }, 1000);
            }
          }, 1000);
        } else {
          alert('Failed to get Gmail authorization URL. Please check your Gmail OAuth2 configuration.');
        }
      },
      error: (err) => {
        console.error('Failed to initiate Gmail connection:', err);
        alert('Failed to connect to Gmail. Please check your Gmail OAuth2 setup in the backend configuration.');
      }
    });
  }

  performGmailSync(): void {
    console.log('Calling Gmail sync API with date range:', this.startDate, 'to', this.endDate);
    
    // Use Gmail sync with date range
    this.emailService.syncEmails(this.startDate, this.endDate).subscribe({
      next: (response) => {
        console.log('Gmail sync response:', response);
        console.log(`Synced ${response.processed_count} emails from ${response.total_fetched} fetched`);
        console.log(`Date range: ${response.start_date} to ${response.end_date}`);
        
        // Refresh the email list
        this.loadEmails();
        
        // Show success message
        alert(`Successfully synced ${response.processed_count} emails from Gmail!\nDate range: ${response.start_date} to ${response.end_date}`);
      },
      error: (err) => {
        console.error('Gmail sync failed:', err);
        console.error('Error details:', err.status, err.statusText, err.error);
        
        if (err.status === 401) {
          alert('Authentication required. Please login again.');
          // Optionally redirect to login page
          // this.router.navigate(['/login']);
        } else if (err.status === 400) {
          alert('Invalid request: ' + (err.error?.error || 'Check your date range format'));
        } else if (err.status === 0) {
          alert('Cannot connect to server. Please check if the backend is running on localhost:8080');
        } else {
          alert('Gmail sync failed: ' + (err.error?.error || err.message || 'Unknown error'));
        }
      }
    });
  }

  createTestUserAndLogin(): void {
    console.log('Creating test user and logging in...');
    
    // Create test user
    const testUser = {
      email: 'test@example.com',
      password: 'test123456',
      name: 'Test User'
    };
    
    // First try to register
    fetch('http://localhost:8080/api/v1/auth/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(testUser)
    }).then(response => response.json())
    .then(data => {
      console.log('Registration response:', data);
      
      // Now login
      return fetch('http://localhost:8080/api/v1/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          email: testUser.email,
          password: testUser.password
        })
      });
    })
    .then(response => response.json())
    .then(data => {
      console.log('Login response:', data);
      
      if (data.token) {
        localStorage.setItem('token', data.token);
        console.log('Test user logged in successfully');
        alert('Test user created and logged in! Now you can sync emails.');
        this.loadEmails(); // Refresh emails
      } else {
        alert('Login failed: ' + (data.error || 'Unknown error'));
      }
    })
    .catch(err => {
      console.error('Test user creation/login failed:', err);
      alert('Test user creation failed: ' + err.message);
    });
  }

  classifyEmail(email: any): void {
    this.emailService.classifyEmail(email.id, email.body).subscribe({
      next: () => this.loadEmails(),
      error: (err) => alert('Classification failed')
    });
  }

  recallEmail(email: any): void {
    if (confirm(`Are you sure you want to recall this email from "${email.sender_email}"? This will mark the email as recalled and remove it from your inbox.`)) {
      // For now, we'll implement a basic recall functionality
      // In a real implementation, this would call an API endpoint to handle the recall
      console.log('Recalling email:', email.id);
      
      // You could implement an API call here
      // this.emailService.recallEmail(email.id).subscribe({
      //   next: () => this.loadEmails(),
      //   error: (err) => alert('Recall failed')
      // });
      
      // For demonstration, we'll just show an alert
      alert(`Email "${email.subject}" has been marked for recall. In a production environment, this would trigger a recall process.`);
      
      // Optionally, you could remove the email from the list immediately
      // this.emails = this.emails.filter(e => e.id !== email.id);
    }
  }

  generateReply(email: any): void {
    if (email.has_reply && email.ai_reply) {
      email.showReply = !email.showReply;
      return;
    }

    console.log('Generating reply for email:', email.id);

    // Check if email body exists and is not empty
    if (!email.body || email.body.trim() === '') {
      console.warn('Cannot generate reply: Email body is empty for ID:', email.id);
      alert('Cannot generate reply: This email has no content.');
      return;
    }

    this.isGenerating = true;
    email.generating = true; // Add loading state for specific email
    
    // Set a timeout to show a warning after 60 seconds
    const timeoutWarning = setTimeout(() => {
      if (email.generating) {
        console.warn('Reply generation taking longer than expected...');
        // Could show a toast notification here
      }
    }, 60000);

    this.emailService.generateReply(email.id, email.body, this.selectedModel).subscribe({
      next: (res: any) => {
        clearTimeout(timeoutWarning);
        this.isGenerating = false;
        email.generating = false;
        email.ai_reply = res.reply;
        email.has_reply = true;
        email.showReply = true;
        
        // Store GPU information
        email.gpu_enabled = res.gpu_enabled;
        email.gpu_type = res.gpu_type;
        
        // Log GPU usage
        if (res.gpu_enabled) {
          console.log(`Reply generated using GPU (${res.gpu_type}) in ${res.response_time}ms`);
        } else {
          console.log(`Reply generated using CPU in ${res.response_time}ms`);
        }
      },
      error: (err) => {
        clearTimeout(timeoutWarning);
        this.isGenerating = false;
        email.generating = false;
        alert('Reply generation failed: ' + (err.error?.error || 'Timeout'));
      }
    });
  }

  // Helper method to check if email has valid body content
  hasValidBody(email: any): boolean {
    const isValid = !!(email.body && email.body.trim().length > 0);
    return isValid;
  }

  onSearch(event: any): void {
    // Basic local search for simplicity in this demo
    const term = (event.target as HTMLInputElement).value.toLowerCase();
    if (!term) {
      this.loadEmails();
      return;
    }
    this.emails = this.emails.filter(e => 
      (e.subject && e.subject.toLowerCase().includes(term)) || 
      (e.sender_email && e.sender_email.toLowerCase().includes(term))
    );
  }

  goToDashboard(): void {
    this.router.navigate(['/dashboard']);
  }
}
