package main

import (
	"fmt"
	"sync"
	"time"
)

// Circuit Breaker Pattern

// CircuitBreaker interface
type CircuitBreaker interface {
	Execute(func() error) error
	GetState() CircuitBreakerState
	GetStats() CircuitBreakerStats
	Reset()
}

type CircuitBreakerState int

const (
	Closed CircuitBreakerState = iota
	Open
	HalfOpen
)

type CircuitBreakerStats struct {
	State         CircuitBreakerState
	FailureCount  int
	SuccessCount  int
	TotalRequests int
	LastFailure   time.Time
}

// Basic Circuit Breaker
type BasicCircuitBreaker struct {
	maxFailures   int
	resetTimeout  time.Duration
	state         CircuitBreakerState
	failureCount  int
	successCount  int
	totalRequests int
	lastFailure   time.Time
	mu            sync.RWMutex
}

func NewBasicCircuitBreaker(maxFailures int, resetTimeout time.Duration) *BasicCircuitBreaker {
	return &BasicCircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        Closed,
	}
}

func (bcb *BasicCircuitBreaker) Execute(task func() error) error {
	bcb.mu.Lock()
	defer bcb.mu.Unlock()
	
	bcb.totalRequests++
	
	// Check if circuit should be reset
	if bcb.state == Open && time.Since(bcb.lastFailure) > bcb.resetTimeout {
		bcb.state = HalfOpen
		bcb.failureCount = 0
		fmt.Println("Circuit Breaker: Transitioning to Half-Open")
	}
	
	// Reject if circuit is open
	if bcb.state == Open {
		return fmt.Errorf("circuit breaker: circuit open")
	}
	
	// Execute task
	err := task()
	
	if err != nil {
		bcb.failureCount++
		bcb.lastFailure = time.Now()
		
		if bcb.failureCount >= bcb.maxFailures {
			bcb.state = Open
			fmt.Println("Circuit Breaker: Transitioning to Open")
		}
	} else {
		bcb.successCount++
		
		if bcb.state == HalfOpen {
			bcb.state = Closed
			bcb.failureCount = 0
			fmt.Println("Circuit Breaker: Transitioning to Closed")
		}
	}
	
	return err
}

func (bcb *BasicCircuitBreaker) GetState() CircuitBreakerState {
	bcb.mu.RLock()
	defer bcb.mu.RUnlock()
	return bcb.state
}

func (bcb *BasicCircuitBreaker) GetStats() CircuitBreakerStats {
	bcb.mu.RLock()
	defer bcb.mu.RUnlock()
	return CircuitBreakerStats{
		State:         bcb.state,
		FailureCount:  bcb.failureCount,
		SuccessCount:  bcb.successCount,
		TotalRequests: bcb.totalRequests,
		LastFailure:   bcb.lastFailure,
	}
}

func (bcb *BasicCircuitBreaker) Reset() {
	bcb.mu.Lock()
	defer bcb.mu.Unlock()
	
	bcb.state = Closed
	bcb.failureCount = 0
	bcb.successCount = 0
	bcb.totalRequests = 0
	bcb.lastFailure = time.Time{}
	fmt.Println("Circuit Breaker: Reset to Closed")
}

// Advanced Circuit Breaker with metrics
type AdvancedCircuitBreaker struct {
	maxFailures        int
	resetTimeout       time.Duration
	halfOpenMaxCalls   int
	state              CircuitBreakerState
	failureCount       int
	successCount       int
	totalRequests      int
	halfOpenCalls      int
	lastFailure        time.Time
	mu                 sync.RWMutex
	onStateChange      func(from, to CircuitBreakerState)
	onSuccess          func()
	onFailure          func(err error)
	metrics            *CircuitBreakerMetrics
}

type CircuitBreakerMetrics struct {
	RequestCount    int64
	SuccessCount    int64
	FailureCount    int64
	TimeoutCount    int64
	AverageResponse time.Duration
	LastResponse    time.Duration
}

