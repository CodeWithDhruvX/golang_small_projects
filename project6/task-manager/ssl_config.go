package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// SSLConfig holds SSL/TLS configuration
type SSLConfig struct {
	CertFile    string
	KeyFile     string
	CAFile      string
	UseHTTPS    bool
	Port        int
	HTTPSPort   int
	RedirectToHTTPS bool
}

// SSLManager manages SSL/TLS certificates and configuration
type SSLManager struct {
	config *SSLConfig
	logger *Logger
}

// NewSSLManager creates a new SSL manager
func NewSSLManager(config *SSLConfig, logger *Logger) *SSLManager {
	return &SSLManager{
		config: config,
		logger: logger,
	}
}

// GenerateSelfSignedCertificate generates a self-signed certificate for development
func (sm *SSLManager) GenerateSelfSignedCertificate() error {
	sm.logger.Info("Generating self-signed certificate", nil)

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Secure API Development"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 1 year
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost", "127.0.0.1"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Save certificate
	certFile, err := os.Create(sm.config.CertFile)
	if err != nil {
		return fmt.Errorf("failed to create cert file: %w", err)
	}
	defer certFile.Close()

	if err := pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}); err != nil {
		return fmt.Errorf("failed to encode certificate: %w", err)
	}

	// Save private key
	keyFile, err := os.Create(sm.config.KeyFile)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyFile.Close()

	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	if err := pem.Encode(keyFile, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyDER,
	}); err != nil {
		return fmt.Errorf("failed to encode private key: %w", err)
	}

	sm.logger.Info("Self-signed certificate generated successfully", map[string]interface{}{
		"cert_file": sm.config.CertFile,
		"key_file":  sm.config.KeyFile,
	})

	return nil
}

// LoadTLSCertificates loads TLS certificates
func (sm *SSLManager) LoadTLSCertificates() (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(sm.config.CertFile, sm.config.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificates: %w", err)
	}

	sm.logger.Info("TLS certificates loaded successfully", nil)
	return &cert, nil
}

// CreateTLSConfig creates a TLS configuration
func (sm *SSLManager) CreateTLSConfig() (*tls.Config, error) {
	cert, err := sm.LoadTLSCertificates()
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		PreferServerCipherSuites: true,
	}

	sm.logger.Info("TLS configuration created", map[string]interface{}{
		"min_version": tlsConfig.MinVersion,
		"cipher_suites": len(tlsConfig.CipherSuites),
	})

	return tlsConfig, nil
}

// ValidateCertificates validates the SSL certificates
func (sm *SSLManager) ValidateCertificates() error {
	// Check if certificate files exist
	if _, err := os.Stat(sm.config.CertFile); os.IsNotExist(err) {
		return fmt.Errorf("certificate file not found: %s", sm.config.CertFile)
	}

	if _, err := os.Stat(sm.config.KeyFile); os.IsNotExist(err) {
		return fmt.Errorf("key file not found: %s", sm.config.KeyFile)
	}

	// Try to load and parse certificate
	cert, err := sm.LoadTLSCertificates()
	if err != nil {
		return fmt.Errorf("failed to load certificates: %w", err)
	}

	// Parse certificate to check expiration
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Check expiration
	if time.Now().After(x509Cert.NotAfter) {
		return fmt.Errorf("certificate has expired on %s", x509Cert.NotAfter.Format(time.RFC3339))
	}

	// Check if certificate is about to expire (within 30 days)
	if time.Now().Add(30 * 24 * time.Hour).After(x509Cert.NotAfter) {
		sm.logger.Warn("Certificate is about to expire", map[string]interface{}{
			"expires_on": x509Cert.NotAfter.Format(time.RFC3339),
			"days_left":  int(x509Cert.NotAfter.Sub(time.Now()).Hours() / 24),
		})
	}

	sm.logger.Info("Certificate validation successful", map[string]interface{}{
		"expires_on": x509Cert.NotAfter.Format(time.RFC3339),
		"issuer":     x509Cert.Issuer.CommonName,
		"subject":    x509Cert.Subject.CommonName,
	})

	return nil
}

// SetupSSL sets up SSL/TLS for the server
func (sm *SSLManager) SetupSSL() error {
	if !sm.config.UseHTTPS {
		sm.logger.Info("HTTPS disabled, running in HTTP mode", nil)
		return nil
	}

	// Check if certificates exist, if not generate them
	if err := sm.ValidateCertificates(); err != nil {
		if os.IsNotExist(err) || err.Error() == "certificate file not found" || err.Error() == "key file not found" {
			sm.logger.Info("Certificates not found, generating self-signed certificates", nil)
			if err := sm.GenerateSelfSignedCertificate(); err != nil {
				return fmt.Errorf("failed to generate certificates: %w", err)
			}
		} else {
			return fmt.Errorf("certificate validation failed: %w", err)
		}
	}

	return nil
}

// GetDefaultSSLConfig returns default SSL configuration
func GetDefaultSSLConfig() *SSLConfig {
	return &SSLConfig{
		CertFile:         "server.crt",
		KeyFile:          "server.key",
		CAFile:           "ca.crt",
		UseHTTPS:         true,
		Port:             8080,
		HTTPSPort:        8443,
		RedirectToHTTPS:  true,
	}
}

// HTTPSToHTTPSRedirectMiddleware redirects HTTP to HTTPS
func (sm *SSLManager) HTTPSToHTTPSRedirectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if sm.config.RedirectToHTTPS && c.Request.TLS == nil {
			httpsURL := fmt.Sprintf("https://%s:%d%s", 
				c.Request.Host, 
				sm.config.HTTPSPort, 
				c.Request.URL.Path)
			
			c.Redirect(http.StatusMovedPermanently, httpsURL)
			c.Abort()
			return
		}
		c.Next()
	}
}

// LogSSLConnection logs SSL/TLS connection information
func (sm *SSLManager) LogSSLConnection(state *tls.ConnectionState) {
	if state != nil {
		sm.logger.Info("TLS connection established", map[string]interface{}{
			"version":           tlsVersionToString(state.Version),
			"cipher_suite":      tls.CipherSuiteName(state.CipherSuite),
			"server_name":       state.ServerName,
			"negotiated_protocol": string(state.NegotiatedProtocol),
		})
	}
}

func tlsVersionToString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("Unknown (%d)", version)
	}
}
