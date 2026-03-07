package main

import (
	"fmt"
	"sync"
	"time"
)

// Bulkhead Pattern

// Bulkhead interface
type Bulkhead interface {
	Execute(task func() error) error
	GetStats() BulkheadStats
}

type BulkheadStats struct {
	ActiveTasks    int
	QueuedTasks    int
	RejectedTasks  int
	CompletedTasks int
	MaxConcurrency int
}

// Semaphore Bulkhead
type SemaphoreBulkhead struct {
	maxConcurrency int
	semaphore      chan struct{}
	stats          BulkheadStats
	mu             sync.Mutex
}

func NewSemaphoreBulkhead(maxConcurrency int) *SemaphoreBulkhead {
	return &SemaphoreBulkhead{
		maxConcurrency: maxConcurrency,
		semaphore:      make(chan struct{}, maxConcurrency),
		stats: BulkheadStats{
			MaxConcurrency: maxConcurrency,
		},
	}
}

func (sb *SemaphoreBulkhead) Execute(task func() error) error {
	// Try to acquire semaphore
	select {
	case sb.semaphore <- struct{}{}:
		// Acquired, execute task
		defer func() { <-sb.semaphore }()
		
		sb.mu.Lock()
		sb.stats.ActiveTasks++
		sb.mu.Unlock()
		
		defer func() {
			sb.mu.Lock()
			sb.stats.ActiveTasks--
			sb.stats.CompletedTasks++
			sb.mu.Unlock()
		}()
		
		return task()
		
	default:
		// Could not acquire, reject task
		sb.mu.Lock()
		sb.stats.RejectedTasks++
		sb.mu.Unlock()
		
		return fmt.Errorf("bulkhead: task rejected - concurrency limit reached")
	}
}

func (sb *SemaphoreBulkhead) GetStats() BulkheadStats {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.stats
}

// Thread Pool Bulkhead
type ThreadPoolBulkhead struct {
	maxConcurrency int
	taskQueue       chan func()
	workers         []*Worker
	stats           BulkheadStats
	mu              sync.Mutex
	wg              sync.WaitGroup
	done            chan struct{}
}

type Worker struct {
	id       int
	taskQueue <-chan func()
	done     <-chan struct{}
}

func NewWorker(id int, taskQueue <-chan func(), done <-chan struct{}) *Worker {
	return &Worker{
		id:       id,
		taskQueue: taskQueue,
		done:     done,
	}
}

func (w *Worker) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case task := <-w.taskQueue:
				if task != nil {
					task()
				}
			case <-w.done:
				return
			}
		}
	}()
}

func NewThreadPoolBulkhead(maxConcurrency, queueSize int) *ThreadPoolBulkhead {
	tpb := &ThreadPoolBulkhead{
		maxConcurrency: maxConcurrency,
		taskQueue:       make(chan func(), queueSize),
		workers:         make([]*Worker, maxConcurrency),
		stats: BulkheadStats{
			MaxConcurrency: maxConcurrency,
		},
		done: make(chan struct{}),
	}
	
	// Create and start workers
	for i := 0; i < maxConcurrency; i++ {
		worker := NewWorker(i, tpb.taskQueue, tpb.done)
		tpb.workers[i] = worker
		worker.Start(&tpb.wg)
	}
	
	return tpb
}

func (tpb *ThreadPoolBulkhead) Execute(task func() error) error {
	resultChan := make(chan error, 1)
	
	wrappedTask := func() {
		resultChan <- task()
	}
	
	select {
	case tpb.taskQueue <- wrappedTask:
		tpb.mu.Lock()
		tpb.stats.ActiveTasks++
		tpb.mu.Unlock()
		
		err := <-resultChan
		
		tpb.mu.Lock()
		tpb.stats.ActiveTasks--
		tpb.stats.CompletedTasks++
		tpb.mu.Unlock()
		
		return err
		
	default:
		tpb.mu.Lock()
		tpb.stats.RejectedTasks++
		tpb.mu.Unlock()
		
		return fmt.Errorf("bulkhead: task rejected - queue full")
	}
}

func (tpb *ThreadPoolBulkhead) GetStats() BulkheadStats {
	tpb.mu.Lock()
	defer tpb.mu.Unlock()
	tpb.stats.QueuedTasks = len(tpb.taskQueue)
	return tpb.stats
}