func NewAdvancedCircuitBreaker(maxFailures int, resetTimeout time.Duration, halfOpenMaxCalls int) *AdvancedCircuitBreaker {
	return &AdvancedCircuitBreaker{
		maxFailures:      maxFailures,
		resetTimeout:     resetTimeout,
		halfOpenMaxCalls: halfOpenMaxCalls,
		state:            Closed,
		metrics:          &CircuitBreakerMetrics{},
	}
}

func (acb *AdvancedCircuitBreaker) Execute(task func() error) error {
	start := time.Now()
	
	acb.mu.Lock()
	defer acb.mu.Unlock()
	
	acb.metrics.RequestCount++
	
	// Check if circuit should be reset
	if acb.state == Open && time.Since(acb.lastFailure) > acb.resetTimeout {
		acb.transitionTo(HalfOpen)
	}
	
	// Reject if circuit is open
	if acb.state == Open {
		acb.metrics.TimeoutCount++
		return fmt.Errorf("circuit breaker: circuit open")
	}
	
	// Limit calls in half-open state
	if acb.state == HalfOpen && acb.halfOpenCalls >= acb.halfOpenMaxCalls {
		return fmt.Errorf("circuit breaker: half-open call limit reached")
	}
	
	if acb.state == HalfOpen {
		acb.halfOpenCalls++
	}
	
	// Execute task
	err := task()
	responseTime := time.Since(start)
	acb.metrics.LastResponse = responseTime
	
	if err != nil {
		acb.failureCount++
		acb.metrics.FailureCount++
		acb.lastFailure = time.Now()
		
		if acb.onFailure != nil {
			acb.onFailure(err)
		}
		
		if acb.failureCount >= acb.maxFailures {
			acb.transitionTo(Open)
		}
	} else {
		acb.successCount++
		acb.metrics.SuccessCount++
		
		if acb.onSuccess != nil {
			acb.onSuccess()
		}
		
		if acb.state == HalfOpen {
			acb.transitionTo(Closed)
		}
	}
	
	// Update average response time
	if acb.metrics.RequestCount > 0 {
		acb.metrics.AverageResponse = time.Duration(
			(int64(acb.metrics.AverageResponse)*int64(acb.metrics.RequestCount-1) + int64(responseTime)) / int64(acb.metrics.RequestCount))
	}
	
	return err
}

func (acb *AdvancedCircuitBreaker) transitionTo(newState CircuitBreakerState) {
	oldState := acb.state
	acb.state = newState
	
	if newState == Open {
		acb.halfOpenCalls = 0
	} else if newState == Closed {
		acb.failureCount = 0
		acb.halfOpenCalls = 0
	}
	
	if acb.onStateChange != nil {
		acb.onStateChange(oldState, newState)
	}
	
	fmt.Printf("Circuit Breaker: Transitioned from %v to %v\n", oldState, newState)
}

func (acb *AdvancedCircuitBreaker) GetState() CircuitBreakerState {
	acb.mu.RLock()
	defer acb.mu.RUnlock()
	return acb.state
}

func (acb *AdvancedCircuitBreaker) GetStats() CircuitBreakerStats {
	acb.mu.RLock()
	defer acb.mu.RUnlock()
	return CircuitBreakerStats{
		State:         acb.state,
		FailureCount:  acb.failureCount,
		SuccessCount:  acb.successCount,
		TotalRequests: int(acb.metrics.RequestCount),
		LastFailure:   acb.lastFailure,
	}
}

func (acb *AdvancedCircuitBreaker) GetMetrics() *CircuitBreakerMetrics {
	acb.mu.RLock()
	defer acb.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	return &CircuitBreakerMetrics{
		RequestCount:    acb.metrics.RequestCount,
		SuccessCount:    acb.metrics.SuccessCount,
		FailureCount:    acb.metrics.FailureCount,
		TimeoutCount:    acb.metrics.TimeoutCount,
		AverageResponse: acb.metrics.AverageResponse,
		LastResponse:    acb.metrics.LastResponse,
	}
}

