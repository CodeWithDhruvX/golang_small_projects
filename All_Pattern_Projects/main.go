package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Main application to run all design pattern demonstrations
type PatternRunner struct {
	patterns map[string]PatternInfo
}

type PatternInfo struct {
	name        string
	category    string
	description string
	filePath    string
	packagePath string
}

func NewPatternRunner() *PatternRunner {
	pr := &PatternRunner{
		patterns: make(map[string]PatternInfo),
	}
	
	pr.registerPatterns()
	return pr
}

func (pr *PatternRunner) registerPatterns() {
	// Behavioral Patterns
	pr.patterns["chain_of_responsibility"] = PatternInfo{
		name:        "Chain of Responsibility",
		category:    "Behavioral",
		description: "Passes request along a chain of handlers",
		filePath:    "chain_of_responsibility.go",
		packagePath: "./behavioral/chain_of_responsibility",
	}
	
	pr.patterns["command"] = PatternInfo{
		name:        "Command",
		category:    "Behavioral",
		description: "Encapsulates a request as an object",
		filePath:    "command.go",
		packagePath: "./behavioral/command",
	}
	
	pr.patterns["mediator"] = PatternInfo{
		name:        "Mediator",
		category:    "Behavioral",
		description: "Defines simplified communication between classes",
		filePath:    "mediator.go",
		packagePath: "./behavioral/mediator",
	}
	
	pr.patterns["memento"] = PatternInfo{
		name:        "Memento",
		category:    "Behavioral",
		description: "Captures and restores an object's internal state",
		filePath:    "memento.go",
		packagePath: "./behavioral/memento",
	}
	
	pr.patterns["observer"] = PatternInfo{
		name:        "Observer",
		category:    "Behavioral",
		description: "Defines a subscription mechanism to notify multiple objects",
		filePath:    "observer.go",
		packagePath: "./behavioral/observer",
	}
	
	pr.patterns["state"] = PatternInfo{
		name:        "State",
		category:    "Behavioral",
		description: "Allows an object to alter its behavior when its state changes",
		filePath:    "state.go",
		packagePath: "./behavioral/state",
	}
	
	pr.patterns["strategy"] = PatternInfo{
		name:        "Strategy",
		category:    "Behavioral",
		description: "Defines a family of algorithms and makes them interchangeable",
		filePath:    "strategy.go",
		packagePath: "../strategy",
	}
	
	pr.patterns["template_method"] = PatternInfo{
		name:        "Template Method",
		category:    "Behavioral",
		description: "Defines the skeleton of an algorithm in a method",
		filePath:    "template_method.go",
		packagePath: "../template_method",
	}
	
	// Concurrency Patterns
	pr.patterns["barrier"] = PatternInfo{
		name:        "Barrier",
		category:    "Concurrency",
		description: "Synchronization primitive that blocks until a certain number of threads have reached the barrier",
		filePath:    "barrier.go",
		packagePath: "./concurrency/barrier",
	}
	
	pr.patterns["fan_in_fan_out"] = PatternInfo{
		name:        "Fan-in/Fan-out",
		category:    "Concurrency",
		description: "Concurrency pattern for distributing work to multiple goroutines and collecting results",
		filePath:    "fan_in_fan_out.go",
		packagePath: "./concurrency/fan_in_fan_out",
	}
	
	pr.patterns["generator"] = PatternInfo{
		name:        "Generator",
		category:    "Concurrency",
		description: "Function that returns a channel that produces a sequence of values",
		filePath:    "generator.go",
		packagePath: "./concurrency/generator",
	}
	
	pr.patterns["pipeline"] = PatternInfo{
		name:        "Pipeline",
		category:    "Concurrency",
		description: "Chain of stages connected by channels where output of one is input to next",
		filePath:    "pipeline.go",
		packagePath: "./concurrency/pipeline",
	}
	
	pr.patterns["semaphore"] = PatternInfo{
		name:        "Semaphore",
		category:    "Concurrency",
		description: "Synchronization primitive that controls access to a common resource",
		filePath:    "semaphore.go",
		packagePath: "./concurrency/semaphore",
	}
	
	pr.patterns["worker_pool"] = PatternInfo{
		name:        "Worker Pool",
		category:    "Concurrency",
		description: "Collection of goroutines waiting for tasks to be assigned",
		filePath:    "worker_pool.go",
		packagePath: "./concurrency/worker_pool",
	}
	
	// Creational Patterns
	pr.patterns["abstract_factory"] = PatternInfo{
		name:        "Abstract Factory",
		category:    "Creational",
		description: "Creates families of related objects without specifying their concrete classes",
		filePath:    "abstract_factory.go",
		packagePath: "./behavioral/abstract_factory",
	}
	
	pr.patterns["builder"] = PatternInfo{
		name:        "Builder",
		category:    "Creational",
		description: "Constructs complex objects step by step",
		filePath:    "builder.go",
		packagePath: "./creational/builder",
	}
	
	pr.patterns["factory_method"] = PatternInfo{
		name:        "Factory Method",
		category:    "Creational",
		description: "Defines an interface for creating objects but lets subclasses decide which class to instantiate",
		filePath:    "factory_method.go",
		packagePath: "./creational/factory_method",
	}
	
	pr.patterns["prototype"] = PatternInfo{
		name:        "Prototype",
		category:    "Creational",
		description: "Creates new objects by copying existing objects",
		filePath:    "prototype.go",
		packagePath: "./creational/prototype",
	}
	
	pr.patterns["singleton"] = PatternInfo{
		name:        "Singleton",
		category:    "Creational",
		description: "Ensures a class has only one instance and provides global access to it",
		filePath:    "singleton.go",
		packagePath: "./creational/singleton",
	}
	
	// Structural Patterns
	pr.patterns["adapter"] = PatternInfo{
		name:        "Adapter",
		category:    "Structural",
		description: "Allows incompatible interfaces to work together",
		filePath:    "adapter.go",
		packagePath: "./structural/adapter",
	}
	
	pr.patterns["bridge"] = PatternInfo{
		name:        "Bridge",
		category:    "Structural",
		description: "Decouples an abstraction from its implementation so they can vary independently",
		filePath:    "bridge.go",
		packagePath: "./structural/bridge",
	}
	
	pr.patterns["composite"] = PatternInfo{
		name:        "Composite",
		category:    "Structural",
		description: "Composes objects into tree structures to represent part-whole hierarchies",
		filePath:    "composite.go",
		packagePath: "./structural/composite",
	}
	
	pr.patterns["decorator"] = PatternInfo{
		name:        "Decorator",
		category:    "Structural",
		description: "Adds new functionality to objects dynamically without altering their structure",
		filePath:    "decorator.go",
		packagePath: "./structural/decorator",
	}
	
	pr.patterns["facade"] = PatternInfo{
		name:        "Facade",
		category:    "Structural",
		description: "Provides a simplified interface to a complex subsystem",
		filePath:    "facade.go",
		packagePath: "./structural/facade",
	}
	
	pr.patterns["flyweight"] = PatternInfo{
		name:        "Flyweight",
		category:    "Structural",
		description: "Reduces memory usage by sharing as much data as possible with similar objects",
		filePath:    "flyweight.go",
		packagePath: "./structural/flyweight",
	}
	
	pr.patterns["proxy"] = PatternInfo{
		name:        "Proxy",
		category:    "Structural",
		description: "Provides a surrogate or placeholder for another object to control access to it",
		filePath:    "proxy.go",
		packagePath: "./structural/proxy",
	}
	
	// Microservices Patterns
	pr.patterns["api_gateway"] = PatternInfo{
		name:        "API Gateway",
		category:    "Microservices",
		description: "Single entry point for all requests, routing them to appropriate services",
		filePath:    "api_gateway.go",
		packagePath: "./microservices/api_gateway",
	}
	
	pr.patterns["bulkhead"] = PatternInfo{
		name:        "Bulkhead",
		category:    "Microservices",
		description: "Isolates different parts of the system to prevent cascading failures",
		filePath:    "bulkhead.go",
		packagePath: "./microservices/bulkhead",
	}
	
	pr.patterns["circuit_breaker"] = PatternInfo{
		name:        "Circuit Breaker",
		category:    "Microservices",
		description: "Detects failures and encapsulates logic to prevent them from recurring",
		filePath:    "circuit_breaker.go",
		packagePath: "./microservices/circuit_breaker",
	}
	
	pr.patterns["rate_limiting"] = PatternInfo{
		name:        "Rate Limiting",
		category:    "Microservices",
		description: "Controls the rate of requests to prevent service overload",
		filePath:    "rate_limiting.go",
		packagePath: "./microservices/rate_limiting",
	}
	
	pr.patterns["saga"] = PatternInfo{
		name:        "Saga",
		category:    "Microservices",
		description: "Manages distributed transactions using a sequence of local transactions",
		filePath:    "saga.go",
		packagePath: "./microservices/saga",
	}
	
	pr.patterns["sidecar"] = PatternInfo{
		name:        "Sidecar",
		category:    "Microservices",
		description: "Deploys helper services alongside the main application to enhance functionality",
		filePath:    "sidecar.go",
		packagePath: "./microservices/sidecar",
	}
}