func (tpb *ThreadPoolBulkhead) Shutdown() {
	close(tpb.done)
	tpb.wg.Wait()
}

// Resource Pool Bulkhead
type ResourcePoolBulkhead struct {
	maxConnections int
	pool           chan *Connection
	stats          BulkheadStats
	mu             sync.Mutex
}

type Connection struct {
	id     int
	active bool
}

func NewConnection(id int) *Connection {
	return &Connection{id: id, active: true}
}

func (c *Connection) Close() {
	c.active = false
}

func NewResourcePoolBulkhead(maxConnections int) *ResourcePoolBulkhead {
	rpb := &ResourcePoolBulkhead{
		maxConnections: maxConnections,
		pool:           make(chan *Connection, maxConnections),
		stats: BulkheadStats{
			MaxConcurrency: maxConnections,
		},
	}
	
	// Initialize pool with connections
	for i := 0; i < maxConnections; i++ {
		rpb.pool <- NewConnection(i)
	}
	
	return rpb
}

func (rpb *ResourcePoolBulkhead) Execute(task func(*Connection) error) error {
	select {
	case conn := <-rpb.pool:
		rpb.mu.Lock()
		rpb.stats.ActiveTasks++
		rpb.mu.Unlock()
		
		defer func() {
			rpb.mu.Lock()
			rpb.stats.ActiveTasks--
			rpb.stats.CompletedTasks++
			rpb.mu.Unlock()
			
			// Return connection to pool
			rpb.pool <- conn
		}()
		
		return task(conn)
		
	default:
		rpb.mu.Lock()
		rpb.stats.RejectedTasks++
		rpb.mu.Unlock()
		
		return fmt.Errorf("bulkhead: task rejected - no available connections")
	}
}

func (rpb *ResourcePoolBulkhead) GetStats() BulkheadStats {
	rpb.mu.Lock()
	defer rpb.mu.Unlock()
	rpb.stats.QueuedTasks = len(rpb.pool)
	return rpb.stats
}

// Service with Bulkhead
type UserService struct {
	bulkhead Bulkhead
}

func NewUserService(bulkhead Bulkhead) *UserService {
	return &UserService{bulkhead: bulkhead}
}

func (us *UserService) GetUser(id string) error {
	task := func() error {
		fmt.Printf("UserService: Processing user %s\n", id)
		time.Sleep(100 * time.Millisecond) // Simulate work
		fmt.Printf("UserService: Completed user %s\n", id)
		return nil
	}
	
	return us.bulkhead.Execute(task)
}

type OrderService struct {
	bulkhead Bulkhead
}

func NewOrderService(bulkhead Bulkhead) *OrderService {
	return &OrderService{bulkhead: bulkhead}
}

func (os *OrderService) ProcessOrder(id string) error {
	task := func() error {
		fmt.Printf("OrderService: Processing order %s\n", id)
		time.Sleep(150 * time.Millisecond) // Simulate work
		fmt.Printf("OrderService: Completed order %s\n", id)
		return nil
	}
	
	return os.bulkhead.Execute(task)
}

type PaymentService struct {
	bulkhead Bulkhead
}

func NewPaymentService(bulkhead Bulkhead) *PaymentService {
	return &PaymentService{bulkhead: bulkhead}
}

func (ps *PaymentService) ProcessPayment(id string) error {
	task := func() error {
		fmt.Printf("PaymentService: Processing payment %s\n", id)
		time.Sleep(200 * time.Millisecond) // Simulate work
		fmt.Printf("PaymentService: Completed payment %s\n", id)
		return nil
	}
	
	return ps.bulkhead.Execute(task)
}

// Database Service with Connection Pool
type DatabaseService struct {
	bulkhead *ResourcePoolBulkhead
}

func NewDatabaseService(maxConnections int) *DatabaseService {
	return &DatabaseService{
		bulkhead: NewResourcePoolBulkhead(maxConnections),
	}
}

func (ds *DatabaseService) Query(query string) error {
	task := func(conn *Connection) error {
		fmt.Printf("DatabaseService: Executing query '%s' with connection %d\n", query, conn.id)
		time.Sleep(50 * time.Millisecond) // Simulate database work
		fmt.Printf("DatabaseService: Completed query '%s'\n", query)
		return nil
	}
	
	return ds.bulkhead.Execute(task)
}