func (acb *AdvancedCircuitBreaker) Reset() {
	acb.mu.Lock()
	defer acb.mu.Unlock()
	
	acb.state = Closed
	acb.failureCount = 0
	acb.successCount = 0
	acb.halfOpenCalls = 0
	acb.lastFailure = time.Time{}
	acb.metrics = &CircuitBreakerMetrics{}
	
	fmt.Println("Circuit Breaker: Reset to Closed")
}

func (acb *AdvancedCircuitBreaker) SetStateChangeCallback(callback func(from, to CircuitBreakerState)) {
	acb.mu.Lock()
	defer acb.mu.Unlock()
	acb.onStateChange = callback
}

func (acb *AdvancedCircuitBreaker) SetSuccessCallback(callback func()) {
	acb.mu.Lock()
	defer acb.mu.Unlock()
	acb.onSuccess = callback
}

func (acb *AdvancedCircuitBreaker) SetFailureCallback(callback func(error)) {
	acb.mu.Lock()
	defer acb.mu.Unlock()
	acb.onFailure = callback
}

// Service with Circuit Breaker
type Service struct {
	name           string
	circuitBreaker CircuitBreaker
	failureRate    float64
}

func NewService(name string, circuitBreaker CircuitBreaker, failureRate float64) *Service {
	return &Service{
		name:           name,
		circuitBreaker: circuitBreaker,
		failureRate:    failureRate,
	}
}

func (s *Service) Call() error {
	task := func() error {
		fmt.Printf("Service %s: Processing request\n", s.name)
		time.Sleep(50 * time.Millisecond)
		
		// Simulate failure based on failure rate
		if s.failureRate > 0 && time.Now().UnixNano()%100 < int(s.failureRate*100) {
			return fmt.Errorf("service %s: simulated failure", s.name)
		}
		
		fmt.Printf("Service %s: Request completed successfully\n", s.name)
		return nil
	}
	
	return s.circuitBreaker.Execute(task)
}

func (s *Service) GetStatus() string {
	state := s.circuitBreaker.GetState()
	stats := s.circuitBreaker.GetStats()
	
	return fmt.Sprintf("Service %s: State=%v, Failures=%d, Successes=%d, Total=%d", 
		s.name, state, stats.FailureCount, stats.SuccessCount, stats.TotalRequests)
}

// Circuit Breaker Registry
type CircuitBreakerRegistry struct {
	circuitBreakers map[string]CircuitBreaker
	mu              sync.RWMutex
}

func NewCircuitBreakerRegistry() *CircuitBreakerRegistry {
	return &CircuitBreakerRegistry{
		circuitBreakers: make(map[string]CircuitBreaker),
	}
}

func (cbr *CircuitBreakerRegistry) Register(name string, circuitBreaker CircuitBreaker) {
	cbr.mu.Lock()
	defer cbr.mu.Unlock()
	cbr.circuitBreakers[name] = circuitBreaker
	fmt.Printf("Registry: Registered circuit breaker '%s'\n", name)
}

func (cbr *CircuitBreakerRegistry) Get(name string) CircuitBreaker {
	cbr.mu.RLock()
	defer cbr.mu.RUnlock()
	return cbr.circuitBreakers[name]
}

func (cbr *CircuitBreakerRegistry) GetAllStatus() map[string]string {
	cbr.mu.RLock()
	defer cbr.mu.RUnlock()
	
	status := make(map[string]string)
	for name, cb := range cbr.circuitBreakers {
		stats := cb.GetStats()
		status[name] = fmt.Sprintf("State=%v, Failures=%d, Successes=%d", 
			stats.State, stats.FailureCount, stats.SuccessCount)
	}
	return status
}

func (cbr *CircuitBreakerRegistry) ResetAll() {
	cbr.mu.RLock()
	defer cbr.mu.RUnlock()
	
	for name, cb := range cbr.circuitBreakers {
		cb.Reset()
		fmt.Printf("Registry: Reset circuit breaker '%s'\n", name)
	}
}

