# AI Email Response Generator Chrome Extension

A powerful Chrome extension that integrates with your existing AI Recruiter Assistant backend to generate intelligent email responses directly in Gmail.

## Features

- **AI-Powered Responses**: Generate context-aware email responses using your existing Ollama AI service
- **Seamless Gmail Integration**: Works directly within the Gmail interface
- **Customizable Tones**: Choose from Professional, Casual, Friendly, or Formal response styles
- **Flexible Length**: Generate Short, Medium, or Detailed responses
- **Real-time Generation**: Get instant responses without leaving Gmail
- **Secure Authentication**: Uses OAuth2 for secure Gmail integration
- **User-Friendly Interface**: Clean, modern UI that matches Gmail's design language

## Installation

### Prerequisites

1. **Backend Setup**: Ensure your AI Recruiter Assistant backend is running and accessible
2. **Gmail OAuth2**: Configure Gmail OAuth2 credentials in your backend
3. **Ollama Service**: Make sure Ollama is running with compatible models

### Development Installation

1. Clone this repository or navigate to the chrome-extension directory
2. Open Chrome and go to `chrome://extensions/`
3. Enable "Developer mode" in the top right
4. Click "Load unpacked" and select the chrome-extension directory
5. The extension will appear in your Chrome toolbar

### Production Installation

1. Build the extension for production
2. Submit to Chrome Web Store
3. Users can install from the Chrome Web Store

## Configuration

### Backend Configuration

Update the API URL in `background.js` and `content.js`:

```javascript
const apiUrl = 'https://yourdomain.com/api/v1'; // Production
// or
const apiUrl = 'http://localhost:8080/api/v1'; // Development
```

### Gmail OAuth2 Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable Gmail API
4. Create OAuth2 credentials
5. Add the extension's redirect URI
6. Update backend configuration with client ID and secret

## Usage

### Basic Usage

1. **Connect Gmail**: Click the extension icon and select "Connect Gmail"
2. **Authenticate**: Complete the OAuth2 flow in the popup window
3. **Open Gmail**: Navigate to Gmail and open any email
4. **Generate Response**: Click the "Generate AI Response" button below the email
5. **Review & Edit**: Review the generated response and edit as needed
6. **Send**: Send the email as usual

### Settings Configuration

Access settings by:
- Clicking the ⚙️ button next to the Generate Response button
- Opening the extension popup and clicking "Settings"

**Available Settings:**
- **Response Tone**: Professional, Casual, Friendly, Formal
- **Response Length**: Short, Medium, Detailed
- **Auto-generate**: Automatically generate responses when opening emails

### Extension Popup

Click the extension icon to access:
- Connection status
- Quick settings
- Email sync functionality
- User information
- Help and support links

## API Integration

The extension integrates with your existing backend API:

### Endpoints Used

- `POST /api/v1/gmail/generate-response` - Generate email responses
- `GET /api/v1/gmail/status` - Check Gmail connection status
- `POST /api/v1/gmail/sync` - Sync emails from Gmail
- `GET /api/v1/auth/gmail` - Initiate OAuth2 flow
- `GET /api/v1/auth/gmail/callback` - Handle OAuth2 callback

### Request/Response Format

**Generate Response Request:**
```json
{
  "email_body": "Original email content...",
  "subject": "Email subject",
  "sender": "sender@example.com",
  "tone": "professional",
  "length": "medium"
}
```

**Generate Response Response:**
```json
{
  "success": true,
  "response": "Generated email response...",
  "tone": "professional",
  "length": "medium",
  "tokens_used": 150
}
```

## Development

### File Structure

```
chrome-extension/
├── manifest.json          # Extension manifest
├── background.js          # Background service worker
├── content.js            # Gmail content script
├── content.css           # Styles for Gmail integration
├── popup.html            # Extension popup UI
├── popup.css             # Popup styles
├── popup.js              # Popup functionality
├── welcome.html          # Welcome page
├── icons/                # Extension icons
└── README.md             # This file
```

### Key Components

#### Content Script (`content.js`)
- Injects UI elements into Gmail
- Handles email content extraction
- Manages response generation
- Provides settings interface

#### Background Script (`background.js`)
- Manages API communication
- Handles authentication tokens
- Provides extension-wide state management
- Processes extension messages

#### Popup (`popup.html`, `popup.css`, `popup.js`)
- Extension settings interface
- Connection status display
- Quick actions and user management

### Customization

#### Adding New Tones

1. Update `content.js` tone options
2. Update `popup.js` form controls
3. Update backend AI service prompt generation

#### Adding New Lengths

1. Update `content.js` length options
2. Update `popup.js` form controls
3. Update backend AI service prompt generation

#### UI Customization

- Modify `content.css` for Gmail integration styles
- Update `popup.css` for popup appearance
- Adjust colors and design tokens as needed

## Troubleshooting

### Common Issues

**Extension not loading:**
- Check Chrome developer mode is enabled
- Verify manifest.json syntax
- Check for JavaScript errors in console

**Gmail integration not working:**
- Ensure content script permissions are correct
- Check Gmail UI selectors (Google may update them)
- Verify content script is injected properly

**API connection issues:**
- Check backend is running and accessible
- Verify API URL configuration
- Check CORS settings on backend
- Ensure authentication tokens are valid

**OAuth2 problems:**
- Verify Google Cloud Console configuration
- Check redirect URI matches extension URL
- Ensure OAuth2 credentials are correct

### Debugging

1. **Extension Console**: Right-click extension icon → Inspect popup
2. **Content Script Console**: Open Gmail DevTools → Console tab
3. **Background Script**: Go to `chrome://extensions/` → Details → Inspect service worker
4. **Network Tab**: Monitor API requests and responses

### Logs

Enable detailed logging:
```javascript
// In content.js or background.js
console.log('Debug info:', debugData);
```

Check browser console for error messages and debugging information.

## Security Considerations

- **Token Storage**: Authentication tokens are stored securely in Chrome storage
- **HTTPS Required**: Production use requires HTTPS for all API communications
- **Permissions**: Extension only requests necessary permissions
- **Content Security**: Follows Chrome extension security best practices
- **OAuth2 Security**: Uses secure OAuth2 flow with state parameter

## Performance

- **Lazy Loading**: Components load only when needed
- **Efficient DOM Manipulation**: Minimizes Gmail UI impact
- **Caching**: Stores settings and tokens locally
- **Background Processing**: Heavy operations run in background script

## Browser Compatibility

- **Chrome**: Full support (Manifest V3)
- **Edge**: Full support (Chromium-based)
- **Firefox**: Not supported (different extension API)
- **Safari**: Not supported (different extension API)

## Updates and Maintenance

### Version Updates

1. Update `version` in manifest.json
2. Test thoroughly with current Gmail UI
3. Update documentation
4. Release updated extension

### Gmail UI Changes

Google frequently updates Gmail's UI. When changes occur:
1. Update CSS selectors in `content.js`
2. Test all functionality
3. Update styling in `content.css`
4. Release compatibility update

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This Chrome extension is part of the AI Recruiter Assistant project. See the main project license for details.

## Support

For issues, questions, or feature requests:
1. Check the troubleshooting section
2. Review existing GitHub issues
3. Create a new issue with detailed information
4. Include browser version, Gmail UI version, and error logs

## Changelog

### v1.0.0
- Initial release
- Gmail integration
- AI response generation
- OAuth2 authentication
- Settings management
- Background service worker
- Popup interface
