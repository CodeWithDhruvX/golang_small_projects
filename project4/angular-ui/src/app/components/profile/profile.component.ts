import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { ProfileService } from '../../services/profile.service';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="min-h-screen bg-gray-50">
      <!-- Header -->
      <header class="bg-white shadow">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16">
            <div class="flex items-center">
              <h1 class="text-2xl font-bold text-gray-900">Profile Settings</h1>
            </div>
            <div class="flex items-center space-x-4">
              <button class="btn btn-secondary" (click)="goToDashboard()">Dashboard</button>
            </div>
          </div>
        </div>
      </header>

      <!-- Main content -->
      <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <!-- Profile Information -->
          <div class="lg:col-span-2">
            <div class="card">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Personal Information</h2>
              <form class="space-y-6" (ngSubmit)="saveProfile()">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label class="form-label">Full Name</label>
                    <input type="text" class="form-input" [(ngModel)]="profile.name" name="name">
                  </div>
                  <div>
                    <label class="form-label">Email</label>
                    <input type="email" class="form-input" [(ngModel)]="profile.email" name="email" readonly>
                  </div>
                </div>

                <div>
                  <label class="form-label">Experience</label>
                  <textarea class="form-input" rows="3" [(ngModel)]="profile.experience" name="experience"></textarea>
                </div>

                <div>
                  <label class="form-label">Skills</label>
                  <input type="text" class="form-input" [(ngModel)]="profile.skills_display" name="skills_display">
                  <p class="text-sm text-gray-500 mt-1">Separate skills with commas</p>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label class="form-label">Current Salary</label>
                    <input type="number" class="form-input" [(ngModel)]="profile.current_salary" name="current_salary">
                  </div>
                  <div>
                    <label class="form-label">Expected Salary</label>
                    <input type="number" class="form-input" [(ngModel)]="profile.expected_salary" name="expected_salary">
                  </div>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label class="form-label">Notice Period (days)</label>
                    <input type="number" class="form-input" [(ngModel)]="profile.notice_period" name="notice_period">
                  </div>
                  <div>
                    <label class="form-label">Location</label>
                    <input type="text" class="form-input" [(ngModel)]="profile.location" name="location">
                  </div>
                </div>

                <div class="flex justify-end space-x-3">
                  <button type="submit" class="btn btn-primary" [disabled]="isSaving">
                    {{ isSaving ? 'Saving...' : 'Save Changes' }}
                  </button>
                </div>
              </form>
            </div>
          </div>

          <!-- Resume Upload -->
          <div class="lg:col-span-1">
            <div class="card">
              <h2 class="text-lg font-medium text-gray-900 mb-4">Resume</h2>
              
              <!-- Current Resume (Mock info for now) -->
              <div class="mb-6" *ngIf="profile.has_resume">
                <div class="flex items-center justify-between p-3 border border-gray-200 rounded-lg">
                  <div class="flex items-center">
                    <svg class="h-8 w-8 text-red-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                    </svg>
                    <div>
                      <p class="text-sm font-medium text-gray-900">Current_Resume.pdf</p>
                    </div>
                  </div>
                  <button class="btn btn-secondary btn-sm" (click)="viewResume()">View</button>
                </div>
              </div>

              <!-- Upload New Resume -->
              <div>
                <label class="form-label">Upload New Resume</label>
                <div class="mt-1 flex justify-center px-6 pt-5 pb-6 border-2 border-gray-300 border-dashed rounded-md hover:border-gray-400 transition-colors">
                  <div class="space-y-1 text-center">
                    <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48">
                      <path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                    </svg>
                    <div class="flex text-sm text-gray-600">
                      <label for="file-upload" class="relative cursor-pointer bg-white rounded-md font-medium text-blue-600 hover:text-blue-500 focus-within:outline-none focus-within:ring-2 focus-within:ring-offset-2 focus-within:ring-blue-500">
                        <span>Upload a file</span>
                        <input id="file-upload" name="file-upload" type="file" class="sr-only" accept=".pdf" (change)="onFileSelected($event)">
                      </label>
                      <p class="pl-1">or drag and drop</p>
                    </div>
                    <p class="text-xs text-gray-500">PDF up to 10MB</p>
                  </div>
                </div>
                <div *ngIf="isUploading" class="mt-2 text-sm text-blue-600">Uploading...</div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  `,
  styles: []
})
export class ProfileComponent implements OnInit {
  profile: any = {
    full_name: '',
    email: '',
    experience: '',
    skills: '',
    current_salary: 0,
    expected_salary: 0,
    notice_period: 0,
    location: '',
    has_resume: false
  };

  isSaving = false;
  isUploading = false;

  constructor(
    private router: Router,
    private profileService: ProfileService
  ) {}

  ngOnInit(): void {
    this.loadProfile();
  }

  loadProfile(): void {
    this.profileService.getProfile().subscribe({
      next: (data) => {
        this.profile = data;
        // Convert skills array to comma-separated string for the form
        if (Array.isArray(this.profile.skills)) {
          this.profile.skills_display = this.profile.skills.join(', ');
        } else {
          this.profile.skills_display = '';
        }
      },
      error: (err) => console.error('Error loading profile', err)
    });
  }

  saveProfile(): void {
    this.isSaving = true;
    
    // Create a copy to modify for the API
    const profileToSave = { ...this.profile };
    
    // Convert skills string back to array
    if (profileToSave.skills_display) {
      profileToSave.skills = profileToSave.skills_display
        .split(',')
        .map((s: string) => s.trim())
        .filter((s: string) => s !== '');
    } else {
      profileToSave.skills = [];
    }
    
    // Remove the display-only field
    delete profileToSave.skills_display;

    this.profileService.updateProfile(profileToSave).subscribe({
      next: () => {
        this.isSaving = false;
        alert('Profile updated successfully!');
        this.loadProfile();
      },
      error: (err) => {
        this.isSaving = false;
        alert('Update failed: ' + (err.error?.error || 'Unknown error'));
      }
    });
  }

  onFileSelected(event: any): void {
    const file = event.target.files[0];
    if (file) {
      this.isUploading = true;
      this.profileService.uploadResume(file).subscribe({
        next: () => {
          this.isUploading = false;
          this.profile.has_resume = true;
          alert('Resume uploaded and processed!');
          this.loadProfile();
        },
        error: (err) => {
          this.isUploading = false;
          alert('Upload failed');
        }
      });
    }
  }

  viewResume(): void {
    // Logic for viewing resume (e.g., opening in new tab)
    alert('Resume viewing not implemented in this demo');
  }

  goToDashboard(): void {
    this.router.navigate(['/dashboard']);
  }
}
