package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Sidecar Pattern

// Main Application
type MainApplication struct {
	name    string
	port    int
	running bool
}

func NewMainApplication(name string, port int) *MainApplication {
	return &MainApplication{
		name:    name,
		port:    port,
		running: false,
	}
}

func (ma *MainApplication) Start() error {
	fmt.Printf("Main Application %s: Starting on port %d\n", ma.name, ma.port)
	ma.running = true
	
	// Simulate main application work
	go func() {
		for ma.running {
			time.Sleep(1 * time.Second)
			fmt.Printf("Main Application %s: Processing business logic...\n", ma.name)
		}
	}()
	
	return nil
}

func (ma *MainApplication) Stop() error {
	fmt.Printf("Main Application %s: Stopping\n", ma.name)
	ma.running = false
	return nil
}

func (ma *MainApplication) GetHealth() HealthStatus {
	return HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Details:   fmt.Sprintf("Main application %s is running", ma.name),
	}
}

// Sidecar Container
type Sidecar struct {
	name           string
	mainApp        *MainApplication
	healthChecker  *HealthChecker
	logger         *Logger
	metrics        *MetricsCollector
	configWatcher  *ConfigWatcher
	proxy          *Proxy
	running        bool
}

func NewSidecar(name string, mainApp *MainApplication) *Sidecar {
	return &Sidecar{
		name:          name,
		mainApp:       mainApp,
		healthChecker: NewHealthChecker(mainApp),
		logger:        NewLogger(),
		metrics:       NewMetricsCollector(),
		configWatcher: NewConfigWatcher(),
		proxy:         NewProxy(),
		running:       false,
	}
}

func (sc *Sidecar) Start() error {
	fmt.Printf("Sidecar %s: Starting...\n", sc.name)
	sc.running = true
	
	// Start all sidecar components
	sc.healthChecker.Start()
	sc.logger.Start()
	sc.metrics.Start()
	sc.configWatcher.Start()
	sc.proxy.Start()
	
	// Start monitoring main application
	go sc.monitorMainApp()
	
	// Start HTTP server for sidecar endpoints
	go sc.startHTTPServer()
	
	return nil
}

func (sc *Sidecar) Stop() error {
	fmt.Printf("Sidecar %s: Stopping...\n", sc.name)
	sc.running = false
	
	// Stop all components
	sc.healthChecker.Stop()
	sc.logger.Stop()
	sc.metrics.Stop()
	sc.configWatcher.Stop()
	sc.proxy.Stop()
	
	return nil
}

func (sc *Sidecar) monitorMainApp() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for sc.running {
		select {
		case <-ticker.C:
			health := sc.mainApp.GetHealth()
			sc.logger.Log(fmt.Sprintf("Health check: %s", health.Status))
			sc.metrics.RecordHealthCheck(health.Status == "healthy")
		}
	}
}

func (sc *Sidecar) startHTTPServer() {
	http.HandleFunc("/health", sc.handleHealth)
	http.HandleFunc("/metrics", sc.handleMetrics)
	http.HandleFunc("/logs", sc.handleLogs)
	http.HandleFunc("/config", sc.handleConfig)
	http.HandleFunc("/proxy", sc.handleProxy)
	
	fmt.Printf("Sidecar %s: HTTP server listening on port 8081\n", sc.name)
	http.ListenAndServe(":8081", nil)
}

func (sc *Sidecar) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := sc.mainApp.GetHealth()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "%s", "details": "%s", "timestamp": "%s"}`,
		health.Status, health.Details, health.Timestamp.Format(time.RFC3339))
}

func (sc *Sidecar) handleMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := sc.metrics.GetMetrics()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"requests": %d, "errors": %d, "uptime": "%s"}`,
		metrics.Requests, metrics.Errors, metrics.Uptime)
}

func (sc *Sidecar) handleLogs(w http.ResponseWriter, r *http.Request) {
	logs := sc.logger.GetRecentLogs(10)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"logs": [%s]}`, strings.Join(logs, ","))
}

func (sc *Sidecar) handleConfig(w http.ResponseWriter, r *http.Request) {
	config := sc.configWatcher.GetCurrentConfig()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"config": %s}`, config)
}