// Bulkhead Manager
type BulkheadManager struct {
	bulkheads map[string]Bulkhead
	mu        sync.RWMutex
}

func NewBulkheadManager() *BulkheadManager {
	return &BulkheadManager{
		bulkheads: make(map[string]Bulkhead),
	}
}

func (bm *BulkheadManager) AddBulkhead(name string, bulkhead Bulkhead) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.bulkheads[name] = bulkhead
	fmt.Printf("BulkheadManager: Added bulkhead '%s'\n", name)
}

func (bm *BulkheadManager) GetBulkhead(name string) Bulkhead {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return bm.bulkheads[name]
}

func (bm *BulkheadManager) GetStats(name string) BulkheadStats {
	bulkhead := bm.GetBulkhead(name)
	if bulkhead != nil {
		return bulkhead.GetStats()
	}
	return BulkheadStats{}
}

func (bm *BulkheadManager) GetAllStats() map[string]BulkheadStats {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	
	stats := make(map[string]BulkheadStats)
	for name, bulkhead := range bm.bulkheads {
		stats[name] = bulkhead.GetStats()
	}
	return stats
}

// Circuit Breaker with Bulkhead
type CircuitBreaker struct {
	maxFailures   int
	resetTimeout  time.Duration
	failures      int
	lastFailTime  time.Time
	state         CircuitBreakerState
	mu            sync.Mutex
}

type CircuitBreakerState int

const (
	Closed CircuitBreakerState = iota
	Open
	HalfOpen
)

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        Closed,
	}
}

func (cb *CircuitBreaker) Execute(task func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	// Check if circuit is open and should be reset
	if cb.state == Open && time.Since(cb.lastFailTime) > cb.resetTimeout {
		cb.state = HalfOpen
		cb.failures = 0
	}
	
	// Reject if circuit is open
	if cb.state == Open {
		return fmt.Errorf("circuit breaker: circuit open")
	}
	
	// Execute task
	err := task()
	
	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()
		
		if cb.failures >= cb.maxFailures {
			cb.state = Open
		}
	} else {
		cb.failures = 0
		cb.state = Closed
	}
	
	return err
}

// Service with both Bulkhead and Circuit Breaker
type RobustService struct {
	bulkhead       Bulkhead
	circuitBreaker *CircuitBreaker
}

func NewRobustService(bulkhead Bulkhead, maxFailures int, resetTimeout time.Duration) *RobustService {
	return &RobustService{
		bulkhead:       bulkhead,
		circuitBreaker: NewCircuitBreaker(maxFailures, resetTimeout),
	}
}

func (rs *RobustService) ProcessRequest(id string) error {
	task := func() error {
		fmt.Printf("RobustService: Processing request %s\n", id)
		time.Sleep(100 * time.Millisecond)
		
		// Simulate occasional failure
		if id == "fail" {
			return fmt.Errorf("simulated failure")
		}
		
		fmt.Printf("RobustService: Completed request %s\n", id)
		return nil
	}
	
	// First go through circuit breaker
	circuitTask := func() error {
		return rs.bulkhead.Execute(task)
	}
	
	return rs.circuitBreaker.Execute(circuitTask)
}

func demonstrateSemaphoreBulkhead() {
	fmt.Println("--- Semaphore Bulkhead Demo ---")
	
	bulkhead := NewSemaphoreBulkhead(3)
	userService := NewUserService(bulkhead)
	
	// Execute tasks concurrently
	var wg sync.WaitGroup
	for i := 1; i <= 8; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := userService.GetUser(fmt.Sprintf("user%d", id))
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}(i)
	}
	
	wg.Wait()
	
	stats := bulkhead.GetStats()
	fmt.Printf("Stats: Active=%d, Completed=%d, Rejected=%d\n", 
		stats.ActiveTasks, stats.CompletedTasks, stats.RejectedTasks)
}

