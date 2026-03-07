package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task represents a unit of work
type Task struct {
	ID   int
	Data string
}

// Result represents the outcome of a task
type Result struct {
	TaskID int
	Output string
	Error  error
	Duration time.Duration
}

// Worker represents a worker that processes tasks
type Worker struct {
	id       int
	taskChan <-chan Task
	resultChan chan<- Result
	quit     chan bool
}

func NewWorker(id int, taskChan <-chan Task, resultChan chan<- Result) *Worker {
	return &Worker{
		id:        id,
		taskChan:  taskChan,
		resultChan: resultChan,
		quit:      make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case task := <-w.taskChan:
				result := w.processTask(task)
				w.resultChan <- result
			case <-w.quit:
				fmt.Printf("Worker %d shutting down\n", w.id)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.quit <- true
}

func (w *Worker) processTask(task Task) Result {
	start := time.Now()
	
	// Simulate work with random duration
	workTime := time.Duration(rand.Intn(500)+100) * time.Millisecond
	time.Sleep(workTime)
	
	// Simulate occasional failure
	var err error
	if rand.Intn(10) == 0 {
		err = fmt.Errorf("random failure in worker %d", w.id)
	}
	
	output := fmt.Sprintf("Task %d processed by worker %d: %s", 
		task.ID, w.id, task.Data)
	
	return Result{
		TaskID:   task.ID,
		Output:   output,
		Error:    err,
		Duration: time.Since(start),
	}
}

// WorkerPool manages a pool of workers
type WorkerPool struct {
	workers    []*Worker
	taskChan   chan Task
	resultChan chan Result
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int, taskBufferSize int) *WorkerPool {
	taskChan := make(chan Task, taskBufferSize)
	resultChan := make(chan Result, taskBufferSize)
	
	pool := &WorkerPool{
		taskChan:   taskChan,
		resultChan: resultChan,
		workers:    make([]*Worker, numWorkers),
	}
	
	for i := 0; i < numWorkers; i++ {
		pool.workers[i] = NewWorker(i+1, taskChan, resultChan)
	}
	
	return pool
}

func (wp *WorkerPool) Start() {
	for _, worker := range wp.workers {
		worker.Start()
	}
	fmt.Printf("Worker pool started with %d workers\n", len(wp.workers))
}

func (wp *WorkerPool) Stop() {
	for _, worker := range wp.workers {
		worker.Stop()
	}
	close(wp.taskChan)
	close(wp.resultChan)
}

func (wp *WorkerPool) SubmitTask(task Task) {
	wp.wg.Add(1)
	go func() {
		defer wp.wg.Done()
		wp.taskChan <- task
	}()
}

func (wp *WorkerPool) GetResults() <-chan Result {
	return wp.resultChan
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// Advanced Worker Pool with dynamic sizing
type DynamicWorkerPool struct {
	minWorkers    int
	maxWorkers    int
	currentWorkers int
	taskChan      chan Task
	resultChan    chan Result
	workers       []*Worker
	mu            sync.RWMutex
	quit          chan bool
}

func NewDynamicWorkerPool(minWorkers, maxWorkers, taskBufferSize int) *DynamicWorkerPool {
	return &DynamicWorkerPool{
		minWorkers:     minWorkers,
		maxWorkers:     maxWorkers,
		currentWorkers: minWorkers,
		taskChan:       make(chan Task, taskBufferSize),
		resultChan:     make(chan Result, taskBufferSize),
		workers:        make([]*Worker, 0),
		quit:           make(chan bool),
	}
}

func (dwp *DynamicWorkerPool) Start() {
	dwp.mu.Lock()
	defer dwp.mu.Unlock()
	
	// Start with minimum workers
	for i := 0; i < dwp.minWorkers; i++ {
		worker := NewWorker(i+1, dwp.taskChan, dwp.resultChan)
		dwp.workers = append(dwp.workers, worker)
		worker.Start()
	}
	
	// Start load balancer
	go dwp.loadBalancer()
	
	fmt.Printf("Dynamic worker pool started with %d workers (min: %d, max: %d)\n", 
		dwp.currentWorkers, dwp.minWorkers, dwp.maxWorkers)
}

func (dwp *DynamicWorkerPool) Stop() {
	dwp.quit <- true
	
	dwp.mu.Lock()
	defer dwp.mu.Unlock()
	
	for _, worker := range dwp.workers {
		worker.Stop()
	}
	
	close(dwp.taskChan)
	close(dwp.resultChan)
}

func (dwp *DynamicWorkerPool) SubmitTask(task Task) {
	dwp.taskChan <- task
}

func (dwp *DynamicWorkerPool) GetResults() <-chan Result {
	return dwp.resultChan
}

func (dwp *DynamicWorkerPool) loadBalancer() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	
	queueLength := 0
	
	for {
		select {
		case <-ticker.C:
			newQueueLength := len(dwp.taskChan)
			
			// Scale up if queue is growing and we can add more workers
			if newQueueLength > queueLength && dwp.currentWorkers < dwp.maxWorkers {
				dwp.scaleUp()
			}
			// Scale down if queue is empty and we have more than minimum workers
			else if newQueueLength == 0 && dwp.currentWorkers > dwp.minWorkers {
				dwp.scaleDown()
			}
			
			queueLength = newQueueLength
			
		case <-dwp.quit:
			return
		}
	}
}

func (dwp *DynamicWorkerPool) scaleUp() {
	dwp.mu.Lock()
	defer dwp.mu.Unlock()
	
	if dwp.currentWorkers >= dwp.maxWorkers {
		return
	}
	
	worker := NewWorker(dwp.currentWorkers+1, dwp.taskChan, dwp.resultChan)
	dwp.workers = append(dwp.workers, worker)
	worker.Start()
	dwp.currentWorkers++
	
	fmt.Printf("Scaled up to %d workers\n", dwp.currentWorkers)
}

func (dwp *DynamicWorkerPool) scaleDown() {
	dwp.mu.Lock()
	defer dwp.mu.Unlock()
	
	if dwp.currentWorkers <= dwp.minWorkers {
		return
	}
	
	// Remove last worker
	worker := dwp.workers[len(dwp.workers)-1]
	worker.Stop()
	dwp.workers = dwp.workers[:len(dwp.workers)-1]
	dwp.currentWorkers--
	
	fmt.Printf("Scaled down to %d workers\n", dwp.currentWorkers)
}

// Priority Worker Pool
type PriorityTask struct {
	Task
	Priority int
}

type PriorityWorkerPool struct {
	workers      []*Worker
	highPriority chan PriorityTask
	lowPriority  chan PriorityTask
	resultChan   chan Result
}

func NewPriorityWorkerPool(numWorkers int) *PriorityWorkerPool {
	return &PriorityWorkerPool{
		workers:      make([]*Worker, numWorkers),
		highPriority: make(chan PriorityTask, 100),
		lowPriority:  make(chan PriorityTask, 100),
		resultChan:   make(chan Result, 100),
	}
}

func (pwp *PriorityWorkerPool) Start() {
	for i := 0; i < len(pwp.workers); i++ {
		worker := NewWorker(i+1, pwp.taskRouter(), pwp.resultChan)
		pwp.workers[i] = worker
		worker.Start()
	}
	
	fmt.Printf("Priority worker pool started with %d workers\n", len(pwp.workers))
}

func (pwp *PriorityWorkerPool) taskRouter() <-chan Task {
	taskChan := make(chan Task)
	
	go func() {
		for {
			select {
			case highTask := <-pwp.highPriority:
				taskChan <- highTask.Task
			case lowTask := <-pwp.lowPriority:
				taskChan <- lowTask.Task
			}
		}
	}()
	
	return taskChan
}

func (pwp *PriorityWorkerPool) SubmitHighPriorityTask(task Task) {
	pwp.highPriority <- PriorityTask{Task: task, Priority: 1}
}

func (pwp *PriorityWorkerPool) SubmitLowPriorityTask(task Task) {
	pwp.lowPriority <- PriorityTask{Task: task, Priority: 0}
}

func (pwp *PriorityWorkerPool) GetResults() <-chan Result {
	return pwp.resultChan
}

func (pwp *PriorityWorkerPool) Stop() {
	for _, worker := range pwp.workers {
		worker.Stop()
	}
	close(pwp.highPriority)
	close(pwp.lowPriority)
	close(pwp.resultChan)
}

func main() {
	fmt.Println("=== Worker Pool Pattern Demo ===")
	
	rand.Seed(time.Now().UnixNano())
	
	// Basic Worker Pool example
	fmt.Println("\n--- Basic Worker Pool Example ---")
	basicPool := NewWorkerPool(3, 10)
	basicPool.Start()
	
	// Submit tasks
	for i := 1; i <= 10; i++ {
		task := Task{
			ID:   i,
			Data: fmt.Sprintf("Task data %d", i),
		}
		basicPool.SubmitTask(task)
	}
	
	// Collect results
	go func() {
		for result := range basicPool.GetResults() {
			if result.Error != nil {
				fmt.Printf("❌ %s (took %v)\n", result.Error, result.Duration)
			} else {
				fmt.Printf("✅ %s (took %v)\n", result.Output, result.Duration)
			}
		}
	}()
	
	basicPool.Wait()
	basicPool.Stop()
	
	// Dynamic Worker Pool example
	fmt.Println("\n--- Dynamic Worker Pool Example ---")
	dynamicPool := NewDynamicWorkerPool(2, 6, 20)
	dynamicPool.Start()
	
	// Submit burst of tasks
	go func() {
		for i := 1; i <= 20; i++ {
			task := Task{
				ID:   i,
				Data: fmt.Sprintf("Dynamic task %d", i),
			}
			dynamicPool.SubmitTask(task)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	// Collect results for a while
	go func() {
		for i := 0; i < 20; i++ {
			select {
			case result := <-dynamicPool.GetResults():
				if result.Error != nil {
					fmt.Printf("❌ %s\n", result.Error)
				} else {
					fmt.Printf("✅ %s\n", result.Output)
				}
			case <-time.After(100 * time.Millisecond):
				continue
			}
		}
	}()
	
	time.Sleep(8 * time.Second) // Let it run for a bit to see scaling
	dynamicPool.Stop()
	
	// Priority Worker Pool example
	fmt.Println("\n--- Priority Worker Pool Example ---")
	priorityPool := NewPriorityWorkerPool(2)
	priorityPool.Start()
	
	// Submit mixed priority tasks
	go func() {
		// Submit low priority tasks first
		for i := 1; i <= 5; i++ {
			task := Task{
				ID:   i,
				Data: fmt.Sprintf("Low priority task %d", i),
			}
			priorityPool.SubmitLowPriorityTask(task)
		}
		
		time.Sleep(100 * time.Millisecond)
		
		// Submit high priority tasks
		for i := 6; i <= 8; i++ {
			task := Task{
				ID:   i,
				Data: fmt.Sprintf("High priority task %d", i),
			}
			priorityPool.SubmitHighPriorityTask(task)
		}
		
		// Submit more low priority tasks
		for i := 9; i <= 10; i++ {
			task := Task{
				ID:   i,
				Data: fmt.Sprintf("Low priority task %d", i),
			}
			priorityPool.SubmitLowPriorityTask(task)
		}
	}()
	
	// Collect results
	go func() {
		for i := 0; i < 10; i++ {
			result := <-priorityPool.GetResults()
			if result.Error != nil {
				fmt.Printf("❌ %s\n", result.Error)
			} else {
				fmt.Printf("✅ %s\n", result.Output)
			}
		}
	}()
	
	time.Sleep(3 * time.Second)
	priorityPool.Stop()
	
	fmt.Println("\nAll worker pool patterns demonstrated successfully!")
}