func (sc *Sidecar) handleProxy(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sc.proxy.HandleRequest(w, r)
	}
}

// Health Checker
type HealthChecker struct {
	mainApp *MainApplication
	running bool
}

func NewHealthChecker(mainApp *MainApplication) *HealthChecker {
	return &HealthChecker{mainApp: mainApp}
}

func (hc *HealthChecker) Start() {
	fmt.Println("Health Checker: Starting...")
	hc.running = true
	
	go func() {
		for hc.running {
			time.Sleep(10 * time.Second)
			health := hc.mainApp.GetHealth()
			fmt.Printf("Health Checker: %s - %s\n", health.Status, health.Details)
		}
	}()
}

func (hc *HealthChecker) Stop() {
	fmt.Println("Health Checker: Stopping...")
	hc.running = false
}

// Logger
type Logger struct {
	logs   []string
	running bool
}

func NewLogger() *Logger {
	return &Logger{
		logs: make([]string, 0),
	}
}

func (l *Logger) Start() {
	fmt.Println("Logger: Starting...")
	l.running = true
}

func (l *Logger) Stop() {
	fmt.Println("Logger: Stopping...")
	l.running = false
}

func (l *Logger) Log(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s", timestamp, message)
	
	l.logs = append(l.logs, logEntry)
	
	// Keep only last 100 logs
	if len(l.logs) > 100 {
		l.logs = l.logs[1:]
	}
	
	fmt.Printf("Logger: %s\n", message)
}

func (l *Logger) GetRecentLogs(count int) []string {
	if count > len(l.logs) {
		count = len(l.logs)
	}
	
	return l.logs[len(l.logs)-count:]
}

// Metrics Collector
type MetricsCollector struct {
	requests int
	errors   int
	startTime time.Time
	running  bool
}

type Metrics struct {
	Requests int
	Errors   int
	Uptime   string
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		startTime: time.Now(),
	}
}

func (mc *MetricsCollector) Start() {
	fmt.Println("Metrics Collector: Starting...")
	mc.running = true
}

func (mc *MetricsCollector) Stop() {
	fmt.Println("Metrics Collector: Stopping...")
	mc.running = false
}

func (mc *MetricsCollector) RecordRequest() {
	mc.requests++
}

func (mc *MetricsCollector) RecordError() {
	mc.errors++
}

func (mc *MetricsCollector) RecordHealthCheck(healthy bool) {
	mc.RecordRequest()
	if !healthy {
		mc.RecordError()
	}
}

func (mc *MetricsCollector) GetMetrics() Metrics {
	uptime := time.Since(mc.startTime)
	return Metrics{
		Requests: mc.requests,
		Errors:   mc.errors,
		Uptime:   uptime.String(),
	}
}

// Config Watcher
type ConfigWatcher struct {
	config  map[string]string
	running bool
}

func NewConfigWatcher() *ConfigWatcher {
	return &ConfigWatcher{
		config: map[string]string{
			"log_level":    "info",
			"max_requests": "1000",
			"timeout":      "30s",
		},
	}
}

func (cw *ConfigWatcher) Start() {
	fmt.Println("Config Watcher: Starting...")
	cw.running = true
	
	go func() {
		for cw.running {
			time.Sleep(30 * time.Second)
			cw.checkForConfigChanges()
		}
	}()
}

func (cw *ConfigWatcher) Stop() {
	fmt.Println("Config Watcher: Stopping...")
	cw.running = false
}

func (cw *ConfigWatcher) checkForConfigChanges() {
	// Simulate config change detection
	if time.Now().Unix()%60 == 0 { // Every minute
		cw.config["log_level"] = "debug"
		fmt.Println("Config Watcher: Configuration updated")
	}
}

func (cw *ConfigWatcher) GetCurrentConfig() string {
	return `{"log_level": "` + cw.config["log_level"] + `", "max_requests": "` + cw.config["max_requests"] + `"}`
}

func (cw *ConfigWatcher) GetConfig(key string) string {
	return cw.config[key]
}

// Proxy
type Proxy struct {
	running bool
}

func NewProxy() *Proxy {
	return &Proxy{}
}

func (p *Proxy) Start() {
	fmt.Println("Proxy: Starting...")
	p.running = true
}

