import { Component, OnInit } from '@angular/core';
import { GmailService, GmailStatus, GmailAuthResponse } from '../../services/gmail.service';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

@Component({
  selector: 'app-gmail',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './gmail.component.html',
  styleUrls: ['./gmail.component.scss']
})
export class GmailComponent implements OnInit {
  gmailStatus: GmailStatus = {
    connected: false,
    message: 'Loading...'
  };
  
  isLoading = false;
  syncResult: any = null;
  showSyncResult = false;

  constructor(
    private gmailService: GmailService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.loadGmailStatus();
  }

  loadGmailStatus(): void {
    this.gmailService.getStatus().subscribe({
      next: (status) => {
        this.gmailStatus = status;
      },
      error: (error) => {
        console.error('Failed to load Gmail status:', error);
        this.gmailStatus = {
          connected: false,
          message: 'Failed to load status'
        };
      }
    });
  }

  connectGmail(): void {
    this.isLoading = true;
    this.gmailService.initiateAuth().subscribe({
      next: (response: GmailAuthResponse) => {
        // Store state for callback verification
        sessionStorage.setItem('gmail_oauth_state', response.state);
        
        // Redirect to Google OAuth2 URL
        window.location.href = response.auth_url;
      },
      error: (error) => {
        console.error('Failed to initiate Gmail auth:', error);
        this.isLoading = false;
        alert('Failed to start Gmail authentication. Please try again.');
      }
    });
  }

  syncEmails(): void {
    this.isLoading = true;
    this.showSyncResult = false;
    
    this.gmailService.syncEmails().subscribe({
      next: (result) => {
        this.syncResult = result;
        this.showSyncResult = true;
        this.isLoading = false;
        
        // Refresh status after sync
        setTimeout(() => {
          this.loadGmailStatus();
        }, 1000);
      },
      error: (error) => {
        console.error('Failed to sync emails:', error);
        this.isLoading = false;
        alert('Failed to sync emails. Please check your Gmail connection.');
      }
    });
  }

  disconnectGmail(): void {
    if (confirm('Are you sure you want to disconnect your Gmail account? This will remove all access tokens.')) {
      this.isLoading = true;
      
      this.gmailService.disconnect().subscribe({
        next: () => {
          this.gmailStatus = {
            connected: false,
            message: 'Gmail account disconnected successfully'
          };
          this.isLoading = false;
        },
        error: (error) => {
          console.error('Failed to disconnect Gmail:', error);
          this.isLoading = false;
          alert('Failed to disconnect Gmail account. Please try again.');
        }
      });
    }
  }

  refreshStatus(): void {
    this.loadGmailStatus();
  }

  formatDate(dateString: string): string {
    if (!dateString) return 'Never';
    try {
      return new Date(dateString).toLocaleString();
    } catch (error) {
      return 'Invalid date';
    }
  }
}
