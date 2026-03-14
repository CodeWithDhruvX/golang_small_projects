import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: 'login',
    loadComponent: () => import('./components/auth/login/login.component').then(m => m.LoginComponent)
  },
  {
    path: 'dashboard',
    loadComponent: () => import('./components/dashboard/dashboard.component').then(m => m.DashboardComponent)
  },
  {
    path: 'emails',
    loadComponent: () => import('./components/emails/emails.component').then(m => m.EmailsComponent)
  },
  {
    path: 'applications',
    loadComponent: () => import('./components/applications/applications.component').then(m => m.ApplicationsComponent)
  },
  {
    path: 'profile',
    loadComponent: () => import('./components/profile/profile.component').then(m => m.ProfileComponent)
  },
  {
    path: 'settings/gmail',
    loadComponent: () => import('./components/gmail/gmail.component').then(m => m.GmailComponent)
  },
  {
    path: 'auth/gmail/callback',
    loadComponent: () => import('./components/gmail/gmail-callback.component').then(m => m.GmailCallbackComponent)
  },
  {
    path: '',
    redirectTo: '/dashboard',
    pathMatch: 'full'
  },
  {
    path: '**',
    redirectTo: '/dashboard'
  }
];
