// Popup script for Chrome extension
class PopupController {
  constructor() {
    this.init();
  }

  async init() {
    console.log('Popup initialized');
    
    // Hide loading overlay
    this.hideLoading();
    
    // Load initial state
    await this.loadSettings();
    await this.checkConnectionStatus();
    
    // Set up event listeners
    this.setupEventListeners();
    
    // Update UI based on connection status
    this.updateUI();
  }

  setupEventListeners() {
    // Connect button
    document.getElementById('connect-btn').addEventListener('click', () => {
      this.connectGmail();
    });

    // Save settings button
    document.getElementById('save-settings-btn').addEventListener('click', () => {
      this.saveSettings();
    });

    // Sync emails button
    document.getElementById('sync-emails-btn').addEventListener('click', () => {
      this.syncEmails();
    });

    // Open Gmail button
    document.getElementById('open-gmail-btn').addEventListener('click', () => {
      this.openGmail();
    });

    // View stats button
    document.getElementById('view-stats-btn').addEventListener('click', () => {
      this.viewStats();
    });

    // Footer links
    document.getElementById('help-link').addEventListener('click', (e) => {
      e.preventDefault();
      this.openHelp();
    });

    document.getElementById('settings-link').addEventListener('click', (e) => {
      e.preventDefault();
      this.toggleSettings();
    });

    document.getElementById('disconnect-link').addEventListener('click', (e) => {
      e.preventDefault();
      this.disconnect();
    });

    // Toast close button
    document.getElementById('toast-close').addEventListener('click', () => {
      this.hideToast();
    });

    // Settings form controls
    document.getElementById('tone-select').addEventListener('change', () => {
      this.markSettingsChanged();
    });

    document.getElementById('length-select').addEventListener('change', () => {
      this.markSettingsChanged();
    });

    document.getElementById('auto-generate').addEventListener('change', () => {
      this.markSettingsChanged();
    });
  }

  async loadSettings() {
    try {
      const response = await this.sendMessage({ action: 'getSettings' });
      
      if (response.success) {
        const settings = response.settings;
        
        // Update form controls
        document.getElementById('tone-select').value = settings.tone || 'professional';
        document.getElementById('length-select').value = settings.length || 'medium';
        document.getElementById('auto-generate').checked = settings.autoGenerate || false;
        
        // Store settings for later use
        this.settings = settings;
      }
    } catch (error) {
      console.error('Failed to load settings:', error);
      this.showToast('Failed to load settings', 'error');
    }
  }

  async checkConnectionStatus() {
    try {
      // Check if we have a token
      const settingsResponse = await this.sendMessage({ action: 'getSettings' });
      
      if (settingsResponse.success && settingsResponse.settings.token) {
        // Validate token with backend
        const validationResponse = await this.sendMessage({
          action: 'validateToken',
          data: { token: settingsResponse.settings.token }
        });
        
        if (validationResponse.success && validationResponse.valid) {
          this.connectionStatus = {
            connected: true,
            gmail: true,
            backend: true,
            user: validationResponse.user
          };
        } else {
          this.connectionStatus = {
            connected: false,
            gmail: false,
            backend: true,
            user: null
          };
        }
      } else {
        this.connectionStatus = {
          connected: false,
          gmail: false,
          backend: true,
          user: null
        };
      }
    } catch (error) {
      console.error('Failed to check connection status:', error);
      this.connectionStatus = {
        connected: false,
        gmail: false,
        backend: false,
        user: null
      };
    }
  }

