import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { EmailService } from '../../services/email.service';
import { AuthService } from '../../services/auth.service';
import { GmailService } from '../../services/gmail.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="min-h-screen bg-gray-50">
      <!-- Header -->
      <header class="bg-white shadow">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16">
            <div class="flex items-center">
              <h1 class="text-2xl font-bold text-gray-900">AI Recruiter Assistant</h1>
            </div>
            <div class="flex items-center space-x-4">
              <button class="btn btn-secondary" (click)="goToProfile()">Profile</button>
              <button class="btn btn-danger" (click)="logout()">Logout</button>
            </div>
          </div>
        </div>
      </header>

      <!-- Main content -->
      <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <!-- Stats Cards -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <div class="card">
            <div class="flex items-center">
              <div class="flex-shrink-0 bg-blue-500 rounded-md p-3">
                <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">Total Emails</dt>
                  <dd class="text-lg font-medium text-gray-900">{{ stats.totalEmails }}</dd>
                </dl>
              </div>
            </div>
          </div>

          <div class="card">
            <div class="flex items-center">
              <div class="flex-shrink-0 bg-green-500 rounded-md p-3">
                <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">Applications</dt>
                  <dd class="text-lg font-medium text-gray-900">{{ stats.applications }}</dd>
                </dl>
              </div>
            </div>
          </div>

          <div class="card">
            <div class="flex items-center">
              <div class="flex-shrink-0 bg-yellow-500 rounded-md p-3">
                <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">Pending</dt>
                  <dd class="text-lg font-medium text-gray-900">{{ stats.pending }}</dd>
                </dl>
              </div>
            </div>
          </div>

          <div class="card">
            <div class="flex items-center">
              <div class="flex-shrink-0 bg-purple-500 rounded-md p-3">
                <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <div class="ml-5 w-0 flex-1">
                <dl>
                  <dt class="text-sm font-medium text-gray-500 truncate">AI Replies</dt>
                  <dd class="text-lg font-medium text-gray-900">{{ stats.aiReplies }}</dd>
                </dl>
              </div>
            </div>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          <div class="card">
            <h2 class="text-lg font-medium text-gray-900 mb-4">Quick Actions</h2>
            <div class="space-y-3">
              <button class="btn btn-primary w-full" (click)="importEmails()">Import Emails</button>
              <button class="btn btn-secondary w-full" (click)="goToGmailSettings()">Gmail Integration</button>
              <button class="btn btn-secondary w-full" (click)="goToProfile()">Upload Resume</button>
              <button class="btn btn-secondary w-full" (click)="viewEmails()">Generate AI Reply</button>
            </div>
          </div>

          <div class="card">
            <h2 class="text-lg font-medium text-gray-900 mb-4">Recent Activity</h2>
            <div class="space-y-3">
              <div *ngFor="let activity of recentActivity" class="flex items-center justify-between py-2 border-b last:border-0 border-gray-100">
                <span class="text-sm text-gray-600">{{ activity.message }}</span>
                <span class="text-xs text-gray-500">{{ activity.time }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Emails -->
        <div class="card">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-medium text-gray-900">Recent Emails</h2>
            <button class="btn btn-secondary" (click)="viewEmails()">View All</button>
          </div>
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">From</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Subject</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr *ngFor="let email of recentEmails">
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ email.from }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ email.subject }}</td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full" 
                          [ngClass]="email.status === 'processed' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'">
                      {{ email.status | titlecase }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button class="btn btn-primary btn-sm" (click)="viewEmails()">View Details</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </main>
    </div>
  `,
  styles: []
})
export class DashboardComponent implements OnInit {
  stats = {
    totalEmails: 0,
    applications: 0,
    pending: 0,
    aiReplies: 0
  };

  recentActivity = [
    { message: 'New email from Google', time: '2 min ago' },
    { message: 'AI reply generated', time: '15 min ago' },
    { message: 'Application submitted', time: '1 hour ago' }
  ];

  recentEmails: any[] = [];

  constructor(
    private router: Router,
    private emailService: EmailService,
    private authService: AuthService,
    private gmailService: GmailService
  ) {}

  ngOnInit(): void {
    this.autoLoginIfNeeded();
  }

  autoLoginIfNeeded(): void {
    if (!this.authService.isLoggedIn()) {
      // Auto-login with test credentials for demo
      this.authService.login({
        email: 'test@example.com',
        password: 'password123'
      }).subscribe({
        next: () => {
          console.log('Auto-logged in successfully');
          this.loadDashboardData();
        },
        error: (err) => {
          console.error('Auto-login failed, redirecting to login', err);
          this.router.navigate(['/login']);
        }
      });
    } else {
      this.loadDashboardData();
    }
  }

  loadDashboardData(): void {
    this.emailService.getEmails().subscribe({
      next: (res: any) => {
        const emails = res.emails || [];
        this.recentEmails = emails.slice(0, 5);
        this.stats.totalEmails = emails.length;
        this.stats.pending = emails.filter((e: any) => e.status === 'pending').length;
        this.stats.aiReplies = emails.filter((e: any) => e.has_reply).length;
      },
      error: (err) => console.error('Error loading dashboard data', err)
    });
  }

  importEmails(): void {
    const status = this.gmailService.getCurrentStatus();
    if (!status.connected) {
      if (confirm('Gmail not connected. Would you like to connect now?')) {
        this.goToGmailSettings();
      }
      return;
    }

    this.gmailService.syncEmails().subscribe({
      next: (res) => {
        alert(`Successfully synced ${res.processed_count} emails.`);
        this.loadDashboardData();
      },
      error: (err) => {
        console.error('Error syncing emails', err);
        alert('Failed to sync emails. Please check your Gmail connection.');
      }
    });
  }

  viewEmails(): void {
    this.router.navigate(['/emails']);
  }

  goToProfile(): void {
    this.router.navigate(['/profile']);
  }

  goToGmailSettings(): void {
    this.router.navigate(['/settings/gmail']);
  }

  logout(): void {
    this.authService.logout();
  }
}
