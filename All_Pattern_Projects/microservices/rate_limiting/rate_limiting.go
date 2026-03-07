package main

import (
	"fmt"
	"sync"
	"time"
)

// Rate Limiting Pattern

// RateLimiter interface
type RateLimiter interface {
	Allow(key string) bool
	AllowN(key string, n int) bool
	GetStats(key string) RateLimiterStats
	Reset(key string)
}

type RateLimiterStats struct {
	AllowedRequests int
	DeniedRequests  int
	TotalRequests   int
	CurrentTokens   int
	MaxTokens       int
}

// Token Bucket Rate Limiter
type TokenBucketRateLimiter struct {
	buckets map[string]*TokenBucket
	mu      sync.RWMutex
	config  TokenBucketConfig
}

type TokenBucketConfig struct {
	MaxTokens   int
	RefillRate  int // tokens per second
	RefillInterval time.Duration
}

type TokenBucket struct {
	tokens       int
	maxTokens    int
	lastRefill   time.Time
	refillRate   int
	allowedCount int
	deniedCount  int
	totalCount   int
	mu           sync.Mutex
}

func NewTokenBucketRateLimiter(config TokenBucketConfig) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		buckets: make(map[string]*TokenBucket),
		config:  config,
	}
}

func (tbl *TokenBucketRateLimiter) getOrCreateBucket(key string) *TokenBucket {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()
	
	if bucket, exists := tbl.buckets[key]; exists {
		return bucket
	}
	
	bucket := &TokenBucket{
		tokens:     tbl.config.MaxTokens,
		maxTokens:  tbl.config.MaxTokens,
		lastRefill: time.Now(),
		refillRate:  tbl.config.RefillRate,
	}
	
	tbl.buckets[key] = bucket
	return bucket
}

func (tbl *TokenBucketRateLimiter) Allow(key string) bool {
	return tbl.AllowN(key, 1)
}

func (tbl *TokenBucketRateLimiter) AllowN(key string, n int) bool {
	bucket := tbl.getOrCreateBucket(key)
	
	bucket.mu.Lock()
	defer bucket.mu.Unlock()
	
	bucket.totalCount++
	
	// Refill tokens
	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill)
	tokensToAdd := int(elapsed.Seconds() * float64(bucket.refillRate))
	
	if tokensToAdd > 0 {
		bucket.tokens = min(bucket.tokens+tokensToAdd, bucket.maxTokens)
		bucket.lastRefill = now
	}
	
	// Check if enough tokens
	if bucket.tokens >= n {
		bucket.tokens -= n
		bucket.allowedCount++
		return true
	}
	
	bucket.deniedCount++
	return false
}

func (tbl *TokenBucketRateLimiter) GetStats(key string) RateLimiterStats {
	tbl.mu.RLock()
	bucket, exists := tbl.buckets[key]
	tbl.mu.RUnlock()
	
	if !exists {
		return RateLimiterStats{}
	}
	
	bucket.mu.Lock()
	defer bucket.mu.Unlock()
	
	return RateLimiterStats{
		AllowedRequests: bucket.allowedCount,
		DeniedRequests:  bucket.deniedCount,
		TotalRequests:   bucket.totalCount,
		CurrentTokens:   bucket.tokens,
		MaxTokens:       bucket.maxTokens,
	}
}

func (tbl *TokenBucketRateLimiter) Reset(key string) {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()
	
	if bucket, exists := tbl.buckets[key]; exists {
		bucket.mu.Lock()
		bucket.tokens = bucket.maxTokens
		bucket.allowedCount = 0
		bucket.deniedCount = 0
		bucket.totalCount = 0
		bucket.lastRefill = time.Now()
		bucket.mu.Unlock()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Sliding Window Rate Limiter
type SlidingWindowRateLimiter struct {
	windows map[string]*SlidingWindow
	mu      sync.RWMutex
	config  SlidingWindowConfig
}

type SlidingWindowConfig struct {
	WindowDuration time.Duration
	MaxRequests    int
}

type SlidingWindow struct {
	requests      []time.Time
	windowDuration time.Duration
	maxRequests    int
	allowedCount   int
	deniedCount    int
	totalCount     int
	mu             sync.Mutex
}

func NewSlidingWindowRateLimiter(config SlidingWindowConfig) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		windows: make(map[string]*SlidingWindow),
		config:  config,
	}
}