func (p *Proxy) Stop() {
	fmt.Println("Proxy: Stopping...")
	p.running = false
}

func (p *Proxy) HandleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Proxy: Forwarding request to main application: %s %s\n", r.Method, r.URL.Path)
	
	// Simulate proxying to main application
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Proxied response from main application", "path": "%s"}`, r.URL.Path)
}

// Health Status
type HealthStatus struct {
	Status    string
	Timestamp time.Time
	Details   string
}

// Ambassador Sidecar (for external communication)
type AmbassadorSidecar struct {
	name           string
	mainApp        *MainApplication
	externalAPI    *ExternalAPI
	auth           *Authenticator
	rateLimiter    *RateLimiter
	running        bool
}

func NewAmbassadorSidecar(name string, mainApp *MainApplication) *AmbassadorSidecar {
	return &AmbassadorSidecar{
		name:        name,
		mainApp:     mainApp,
		externalAPI: NewExternalAPI(),
		auth:        NewAuthenticator(),
		rateLimiter: NewRateLimiter(100), // 100 requests per minute
	}
}

func (as *AmbassadorSidecar) Start() error {
	fmt.Printf("Ambassador Sidecar %s: Starting...\n", as.name)
	as.running = true
	
	as.externalAPI.Start()
	as.auth.Start()
	as.rateLimiter.Start()
	
	return nil
}

func (as *AmbassadorSidecar) Stop() error {
	fmt.Printf("Ambassador Sidecar %s: Stopping...\n", as.name)
	as.running = false
	
	as.externalAPI.Stop()
	as.auth.Stop()
	as.rateLimiter.Stop()
	
	return nil
}

func (as *AmbassadorSidecar) CallExternalAPI(endpoint string) (string, error) {
	if !as.rateLimiter.Allow() {
		return "", fmt.Errorf("rate limit exceeded")
	}
	
	token, err := as.auth.GetToken()
	if err != nil {
		return "", fmt.Errorf("authentication failed: %v", err)
	}
	
	return as.externalAPI.Call(endpoint, token)
}

// External API
type ExternalAPI struct {
	baseURL string
	running bool
}

func NewExternalAPI() *ExternalAPI {
	return &ExternalAPI{
		baseURL: "https://api.example.com",
	}
}

func (ea *ExternalAPI) Start() {
	fmt.Println("External API: Starting...")
	ea.running = true
}

func (ea *ExternalAPI) Stop() {
	fmt.Println("External API: Stopping...")
	ea.running = false
}

func (ea *ExternalAPI) Call(endpoint, token string) (string, error) {
	fmt.Printf("External API: Calling %s with token %s\n", endpoint, token[:8]+"...")
	time.Sleep(100 * time.Millisecond) // Simulate network latency
	return fmt.Sprintf(`{"data": "Response from %s"}`, endpoint), nil
}

// Authenticator
type Authenticator struct {
	token string
	running bool
}

func NewAuthenticator() *Authenticator {
	return &Authenticator{
		token: "secret_token_12345",
	}
}

func (a *Authenticator) Start() {
	fmt.Println("Authenticator: Starting...")
	a.running = true
}

func (a *Authenticator) Stop() {
	fmt.Println("Authenticator: Stopping...")
	a.running = false
}

func (a *Authenticator) GetToken() (string, error) {
	if a.running {
		return a.token, nil
	}
	return "", fmt.Errorf("authenticator not running")
}

// Rate Limiter
type RateLimiter struct {
	maxRequests int
	requests    []time.Time
	running     bool
}

func NewRateLimiter(maxRequests int) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		requests:    make([]time.Time, 0),
	}
}

func (rl *RateLimiter) Start() {
	fmt.Println("Rate Limiter: Starting...")
	rl.running = true
}

func (rl *RateLimiter) Stop() {
	fmt.Println("Rate Limiter: Stopping...")
	rl.running = false
}

func (rl *RateLimiter) Allow() bool {
	now := time.Now()
	
	// Remove old requests (older than 1 minute)
	validRequests := make([]time.Time, 0)
	for _, req := range rl.requests {
		if now.Sub(req) <= time.Minute {
			validRequests = append(validRequests, req)
		}
	}
	rl.requests = validRequests
	
	// Check if we can allow this request
	if len(rl.requests) < rl.maxRequests {
		rl.requests = append(rl.requests, now)
		return true
	}
	
	return false
}

