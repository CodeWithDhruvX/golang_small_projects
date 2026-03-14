import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ProfileService {
  constructor(private api: ApiService) {}

  getProfile(): Observable<any> {
    return this.api.get('/profile');
  }

  updateProfile(profile: any): Observable<any> {
    return this.api.put('/profile', profile);
  }

  uploadResume(file: File): Observable<any> {
    const formData = new FormData();
    formData.append('resume', file);
    // Custom post for multipart/form-data
    return this.api.post('/profile/resume', formData);
  }

  getResume(): Observable<any> {
    return this.api.get('/profile/resume');
  }
}
