// Simplified popup script for debugging
console.log('Popup script loaded');

// Wait for DOM to be loaded
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded, initializing popup...');
    
    try {
        initializePopup();
    } catch (error) {
        console.error('Error initializing popup:', error);
        showError('Failed to load extension');
    }
});

function initializePopup() {
    console.log('Initializing popup controller...');
    
    // Show basic content immediately
    showBasicContent();
    
    // Initialize controller
    const controller = new PopupController();
    console.log('Popup controller initialized');
}

function showBasicContent() {
    // Ensure basic content is visible
    const container = document.querySelector('.popup-container');
    if (container) {
        container.style.display = 'block';
        console.log('Container shown');
    }
    
    // Update status text
    const statusText = document.querySelector('.status-text');
    if (statusText) {
        statusText.textContent = 'Loading...';
    }
}

function showError(message) {
    const container = document.querySelector('.popup-container');
    if (container) {
        container.innerHTML = `
            <div style="padding: 20px; text-align: center;">
                <h3>Error</h3>
                <p>${message}</p>
                <button onclick="location.reload()">Reload</button>
            </div>
        `;
    }
}

class PopupController {
    constructor() {
        console.log('PopupController constructor');
        this.settings = {};
        this.connectionStatus = {
            connected: false,
            gmail: false,
            backend: false,
            user: null
        };
        
        this.init();
    }

    async init() {
        console.log('PopupController init');
        try {
            await this.loadSettings();
            await this.checkConnectionStatus();
            this.setupEventListeners();
            this.updateUI();
            console.log('PopupController initialized successfully');
        } catch (error) {
            console.error('PopupController init error:', error);
            this.showToast('Initialization failed', 'error');
        }
    }

    setupEventListeners() {
        console.log('Setting up event listeners');
        
        const connectBtn = document.getElementById('connect-btn');
        if (connectBtn) {
            connectBtn.addEventListener('click', () => this.connectGmail());
            console.log('Connect button listener added');
        }

        const saveSettingsBtn = document.getElementById('save-settings-btn');
        if (saveSettingsBtn) {
            saveSettingsBtn.addEventListener('click', () => this.saveSettings());
        }

        const syncEmailsBtn = document.getElementById('sync-emails-btn');
        if (syncEmailsBtn) {
            syncEmailsBtn.addEventListener('click', () => this.syncEmails());
        }

        const openGmailBtn = document.getElementById('open-gmail-btn');
        if (openGmailBtn) {
            openGmailBtn.addEventListener('click', () => this.openGmail());
        }

        const helpLink = document.getElementById('help-link');
        if (helpLink) {
            helpLink.addEventListener('click', (e) => {
                e.preventDefault();
                this.openHelp();
            });
        }

        const settingsLink = document.getElementById('settings-link');
        if (settingsLink) {
            settingsLink.addEventListener('click', (e) => {
                e.preventDefault();
                this.toggleSettings();
            });
        }
    }

    async loadSettings() {
        try {
            const response = await this.sendMessage({ action: 'getSettings' });
            if (response && response.success) {
                this.settings = response.settings;
                console.log('Settings loaded:', this.settings);
                
                // Update form controls
                const toneSelect = document.getElementById('tone-select');
                if (toneSelect) {
                    toneSelect.value = this.settings.tone || 'professional';
                }
                
                const lengthSelect = document.getElementById('length-select');
                if (lengthSelect) {
                    lengthSelect.value = this.settings.length || 'medium';
                }
                
                const autoGenerate = document.getElementById('auto-generate');
                if (autoGenerate) {
                    autoGenerate.checked = this.settings.autoGenerate || false;
                }
            }
        } catch (error) {
            console.error('Failed to load settings:', error);
        }
    }