func (pr *PatternRunner) RunPattern(key string) error {
	pattern, exists := pr.patterns[key]
	if !exists {
		return fmt.Errorf("pattern '%s' not found", key)
	}
	
	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("Running %s Pattern (%s)\n", pattern.name, pattern.category)
	fmt.Printf("Description: %s\n", pattern.description)
	fmt.Printf(strings.Repeat("=", 60) + "\n")
	
	// Run the pattern
	cmd := exec.Command("go", "run", pattern.filePath)
	cmd.Dir = pattern.packagePath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run pattern %s: %v\nOutput: %s", key, err, string(output))
	}
	
	fmt.Print(string(output))
	return nil
}

func (pr *PatternRunner) ListPatterns() {
	fmt.Println("\nAvailable Design Patterns:")
	fmt.Println(strings.Repeat("-", 80))
	
	categories := make(map[string][]PatternInfo)
	for _, pattern := range pr.patterns {
		categories[pattern.category] = append(categories[pattern.category], pattern)
	}
	
	order := []string{"Behavioral", "Concurrency", "Creational", "Structural", "Microservices"}
	
	for _, category := range order {
		patterns := categories[category]
		if len(patterns) == 0 {
			continue
		}
		
		fmt.Printf("\n%s Patterns (%d):\n", category, len(patterns))
		fmt.Println(strings.Repeat("-", 40))
		
		for _, pattern := range patterns {
			key := ""
			for k, p := range pr.patterns {
				if p.name == pattern.name {
					key = k
					break
				}
			}
			fmt.Printf("  %-20s - %s\n", key, pattern.description)
		}
	}
	
	fmt.Printf("\nTotal: %d patterns\n", len(pr.patterns))
}