func (swl *SlidingWindowRateLimiter) getOrCreateWindow(key string) *SlidingWindow {
	swl.mu.Lock()
	defer swl.mu.Unlock()
	
	if window, exists := swl.windows[key]; exists {
		return window
	}
	
	window := &SlidingWindow{
		requests:      make([]time.Time, 0),
		windowDuration: swl.config.WindowDuration,
		maxRequests:    swl.config.MaxRequests,
	}
	
	swl.windows[key] = window
	return window
}

func (swl *SlidingWindowRateLimiter) Allow(key string) bool {
	return swl.AllowN(key, 1)
}

func (swl *SlidingWindowRateLimiter) AllowN(key string, n int) bool {
	window := swl.getOrCreateWindow(key)
	
	window.mu.Lock()
	defer window.mu.Unlock()
	
	window.totalCount++
	
	now := time.Now()
	
	// Remove old requests outside the window
	validRequests := make([]time.Time, 0)
	for _, req := range window.requests {
		if now.Sub(req) <= window.windowDuration {
			validRequests = append(validRequests, req)
		}
	}
	window.requests = validRequests
	
	// Check if adding n requests would exceed the limit
	if len(window.requests)+n <= window.maxRequests {
		for i := 0; i < n; i++ {
			window.requests = append(window.requests, now)
		}
		window.allowedCount++
		return true
	}
	
	window.deniedCount++
	return false
}

func (swl *SlidingWindowRateLimiter) GetStats(key string) RateLimiterStats {
	swl.mu.RLock()
	window, exists := swl.windows[key]
	swl.mu.RUnlock()
	
	if !exists {
		return RateLimiterStats{}
	}
	
	window.mu.Lock()
	defer window.mu.Unlock()
	
	return RateLimiterStats{
		AllowedRequests: window.allowedCount,
		DeniedRequests:  window.deniedCount,
		TotalRequests:   window.totalCount,
		CurrentTokens:   window.maxRequests - len(window.requests),
		MaxTokens:       window.maxRequests,
	}
}

func (swl *SlidingWindowRateLimiter) Reset(key string) {
	swl.mu.Lock()
	defer swl.mu.Unlock()
	
	if window, exists := swl.windows[key]; exists {
		window.mu.Lock()
		window.requests = window.requests[:0]
		window.allowedCount = 0
		window.deniedCount = 0
		window.totalCount = 0
		window.mu.Unlock()
	}
}

// Fixed Window Rate Limiter
type FixedWindowRateLimiter struct {
	windows map[string]*FixedWindow
	mu      sync.RWMutex
	config  FixedWindowConfig
}

type FixedWindowConfig struct {
	WindowDuration time.Duration
	MaxRequests    int
}

type FixedWindow struct {
	requestCount   int
	windowStart    time.Time
	windowDuration time.Duration
	maxRequests    int
	allowedCount   int
	deniedCount    int
	totalCount     int
	mu             sync.Mutex
}

func NewFixedWindowRateLimiter(config FixedWindowConfig) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		windows: make(map[string]*FixedWindow),
		config:  config,
	}
}

func (fwl *FixedWindowRateLimiter) getOrCreateWindow(key string) *FixedWindow {
	fwl.mu.Lock()
	defer fwl.mu.Unlock()
	
	if window, exists := fwl.windows[key]; exists {
		return window
	}
	
	window := &FixedWindow{
		windowStart:    time.Now(),
		windowDuration: fwl.config.WindowDuration,
		maxRequests:    fwl.config.MaxRequests,
	}
	
	fwl.windows[key] = window
	return window
}