// Timeout Circuit Breaker
type TimeoutCircuitBreaker struct {
	*AdvancedCircuitBreaker
	timeout time.Duration
}

func NewTimeoutCircuitBreaker(maxFailures, resetTimeout, timeout, halfOpenMaxCalls int) *TimeoutCircuitBreaker {
	return &TimeoutCircuitBreaker{
		AdvancedCircuitBreaker: NewAdvancedCircuitBreaker(maxFailures, 
			time.Duration(resetTimeout)*time.Second, halfOpenMaxCalls),
		timeout: time.Duration(timeout) * time.Millisecond,
	}
}

func (tcb *TimeoutCircuitBreaker) Execute(task func() error) error {
	wrappedTask := func() error {
		done := make(chan error, 1)
		
		go func() {
			done <- task()
		}()
		
		select {
		case err := <-done:
			return err
		case <-time.After(tcb.timeout):
			return fmt.Errorf("circuit breaker: timeout after %v", tcb.timeout)
		}
	}
	
	return tcb.AdvancedCircuitBreaker.Execute(wrappedTask)
}

// Retry Circuit Breaker
type RetryCircuitBreaker struct {
	*AdvancedCircuitBreaker
	maxRetries int
}

func NewRetryCircuitBreaker(maxFailures, resetTimeout, halfOpenMaxCalls, maxRetries int) *RetryCircuitBreaker {
	return &RetryCircuitBreaker{
		AdvancedCircuitBreaker: NewAdvancedCircuitBreaker(maxFailures, 
			time.Duration(resetTimeout)*time.Second, halfOpenMaxCalls),
		maxRetries: maxRetries,
	}
}

func (rcb *RetryCircuitBreaker) Execute(task func() error) error {
	wrappedTask := func() error {
		var lastErr error
		
		for i := 0; i <= rcb.maxRetries; i++ {
			if i > 0 {
				fmt.Printf("Circuit Breaker: Retry attempt %d/%d\n", i, rcb.maxRetries)
				time.Sleep(time.Duration(i*100) * time.Millisecond) // Exponential backoff
			}
			
			err := task()
			if err == nil {
				return nil
			}
			
			lastErr = err
		}
		
		return lastErr
	}
	
	return rcb.AdvancedCircuitBreaker.Execute(wrappedTask)
}

func demonstrateBasicCircuitBreaker() {
	fmt.Println("--- Basic Circuit Breaker Demo ---")
	
	cb := NewBasicCircuitBreaker(3, 2*time.Second)
	service := NewService("UserService", cb, 0.8) // 80% failure rate
	
	// Execute requests
	for i := 1; i <= 10; i++ {
		fmt.Printf("\nRequest %d:\n", i)
		err := service.Call()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println(service.GetStatus())
		
		time.Sleep(500 * time.Millisecond)
	}
	
	// Wait for reset
	fmt.Println("\nWaiting for circuit breaker to reset...")
	time.Sleep(3 * time.Second)
	
	// Try again
	fmt.Println("\nAfter reset:")
	err := service.Call()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println(service.GetStatus())
}

func demonstrateAdvancedCircuitBreaker() {
	fmt.Println("\n--- Advanced Circuit Breaker Demo ---")
	
	cb := NewAdvancedCircuitBreaker(3, 2*time.Second, 2)
	
	// Set callbacks
	cb.SetStateChangeCallback(func(from, to CircuitBreakerState) {
		fmt.Printf("Callback: State changed from %v to %v\n", from, to)
	})
	
	cb.SetSuccessCallback(func() {
		fmt.Println("Callback: Success!")
	})
	
	cb.SetFailureCallback(func(err error) {
		fmt.Printf("Callback: Failure - %v\n", err)
	})
	
	service := NewService("OrderService", cb, 0.7) // 70% failure rate
	
	// Execute requests
	for i := 1; i <= 8; i++ {
		fmt.Printf("\nRequest %d:\n", i)
		err := service.Call()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		
		metrics := cb.GetMetrics()
		fmt.Printf("Metrics: Requests=%d, Success=%d, Failures=%d, Timeouts=%d, Avg Response=%v\n",
			metrics.RequestCount, metrics.SuccessCount, metrics.FailureCount, 
			metrics.TimeoutCount, metrics.AverageResponse)
		
		time.Sleep(300 * time.Millisecond)
	}
}