    async checkConnectionStatus() {
        try {
            console.log('Checking connection status...');
            
            // Check if we have a token
            const settingsResponse = await this.sendMessage({ action: 'getSettings' });
            
            if (settingsResponse && settingsResponse.success && settingsResponse.settings.token) {
                console.log('Token found, validating...');
                
                // For now, assume connected if token exists
                this.connectionStatus = {
                    connected: true,
                    gmail: true,
                    backend: true,
                    user: { email: 'user@example.com' } // Mock user
                };
            } else {
                console.log('No token found');
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
        console.log('Updating UI with status:', this.connectionStatus);
        
        // Update status indicator
        const statusIndicator = document.getElementById('status-indicator');
        const statusDot = statusIndicator ? statusIndicator.querySelector('.status-dot') : null;
        const statusText = statusIndicator ? statusIndicator.querySelector('.status-text') : null;
        
        if (statusDot && statusText) {
            if (this.connectionStatus.connected) {
                statusDot.classList.remove('disconnected');
                statusText.textContent = 'Connected';
            } else {
                statusDot.classList.add('disconnected');
                statusText.textContent = 'Disconnected';
            }
        }

        // Update connection status card
        const gmailStatus = document.getElementById('gmail-status');
        const backendStatus = document.getElementById('backend-status');
        const connectBtn = document.getElementById('connect-btn');
        
        if (gmailStatus) {
            gmailStatus.textContent = this.connectionStatus.gmail ? 'Connected' : 'Not connected';
            gmailStatus.className = this.connectionStatus.gmail ? 'value connected' : 'value disconnected';
        }
        
        if (backendStatus) {
            backendStatus.textContent = this.connectionStatus.backend ? 'Online' : 'Offline';
            backendStatus.className = this.connectionStatus.backend ? 'value connected' : 'value disconnected';
        }
        
        if (connectBtn) {
            connectBtn.textContent = this.connectionStatus.connected ? 'Disconnect' : 'Connect Gmail';
        }

        // Show/hide sections based on connection status
        const settingsSection = document.getElementById('settings-section');
        const quickActions = document.getElementById('quick-actions');
        const userInfo = document.getElementById('user-info');
        const disconnectLink = document.getElementById('disconnect-link');
        
        if (this.connectionStatus.connected) {
            if (settingsSection) settingsSection.style.display = 'block';
            if (quickActions) quickActions.style.display = 'block';
            if (disconnectLink) disconnectLink.style.display = 'inline';
            if (userInfo) {
                userInfo.style.display = 'block';
                this.updateUserInfo(this.connectionStatus.user);
            }
        } else {
            if (settingsSection) settingsSection.style.display = 'none';
            if (quickActions) quickActions.style.display = 'none';
            if (userInfo) userInfo.style.display = 'none';
            if (disconnectLink) disconnectLink.style.display = 'none';
        }
    }

    updateUserInfo(user) {
        if (!user) return;
        
        const userInitial = document.getElementById('user-initial');
        const userEmail = document.getElementById('user-email');
        const lastSyncTime = document.getElementById('last-sync-time');
        
        if (userInitial && user.email) {
            userInitial.textContent = user.email.charAt(0).toUpperCase();
        }
        
        if (userEmail) {
            userEmail.textContent = user.email;
        }
        
        if (lastSyncTime) {
            lastSyncTime.textContent = this.settings.lastSync ? 
                this.formatRelativeTime(new Date(this.settings.lastSync)) : 'Never';
        }
    }

    async connectGmail() {
        if (this.connectionStatus.connected) {
            await this.disconnect();
            return;
        }
        
        try {
            this.showToast('Opening Gmail authentication...', 'info');
            
            // Open Gmail OAuth flow
            const authUrl = 'http://localhost:8080/api/v1/auth/gmail';
            chrome.tabs.create({ url: authUrl });
            
        } catch (error) {
            console.error('Failed to connect Gmail:', error);
            this.showToast('Failed to connect Gmail', 'error');
        }
    }

    async disconnect() {
        try {
            await this.sendMessage({
                action: 'saveSettings',
                data: { token: '', refreshToken: '' }
            });
            
            this.showToast('Gmail disconnected', 'success');
            
            this.connectionStatus.connected = false;
            this.connectionStatus.gmail = false;
            this.updateUI();
            
        } catch (error) {
            console.error('Failed to disconnect:', error);
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
            
            const response = await this.sendMessage({
                action: 'saveSettings',
                data: settings
            });
            
            if (response && response.success) {
                this.showToast('Settings saved successfully!', 'success');
            } else {
                throw new Error(response ? response.error : 'Unknown error');
            }
            
        } catch (error) {
            console.error('Failed to save settings:', error);
            this.showToast('Failed to save settings', 'error');
        }
    }

    async syncEmails() {
        try {
            this.showToast('Syncing emails...', 'info');
            
            const response = await this.sendMessage({
                action: 'syncEmails',
                data: { maxResults: 50 }
            });
            
            if (response && response.success) {
                this.showToast(`Synced ${response.processed_count} emails successfully!`, 'success');
            } else {
                throw new Error(response ? response.error : 'Unknown error');
            }
            
        } catch (error) {
            console.error('Failed to sync emails:', error);
            this.showToast('Failed to sync emails', 'error');
        }
    }

    openGmail() {
        chrome.tabs.create({ url: 'https://mail.google.com' });
    }

    openHelp() {
        chrome.tabs.create({ 
            url: chrome.runtime.getURL('welcome.html') 
        });
    }

    toggleSettings() {
        const settingsSection = document.getElementById('settings-section');
        if (settingsSection) {
            const isVisible = settingsSection.style.display !== 'none';
            settingsSection.style.display = isVisible ? 'none' : 'block';
        }
    }

    async sendMessage(message) {
        return new Promise((resolve) => {
            chrome.runtime.sendMessage(message, (response) => {
                if (chrome.runtime.lastError) {
                    console.error('Chrome runtime error:', chrome.runtime.lastError);
                    resolve({ success: false, error: chrome.runtime.lastError.message });
                } else {
                    resolve(response);
                }
            });
        });
    }

    showToast(message, type = 'info') {
        const toast = document.getElementById('toast');
        const toastIcon = document.getElementById('toast-icon');
        const toastMessage = document.getElementById('toast-message');
        
        if (!toast || !toastIcon || !toastMessage) {
            console.log('Toast elements not found, showing alert instead:', message);
            alert(message);
            return;
        }
        
        toast.className = `toast ${type}`;
        
        const icons = {
            success: '✓',
            error: '✗',
            warning: '⚠',
            info: 'ℹ'
        };
        
        toastIcon.textContent = icons[type] || icons.info;
        toastMessage.textContent = message;
        
        toast.classList.add('show');
        
        setTimeout(() => {
            toast.classList.remove('show');
        }, 3000);
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