func (pr *PatternRunner) RunCategory(category string) error {
	var patternsToRun []string
	
	for key, pattern := range pr.patterns {
		if pattern.category == category {
			patternsToRun = append(patternsToRun, key)
		}
	}
	
	if len(patternsToRun) == 0 {
		return fmt.Errorf("no patterns found for category '%s'", category)
	}
	
	fmt.Printf("\nRunning all %s patterns...\n", category)
	
	for _, key := range patternsToRun {
		if err := pr.RunPattern(key); err != nil {
			fmt.Printf("Error running pattern %s: %v\n", key, err)
		}
		time.Sleep(1 * time.Second) // Brief pause between patterns
	}
	
	return nil
}

func (pr *PatternRunner) RunAll() error {
	fmt.Println("\nRunning ALL design patterns...")
	fmt.Println("This will take some time. Press Ctrl+C to stop.")
	
	// Run in order: Behavioral, Concurrency, Creational, Structural, Microservices
	order := []string{"Behavioral", "Concurrency", "Creational", "Structural", "Microservices"}
	
	for _, category := range order {
		if err := pr.RunCategory(category); err != nil {
			return fmt.Errorf("error running category %s: %v", category, err)
		}
		
		fmt.Printf("\nCompleted %s patterns. Press Enter to continue to next category...", category)
		fmt.Scanln() // Wait for user input
	}
	
	fmt.Println("\nAll patterns completed successfully!")
	return nil
}

