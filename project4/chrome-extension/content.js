// Content script for Gmail integration
class GmailEmailResponseGenerator {
  constructor() {
    this.apiUrl = 'http://localhost:8080/api/v1';
    this.isGenerating = false;
    this.init();
  }

  init() {
    console.log('AI Email Response Generator initialized');
    this.observeGmailChanges();
    this.addResponseButtons();
  }

  // Observe Gmail DOM changes for dynamic content
  observeGmailChanges() {
    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.addedNodes.length) {
          this.addResponseButtons();
        }
      });
    });

    observer.observe(document.body, {
      childList: true,
      subtree: true
    });
  }

  // Add response generation buttons to email threads
  addResponseButtons() {
    // Find email open views
    const emailViews = document.querySelectorAll('.nH.if, .nH.nn');
    
    emailViews.forEach(emailView => {
      // Check if we've already added buttons to this email
      if (emailView.querySelector('.ai-response-btn')) {
        return;
      }

      // Find the reply area
      const replyArea = this.findReplyArea(emailView);
      if (replyArea) {
        this.addResponseButton(replyArea, emailView);
      }
    });
  }

  // Find the appropriate reply area in Gmail's UI
  findReplyArea(emailView) {
    // Try different selectors for Gmail's reply area
    const selectors = [
      '.Am.Al.editable',
      '.LW-avf',
      '[g_editable="true"]',
      '.Am',
      '.aO7'
    ];

    for (const selector of selectors) {
      const element = emailView.querySelector(selector);
      if (element && element.getAttribute('contenteditable') === 'true') {
        return element;
      }
    }

    return null;
  }

  // Add response generation button
  addResponseButton(replyArea, emailView) {
    const buttonContainer = document.createElement('div');
    buttonContainer.className = 'ai-response-container';
    buttonContainer.style.cssText = `
      margin: 8px 0;
      display: flex;
      gap: 8px;
      align-items: center;
    `;

    // Generate response button
    const generateBtn = document.createElement('button');
    generateBtn.className = 'ai-response-btn';
    generateBtn.textContent = 'Generate AI Response';
    generateBtn.style.cssText = `
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      border: none;
      padding: 8px 16px;
      border-radius: 20px;
      cursor: pointer;
      font-size: 13px;
      font-weight: 500;
      transition: all 0.2s ease;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    `;

    generateBtn.addEventListener('mouseenter', () => {
      generateBtn.style.transform = 'translateY(-1px)';
      generateBtn.style.boxShadow = '0 4px 8px rgba(0,0,0,0.15)';
    });

    generateBtn.addEventListener('mouseleave', () => {
      generateBtn.style.transform = 'translateY(0)';
      generateBtn.style.boxShadow = '0 2px 4px rgba(0,0,0,0.1)';
    });

    generateBtn.addEventListener('click', () => {
      this.generateResponse(replyArea, emailView);
    });

    // Settings button
    const settingsBtn = document.createElement('button');
    settingsBtn.className = 'ai-settings-btn';
    settingsBtn.textContent = '⚙️';
    settingsBtn.title = 'AI Response Settings';
    settingsBtn.style.cssText = `
      background: #f1f3f4;
      color: #5f6368;
      border: none;
      padding: 8px;
      border-radius: 50%;
      cursor: pointer;
      font-size: 14px;
      width: 32px;
      height: 32px;
      transition: all 0.2s ease;
    `;

    settingsBtn.addEventListener('click', (e) => {
      e.stopPropagation();
      this.showSettings(replyArea);
    });

    buttonContainer.appendChild(generateBtn);
    buttonContainer.appendChild(settingsBtn);

    // Insert the button container before the reply area
    replyArea.parentNode.insertBefore(buttonContainer, replyArea);
  }

  // Generate AI response for the current email
  async generateResponse(replyArea, emailView) {
    if (this.isGenerating) {
      this.showNotification('Already generating a response...', 'warning');
      return;
    }

    this.isGenerating = true;
    const generateBtn = emailView.querySelector('.ai-response-btn');
    const originalText = generateBtn.textContent;

    try {
      generateBtn.textContent = 'Generating...';
      generateBtn.disabled = true;

      // Extract email content
      const emailContent = this.extractEmailContent(emailView);
      
      if (!emailContent.body) {
        throw new Error('Could not extract email content');
      }

      // Get user settings
      const settings = await this.getSettings();
      
      // Generate response via API
      const response = await this.callAPI('/gmail/generate-response', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${settings.token}`
        },
        body: JSON.stringify({
          email_body: emailContent.body,
          subject: emailContent.subject,
          sender: emailContent.sender,
          tone: settings.tone || 'professional',
          length: settings.length || 'medium'
        })
      });

      if (response.success) {
        // Insert the generated response
        this.insertResponse(replyArea, response.response);
        this.showNotification('Response generated successfully!', 'success');
      } else {
        throw new Error(response.error || 'Failed to generate response');
      }

    } catch (error) {
      console.error('Error generating response:', error);
      this.showNotification(`Error: ${error.message}`, 'error');
    } finally {
      generateBtn.textContent = originalText;
      generateBtn.disabled = false;
      this.isGenerating = false;
    }
  }

  // Extract email content from Gmail UI
  extractEmailContent(emailView) {
    const content = {
      subject: '',
      sender: '',
      body: ''
    };

    // Extract subject
    const subjectElement = emailView.querySelector('.hP, .ha');
    if (subjectElement) {
      content.subject = subjectElement.textContent.trim();
    }

    // Extract sender
    const senderElement = emailView.querySelector('.gD, .go');
    if (senderElement) {
      content.sender = senderElement.textContent.trim();
    }

    // Extract email body
    const bodyElement = emailView.querySelector('.a3s, .ii.gt');
    if (bodyElement) {
      // Remove quoted text and signatures
      const bodyClone = bodyElement.cloneNode(true);
      
      // Remove common Gmail quote elements
      const quotes = bodyClone.querySelectorAll('.gmail_quote, .gmail_extra');
      quotes.forEach(quote => quote.remove());
      
      content.body = bodyClone.textContent.trim();
    }

    return content;
  }

  // Insert generated response into reply area
  insertResponse(replyArea, response) {
    // Clear existing content
    replyArea.innerHTML = '';
    
    // Create a temporary div to parse HTML response
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = response;
    
    // Insert the response
    replyArea.innerHTML = tempDiv.textContent;
    
    // Trigger input event to let Gmail know content changed
    const event = new Event('input', { bubbles: true });
    replyArea.dispatchEvent(event);
  }

  // Call backend API
  async callAPI(endpoint, options) {
    try {
      const response = await fetch(this.apiUrl + endpoint, options);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      return await response.json();
    } catch (error) {
      console.error('API call failed:', error);
      throw error;
    }
  }

  // Get extension settings
  async getSettings() {
    return new Promise((resolve) => {
      chrome.storage.sync.get({
        token: '',
        tone: 'professional',
        length: 'medium',
        autoGenerate: false
      }, resolve);
    });
  }

  // Show settings modal
  showSettings(replyArea) {
    // Remove existing settings modal
    const existingModal = document.querySelector('.ai-settings-modal');
    if (existingModal) {
      existingModal.remove();
    }

    const modal = document.createElement('div');
    modal.className = 'ai-settings-modal';
    modal.style.cssText = `
      position: fixed;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      background: white;
      padding: 24px;
      border-radius: 12px;
      box-shadow: 0 20px 40px rgba(0,0,0,0.15);
      z-index: 10000;
      min-width: 400px;
      font-family: 'Google Sans', Roboto, Arial, sans-serif;
    `;

    modal.innerHTML = `
      <h3 style="margin: 0 0 20px 0; color: #202124;">AI Response Settings</h3>
      
      <div style="margin-bottom: 16px;">
        <label style="display: block; margin-bottom: 8px; color: #5f6368; font-size: 14px;">Response Tone</label>
        <select id="ai-tone" style="width: 100%; padding: 8px; border: 1px solid #dadce0; border-radius: 4px;">
          <option value="professional">Professional</option>
          <option value="casual">Casual</option>
          <option value="friendly">Friendly</option>
          <option value="formal">Formal</option>
        </select>
      </div>
      
      <div style="margin-bottom: 16px;">
        <label style="display: block; margin-bottom: 8px; color: #5f6368; font-size: 14px;">Response Length</label>
        <select id="ai-length" style="width: 100%; padding: 8px; border: 1px solid #dadce0; border-radius: 4px;">
          <option value="short">Short</option>
          <option value="medium">Medium</option>
          <option value="detailed">Detailed</option>
        </select>
      </div>
      
      <div style="margin-bottom: 20px;">
        <label style="display: flex; align-items: center; cursor: pointer;">
          <input type="checkbox" id="auto-generate" style="margin-right: 8px;">
          <span style="color: #5f6368; font-size: 14px;">Auto-generate on email open</span>
        </label>
      </div>
      
      <div style="display: flex; gap: 12px; justify-content: flex-end;">
        <button id="cancel-settings" style="padding: 8px 16px; border: 1px solid #dadce0; background: white; border-radius: 4px; cursor: pointer;">Cancel</button>
        <button id="save-settings" style="padding: 8px 16px; background: #1a73e8; color: white; border: none; border-radius: 4px; cursor: pointer;">Save</button>
      </div>
    `;

    // Add backdrop
    const backdrop = document.createElement('div');
    backdrop.className = 'ai-settings-backdrop';
    backdrop.style.cssText = `
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0,0,0,0.5);
      z-index: 9999;
    `;

    document.body.appendChild(backdrop);
    document.body.appendChild(modal);

    // Load current settings
    this.getSettings().then(settings => {
      document.getElementById('ai-tone').value = settings.tone;
      document.getElementById('ai-length').value = settings.length;
      document.getElementById('auto-generate').checked = settings.autoGenerate;
    });

    // Handle save
    document.getElementById('save-settings').addEventListener('click', () => {
      const newSettings = {
        tone: document.getElementById('ai-tone').value,
        length: document.getElementById('ai-length').value,
        autoGenerate: document.getElementById('auto-generate').checked
      };

      chrome.storage.sync.set(newSettings, () => {
        this.showNotification('Settings saved!', 'success');
        modal.remove();
        backdrop.remove();
      });
    });

    // Handle cancel
    document.getElementById('cancel-settings').addEventListener('click', () => {
      modal.remove();
      backdrop.remove();
    });

    // Close on backdrop click
    backdrop.addEventListener('click', () => {
      modal.remove();
      backdrop.remove();
    });
  }

  // Show notification
  showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = 'ai-notification';
    
    const colors = {
      success: '#34a853',
      error: '#ea4335',
      warning: '#fbbc04',
      info: '#1a73e8'
    };

    notification.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      background: ${colors[type]};
      color: white;
      padding: 12px 20px;
      border-radius: 8px;
      box-shadow: 0 4px 12px rgba(0,0,0,0.15);
      z-index: 10001;
      font-family: 'Google Sans', Roboto, Arial, sans-serif;
      font-size: 14px;
      max-width: 300px;
      animation: slideIn 0.3s ease;
    `;

    notification.textContent = message;
    document.body.appendChild(notification);

    // Auto remove after 3 seconds
    setTimeout(() => {
      notification.style.animation = 'slideOut 0.3s ease';
      setTimeout(() => notification.remove(), 300);
    }, 3000);
  }
}

// Add CSS animations
const style = document.createElement('style');
style.textContent = `
  @keyframes slideIn {
    from {
      transform: translateX(100%);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }

  @keyframes slideOut {
    from {
      transform: translateX(0);
      opacity: 1;
    }
    to {
      transform: translateX(100%);
      opacity: 0;
    }
  }
`;
document.head.appendChild(style);

// Initialize the extension
const aiResponder = new GmailEmailResponseGenerator();
