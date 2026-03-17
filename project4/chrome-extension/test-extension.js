/**
 * Chrome Extension Testing Script
 * 
 * This script helps test the Chrome extension functionality
 * Run this in the browser console to verify extension is working correctly
 */

class ExtensionTester {
  constructor() {
    this.testResults = [];
    this.apiUrl = 'http://localhost:8080/api/v1';
  }

  // Run all tests
  async runAllTests() {
    console.log('🧪 Starting Chrome Extension Tests...\n');

    await this.testManifest();
    await this.testContentScript();
    await this.testBackgroundScript();
    await this.testAPIConnection();
    await this.testGmailIntegration();

    this.displayResults();
  }

  // Test manifest.json
  async testManifest() {
    const testName = 'Manifest Validation';
    try {
      // Check if manifest is properly loaded
      if (typeof chrome !== 'undefined' && chrome.runtime) {
        const manifest = chrome.runtime.getManifest();
        
        const requiredFields = ['name', 'version', 'manifest_version', 'permissions'];
        const missingFields = requiredFields.filter(field => !manifest[field]);
        
        if (missingFields.length === 0) {
          this.addResult(testName, true, 'Manifest loaded successfully');
        } else {
          this.addResult(testName, false, `Missing manifest fields: ${missingFields.join(', ')}`);
        }
      } else {
        this.addResult(testName, false, 'Chrome APIs not available');
      }
    } catch (error) {
      this.addResult(testName, false, `Error: ${error.message}`);
    }
  }

  // Test content script injection
  async testContentScript() {
    const testName = 'Content Script';
    try {
      // Check if we're on Gmail
      if (!window.location.hostname.includes('mail.google.com')) {
        this.addResult(testName, 'warning', 'Not on Gmail - cannot test content script');
        return;
      }

      // Look for extension elements
      const extensionElements = document.querySelectorAll('.ai-response-btn, .ai-response-container');
      if (extensionElements.length > 0) {
        this.addResult(testName, true, `Found ${extensionElements.length} extension elements`);
      } else {
        this.addResult(testName, false, 'No extension elements found on page');
      }

      // Check if GmailEmailResponseGenerator class exists
      if (typeof GmailEmailResponseGenerator !== 'undefined') {
        this.addResult(testName, true, 'Content script loaded successfully');
      } else {
        this.addResult(testName, false, 'Content script not loaded');
      }
    } catch (error) {
      this.addResult(testName, false, `Error: ${error.message}`);
    }
  }

  // Test background script
  async testBackgroundScript() {
    const testName = 'Background Script';
    try {
      // Send a test message to background script
      const response = await new Promise((resolve) => {
        chrome.runtime.sendMessage({ action: 'getSettings' }, resolve);
      });

      if (response && response.success) {
        this.addResult(testName, true, 'Background script responding');
      } else {
        this.addResult(testName, false, 'Background script not responding');
      }
    } catch (error) {
      this.addResult(testName, false, `Error: ${error.message}`);
    }
  }

