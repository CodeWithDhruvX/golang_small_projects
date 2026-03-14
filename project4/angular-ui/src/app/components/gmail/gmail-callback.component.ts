import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { GmailService } from '../../services/gmail.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-gmail-callback',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './gmail-callback.component.html',
  styleUrls: ['./gmail-callback.component.scss']
})
export class GmailCallbackComponent implements OnInit {
  isLoading = true;
  isSuccess = false;
  errorMessage = '';
  successMessage = '';

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private gmailService: GmailService
  ) {}

  ngOnInit(): void {
    this.handleCallback();
  }

  private handleCallback(): void {
    // Get query parameters from the URL
    const code = this.route.snapshot.queryParamMap.get('code');
    const state = this.route.snapshot.queryParamMap.get('state');
    const error = this.route.snapshot.queryParamMap.get('error');

    // Check for direct token injection (from backend redirect)
    const token = this.route.snapshot.queryParamMap.get('token');
    const email = this.route.snapshot.queryParamMap.get('email');
    
    if (token) {
      localStorage.setItem('token', token);
      if (email) localStorage.setItem('user_email', email);
      this.checkAuthenticationStatus();
      return;
    }

    // Check for errors
    if (error) {
      this.handleError(error);
      return;
    }

    // Check for authorization code
    if (!code) {
      this.handleError('No authorization code received');
      return;
    }

    // Verify state parameter (prevent CSRF attacks)
    const storedState = sessionStorage.getItem('gmail_oauth_state');
    if (!storedState || storedState !== state) {
      this.handleError('Invalid state parameter');
      return;
    }

    // Clear stored state
    sessionStorage.removeItem('gmail_oauth_state');

    // The callback is actually handled by the backend
    // We just need to check if the authentication was successful
    this.checkAuthenticationStatus();
  }

  private checkAuthenticationStatus(): void {
    this.gmailService.getStatus().subscribe({
      next: (status) => {
        this.isLoading = false;
        if (status.connected) {
          this.isSuccess = true;
          
          // Check what type of flow this was
          const isLoginFlow = sessionStorage.getItem('google_login_flow') === 'true';
          const isSignupFlow = sessionStorage.getItem('google_signup_flow') === 'true';
          
          // Clean up session storage
          sessionStorage.removeItem('google_login_flow');
          sessionStorage.removeItem('google_signup_flow');
          
          if (isSignupFlow) {
            // For signup flow, show welcome message and redirect to dashboard
            this.successMessage = 'Account created successfully! Welcome to AI Recruiter Assistant.';
            setTimeout(() => {
              this.router.navigate(['/dashboard']);
            }, 3000);
          } else if (isLoginFlow) {
            // For login flow, redirect to dashboard
            setTimeout(() => {
              this.router.navigate(['/dashboard']);
            }, 2000);
          } else {
            // For settings flow, redirect to Gmail settings
            setTimeout(() => {
              this.router.navigate(['/settings/gmail']);
            }, 2000);
          }
        } else {
          this.errorMessage = 'Authentication failed. Please try again.';
        }
      },
      error: (error) => {
        this.isLoading = false;
        this.errorMessage = 'Failed to verify authentication status. Please try again.';
        console.error('Gmail callback error:', error);
      }
    });
  }

  private handleError(error: string): void {
    this.isLoading = false;
    this.isSuccess = false;
    
    if (error === 'access_denied') {
      this.errorMessage = 'You denied access to your Gmail account. Please try again and grant the necessary permissions.';
    } else {
      this.errorMessage = `Authentication error: ${error}`;
    }
  }

  retryConnection(): void {
    const isLoginFlow = sessionStorage.getItem('google_login_flow') === 'true';
    const isSignupFlow = sessionStorage.getItem('google_signup_flow') === 'true';
    
    if (isSignupFlow || isLoginFlow) {
      this.router.navigate(['/login']);
    } else {
      this.router.navigate(['/settings/gmail']);
    }
  }
}
