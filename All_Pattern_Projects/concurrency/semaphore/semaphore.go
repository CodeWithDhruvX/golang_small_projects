package main

import (
	"fmt"
	"sync"
	"time"
)

// Semaphore implementation using channels
type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(capacity int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, capacity),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s *Semaphore) AcquireWithTimeout(timeout time.Duration) bool {
	select {
	case s.ch <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Worker function that uses semaphore
func worker(id int, semaphore *Semaphore, wg *sync.WaitGroup) {
	defer wg.Done()
	
	fmt.Printf("Worker %d: Waiting to acquire semaphore\n", id)
	semaphore.Acquire()
	
	fmt.Printf("Worker %d: Acquired semaphore, starting work\n", id)
	time.Sleep(time.Duration(500+id*100) * time.Millisecond)
	fmt.Printf("Worker %d: Completed work\n", id)
	
	semaphore.Release()
	fmt.Printf("Worker %d: Released semaphore\n", id)
}

// Resource Pool example using semaphore
type ResourcePool struct {
	resources []string
	semaphore *Semaphore
	mu        sync.Mutex
}

func NewResourcePool(resources []string) *ResourcePool {
	return &ResourcePool{
		resources: resources,
		semaphore: NewSemaphore(len(resources)),
	}
}

func (rp *ResourcePool) AcquireResource() (string, bool) {
	rp.semaphore.Acquire()
	
	rp.mu.Lock()
	defer rp.mu.Unlock()
	
	if len(rp.resources) == 0 {
		rp.semaphore.Release()
		return "", false
	}
	
	resource := rp.resources[0]
	rp.resources = rp.resources[1:]
	return resource, true
}

func (rp *ResourcePool) ReleaseResource(resource string) {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	
	rp.resources = append(rp.resources, resource)
	rp.semaphore.Release()
}

func resourceUser(id int, pool *ResourcePool, wg *sync.WaitGroup) {
	defer wg.Done()
	
	fmt.Printf("User %d: Trying to acquire resource\n", id)
	resource, ok := pool.AcquireResource()
	
	if !ok {
		fmt.Printf("User %d: No resources available\n", id)
		return
	}
	
	fmt.Printf("User %d: Acquired resource %s\n", id, resource)
	time.Sleep(time.Duration(300+id*100) * time.Millisecond)
	
	pool.ReleaseResource(resource)
	fmt.Printf("User %d: Released resource %s\n", id, resource)
}

// Rate Limiter using semaphore
type RateLimiter struct {
	semaphore *Semaphore
	ticker    *time.Ticker
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		semaphore: NewSemaphore(rate),
		ticker:    time.NewTicker(window),
	}
	
	// Reset semaphore periodically
	go func() {
		for range rl.ticker.C {
			// Drain and refill semaphore
			for {
				if !rl.semaphore.TryAcquire() {
					break
				}
			}
			// Refill to capacity
			for i := 0; i < rate; i++ {
				rl.semaphore.Release()
			}
		}
	}()
	
	return rl
}

func (rl *RateLimiter) Allow() bool {
	return rl.semaphore.TryAcquire()
}

func rateLimitedRequest(id int, limiter *RateLimiter, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for i := 0; i < 3; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d-%d: Allowed at %s\n", id, i, time.Now().Format("15:04:05.000"))
		} else {
			fmt.Printf("Request %d-%d: Rate limited at %s\n", id, i, time.Now().Format("15:04:05.000"))
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// Binary Semaphore (Mutex)
type BinarySemaphore struct {
	ch chan struct{}
}

func NewBinarySemaphore() *BinarySemaphore {
	bs := &BinarySemaphore{
		ch: make(chan struct{}, 1),
	}
	bs.ch <- struct{}{} // Initially available
	return bs
}

func (bs *BinarySemaphore) Lock() {
	<-bs.ch
}

func (bs *BinarySemaphore) Unlock() {
	bs.ch <- struct{}{}
}

func (bs *BinarySemaphore) TryLock() bool {
	select {
	case <-bs.ch:
		return true
	default:
		return false
	}
}

func criticalSection(id int, mutex *BinarySemaphore, wg *sync.WaitGroup) {
	defer wg.Done()
	
	fmt.Printf("Goroutine %d: Attempting to enter critical section\n", id)
	
	if mutex.TryLock() {
		fmt.Printf("Goroutine %d: Entered critical section\n", id)
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("Goroutine %d: Leaving critical section\n", id)
		mutex.Unlock()
	} else {
		fmt.Printf("Goroutine %d: Could not acquire lock, skipping\n", id)
	}
}

// Counting Semaphore with timeout
type CountingSemaphore struct {
	ch chan struct{}
}

func NewCountingSemaphore(capacity int) *CountingSemaphore {
	return &CountingSemaphore{
		ch: make(chan struct{}, capacity),
	}
}

func (cs *CountingSemaphore) Acquire() {
	cs.ch <- struct{}{}
}

func (cs *CountingSemaphore) Release() {
	<-cs.ch
}

func (cs *CountingSemaphore) AcquireWithTimeout(timeout time.Duration) bool {
	select {
	case cs.ch <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

func timeoutWorker(id int, semaphore *CountingSemaphore, wg *sync.WaitGroup) {
	defer wg.Done()
	
	fmt.Printf("Timeout Worker %d: Waiting to acquire semaphore\n", id)
	
	if semaphore.AcquireWithTimeout(300 * time.Millisecond) {
		fmt.Printf("Timeout Worker %d: Acquired semaphore, working\n", id)
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("Timeout Worker %d: Completed work\n", id)
		semaphore.Release()
	} else {
		fmt.Printf("Timeout Worker %d: Timeout waiting for semaphore\n", id)
	}
}

func main() {
	fmt.Println("=== Semaphore Pattern Demo ===")
	
	// Basic Semaphore example
	fmt.Println("\n--- Basic Semaphore Example ---")
	const maxWorkers = 3
	const totalWorkers = 8
	
	semaphore := NewSemaphore(maxWorkers)
	
	var wg sync.WaitGroup
	wg.Add(totalWorkers)
	
	for i := 1; i <= totalWorkers; i++ {
		go worker(i, semaphore, &wg)
	}
	
	wg.Wait()
	
	// Resource Pool example
	fmt.Println("\n--- Resource Pool Example ---")
	resources := []string{"Resource1", "Resource2", "Resource3"}
	pool := NewResourcePool(resources)
	
	var poolWg sync.WaitGroup
	poolWg.Add(6)
	
	for i := 1; i <= 6; i++ {
		go resourceUser(i, pool, &poolWg)
	}
	
	poolWg.Wait()
	
	// Rate Limiter example
	fmt.Println("\n--- Rate Limiter Example ---")
	limiter := NewRateLimiter(3, time.Second)
	
	var rateWg sync.WaitGroup
	rateWg.Add(3)
	
	for i := 1; i <= 3; i++ {
		go rateLimitedRequest(i, limiter, &rateWg)
	}
	
	rateWg.Wait()
	limiter.ticker.Stop()
	
	// Binary Semaphore (Mutex) example
	fmt.Println("\n--- Binary Semaphore (Mutex) Example ---")
	mutex := NewBinarySemaphore()
	
	var mutexWg sync.WaitGroup
	mutexWg.Add(5)
	
	for i := 1; i <= 5; i++ {
		go criticalSection(i, mutex, &mutexWg)
	}
	
	mutexWg.Wait()
	
	// Counting Semaphore with timeout example
	fmt.Println("\n--- Counting Semaphore with Timeout Example ---")
	countSemaphore := NewCountingSemaphore(2)
	
	var timeoutWg sync.WaitGroup
	timeoutWg.Add(5)
	
	for i := 1; i <= 5; i++ {
		go timeoutWorker(i, countSemaphore, &timeoutWg)
	}
	
	timeoutWg.Wait()
	
	fmt.Println("\nAll semaphore patterns demonstrated successfully!")
}