  // Test API connection
  async testAPIConnection() {
    const testName = 'API Connection';
    try {
      const response = await fetch(`${this.apiUrl}/health`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (response.ok) {
        this.addResult(testName, true, 'API connection successful');
      } else {
        this.addResult(testName, false, `API returned status: ${response.status}`);
      }
    } catch (error) {
      this.addResult(testName, false, `API connection failed: ${error.message}`);
    }
  }

  // Test Gmail integration
  async testGmailIntegration() {
    const testName = 'Gmail Integration';
    try {
      if (!window.location.hostname.includes('mail.google.com')) {
        this.addResult(testName, 'warning', 'Not on Gmail - cannot test integration');
        return;
      }

      // Look for Gmail email elements
      const emailElements = document.querySelectorAll('.nH.if, .nH.nn');
      if (emailElements.length > 0) {
        this.addResult(testName, true, `Found ${emailElements.length} email views`);
        
        // Check if reply areas are found
        const replyAreas = document.querySelectorAll('[g_editable="true"], .Am.Al.editable');
        if (replyAreas.length > 0) {
          this.addResult(testName, true, `Found ${replyAreas.length} reply areas`);
        } else {
          this.addResult(testName, 'warning', 'No reply areas found (open an email to test)');
        }
      } else {
        this.addResult(testName, 'warning', 'No email views found (open Gmail to test)');
      }
    } catch (error) {
      this.addResult(testName, false, `Error: ${error.message}`);
    }
  }

  // Add test result
  addResult(testName, status, message) {
    this.testResults.push({
      test: testName,
      status: status, // true, false, or 'warning'
      message: message
    });
  }

  // Display all test results
  displayResults() {
    console.log('\n📊 Test Results:');
    console.log('================');

    let passed = 0;
    let failed = 0;
    let warnings = 0;

    this.testResults.forEach(result => {
      const status = result.status === true ? '✅' : 
                    result.status === false ? '❌' : '⚠️';
      
      console.log(`${status} ${result.test}: ${result.message}`);
      
      if (result.status === true) passed++;
      else if (result.status === false) failed++;
      else warnings++;
    });

    console.log('\n📈 Summary:');
    console.log(`✅ Passed: ${passed}`);
    console.log(`❌ Failed: ${failed}`);
    console.log(`⚠️  Warnings: ${warnings}`);
    console.log(`📊 Total: ${this.testResults.length}`);

    if (failed === 0) {
      console.log('\n🎉 All critical tests passed! Extension should work correctly.');
    } else {
      console.log('\n⚠️  Some tests failed. Please check the issues above.');
    }
  }

  // Test specific functionality
  async testResponseGeneration() {
    console.log('🔄 Testing response generation...');
    
    const testData = {
      email_body: "Hello, I wanted to follow up on our conversation from last week.",
      subject: "Follow up",
      sender: "test@example.com",
      tone: "professional",
      length: "medium"
    };

    try {
      const response = await fetch(`${this.apiUrl}/gmail/generate-response`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer test-token'
        },
        body: JSON.stringify(testData)
      });

      if (response.ok) {
        const result = await response.json();
        console.log('✅ Response generation successful:', result);
        return true;
      } else {
        console.log('❌ Response generation failed:', response.status);
        return false;
      }
    } catch (error) {
      console.log('❌ Response generation error:', error.message);
      return false;
    }
  }

  // Test email content extraction
  testEmailExtraction() {
    console.log('🔍 Testing email content extraction...');
    
    if (!window.location.hostname.includes('mail.google.com')) {
      console.log('⚠️  Not on Gmail - cannot test email extraction');
      return false;
    }

    try {
      // Simulate email content extraction
      const subject = document.querySelector('.hP, .ha');
      const sender = document.querySelector('.gD, .go');
      const body = document.querySelector('.a3s, .ii.gt');

      const extracted = {
        subject: subject ? subject.textContent.trim() : 'Not found',
        sender: sender ? sender.textContent.trim() : 'Not found',
        body: body ? body.textContent.substring(0, 100) + '...' : 'Not found'
      };

      console.log('📧 Extracted email content:', extracted);
      
      const allFound = subject && sender && body;
      if (allFound) {
        console.log('✅ Email extraction successful');
      } else {
        console.log('⚠️  Some email elements not found (open an email to test)');
      }
      
      return allFound;
    } catch (error) {
      console.log('❌ Email extraction error:', error.message);
      return false;
    }
  }
}

// Auto-run tests if on Gmail
if (window.location.hostname.includes('mail.google.com')) {
  console.log('📧 Detected Gmail - running extension tests...');
  const tester = new ExtensionTester();
  
  // Run basic tests
  setTimeout(() => {
    tester.runAllTests();
  }, 2000); // Wait for page to load
  
  // Make tester available globally for manual testing
  window.extensionTester = tester;
  
  console.log('💡 Available commands:');
  console.log('- extensionTester.runAllTests() - Run all tests');
  console.log('- extensionTester.testResponseGeneration() - Test API response generation');
  console.log('- extensionTester.testEmailExtraction() - Test email content extraction');
} else {
  console.log('ℹ️  Extension testing script loaded. Navigate to Gmail to run tests.');
  console.log('📧 Go to https://mail.google.com to test the extension');
}

// Export for use in other scripts
if (typeof module !== 'undefined' && module.exports) {
  module.exports = ExtensionTester;
}