  updateUI() {
    const status = this.connectionStatus;
    
    // Update status indicator
    const statusIndicator = document.getElementById('status-indicator');
    const statusDot = statusIndicator.querySelector('.status-dot');
    const statusText = statusIndicator.querySelector('.status-text');
    
    if (status.connected) {
      statusDot.classList.remove('disconnected');
      statusText.textContent = 'Connected';
    } else {
      statusDot.classList.add('disconnected');
      statusText.textContent = 'Disconnected';
    }
    
    // Update connection status card
    const gmailStatus = document.getElementById('gmail-status');
    const backendStatus = document.getElementById('backend-status');
    const connectBtn = document.getElementById('connect-btn');
    
    if (status.gmail) {
      gmailStatus.textContent = 'Connected';
      gmailStatus.className = 'value connected';
      connectBtn.textContent = 'Disconnect';
    } else {
      gmailStatus.textContent = 'Not connected';
      gmailStatus.className = 'value disconnected';
      connectBtn.textContent = 'Connect Gmail';
    }
    
    if (status.backend) {
      backendStatus.textContent = 'Online';
      backendStatus.className = 'value connected';
    } else {
      backendStatus.textContent = 'Offline';
      backendStatus.className = 'value disconnected';
    }
    
    // Show/hide sections based on connection status
    const settingsSection = document.getElementById('settings-section');
    const quickActions = document.getElementById('quick-actions');
    const userInfo = document.getElementById('user-info');
    const disconnectLink = document.getElementById('disconnect-link');
    
    if (status.connected) {
      settingsSection.style.display = 'block';
      quickActions.style.display = 'block';
      disconnectLink.style.display = 'inline';
      
      if (status.user) {
        userInfo.style.display = 'block';
        this.updateUserInfo(status.user);
      }
    } else {
      settingsSection.style.display = 'none';
      quickActions.style.display = 'none';
      userInfo.style.display = 'none';
      disconnectLink.style.display = 'none';
    }
  }

  updateUserInfo(user) {
    const userInitial = document.getElementById('user-initial');
    const userEmail = document.getElementById('user-email');
    const lastSyncTime = document.getElementById('last-sync-time');
    
    if (user.email) {
      userInitial.textContent = user.email.charAt(0).toUpperCase();
      userEmail.textContent = user.email;
    }
    
    // Update last sync time if available
    if (this.settings && this.settings.lastSync) {
      const lastSync = new Date(this.settings.lastSync);
      lastSyncTime.textContent = this.formatRelativeTime(lastSync);
    } else {
      lastSyncTime.textContent = 'Never';
    }
  }

  async connectGmail() {
    if (this.connectionStatus.connected) {
      // Disconnect
      await this.disconnect();
      return;
    }
    
    try {
      this.showLoading('Connecting to Gmail...');
      
      // Open Gmail OAuth flow
      const authUrl = 'http://localhost:8080/api/v1/auth/gmail';
      
      // Open new window for OAuth
      const authWindow = await chrome.windows.create({
        url: authUrl,
        width: 500,
        height: 600,
        type: 'popup'
      });
      
      // Listen for window close (user cancelled)
      const checkWindow = setInterval(() => {
        chrome.windows.get(authWindow.id, (window) => {
          if (chrome.runtime.lastError || !window) {
            clearInterval(checkWindow);
            this.hideLoading();
            this.showToast('Gmail connection cancelled', 'warning');
          }
        });
      }, 1000);
      
      // Listen for OAuth success message
      chrome.runtime.onMessage.addListener((message) => {
        if (message.action === 'oauthSuccess') {
          clearInterval(checkWindow);
          chrome.windows.remove(authWindow.id);
          this.hideLoading();
          this.showToast('Gmail connected successfully!', 'success');
          
          // Reload status
          setTimeout(() => {
            this.checkConnectionStatus().then(() => this.updateUI());
          }, 1000);
        }
      });
      
    } catch (error) {
      console.error('Failed to connect Gmail:', error);
      this.hideLoading();
      this.showToast('Failed to connect Gmail', 'error');
    }
  }

  async disconnect() {
    try {
      this.showLoading('Disconnecting...');
      
      // Clear stored token
      await this.sendMessage({
        action: 'saveSettings',
        data: { token: '', refreshToken: '' }
      });
      
      this.hideLoading();
      this.showToast('Gmail disconnected', 'success');
      
      // Update UI
      this.connectionStatus.connected = false;
      this.connectionStatus.gmail = false;
      this.updateUI();
      
    } catch (error) {
      console.error('Failed to disconnect:', error);
      this.hideLoading();
      this.showToast('Failed to disconnect', 'error');
    }
  }