// Adapter Sidecar (for protocol adaptation)
type AdapterSidecar struct {
	name        string
	mainApp     *MainApplication
	legacyAPI   *LegacyAPI
	transformer *DataTransformer
	running     bool
}

func NewAdapterSidecar(name string, mainApp *MainApplication) *AdapterSidecar {
	return &AdapterSidecar{
		name:        name,
		mainApp:     mainApp,
		legacyAPI:   NewLegacyAPI(),
		transformer: NewDataTransformer(),
	}
}

func (as *AdapterSidecar) Start() error {
	fmt.Printf("Adapter Sidecar %s: Starting...\n", as.name)
	as.running = true
	
	as.legacyAPI.Start()
	as.transformer.Start()
	
	return nil
}

func (as *AdapterSidecar) Stop() error {
	fmt.Printf("Adapter Sidecar %s: Stopping...\n", as.name)
	as.running = false
	
	as.legacyAPI.Stop()
	as.transformer.Stop()
	
	return nil
}

func (as *AdapterSidecar) HandleLegacyRequest(legacyData string) (string, error) {
	// Transform legacy data to modern format
	modernData := as.transformer.LegacyToModern(legacyData)
	
	// Process with main application (simulated)
	fmt.Printf("Adapter Sidecar: Processing transformed data: %s\n", modernData)
	
	// Transform response back to legacy format
	response := as.transformer.ModernToLegacy(`{"status": "success", "data": "processed"}`)
	
	return response, nil
}

// Legacy API
type LegacyAPI struct {
	running bool
}

func NewLegacyAPI() *LegacyAPI {
	return &LegacyAPI{}
}

func (la *LegacyAPI) Start() {
	fmt.Println("Legacy API: Starting...")
	la.running = true
}

func (la *LegacyAPI) Stop() {
	fmt.Println("Legacy API: Stopping...")
	la.running = false
}

// Data Transformer
type DataTransformer struct {
	running bool
}

func NewDataTransformer() *DataTransformer {
	return &DataTransformer{}
}

func (dt *DataTransformer) Start() {
	fmt.Println("Data Transformer: Starting...")
	dt.running = true
}

func (dt *DataTransformer) Stop() {
	fmt.Println("Data Transformer: Stopping...")
	dt.running = false
}

func (dt *DataTransformer) LegacyToModern(legacyData string) string {
	// Simple transformation: XML to JSON
	return fmt.Sprintf(`{"legacy_data": "%s", "format": "xml_to_json"}`, legacyData)
}

func (dt *DataTransformer) ModernToLegacy(modernData string) string {
	// Simple transformation: JSON to XML
	return fmt.Sprintf("<response><data>%s</data></response>", modernData)
}

// Container Orchestrator
type ContainerOrchestrator struct {
	mainApp     *MainApplication
	sidecar     *Sidecar
	ambassador  *AmbassadorSidecar
	adapter     *AdapterSidecar
}

func NewContainerOrchestrator() *ContainerOrchestrator {
	mainApp := NewMainApplication("WebApp", 8080)
	
	return &ContainerOrchestrator{
		mainApp:    mainApp,
		sidecar:    NewSidecar("Sidecar", mainApp),
		ambassador: NewAmbassadorSidecar("Ambassador", mainApp),
		adapter:    NewAdapterSidecar("Adapter", mainApp),
	}
}

func (co *ContainerOrchestrator) StartAll() error {
	fmt.Println("Container Orchestrator: Starting all containers...")
	
	// Start main application
	if err := co.mainApp.Start(); err != nil {
		return fmt.Errorf("failed to start main application: %v", err)
	}
	
	// Start sidecar
	if err := co.sidecar.Start(); err != nil {
		return fmt.Errorf("failed to start sidecar: %v", err)
	}
	
	// Start ambassador sidecar
	if err := co.ambassador.Start(); err != nil {
		return fmt.Errorf("failed to start ambassador sidecar: %v", err)
	}
	
	// Start adapter sidecar
	if err := co.adapter.Start(); err != nil {
		return fmt.Errorf("failed to start adapter sidecar: %v", err)
	}
	
	fmt.Println("Container Orchestrator: All containers started successfully")
	return nil
}