func (fwl *FixedWindowRateLimiter) Allow(key string) bool {
	return fwl.AllowN(key, 1)
}

func (fwl *FixedWindowRateLimiter) AllowN(key string, n int) bool {
	window := fwl.getOrCreateWindow(key)
	
	window.mu.Lock()
	defer window.mu.Unlock()
	
	window.totalCount++
	
	now := time.Now()
	
	// Reset window if needed
	if now.Sub(window.windowStart) >= window.windowDuration {
		window.requestCount = 0
		window.windowStart = now
	}
	
	// Check if adding n requests would exceed the limit
	if window.requestCount+n <= window.maxRequests {
		window.requestCount += n
		window.allowedCount++
		return true
	}
	
	window.deniedCount++
	return false
}

func (fwl *FixedWindowRateLimiter) GetStats(key string) RateLimiterStats {
	fwl.mu.RLock()
	window, exists := fwl.windows[key]
	fwl.mu.RUnlock()
	
	if !exists {
		return RateLimiterStats{}
	}
	
	window.mu.Lock()
	defer window.mu.Unlock()
	
	return RateLimiterStats{
		AllowedRequests: window.allowedCount,
		DeniedRequests:  window.deniedCount,
		TotalRequests:   window.totalCount,
		CurrentTokens:   window.maxRequests - window.requestCount,
		MaxTokens:       window.maxRequests,
	}
}

func (fwl *FixedWindowRateLimiter) Reset(key string) {
	fwl.mu.Lock()
	defer fwl.mu.Unlock()
	
	if window, exists := fwl.windows[key]; exists {
		window.mu.Lock()
		window.requestCount = 0
		window.windowStart = time.Now()
		window.allowedCount = 0
		window.deniedCount = 0
		window.totalCount = 0
		window.mu.Unlock()
	}
}

// API Service with Rate Limiting
type APIService struct {
	name       string
	rateLimiter RateLimiter
}

func NewAPIService(name string, rateLimiter RateLimiter) *APIService {
	return &APIService{
		name:        name,
		rateLimiter: rateLimiter,
	}
}

func (api *APIService) HandleRequest(clientID string) error {
	if api.rateLimiter.Allow(clientID) {
		fmt.Printf("API Service %s: Request from %s allowed\n", api.name, clientID)
		return nil
	}
	
	fmt.Printf("API Service %s: Request from %s denied (rate limit exceeded)\n", api.name, clientID)
	return fmt.Errorf("rate limit exceeded")
}

func (api *APIService) GetStats(clientID string) {
	stats := api.rateLimiter.GetStats(clientID)
	fmt.Printf("API Service %s Stats for %s: Allowed=%d, Denied=%d, Total=%d\n", 
		api.name, clientID, stats.AllowedRequests, stats.DeniedRequests, stats.TotalRequests)
}

// Multi-tier Rate Limiter
type MultiTierRateLimiter struct {
	limiters []RateLimiter
}

func NewMultiTierRateLimiter(limiters ...RateLimiter) *MultiTierRateLimiter {
	return &MultiTierRateLimiter{limiters: limiters}
}

func (mtrl *MultiTierRateLimiter) Allow(key string) bool {
	return mtrl.AllowN(key, 1)
}

func (mtrl *MultiTierRateLimiter) AllowN(key string, n int) bool {
	for _, limiter := range mtrl.limiters {
		if !limiter.AllowN(key, n) {
			return false
		}
	}
	return true
}

func (mtrl *MultiTierRateLimiter) GetStats(key string) RateLimiterStats {
	// Return stats from the first limiter
	if len(mtrl.limiters) > 0 {
		return mtrl.limiters[0].GetStats(key)
	}
	return RateLimiterStats{}
}

func (mtrl *MultiTierRateLimiter) Reset(key string) {
	for _, limiter := range mtrl.limiters {
		limiter.Reset(key)
	}
}

// Rate Limiter Manager
type RateLimiterManager struct {
	limiters map[string]RateLimiter
	mu       sync.RWMutex
}