func demonstrateCircuitBreakerRegistry() {
	fmt.Println("\n--- Circuit Breaker Registry Demo ---")
	
	registry := NewCircuitBreakerRegistry()
	
	// Register multiple circuit breakers
	registry.Register("UserService", NewBasicCircuitBreaker(2, 1*time.Second))
	registry.Register("OrderService", NewBasicCircuitBreaker(3, 2*time.Second))
	registry.Register("PaymentService", NewBasicCircuitBreaker(1, 500*time.Millisecond))
	
	// Create services
	userService := NewService("UserService", registry.Get("UserService"), 0.6)
	orderService := NewService("OrderService", registry.Get("OrderService"), 0.7)
	paymentService := NewService("PaymentService", registry.Get("PaymentService"), 0.8)
	
	// Execute requests across services
	services := []*Service{userService, orderService, paymentService}
	
	for round := 1; round <= 3; round++ {
		fmt.Printf("\nRound %d:\n", round)
		
		for _, service := range services {
			err := service.Call()
			if err != nil {
				fmt.Printf("Error in %s: %v\n", service.name, err)
			}
		}
		
		// Show registry status
		fmt.Println("\nRegistry Status:")
		for name, status := range registry.GetAllStatus() {
			fmt.Printf("%s: %s\n", name, status)
		}
		
		time.Sleep(1 * time.Second)
	}
}

func demonstrateTimeoutCircuitBreaker() {
	fmt.Println("\n--- Timeout Circuit Breaker Demo ---")
	
	cb := NewTimeoutCircuitBreaker(2, 1, 100, 1) // 100ms timeout
	
	service := NewService("SlowService", cb, 0.3)
	
	// Override service call to simulate slow response
	slowService := struct {
		name           string
		circuitBreaker CircuitBreaker
	}{
		name:           "SlowService",
		circuitBreaker: cb,
	}
	
	slowService.Call = func() error {
		task := func() error {
			fmt.Printf("SlowService: Processing request (will take 150ms)\n")
			time.Sleep(150 * time.Millisecond) // Longer than timeout
			return nil
		}
		return slowService.circuitBreaker.Execute(task)
	}
	
	// Execute requests
	for i := 1; i <= 5; i++ {
		fmt.Printf("\nRequest %d:\n", i)
		err := slowService.Call()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		
		stats := cb.GetStats()
		fmt.Printf("Stats: State=%v, Failures=%d, Successes=%d\n", 
			stats.State, stats.FailureCount, stats.SuccessCount)
		
		time.Sleep(200 * time.Millisecond)
	}
}

func demonstrateRetryCircuitBreaker() {
	fmt.Println("\n--- Retry Circuit Breaker Demo ---")
	
	cb := NewRetryCircuitBreaker(2, 1, 1, 3) // 3 retries
	
	service := NewService("FlakyService", cb, 0.8) // High failure rate, but retries help
	
	// Execute requests
	for i := 1; i <= 6; i++ {
		fmt.Printf("\nRequest %d:\n", i)
		err := service.Call()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println("Success!")
		}
		
		stats := cb.GetStats()
		fmt.Printf("Stats: State=%v, Failures=%d, Successes=%d\n", 
			stats.State, stats.FailureCount, stats.SuccessCount)
		
		time.Sleep(300 * time.Millisecond)
	}
}

func main() {
	fmt.Println("=== Circuit Breaker Pattern Demo ===")
	
	demonstrateBasicCircuitBreaker()
	demonstrateAdvancedCircuitBreaker()
	demonstrateCircuitBreakerRegistry()
	demonstrateTimeoutCircuitBreaker()
	demonstrateRetryCircuitBreaker()
	
	fmt.Println("\nAll circuit breaker patterns demonstrated successfully!")
}
