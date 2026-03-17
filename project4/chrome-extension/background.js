// Background script for Chrome extension
class BackgroundService {
  constructor() {
    this.apiUrl = 'http://localhost:8080/api/v1';
    this.init();
  }

  init() {
    console.log('AI Email Response Generator background service started');
    
    // Set up message listeners
    chrome.runtime.onMessage.addListener(this.handleMessage.bind(this));
    
    // Set up storage listener for settings changes
    chrome.storage.onChanged.addListener(this.handleStorageChange.bind(this));
    
    // Initialize default settings
    this.initializeSettings();
    
    // Set up alarm for periodic token refresh
    chrome.alarms.create('tokenRefresh', { periodInMinutes: 30 });
    chrome.alarms.onAlarm.addListener(this.handleAlarm.bind(this));
  }

  // Handle messages from content script and popup
  async handleMessage(request, sender, sendResponse) {
    try {
      switch (request.action) {
        case 'generateResponse':
          return await this.generateResponse(request.data, sendResponse);
        
        case 'validateToken':
          return await this.validateToken(request.data, sendResponse);
        
        case 'saveSettings':
          return await this.saveSettings(request.data, sendResponse);
        
        case 'getSettings':
          return await this.getSettings(sendResponse);
        
        case 'syncEmails':
          return await this.syncEmails(request.data, sendResponse);
        
        default:
          throw new Error(`Unknown action: ${request.action}`);
      }
    } catch (error) {
      console.error('Background script error:', error);
      sendResponse({ 
        success: false, 
        error: error.message 
      });
    }
    
    return true; // Keep message channel open for async response
  }

  // Generate AI response for email
  async generateResponse(data, sendResponse) {
    try {
      const settings = await this.getSettingsObject();
      
      if (!settings.token) {
        throw new Error('No authentication token found. Please connect your Gmail account first.');
      }

      const response = await fetch(`${this.apiUrl}/gmail/generate-response`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${settings.token}`
        },
        body: JSON.stringify({
          email_body: data.email_body,
          subject: data.subject,
          sender: data.sender,
          tone: data.tone || settings.tone || 'professional',
          length: data.length || settings.length || 'medium'
        })
      });

      if (!response.ok) {
        if (response.status === 401) {
          // Token expired, try to refresh
          const refreshResult = await this.refreshToken();
          if (refreshResult.success) {
            // Retry with new token
            const retryResponse = await fetch(`${this.apiUrl}/gmail/generate-response`, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${refreshResult.token}`
              },
              body: JSON.stringify({
                email_body: data.email_body,
                subject: data.subject,
                sender: data.sender,
                tone: data.tone || settings.tone || 'professional',
                length: data.length || settings.length || 'medium'
              })
            });

            if (!retryResponse.ok) {
              throw new Error(`API error: ${retryResponse.status}`);
            }

