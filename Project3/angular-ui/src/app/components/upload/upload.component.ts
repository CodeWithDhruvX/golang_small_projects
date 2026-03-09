import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { DocumentService } from '../../services/document.service';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-upload',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.scss']
})
export class UploadComponent {
  selectedFile: File | null = null;
  isDragging = false;
  isUploading = false;
  uploadProgress = 0;
  errorMessage = '';
  successMessage = '';

  constructor(
    private documentService: DocumentService,
    private authService: AuthService,
    private router: Router
  ) {}

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      this.selectFile(input.files[0]);
    }
  }

  onDragOver(event: DragEvent): void {
    event.preventDefault();
    this.isDragging = true;
  }

  onDragLeave(event: DragEvent): void {
    event.preventDefault();
    this.isDragging = false;
  }

  onDrop(event: DragEvent): void {
    event.preventDefault();
    this.isDragging = false;
    
    if (event.dataTransfer?.files && event.dataTransfer.files.length > 0) {
      this.selectFile(event.dataTransfer.files[0]);
    }
  }

  private selectFile(file: File): void {
    // Validate file size (50MB max)
    const maxSize = 50 * 1024 * 1024; // 50MB in bytes
    if (file.size > maxSize) {
      this.errorMessage = 'File size exceeds maximum limit of 50MB';
      this.selectedFile = null;
      return;
    }

    // Validate file type
    const allowedTypes = [
      'application/pdf',
      'text/plain',
      'text/markdown',
      'text/x-markdown',
      'application/msword',
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
      'application/vnd.ms-excel',
      'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      'application/vnd.ms-powerpoint',
      'application/vnd.openxmlformats-officedocument.presentationml.presentation',
      'text/html'
    ];

    if (!allowedTypes.includes(file.type) && !file.name.endsWith('.md') && !file.name.endsWith('.txt')) {
      this.errorMessage = 'File type not supported. Please upload PDF, Word, Excel, PowerPoint, text, or markdown files.';
      this.selectedFile = null;
      return;
    }

    this.selectedFile = file;
    this.errorMessage = '';
    this.successMessage = '';
  }

  uploadFile(): void {
    if (!this.selectedFile) {
      this.errorMessage = 'Please select a file to upload';
      return;
    }

    // Check authentication before making API call
    if (!this.authService.isAuthenticated) {
      this.router.navigate(['/login']);
      return;
    }

    this.isUploading = true;
    this.uploadProgress = 0;
    this.errorMessage = '';
    this.successMessage = '';

    const headers = this.authService.getAuthHeadersForUpload();

    this.documentService.uploadDocument(this.selectedFile, headers).subscribe({
      next: (response) => {
        this.successMessage = `Document "${response.filename}" uploaded successfully! Processing started.`;
        this.selectedFile = null;
        this.isUploading = false;
        this.uploadProgress = 100;
        
        // Reset file input
        const fileInput = document.getElementById('file-input') as HTMLInputElement;
        if (fileInput) {
          fileInput.value = '';
        }

        setTimeout(() => {
          this.uploadProgress = 0;
          this.router.navigate(['/documents']);
        }, 1500);
      },
      error: (error) => {
        this.errorMessage = error.message || 'Upload failed. Please try again.';
        this.isUploading = false;
        this.uploadProgress = 0;
      }
    });
  }

  removeFile(): void {
    this.selectedFile = null;
    this.errorMessage = '';
    this.successMessage = '';
    
    // Reset file input
    const fileInput = document.getElementById('file-input') as HTMLInputElement;
    if (fileInput) {
      fileInput.value = '';
    }
  }

  formatFileSize(bytes: number): string {
    return this.documentService.formatFileSize(bytes);
  }

  getFileIcon(): string {
    if (!this.selectedFile) return '';
    
    const fileName = this.selectedFile.name.toLowerCase();
    const fileType = this.selectedFile.type.toLowerCase();
    
    if (fileType.includes('pdf') || fileName.endsWith('.pdf')) return '📄';
    if (fileType.includes('word') || fileName.endsWith('.doc') || fileName.endsWith('.docx')) return '📝';
    if (fileType.includes('excel') || fileName.endsWith('.xls') || fileName.endsWith('.xlsx')) return '📊';
    if (fileType.includes('powerpoint') || fileName.endsWith('.ppt') || fileName.endsWith('.pptx')) return '📈';
    if (fileType.includes('text') || fileName.endsWith('.txt')) return '📃';
    if (fileName.endsWith('.md')) return '📋';
    if (fileType.includes('html')) return '🌐';
    
    return '📄';
  }
}
