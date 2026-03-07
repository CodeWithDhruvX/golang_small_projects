package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Fan-out: Distribute work to multiple goroutines
func fanOut(jobs <-chan int, workers int, processor func(int) int) <-chan int {
	results := make(chan int, workers)
	
	var wg sync.WaitGroup
	wg.Add(workers)
	
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for job := range jobs {
				result := processor(job)
				results <- result
			}
		}()
	}
	
	go func() {
		wg.Wait()
		close(results)
	}()
	
	return results
}

// Fan-in: Collect results from multiple goroutines into one channel
func fanIn(inputs ...<-chan int) <-chan int {
	output := make(chan int)
	
	var wg sync.WaitGroup
	wg.Add(len(inputs))
	
	for _, input := range inputs {
		go func(ch <-chan int) {
			defer wg.Done()
			for value := range ch {
				output <- value
			}
		}(input)
	}
	
	go func() {
		wg.Wait()
		close(output)
	}()
	
	return output
}

// Example 1: Basic Fan-in/Fan-out
func basicFanInFanOut() {
	fmt.Println("--- Basic Fan-in/Fan-out Example ---")
	
	jobs := make(chan int, 10)
	
	// Generate jobs
	go func() {
		for i := 1; i <= 10; i++ {
			jobs <- i
		}
		close(jobs)
	}()
	
	// Define processor function
	processor := func(job int) int {
		fmt.Printf("Processing job %d\n", job)
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		return job * 2 // Double the input
	}
	
	// Fan-out to 3 workers
	results := fanOut(jobs, 3, processor)
	
	// Collect results
	for result := range results {
		fmt.Printf("Received result: %d\n", result)
	}
}

// Example 2: Multiple producers with Fan-in
func multipleProducersFanIn() {
	fmt.Println("\n--- Multiple Producers Fan-in Example ---")
	
	// Producer 1
	producer1 := func() <-chan int {
		output := make(chan int)
		go func() {
			defer close(output)
			for i := 1; i <= 5; i++ {
				output <- i * 10
				time.Sleep(200 * time.Millisecond)
			}
		}()
		return output
	}
	
	// Producer 2
	producer2 := func() <-chan int {
		output := make(chan int)
		go func() {
			defer close(output)
			for i := 1; i <= 5; i++ {
				output <- i * 100
				time.Sleep(300 * time.Millisecond)
			}
		}()
		return output
	}
	
	// Producer 3
	producer3 := func() <-chan int {
		output := make(chan int)
		go func() {
			defer close(output)
			for i := 1; i <= 5; i++ {
				output <- i * 1000
				time.Sleep(400 * time.Millisecond)
			}
		}()
		return output
	}
	
	// Fan-in from multiple producers
	combined := fanIn(producer1(), producer2(), producer3())
	
	// Process combined results
	for result := range combined {
		fmt.Printf("Combined result: %d\n", result)
	}
}

// Example 3: Pipeline with Fan-in/Fan-out
func pipelineFanInFanOut() {
	fmt.Println("\n--- Pipeline Fan-in/Fan-out Example ---")
	
	// Stage 1: Generate numbers
	generator := func() <-chan int {
		output := make(chan int)
		go func() {
			defer close(output)
			for i := 1; i <= 20; i++ {
				output <- i
				time.Sleep(50 * time.Millisecond)
			}
		}()
		return output
	}
	
	// Stage 2: Square numbers (fan-out)
	squarer := func(input <-chan int) <-chan int {
		return fanOut(input, 4, func(x int) int {
			result := x * x
			fmt.Printf("Squared %d to %d\n", x, result)
			return result
		})
	}
	
	// Stage 3: Filter even numbers (fan-out)
	filter := func(input <-chan int) <-chan int {
		return fanOut(input, 2, func(x int) int {
			if x%2 == 0 {
				return x
			}
			return -1 // Signal to filter out
		})
	}
	
	// Stage 4: Process valid results
	process := func(input <-chan int) {
		for result := range input {
			if result != -1 {
				fmt.Printf("Final result: %d\n", result)
			}
		}
	}
	
	// Build pipeline
	numbers := generator()
	squared := squarer(numbers)
	filtered := filter(squared)
	process(filtered)
}

// Example 4: Worker Pool with Fan-in/Fan-out
type Task struct {
	ID   int
	Data string
}

type Result struct {
	TaskID int
	Output string
	Error  error
}

func workerPoolFanInFanOut() {
	fmt.Println("\n--- Worker Pool Fan-in/Fan-out Example ---")
	
	const numWorkers = 3
	const numTasks = 10
	
	tasks := make(chan Task, numTasks)
	results := make(chan Result, numTasks)
	
	// Fan-out: Create worker pool
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for task := range tasks {
				result := processTask(task, workerID)
				results <- result
			}
		}(i)
	}
	
	// Fan-in: Collect results
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Generate tasks
	go func() {
		defer close(tasks)
		for i := 1; i <= numTasks; i++ {
			tasks <- Task{
				ID:   i,
				Data: fmt.Sprintf("Task %d data", i),
			}
		}
	}()
	
	// Process results
	for result := range results {
		if result.Error != nil {
			fmt.Printf("Task %d failed: %v\n", result.TaskID, result.Error)
		} else {
			fmt.Printf("Task %d completed: %s\n", result.TaskID, result.Output)
		}
	}
}

func processTask(task Task, workerID int) Result {
	// Simulate work
	time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
	
	// Simulate occasional failure
	if rand.Intn(10) == 0 {
		return Result{
			TaskID: task.ID,
			Error:  fmt.Errorf("random failure in worker %d", workerID),
		}
	}
	
	return Result{
		TaskID: task.ID,
		Output: fmt.Sprintf("Processed '%s' by worker %d", task.Data, workerID),
	}
}

func main() {
	fmt.Println("=== Fan-in/Fan-out Pattern Demo ===")
	
	rand.Seed(time.Now().UnixNano())
	
	basicFanInFanOut()
	multipleProducersFanIn()
	pipelineFanInFanOut()
	workerPoolFanInFanOut()
	
	fmt.Println("\nAll fan-in/fan-out patterns demonstrated successfully!")
}