            const result = await retryResponse.json();
            sendResponse({ success: true, ...result });
            return;
          }
        }
        throw new Error(`API error: ${response.status}`);
      }

      const result = await response.json();
      sendResponse({ success: true, ...result });

    } catch (error) {
      console.error('Generate response error:', error);
      sendResponse({ 
        success: false, 
        error: error.message 
      });
    }
  }

  // Validate authentication token
  async validateToken(data, sendResponse) {
    try {
      const response = await fetch(`${this.apiUrl}/auth/validate`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${data.token}`
        }
      });

      if (!response.ok) {
        throw new Error('Invalid token');
      }

      const result = await response.json();
      sendResponse({ success: true, valid: true, user: result.user });

    } catch (error) {
      sendResponse({ 
        success: false, 
        valid: false, 
        error: error.message 
      });
    }
  }

  // Save extension settings
  async saveSettings(data, sendResponse) {
    try {
      await chrome.storage.sync.set(data);
      sendResponse({ success: true });
    } catch (error) {
      sendResponse({ 
        success: false, 
        error: error.message 
      });
    }
  }

  // Get extension settings
  async getSettings(sendResponse) {
    try {
      const settings = await this.getSettingsObject();
      sendResponse({ success: true, settings });
    } catch (error) {
      sendResponse({ 
        success: false, 
        error: error.message 
      });
    }
  }

  // Helper to get settings object
  async getSettingsObject() {
    return new Promise((resolve) => {
      chrome.storage.sync.get({
        token: '',
        tone: 'professional',
        length: 'medium',
        autoGenerate: false,
        userEmail: '',
        lastSync: null
      }, resolve);
    });
  }

  // Sync emails from Gmail
  async syncEmails(data, sendResponse) {
    try {
      const settings = await this.getSettingsObject();
      
      if (!settings.token) {
        throw new Error('No authentication token found');
      }

      const response = await fetch(`${this.apiUrl}/gmail/sync`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${settings.token}`
        },
        body: JSON.stringify({
          max_results: data.maxResults || 50,
          start_date: data.startDate,
          end_date: data.endDate
        })
      });

      if (!response.ok) {
        throw new Error(`Sync error: ${response.status}`);
      }

      const result = await response.json();
      
      // Update last sync time
      await chrome.storage.sync.set({ lastSync: new Date().toISOString() });
      
      sendResponse({ success: true, ...result });

    } catch (error) {
      sendResponse({ 
        success: false, 
        error: error.message 
      });
    }
  }

  // Refresh authentication token
  async refreshToken() {
    try {
      const settings = await this.getSettingsObject();
      
      if (!settings.refreshToken) {
        throw new Error('No refresh token available');
      }

      const response = await fetch(`${this.apiUrl}/auth/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          refresh_token: settings.refreshToken
        })
      });

      if (!response.ok) {
        throw new Error('Token refresh failed');
      }

      const result = await response.json();
      
      // Update stored token
      await chrome.storage.sync.set({ 
        token: result.token,
        lastTokenRefresh: new Date().toISOString()
      });

      return { success: true, token: result.token };

    } catch (error) {
      console.error('Token refresh error:', error);
      return { success: false, error: error.message };
    }
  }

  // Handle storage changes
  handleStorageChange(changes, namespace) {
    if (namespace === 'sync') {
      if (changes.token) {
        console.log('Token updated');
        // Notify content scripts about token change
        this.notifyContentScripts({ 
          action: 'tokenUpdated', 
          token: changes.token.newValue 
        });
      }
    }
  }

  // Notify all content scripts
  async notifyContentScripts(message) {
    try {
      const tabs = await chrome.tabs.query({ url: 'https://mail.google.com/*' });
      
      for (const tab of tabs) {
        try {
          await chrome.tabs.sendMessage(tab.id, message);
        } catch (error) {
          // Ignore errors for tabs that don't have content script loaded
          console.log('Could not send message to tab:', tab.id);
        }
      }
    } catch (error) {
      console.error('Error notifying content scripts:', error);
    }
  }

  // Handle periodic alarms
  async handleAlarm(alarm) {
    if (alarm.name === 'tokenRefresh') {
      console.log('Checking token refresh...');
      const settings = await this.getSettingsObject();
      
      if (settings.token) {
        // Validate current token
        try {
          const response = await fetch(`${this.apiUrl}/auth/validate`, {
            method: 'GET',
            headers: {
              'Authorization': `Bearer ${settings.token}`
            }
          });

          if (!response.ok) {
            console.log('Token expired, attempting refresh...');
            await this.refreshToken();
          }
        } catch (error) {
          console.error('Token validation error:', error);
        }
      }
    }
  }

  // Initialize default settings
  async initializeSettings() {
    const settings = await this.getSettingsObject();
    
    // Set defaults if not present
    const defaults = {
      tone: 'professional',
      length: 'medium',
      autoGenerate: false
    };

    let needsUpdate = false;
    const updates = {};

    for (const [key, value] of Object.entries(defaults)) {
      if (settings[key] === undefined || settings[key] === '') {
        updates[key] = value;
        needsUpdate = true;
      }
    }

    if (needsUpdate) {
      await chrome.storage.sync.set(updates);
      console.log('Default settings initialized');
    }
  }

  // Handle extension installation
  async handleInstall(details) {
    if (details.reason === 'install') {
      console.log('Extension installed');
      
      // Open welcome page
      chrome.tabs.create({
        url: chrome.runtime.getURL('welcome.html')
      });
      
    } else if (details.reason === 'update') {
      console.log('Extension updated');
    }
  }
}

// Initialize background service
const backgroundService = new BackgroundService();

// Handle extension install/update
chrome.runtime.onInstalled.addListener(backgroundService.handleInstall.bind(backgroundService));
