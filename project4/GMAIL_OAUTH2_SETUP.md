# Gmail OAuth2 Setup Guide

This guide will help you set up Gmail OAuth2 integration for the AI Recruiter Assistant application.

## Overview

The application now supports Gmail OAuth2 authentication, which is the recommended approach for connecting to Gmail accounts. This provides:
- **Secure authentication** without storing passwords
- **Full Gmail API access** for reading and sending emails
- **Automatic token refresh** for long-term access
- **Google's security protections**

## Prerequisites

1. Google Cloud Platform account
2. Gmail account for testing
3. Docker and Docker Compose installed
4. Go 1.21+ installed

## Step 1: Create Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Gmail API:
   - Go to "APIs & Services" → "Library"
   - Search for "Gmail API"
   - Click "Enable"

## Step 2: Create OAuth2 Credentials

1. Go to "APIs & Services" → "Credentials"
2. Click "Create Credentials" → "OAuth client ID"
3. Configure the consent screen first if prompted:
   - Choose "External" for user type
   - Fill in required app information
   - Add scopes:
     - `https://www.googleapis.com/auth/gmail.readonly`
     - `https://www.googleapis.com/auth/gmail.send`
     - `https://www.googleapis.com/auth/userinfo.email`
     - `https://www.googleapis.com/auth/userinfo.profile`
4. Create OAuth2 client ID:
   - Application type: "Web application"
   - Name: "AI Recruiter Assistant"
   - Authorized redirect URIs:
     - `http://localhost:8080/api/v1/auth/gmail/callback` (development)
     - `https://yourdomain.com/api/v1/auth/gmail/callback` (production)
5. Save the **Client ID** and **Client Secret**

## Step 3: Configure Environment Variables

Create a `.env` file in the project root:

```bash
# Gmail OAuth2 Configuration
GMAIL_CLIENT_ID=your_google_client_id_here
GMAIL_CLIENT_SECRET=your_google_client_secret_here
GMAIL_REDIRECT_URL=http://localhost:8080/api/v1/auth/gmail/callback

# Database Configuration
DATABASE_URL=postgres://postgres:postgres@localhost:5432/ai_recruiter?sslmode=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production

# Other configurations...
REDIS_URL=redis://localhost:6379
OLLAMA_URL=http://localhost:11434
PORT=8080
```

## Step 4: Start the Application

1. Start the database and services:
```bash
docker-compose up -d
```

2. Install Go dependencies:
```bash
cd go-backend
go mod tidy
```

3. Run the application:
```bash
go run cmd/server/main.go
```

The application will be available at `http://localhost:8080`

## Step 5: Test Gmail Integration

### 1. Register/Login
First, create an account or login to get a JWT token.

### 2. Initiate Gmail OAuth
```bash
curl -X GET "http://localhost:8080/api/v1/auth/gmail" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response:
```json
{
  "auth_url": "https://accounts.google.com/o/oauth2/v2/auth?...",
  "state": "random_state_string",
  "message": "Visit the authorization URL to connect your Gmail account"
}
```

### 3. Complete OAuth Flow
1. Visit the `auth_url` in your browser
2. Sign in with your Google account
3. Grant the requested permissions
4. You'll be redirected to the callback URL

### 4. Check Integration Status
```bash
curl -X GET "http://localhost:8080/api/v1/gmail/status" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. Sync Emails
```bash
curl -X POST "http://localhost:8080/api/v1/gmail/sync" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"max_results": 50}'
```

### 6. Send Email
```bash
curl -X POST "http://localhost:8080/api/v1/gmail/send" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "to": "recipient@example.com",
    "subject": "Test Email from AI Recruiter Assistant",
    "body": "This is a test email sent via Gmail API."
  }'
```

## API Endpoints

### Authentication
- `GET /api/v1/auth/gmail` - Initiate Gmail OAuth2 flow
- `GET /api/v1/auth/gmail/callback` - OAuth2 callback (handled by Google)

### Gmail Management
- `GET /api/v1/gmail/status` - Get Gmail integration status
- `POST /api/v1/gmail/sync` - Sync emails from Gmail
- `POST /api/v1/gmail/send` - Send email via Gmail
- `DELETE /api/v1/gmail/disconnect` - Disconnect Gmail account

## Security Considerations

1. **Environment Variables**: Never commit `.env` files to version control
2. **HTTPS**: Always use HTTPS in production for OAuth2 callbacks
3. **Token Storage**: Tokens are encrypted in the database
4. **Scope Limitation**: Only request necessary Gmail scopes
5. **Token Refresh**: Automatic token refresh prevents access expiration

## Troubleshooting

### Common Issues

1. **Redirect URI Mismatch**
   - Ensure the redirect URI in Google Console matches exactly
   - Check for trailing slashes and http/https

2. **Invalid Client**
   - Verify Client ID and Client Secret are correct
   - Check if the Gmail API is enabled

3. **Access Denied**
   - User denied permission during OAuth flow
   - Some scopes may require additional verification

4. **Token Expired**
   - The application handles automatic token refresh
   - If issues persist, disconnect and reconnect Gmail account

### Debug Mode

Enable debug logging by setting:
```bash
ENVIRONMENT=development
```

## Production Deployment

For production deployment:

1. **Domain Setup**: Use your actual domain in redirect URIs
2. **HTTPS**: Configure SSL/TLS certificates
3. **Environment Variables**: Use secure secret management
4. **Rate Limiting**: Implement rate limiting for Gmail API calls
5. **Monitoring**: Set up monitoring for token refresh failures

## Migration from IMAP

If you were previously using IMAP:

1. The Gmail OAuth2 integration replaces IMAP functionality
2. Existing emails in the database remain intact
3. New emails will be fetched via Gmail API
4. IMAP configuration can be removed from environment variables

## Support

For issues with:
- **Google Cloud Console**: Google Cloud Support
- **Gmail API**: Gmail API Documentation
- **Application Issues**: Create GitHub issues

## Next Steps

1. Set up Google Cloud project and credentials
2. Configure environment variables
3. Test the OAuth2 flow
4. Integrate with your frontend application
5. Deploy to production with HTTPS

The Gmail OAuth2 integration is now fully implemented and ready for use!