func NewRateLimiterManager() *RateLimiterManager {
	return &RateLimiterManager{
		limiters: make(map[string]RateLimiter),
	}
}

func (rlm *RateLimiterManager) AddLimiter(name string, limiter RateLimiter) {
	rlm.mu.Lock()
	defer rlm.mu.Unlock()
	rlm.limiters[name] = limiter
	fmt.Printf("Rate Limiter Manager: Added limiter '%s'\n", name)
}

func (rlm *RateLimiterManager) GetLimiter(name string) RateLimiter {
	rlm.mu.RLock()
	defer rlm.mu.RUnlock()
	return rlm.limiters[name]
}

func (rlm *RateLimiterManager) GetAllStats() map[string]map[string]RateLimiterStats {
	rlm.mu.RLock()
	defer rlm.mu.RUnlock()
	
	allStats := make(map[string]map[string]RateLimiterStats)
	
	// For demo, we'll use some sample client IDs
	clientIDs := []string{"client1", "client2", "client3"}
	
	for name, limiter := range rlm.limiters {
		allStats[name] = make(map[string]RateLimiterStats)
		for _, clientID := range clientIDs {
			allStats[name][clientID] = limiter.GetStats(clientID)
		}
	}
	
	return allStats
}

func demonstrateTokenBucket() {
	fmt.Println("--- Token Bucket Rate Limiter Demo ---")
	
	config := TokenBucketConfig{
		MaxTokens:     10,
		RefillRate:    2,  // 2 tokens per second
		RefillInterval: time.Second,
	}
	
	limiter := NewTokenBucketRateLimiter(config)
	api := NewAPIService("UserService", limiter)
	
	// Simulate rapid requests
	clientID := "user123"
	
	for i := 1; i <= 15; i++ {
		fmt.Printf("Request %d: ", i)
		err := api.HandleRequest(clientID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		
		if i%5 == 0 {
			api.GetStats(clientID)
		}
		
		time.Sleep(100 * time.Millisecond)
	}
	
	// Wait for refill
	fmt.Println("\nWaiting for token refill...")
	time.Sleep(2 * time.Second)
	
	fmt.Println("\nAfter refill:")
	err := api.HandleRequest(clientID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	api.GetStats(clientID)
}

func demonstrateSlidingWindow() {
	fmt.Println("\n--- Sliding Window Rate Limiter Demo ---")
	
	config := SlidingWindowConfig{
		WindowDuration: 5 * time.Second,
		MaxRequests:    5,
	}
	
	limiter := NewSlidingWindowRateLimiter(config)
	api := NewAPIService("OrderService", limiter)
	
	clientID := "client456"
	
	// Make requests
	for i := 1; i <= 8; i++ {
		fmt.Printf("Request %d: ", i)
		err := api.HandleRequest(clientID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		
		time.Sleep(500 * time.Millisecond)
	}
	
	api.GetStats(clientID)
	
	// Wait and try again
	fmt.Println("\nWaiting 3 seconds...")
	time.Sleep(3 * time.Second)
	
	fmt.Println("\nAfter wait:")
	err := api.HandleRequest(clientID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	api.GetStats(clientID)
}

func demonstrateFixedWindow() {
	fmt.Println("\n--- Fixed Window Rate Limiter Demo ---")
	
	config := FixedWindowConfig{
		WindowDuration: 3 * time.Second,
		MaxRequests:    4,
	}
	
	limiter := NewFixedWindowRateLimiter(config)
	api := NewAPIService("PaymentService", limiter)
	
	clientID := "client789"
	
	// Make requests
	for i := 1; i <= 6; i++ {
		fmt.Printf("Request %d: ", i)
		err := api.HandleRequest(clientID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		
		time.Sleep(200 * time.Millisecond)
	}
	
	api.GetStats(clientID)
	
	// Wait for window reset
	fmt.Println("\nWaiting for window reset...")
	time.Sleep(4 * time.Second)
	
	fmt.Println("\nAfter window reset:")
	err := api.HandleRequest(clientID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	api.GetStats(clientID)
}

func demonstrateMultiTier() {
	fmt.Println("\n--- Multi-tier Rate Limiter Demo ---")
	
	// Create multiple rate limiters
	tokenBucketConfig := TokenBucketConfig{
		MaxTokens:     5,
		RefillRate:    1,
		RefillInterval: time.Second,
	}
	
	slidingWindowConfig := SlidingWindowConfig{
		WindowDuration: 10 * time.Second,
		MaxRequests:    8,
	}
	
	tokenBucketLimiter := NewTokenBucketRateLimiter(tokenBucketConfig)
	slidingWindowLimiter := NewSlidingWindowRateLimiter(slidingWindowConfig)
	
	// Combine them - both must allow the request
	multiTierLimiter := NewMultiTierRateLimiter(tokenBucketLimiter, slidingWindowLimiter)
	api := NewAPIService("PremiumService", multiTierLimiter)
	
	clientID := "premium_client"
	
	// Make requests
	for i := 1; i <= 10; i++ {
		fmt.Printf("Request %d: ", i)
		err := api.HandleRequest(clientID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		
		time.Sleep(300 * time.Millisecond)
	}
	
	fmt.Printf("\nToken Bucket Stats: %v\n", tokenBucketLimiter.GetStats(clientID))
	fmt.Printf("Sliding Window Stats: %v\n", slidingWindowLimiter.GetStats(clientID))
}

func demonstrateRateLimiterManager() {
	fmt.Println("\n--- Rate Limiter Manager Demo ---")
	
	manager := NewRateLimiterManager()
	
	// Add different rate limiters for different services
	tokenBucketConfig := TokenBucketConfig{
		MaxTokens:     10,
		RefillRate:    3,
		RefillInterval: time.Second,
	}
	
	slidingWindowConfig := SlidingWindowConfig{
		WindowDuration: 5 * time.Second,
		MaxRequests:    7,
	}
	
	fixedWindowConfig := FixedWindowConfig{
		WindowDuration: 4 * time.Second,
		MaxRequests:    6,
	}
	
	manager.AddLimiter("UserService", NewTokenBucketRateLimiter(tokenBucketConfig))
	manager.AddLimiter("OrderService", NewSlidingWindowRateLimiter(slidingWindowConfig))
	manager.AddLimiter("PaymentService", NewFixedWindowRateLimiter(fixedWindowConfig))
	
	// Create services
	userService := NewAPIService("UserService", manager.GetLimiter("UserService"))
	orderService := NewAPIService("OrderService", manager.GetLimiter("OrderService"))
	paymentService := NewAPIService("PaymentService", manager.GetLimiter("PaymentService"))
	
	// Simulate requests across services
	services := []*APIService{userService, orderService, paymentService}
	clientIDs := []string{"client1", "client2", "client3"}
	
	for round := 1; round <= 3; round++ {
		fmt.Printf("\nRound %d:\n", round)
		
		for _, service := range services {
			for _, clientID := range clientIDs {
				err := service.HandleRequest(clientID)
				if err != nil {
					fmt.Printf("Error in %s for %s: %v\n", service.name, clientID, err)
				}
			}
		}
		
		time.Sleep(1 * time.Second)
	}
	
	// Show all stats
	fmt.Println("\nFinal Stats:")
	for serviceName, clientStats := range manager.GetAllStats() {
		fmt.Printf("\n%s:\n", serviceName)
		for clientID, stats := range clientStats {
			fmt.Printf("  %s: Allowed=%d, Denied=%d, Total=%d\n", 
				clientID, stats.AllowedRequests, stats.DeniedRequests, stats.TotalRequests)
		}
	}
}

func main() {
	fmt.Println("=== Rate Limiting Pattern Demo ===")
	
	demonstrateTokenBucket()
	demonstrateSlidingWindow()
	demonstrateFixedWindow()
	demonstrateMultiTier()
	demonstrateRateLimiterManager()
	
	fmt.Println("\nAll rate limiting patterns demonstrated successfully!")
}
