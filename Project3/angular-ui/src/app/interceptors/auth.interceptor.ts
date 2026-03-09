import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpRequest, HttpHandler } from '@angular/common/http';
import { AuthService } from '../services/auth.service';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  constructor(private authService: AuthService) {}

  intercept(req: HttpRequest<any>, next: HttpHandler) {
    // Skip adding auth headers for login endpoint
    if (req.url.includes('/auth/login')) {
      return next.handle(req);
    }

    // Skip if request already has auth headers (like file uploads)
    if (req.headers.has('Authorization')) {
      return next.handle(req);
    }

    const token = this.authService.getStoredToken();
    if (token) {
      const authReq = req.clone({
        headers: req.headers.set('Authorization', `Bearer ${token}`)
      });
      return next.handle(authReq);
    }

    return next.handle(req);
  }
}
