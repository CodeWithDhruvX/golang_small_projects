package main

import (
	"fmt"
	"sync"
	"time"
)

// Barrier allows multiple goroutines to wait for each other to reach a certain point
type Barrier struct {
	count    int
	waiting  int
	mu       sync.Mutex
	cond     *sync.Cond
	done     bool
}

func NewBarrier(count int) *Barrier {
	b := &Barrier{
		count:   count,
		waiting: 0,
	}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (b *Barrier) Wait() {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.waiting++
	
	if b.waiting == b.count {
		// Last goroutine to arrive releases all waiting goroutines
		b.done = true
		b.cond.Broadcast()
	} else {
		// Wait for the last goroutine to arrive
		for !b.done {
			b.cond.Wait()
		}
	}
	
	b.waiting--
	if b.waiting == 0 {
		b.done = false
	}
}

// Worker function that uses the barrier
func worker(id int, barrier *Barrier, wg *sync.WaitGroup) {
	defer wg.Done()
	
	fmt.Printf("Worker %d: Starting work\n", id)
	time.Sleep(time.Duration(id*200) * time.Millisecond) // Simulate different work durations
	
	fmt.Printf("Worker %d: Reached barrier, waiting for others\n", id)
	barrier.Wait()
	
	fmt.Printf("Worker %d: Passed barrier, continuing work\n", id)
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("Worker %d: Completed work\n", id)
}

// Alternative implementation using channels
type ChannelBarrier struct {
	count   int
	ready   chan struct{}
	done    chan struct{}
	waiting int
	mu      sync.Mutex
}

func NewChannelBarrier(count int) *ChannelBarrier {
	return &ChannelBarrier{
		count: count,
		ready: make(chan struct{}),
		done:  make(chan struct{}),
	}
}

func (cb *ChannelBarrier) Wait() {
	cb.mu.Lock()
	cb.waiting++
	
	if cb.waiting == cb.count {
		// Last goroutine to arrive
		cb.mu.Unlock()
		close(cb.ready) // Signal that all are ready
	} else {
		cb.mu.Unlock()
		<-cb.ready // Wait for all to be ready
	}
	
	// After passing barrier, signal completion
	cb.done <- struct{}{}
	
	// Reset for next use
	cb.mu.Lock()
	cb.waiting--
	if cb.waiting == 0 {
		cb.ready = make(chan struct{})
		cb.done = make(chan struct{})
	}
	cb.mu.Unlock()
}

func channelWorker(id int, barrier *ChannelBarrier, wg *sync.WaitGroup) {
	defer wg.Done()
	
	fmt.Printf("Channel Worker %d: Starting work\n", id)
	time.Sleep(time.Duration(id*150) * time.Millisecond)
	
	fmt.Printf("Channel Worker %d: Reached barrier, waiting for others\n", id)
	barrier.Wait()
	
	fmt.Printf("Channel Worker %d: Passed barrier, continuing work\n", id)
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("Channel Worker %d: Completed work\n", id)
}

// Phase Barrier - allows multiple phases of synchronization
type PhaseBarrier struct {
	count    int
	phase    int
	waiting  int
	mu       sync.Mutex
	cond     *sync.Cond
}

func NewPhaseBarrier(count int) *PhaseBarrier {
	b := &PhaseBarrier{
		count: count,
		phase: 0,
	}
	b.cond = sync.NewCond(&b.mu)
	return b
}

func (pb *PhaseBarrier) WaitPhase(phase int) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	
	if phase != pb.phase {
		panic(fmt.Sprintf("Phase mismatch: expected %d, got %d", pb.phase, phase))
	}
	
	pb.waiting++
	
	if pb.waiting == pb.count {
		// All goroutines have reached this phase
		pb.phase++
		pb.waiting = 0
		pb.cond.Broadcast()
	} else {
		// Wait for others
		pb.cond.Wait()
	}
}

func phaseWorker(id int, barrier *PhaseBarrier, wg *sync.WaitGroup) {
	defer wg.Done()
	
	// Phase 1: Initial setup
	fmt.Printf("Phase Worker %d: Phase 1 - Initial setup\n", id)
	time.Sleep(time.Duration(id*100) * time.Millisecond)
	barrier.WaitPhase(0)
	
	// Phase 2: Processing
	fmt.Printf("Phase Worker %d: Phase 2 - Processing\n", id)
	time.Sleep(time.Duration((id+1)*100) * time.Millisecond)
	barrier.WaitPhase(1)
	
	// Phase 3: Cleanup
	fmt.Printf("Phase Worker %d: Phase 3 - Cleanup\n", id)
	time.Sleep(time.Duration((3-id)*100) * time.Millisecond)
	barrier.WaitPhase(2)
	
	fmt.Printf("Phase Worker %d: All phases completed\n", id)
}

func main() {
	fmt.Println("=== Barrier Pattern Demo ===")
	
	// Basic Barrier example
	fmt.Println("\n--- Basic Barrier Example ---")
	const numWorkers = 4
	barrier := NewBarrier(numWorkers)
	
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	
	for i := 1; i <= numWorkers; i++ {
		go worker(i, barrier, &wg)
	}
	
	wg.Wait()
	
	// Channel Barrier example
	fmt.Println("\n--- Channel Barrier Example ---")
	channelBarrier := NewChannelBarrier(numWorkers)
	
	var channelWg sync.WaitGroup
	channelWg.Add(numWorkers)
	
	for i := 1; i <= numWorkers; i++ {
		go channelWorker(i, channelBarrier, &channelWg)
	}
	
	channelWg.Wait()
	
	// Phase Barrier example
	fmt.Println("\n--- Phase Barrier Example ---")
	phaseBarrier := NewPhaseBarrier(numWorkers)
	
	var phaseWg sync.WaitGroup
	phaseWg.Add(numWorkers)
	
	for i := 1; i <= numWorkers; i++ {
		go phaseWorker(i, phaseBarrier, &phaseWg)
	}
	
	phaseWg.Wait()
	
	fmt.Println("\nAll barrier patterns demonstrated successfully!")
}