func (co *ContainerOrchestrator) StopAll() error {
	fmt.Println("Container Orchestrator: Stopping all containers...")
	
	// Stop in reverse order
	co.adapter.Stop()
	co.ambassador.Stop()
	co.sidecar.Stop()
	co.mainApp.Stop()
	
	fmt.Println("Container Orchestrator: All containers stopped")
	return nil
}

func (co *ContainerOrchestrator) WaitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	<-sigChan
	fmt.Println("\nContainer Orchestrator: Shutdown signal received")
}

func demonstrateBasicSidecar() {
	fmt.Println("--- Basic Sidecar Demo ---")
	
	orchestrator := NewContainerOrchestrator()
	
	// Start all containers
	err := orchestrator.StartAll()
	if err != nil {
		fmt.Printf("Error starting containers: %v\n", err)
		return
	}
	
	// Let them run for a bit
	time.Sleep(5 * time.Second)
	
	// Stop all containers
	err = orchestrator.StopAll()
	if err != nil {
		fmt.Printf("Error stopping containers: %v\n", err)
	}
}

func demonstrateAmbassadorSidecar() {
	fmt.Println("\n--- Ambassador Sidecar Demo ---")
	
	mainApp := NewMainApplication("APIApp", 8080)
	ambassador := NewAmbassadorSidecar("APIAmbassador", mainApp)
	
	mainApp.Start()
	ambassador.Start()
	
	// Make external API calls through ambassador
	for i := 1; i <= 3; i++ {
		fmt.Printf("Making external API call %d...\n", i)
		response, err := ambassador.CallExternalAPI("/api/users")
		if err != nil {
			fmt.Printf("Call %d failed: %v\n", i, err)
		} else {
			fmt.Printf("Call %d succeeded: %s\n", i, response)
		}
		time.Sleep(500 * time.Millisecond)
	}
	
	ambassador.Stop()
	mainApp.Stop()
}

func demonstrateAdapterSidecar() {
	fmt.Println("\n--- Adapter Sidecar Demo ---")
	
	mainApp := NewMainApplication("LegacyApp", 8080)
	adapter := NewAdapterSidecar("LegacyAdapter", mainApp)
	
	mainApp.Start()
	adapter.Start()
	
	// Handle legacy requests
	legacyRequests := []string{
		"<user><name>John</name><age>30</age></user>",
		"<order><id>123</id><amount>99.99</amount></order>",
	}
	
	for i, request := range legacyRequests {
		fmt.Printf("Handling legacy request %d: %s\n", i+1, request)
		response, err := adapter.HandleLegacyRequest(request)
		if err != nil {
			fmt.Printf("Request %d failed: %v\n", i+1, err)
		} else {
			fmt.Printf("Request %d succeeded: %s\n", i+1, response)
		}
		time.Sleep(200 * time.Millisecond)
	}
	
	adapter.Stop()
	mainApp.Stop()
}

func demonstrateSidecarFeatures() {
	fmt.Println("\n--- Sidecar Features Demo ---")
	
	mainApp := NewMainApplication("FeatureApp", 8080)
	sidecar := NewSidecar("FeatureSidecar", mainApp)
	
	mainApp.Start()
	sidecar.Start()
	
	// Simulate some activity
	for i := 1; i <= 5; i++ {
		fmt.Printf("Activity cycle %d\n", i)
		time.Sleep(1 * time.Second)
	}
	
	// Show sidecar functionality
	fmt.Println("\nSidecar Features:")
	fmt.Println("- Health monitoring: Active")
	fmt.Println("- Logging: Active")
	fmt.Println("- Metrics collection: Active")
	fmt.Println("- Configuration watching: Active")
	fmt.Println("- Proxy functionality: Active")
	
	sidecar.Stop()
	mainApp.Stop()
}

func main() {
	fmt.Println("=== Sidecar Pattern Demo ===")
	
	demonstrateBasicSidecar()
	demonstrateAmbassadorSidecar()
	demonstrateAdapterSidecar()
	demonstrateSidecarFeatures()
	
	fmt.Println("\nAll sidecar patterns demonstrated successfully!")
}
