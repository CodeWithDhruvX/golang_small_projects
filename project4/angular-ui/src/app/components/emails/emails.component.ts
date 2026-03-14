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
              <button class="btn btn-secondary" (click)="goToDashboard()">Dashboard</button>
            </div>
          </div>
        </div>
      </header>

      <!-- Main content -->
      <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <!-- Filters (Mocked for now) -->
        <div class="bg-white rounded-lg shadow p-4 mb-6">
          <div class="flex flex-wrap gap-4 items-center">
            <div class="flex-1">
              <input type="text" placeholder="Search emails..." class="form-input" (input)="onSearch($event)">
            </div>
            <div class="flex items-center space-x-2">
              <label class="text-sm font-medium text-gray-700">AI Model:</label>
              <select [(ngModel)]="selectedModel" class="form-input py-1 text-sm bg-white border-gray-300 rounded-md focus:ring-purple-500 focus:border-purple-500">
                <option *ngFor="let model of availableModels" [value]="model.id">
                  {{ model.name }}
                </option>
              </select>
            </div>
          </div>
        </div>

        <!-- Email List -->
        <div class="bg-white rounded-lg shadow overflow-hidden">
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
              No emails found. Try importing some.
            </div>
          </ng-template>
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
  `]
})
export class EmailsComponent implements OnInit {
  emails: any[] = [];
  isGenerating = false;
  gpuStatus: any = null;
  
  selectedModel = 'llama3.1:8b';
  availableModels = [
    { id: 'llama3.1:8b', name: 'Llama 3.1 8B' },
    { id: 'phi3:latest', name: 'Phi-3' },
    { id: 'qwen2.5-coder:3b', name: 'Qwen 2.5 Coder 3B' }
  ];

  constructor(
    private router: Router,
    private emailService: EmailService
  ) {}

  ngOnInit(): void {
    this.loadEmails();
    this.loadGPUStatus();
  }

  loadEmails(): void {
    this.emailService.getEmails().subscribe({
      next: (res: any) => {
        console.log('Loaded emails:', res.emails);
        this.emails = res.emails || [];
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

  importEmails(): void {
    this.emailService.importEmails().subscribe({
      next: () => this.loadEmails(),
      error: (err) => alert('Import failed: ' + (err.error?.error || 'Unknown error'))
    });
  }

  classifyEmail(email: any): void {
    this.emailService.classifyEmail(email.id, email.body).subscribe({
      next: () => this.loadEmails(),
      error: (err) => alert('Classification failed')
    });
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
      (e.from && e.from.toLowerCase().includes(term))
    );
  }

  goToDashboard(): void {
    this.router.navigate(['/dashboard']);
  }
}