  async saveSettings() {
    try {
      const settings = {
        tone: document.getElementById('tone-select').value,
        length: document.getElementById('length-select').value,
        autoGenerate: document.getElementById('auto-generate').checked
      };
      
      this.showLoading('Saving settings...');
      
      const response = await this.sendMessage({
        action: 'saveSettings',
        data: settings
      });
      
      this.hideLoading();
      
      if (response.success) {
        this.showToast('Settings saved successfully!', 'success');
        this.settingsChanged = false;
      } else {
        throw new Error(response.error);
      }
      
    } catch (error) {
      console.error('Failed to save settings:', error);
      this.hideLoading();
      this.showToast('Failed to save settings', 'error');
    }
  }

  async syncEmails() {
    try {
      this.showLoading('Syncing emails...');
      
      const response = await this.sendMessage({
        action: 'syncEmails',
        data: { maxResults: 50 }
      });
      
      this.hideLoading();
      
      if (response.success) {
        this.showToast(
          `Synced ${response.processed_count} emails successfully!`,
          'success'
        );
        
        // Update last sync time
        await this.loadSettings();
        if (this.connectionStatus.user) {
          this.updateUserInfo(this.connectionStatus.user);
        }
      } else {
        throw new Error(response.error);
      }
      
    } catch (error) {
      console.error('Failed to sync emails:', error);
      this.hideLoading();
      this.showToast('Failed to sync emails', 'error');
    }
  }

  openGmail() {
    chrome.tabs.create({ url: 'https://mail.google.com' });
  }

  viewStats() {
    // For now, just open a simple stats page
    chrome.tabs.create({ 
      url: chrome.runtime.getURL('stats.html') 
    });
  }

  openHelp() {
    chrome.tabs.create({ 
      url: 'https://github.com/your-repo/ai-email-assistant#readme' 
    });
  }

  toggleSettings() {
    const settingsSection = document.getElementById('settings-section');
    const isVisible = settingsSection.style.display !== 'none';
    
    settingsSection.style.display = isVisible ? 'none' : 'block';
  }

  markSettingsChanged() {
    this.settingsChanged = true;
  }

  // Utility functions
  async sendMessage(message) {
    return new Promise((resolve) => {
      chrome.runtime.sendMessage(message, resolve);
    });
  }

  showLoading(text = 'Loading...') {
    const overlay = document.getElementById('loading-overlay');
    const loadingText = overlay.querySelector('.loading-text');
    loadingText.textContent = text;
    overlay.classList.remove('hidden');
  }

  hideLoading() {
    const overlay = document.getElementById('loading-overlay');
    overlay.classList.add('hidden');
  }

  showToast(message, type = 'info') {
    const toast = document.getElementById('toast');
    const toastIcon = document.getElementById('toast-icon');
    const toastMessage = document.getElementById('toast-message');
    
    // Set toast type
    toast.className = `toast ${type}`;
    
    // Set icon based on type
    const icons = {
      success: '✓',
      error: '✗',
      warning: '⚠',
      info: 'ℹ'
    };
    
    toastIcon.textContent = icons[type] || icons.info;
    toastMessage.textContent = message;
    
    // Show toast
    toast.classList.add('show');
    
    // Auto hide after 3 seconds
    setTimeout(() => {
      this.hideToast();
    }, 3000);
  }

  hideToast() {
    const toast = document.getElementById('toast');
    toast.classList.remove('show');
  }

  formatRelativeTime(date) {
    const now = new Date();
    const diff = now - date;
    const seconds = Math.floor(diff / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);
    
    if (days > 0) {
      return `${days} day${days > 1 ? 's' : ''} ago`;
    } else if (hours > 0) {
      return `${hours} hour${hours > 1 ? 's' : ''} ago`;
    } else if (minutes > 0) {
      return `${minutes} minute${minutes > 1 ? 's' : ''} ago`;
    } else {
      return 'Just now';
    }
  }
}

// Initialize popup when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  new PopupController();
});