func (pr *PatternRunner) GetStats() PatternStats {
	categoryCount := make(map[string]int)
	
	for _, pattern := range pr.patterns {
		categoryCount[pattern.category]++
	}
	
	return PatternStats{
		TotalPatterns: len(pr.patterns),
		Categories:    categoryCount,
	}
}

type PatternStats struct {
	TotalPatterns int
	Categories    map[string]int
}

func (ps PatternStats) Display() {
	fmt.Printf("\nDesign Patterns Statistics:\n")
	fmt.Printf("Total Patterns: %d\n", ps.TotalPatterns)
	fmt.Printf("Categories: %d\n", len(ps.Categories))
	
	fmt.Println("\nPatterns by Category:")
	for category, count := range ps.Categories {
		fmt.Printf("  %s: %d\n", category, count)
	}
}

func showBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════╗
║                    GO DESIGN PATTERNS                           ║
║                                                              ║
║    A comprehensive collection of design patterns in Go         ║
║                                                              ║
║    Categories:                                                ║
║    • Behavioral (8 patterns)                                 ║
║    • Concurrency (6 patterns)                                ║
║    • Creational (5 patterns)                                 ║
║    • Structural (7 patterns)                                 ║
║    • Microservices (6 patterns)                              ║
║                                                              ║
║    Total: 32 design patterns                                 ║
╚══════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

func showMenu() {
	fmt.Println("\nMain Menu:")
	fmt.Println("1. List all patterns")
	fmt.Println("2. Run a specific pattern")
	fmt.Println("3. Run all patterns in a category")
	fmt.Println("4. Run all patterns")
	fmt.Println("5. Show statistics")
	fmt.Println("6. Exit")
}

func getUserChoice() string {
	fmt.Print("\nEnter your choice (1-6): ")
	var choice string
	fmt.Scanln(&choice)
	return choice
}

func getPatternChoice(runner *PatternRunner) string {
	fmt.Println("\nAvailable pattern keys:")
	for key := range runner.patterns {
		fmt.Printf("  %s\n", key)
	}
	
	fmt.Print("\nEnter pattern key: ")
	var choice string
	fmt.Scanln(&choice)
	return choice
}

func getCategoryChoice() string {
	fmt.Println("\nAvailable categories:")
	fmt.Println("  behavioral")
	fmt.Println("  concurrency")
	fmt.Println("  creational")
	fmt.Println("  structural")
	fmt.Println("  microservices")
	
	fmt.Print("\nEnter category: ")
	var choice string
	fmt.Scanln(&choice)
	return choice
}

func main() {
	// Uncomment to run tests
	// testPatternRunner()
	// return
	
	showBanner()
	
	runner := NewPatternRunner()
	
	for {
		showMenu()
		choice := getUserChoice()
		
		switch choice {
		case "1":
			runner.ListPatterns()
			
		case "2":
			patternKey := getPatternChoice(runner)
			if err := runner.RunPattern(patternKey); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			
		case "3":
			category := getCategoryChoice()
			if err := runner.RunCategory(category); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			
		case "4":
			fmt.Println("\nStarting to run all patterns...")
			if err := runner.RunAll(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			
		case "5":
			stats := runner.GetStats()
			stats.Display()
			
		case "6":
			fmt.Println("\nThank you for using Go Design Patterns!")
			return
			
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 6.")
		}
		
		fmt.Println("\nPress Enter to continue...")
		fmt.Scanln()
	}
}