func demonstrateThreadPoolBulkhead() {
	fmt.Println("\n--- Thread Pool Bulkhead Demo ---")
	
	bulkhead := NewThreadPoolBulkhead(2, 5) // 2 workers, queue size 5
	orderService := NewOrderService(bulkhead)
	
	// Execute tasks
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := orderService.ProcessOrder(fmt.Sprintf("order%d", id))
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}(i)
	}
	
	wg.Wait()
	
	stats := bulkhead.GetStats()
	fmt.Printf("Stats: Active=%d, Queued=%d, Completed=%d, Rejected=%d\n", 
		stats.ActiveTasks, stats.QueuedTasks, stats.CompletedTasks, stats.RejectedTasks)
	
	bulkhead.Shutdown()
}

func demonstrateResourcePoolBulkhead() {
	fmt.Println("\n--- Resource Pool Bulkhead Demo ---")
	
	dbService := NewDatabaseService(2) // 2 connections
	
	// Execute queries
	var wg sync.WaitGroup
	for i := 1; i <= 6; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := dbService.Query(fmt.Sprintf("SELECT * FROM table%d", id))
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}(i)
	}
	
	wg.Wait()
	
	stats := dbService.bulkhead.GetStats()
	fmt.Printf("Stats: Active=%d, Available=%d, Completed=%d, Rejected=%d\n", 
		stats.ActiveTasks, stats.QueuedTasks, stats.CompletedTasks, stats.RejectedTasks)
}

func demonstrateBulkheadManager() {
	fmt.Println("\n--- Bulkhead Manager Demo ---")
	
	manager := NewBulkheadManager()
	
	// Add different bulkheads for different services
	manager.AddBulkhead("UserService", NewSemaphoreBulkhead(2))
	manager.AddBulkhead("OrderService", NewThreadPoolBulkhead(3, 10))
	manager.AddBulkhead("PaymentService", NewSemaphoreBulkhead(1))
	
	userService := NewUserService(manager.GetBulkhead("UserService"))
	orderService := NewOrderService(manager.GetBulkhead("OrderService"))
	paymentService := NewPaymentService(manager.GetBulkhead("PaymentService"))
	
	// Execute tasks across different services
	var wg sync.WaitGroup
	
	// User service tasks
	for i := 1; i <= 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			userService.GetUser(fmt.Sprintf("user%d", id))
		}(i)
	}
	
	// Order service tasks
	for i := 1; i <= 6; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			orderService.ProcessOrder(fmt.Sprintf("order%d", id))
		}(i)
	}
	
	// Payment service tasks
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			paymentService.ProcessPayment(fmt.Sprintf("payment%d", id))
		}(i)
	}
	
	wg.Wait()
	
	// Show stats for all bulkheads
	fmt.Println("\nBulkhead Stats:")
	for name, stats := range manager.GetAllStats() {
		fmt.Printf("%s: Active=%d, Completed=%d, Rejected=%d\n", 
			name, stats.ActiveTasks, stats.CompletedTasks, stats.RejectedTasks)
	}
}

func demonstrateCircuitBreakerWithBulkhead() {
	fmt.Println("\n--- Circuit Breaker with Bulkhead Demo ---")
	
	bulkhead := NewSemaphoreBulkhead(2)
	service := NewRobustService(bulkhead, 3, 2*time.Second)
	
	// Execute some requests
	requests := []string{"req1", "req2", "fail", "fail", "fail", "req3", "req4"}
	
	for _, req := range requests {
		err := service.ProcessRequest(req)
		if err != nil {
			fmt.Printf("Request %s failed: %v\n", req, err)
		} else {
			fmt.Printf("Request %s succeeded\n", req)
		}
		time.Sleep(100 * time.Millisecond)
	}
	
	// Wait for circuit breaker to reset
	fmt.Println("Waiting for circuit breaker to reset...")
	time.Sleep(3 * time.Second)
	
	// Try again
	err := service.ProcessRequest("req5")
	if err != nil {
		fmt.Printf("Request req5 failed: %v\n", err)
	} else {
		fmt.Printf("Request req5 succeeded\n")
	}
}

func main() {
	fmt.Println("=== Bulkhead Pattern Demo ===")
	
	demonstrateSemaphoreBulkhead()
	demonstrateThreadPoolBulkhead()
	demonstrateResourcePoolBulkhead()
	demonstrateBulkheadManager()
	demonstrateCircuitBreakerWithBulkhead()
	
	fmt.Println("\nAll bulkhead patterns demonstrated successfully!")
}
