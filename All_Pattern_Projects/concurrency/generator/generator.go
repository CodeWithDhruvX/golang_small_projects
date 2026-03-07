package main

import (
	"fmt"
	"time"
)

// Generator function that produces a sequence of values
func fibonacciGenerator() func() int {
	a, b := 0, 1
	
	return func() int {
		result := a
		a, b = b, a+b
		return result
	}
}

// Channel-based generator
func channelFibonacci(n int) <-chan int {
	ch := make(chan int)
	
	go func() {
		defer close(ch)
		a, b := 0, 1
		for i := 0; i < n; i++ {
			ch <- a
			a, b = b, a+b
		}
	}()
	
	return ch
}

// Generator with cancellation
func cancellableGenerator(done <-chan struct{}) <-chan int {
	ch := make(chan int)
	
	go func() {
		defer close(ch)
		for i := 0; ; i++ {
			select {
			case <-done:
				fmt.Println("Generator cancelled")
				return
			case ch <- i:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	
	return ch
}

// Prime number generator
func primeGenerator() <-chan int {
	ch := make(chan int)
	
	go func() {
		defer close(ch)
		
		// Sieve of Eratosthenes
		numbers := make(chan int)
		
		// Generate numbers starting from 2
		go func() {
			defer close(numbers)
			for i := 2; ; i++ {
				numbers <- i
			}
		}()
		
		// Sieve process
		for {
			prime := <-numbers
			ch <- prime
			
			// Filter out multiples of the prime
			filtered := make(chan int)
			go func() {
				defer close(filtered)
				for num := range numbers {
					if num%prime != 0 {
						filtered <- num
					}
				}
			}()
			numbers = filtered
		}
	}()
	
	return ch
}

// Limited prime generator (generates first n primes)
func limitedPrimeGenerator(n int) <-chan int {
	ch := make(chan int)
	
	go func() {
		defer close(ch)
		count := 0
		candidate := 2
		
		for count < n {
			if isPrime(candidate) {
				ch <- candidate
				count++
			}
			candidate++
		}
	}()
	
	return ch
}

func isPrime(num int) bool {
	if num <= 1 {
		return false
	}
	if num <= 3 {
		return true
	}
	if num%2 == 0 || num%3 == 0 {
		return false
	}
	
	i := 5
	for i*i <= num {
		if num%i == 0 || num%(i+2) == 0 {
			return false
		}
		i += 6
	}
	
	return true
}

// Random number generator
func randomGenerator(min, max int) <-chan int {
	ch := make(chan int)
	
	go func() {
		defer close(ch)
		for {
			// Simple random number generation
			num := (num*1103515245 + 12345) % (1 << 31)
			result := min + (num % (max - min + 1))
			ch <- result
			time.Sleep(200 * time.Millisecond)
		}
	}()
	
	return ch
}

// Temperature sensor generator
func temperatureSensorGenerator() <-chan float64 {
	ch := make(chan float64)
	
	go func() {
		defer close(ch)
		baseTemp := 20.0
		
		for {
			// Simulate temperature fluctuation
			variation := (float64((baseTemp*100)+12345)%1000 - 500) / 100.0
			temp := baseTemp + variation
			ch <- temp
			time.Sleep(500 * time.Millisecond)
		}
	}()
	
	return ch
}

// Generator pipeline example
func generatorPipeline() {
	fmt.Println("--- Generator Pipeline Example ---")
	
	// Stage 1: Generate numbers
	numbers := numberGenerator(10)
	
	// Stage 2: Square numbers
	squares := squareGenerator(numbers)
	
	// Stage 3: Filter even squares
	evenSquares := filterGenerator(squares, func(x int) bool {
		return x%2 == 0
	})
	
	// Consume final results
	for result := range evenSquares {
		fmt.Printf("Pipeline result: %d\n", result)
	}
}

func numberGenerator(n int) <-chan int {
	ch := make(chan int)
	
	go func() {
		defer close(ch)
		for i := 1; i <= n; i++ {
			ch <- i
		}
	}()
	
	return ch
}

func squareGenerator(input <-chan int) <-chan int {
	ch := make(chan int)
	
	go func() {
		defer close(ch)
		for num := range input {
			ch <- num * num
		}
	}()
	
	return ch
}

func filterGenerator(input <-chan int, predicate func(int) bool) <-chan int {
	ch := make(chan int)
	
	go func() {
		defer close(ch)
		for num := range input {
			if predicate(num) {
				ch <- num
			}
		}
	}()
	
	return ch
}

var num int = 12345 // Simple seed for random number generation

func main() {
	fmt.Println("=== Generator Pattern Demo ===")
	
	// Function-based generator
	fmt.Println("\n--- Function-based Fibonacci Generator ---")
	fib := fibonacciGenerator()
	for i := 0; i < 10; i++ {
		fmt.Printf("Fibonacci %d: %d\n", i, fib())
	}
	
	// Channel-based generator
	fmt.Println("\n--- Channel-based Fibonacci Generator ---")
	for fib := range channelFibonacci(10) {
		fmt.Printf("Fibonacci: %d\n", fib)
	}
	
	// Cancellable generator
	fmt.Println("\n--- Cancellable Generator ---")
	done := make(chan struct{})
	gen := cancellableGenerator(done)
	
	for i := 0; i < 5; i++ {
		fmt.Printf("Generated: %d\n", <-gen)
	}
	close(done) // Cancel the generator
	
	// Prime number generator (limited to first 10 primes)
	fmt.Println("\n--- Prime Number Generator ---")
	for prime := range limitedPrimeGenerator(10) {
		fmt.Printf("Prime: %d\n", prime)
	}
	
	// Random number generator
	fmt.Println("\n--- Random Number Generator ---")
	random := randomGenerator(1, 100)
	for i := 0; i < 5; i++ {
		fmt.Printf("Random: %d\n", <-random)
	}
	
	// Temperature sensor generator
	fmt.Println("\n--- Temperature Sensor Generator ---")
	tempSensor := temperatureSensorGenerator()
	for i := 0; i < 5; i++ {
		temp := <-tempSensor
		fmt.Printf("Temperature: %.2f°C\n", temp)
	}
	
	// Generator pipeline
	generatorPipeline()
	
	fmt.Println("\nAll generator patterns demonstrated successfully!")
}
