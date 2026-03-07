package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Pipeline Stage interface
type Stage interface {
	Process(input interface{}) interface{}
}

// Concrete Stage 1: Data Generator
type DataGenerator struct {
	count int
}

func NewDataGenerator(count int) *DataGenerator {
	return &DataGenerator{count: count}
}

func (dg *DataGenerator) Process(input interface{}) interface{} {
	// Generate a sequence of numbers
	if num, ok := input.(int); ok {
		if num < dg.count {
			return num + 1
		}
	}
	return nil
}

// Concrete Stage 2: Data Transformer
type DataTransformer struct {
	transform func(interface{}) interface{}
}

func NewDataTransformer(transform func(interface{}) interface{}) *DataTransformer {
	return &DataTransformer{transform: transform}
}

func (dt *DataTransformer) Process(input interface{}) interface{} {
	return dt.transform(input)
}

// Concrete Stage 3: Data Filter
type DataFilter struct {
	predicate func(interface{}) bool
}

func NewDataFilter(predicate func(interface{}) bool) *DataFilter {
	return &DataFilter{predicate: predicate}
}

func (df *DataFilter) Process(input interface{}) interface{} {
	if df.predicate(input) {
		return input
	}
	return nil
}

// Concrete Stage 4: Data Aggregator
type DataAggregator struct {
	aggregate func([]interface{}) interface{}
	buffer    []interface{}
	size      int
}

func NewDataAggregator(size int, aggregate func([]interface{}) interface{}) *DataAggregator {
	return &DataAggregator{
		buffer:    make([]interface{}, 0, size),
		size:      size,
		aggregate: aggregate,
	}
}

func (da *DataAggregator) Process(input interface{}) interface{} {
	if input == nil {
		return nil
	}
	
	da.buffer = append(da.buffer, input)
	
	if len(da.buffer) >= da.size {
		result := da.aggregate(da.buffer)
		da.buffer = da.buffer[:0] // Reset buffer
		return result
	}
	
	return nil
}

// Pipeline implementation
type Pipeline struct {
	stages []Stage
}

func NewPipeline() *Pipeline {
	return &Pipeline{stages: make([]Stage, 0)}
}

func (p *Pipeline) AddStage(stage Stage) {
	p.stages = append(p.stages, stage)
}

func (p *Pipeline) Execute(input interface{}) interface{} {
	result := input
	
	for _, stage := range p.stages {
		result = stage.Process(result)
		if result == nil {
			break
		}
	}
	
	return result
}

// Channel-based Pipeline
func ChannelPipeline() {
	fmt.Println("--- Channel-based Pipeline Example ---")
	
	// Stage 1: Generate numbers
	generator := func() <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for i := 1; i <= 20; i++ {
				out <- i
				time.Sleep(100 * time.Millisecond)
			}
		}()
		return out
	}
	
	// Stage 2: Square numbers
	squarer := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for num := range in {
				squared := num * num
				fmt.Printf("Squaring %d -> %d\n", num, squared)
				out <- squared
			}
		}()
		return out
	}
	
	// Stage 3: Filter odd squares
	filter := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for num := range in {
				if num%2 == 0 {
					fmt.Printf("Filtering in %d (even)\n", num)
					out <- num
				} else {
					fmt.Printf("Filtering out %d (odd)\n", num)
				}
			}
		}()
		return out
	}
	
	// Stage 4: Convert to string
	toString := func(in <-chan int) <-chan string {
		out := make(chan string)
		go func() {
			defer close(out)
			for num := range in {
				str := fmt.Sprintf("Result: %d", num)
				fmt.Printf("Converting %d to string\n", num)
				out <- str
			}
		}()
		return out
	}
	
	// Build and run pipeline
	numbers := generator()
	squared := squarer(numbers)
	filtered := filter(squared)
	strings := toString(filtered)
	
	// Consume results
	for result := range strings {
		fmt.Printf("Final result: %s\n", result)
	}
}

// Text Processing Pipeline
func TextProcessingPipeline() {
	fmt.Println("\n--- Text Processing Pipeline Example ---")
	
	text := "The quick brown fox jumps over the lazy dog"
	
	// Stage 1: Split into words
	splitStage := NewDataTransformer(func(input interface{}) interface{} {
		if text, ok := input.(string); ok {
			return strings.Fields(text)
		}
		return nil
	})
	
	// Stage 2: Filter short words
	filterStage := NewDataFilter(func(input interface{}) bool {
		if words, ok := input.([]string); ok {
			var longWords []string
			for _, word := range words {
				if len(word) > 3 {
					longWords = append(longWords, word)
				}
			}
			return len(longWords) > 0
		}
		return false
	})
	
	// Stage 3: Transform to uppercase
	upperStage := NewDataTransformer(func(input interface{}) interface{} {
		if words, ok := input.([]string); ok {
			var upperWords []string
			for _, word := range words {
				if len(word) > 3 {
					upperWords = append(upperWords, strings.ToUpper(word))
				}
			}
			return upperWords
		}
		return nil
	})
	
	// Stage 4: Join with spaces
	joinStage := NewDataTransformer(func(input interface{}) interface{} {
		if words, ok := input.([]string); ok {
			return strings.Join(words, " ")
		}
		return nil
	})
	
	// Build pipeline
	pipeline := NewPipeline()
	pipeline.AddStage(splitStage)
	pipeline.AddStage(filterStage)
	pipeline.AddStage(upperStage)
	pipeline.AddStage(joinStage)
	
	// Execute pipeline
	result := pipeline.Execute(text)
	if result != nil {
		fmt.Printf("Original: %s\n", text)
		fmt.Printf("Processed: %s\n", result)
	}
}

// Data Analysis Pipeline
func DataAnalysisPipeline() {
	fmt.Println("\n--- Data Analysis Pipeline Example ---")
	
	// Generate random data points
	data := make([]float64, 100)
	for i := range data {
		data[i] = rand.Float64() * 100
	}
	
	fmt.Printf("Generated %d data points\n", len(data))
	
	// Stage 1: Filter outliers (> 95th percentile)
	filterStage := NewDataTransformer(func(input interface{}) interface{} {
		if points, ok := input.([]float64); ok {
			var filtered []float64
			
			// Calculate 95th percentile
			sorted := make([]float64, len(points))
			copy(sorted, points)
			
			// Simple sort (bubble sort for demonstration)
			for i := 0; i < len(sorted)-1; i++ {
				for j := 0; j < len(sorted)-i-1; j++ {
					if sorted[j] > sorted[j+1] {
						sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
					}
				}
			}
			
			percentile95 := sorted[int(float64(len(sorted))*0.95)]
			
			for _, point := range points {
				if point <= percentile95 {
					filtered = append(filtered, point)
				}
			}
			
			fmt.Printf("Filtered %d outliers (removed %d points)\n", 
				len(filtered), len(points)-len(filtered))
			return filtered
		}
		return nil
	})
	
	// Stage 2: Calculate moving average
	avgStage := NewDataTransformer(func(input interface{}) interface{} {
		if points, ok := input.([]float64); ok {
			if len(points) < 5 {
				return 0.0
			}
			
			// Calculate average of last 5 points
			sum := 0.0
			for i := len(points) - 5; i < len(points); i++ {
				sum += points[i]
			}
			return sum / 5.0
		}
		return nil
	})
	
	// Stage 3: Format result
	formatStage := NewDataTransformer(func(input interface{}) interface{} {
		if avg, ok := input.(float64); ok {
			return fmt.Sprintf("Moving average: %.2f", avg)
		}
		return nil
	})
	
	// Build pipeline
	pipeline := NewPipeline()
	pipeline.AddStage(filterStage)
	pipeline.AddStage(avgStage)
	pipeline.AddStage(formatStage)
	
	// Execute pipeline
	result := pipeline.Execute(data)
	if result != nil {
		fmt.Printf("Analysis result: %s\n", result)
	}
}

// Pipeline with error handling
type PipelineWithErrors struct {
	stages []Stage
}

func NewPipelineWithErrors() *PipelineWithErrors {
	return &PipelineWithErrors{stages: make([]Stage, 0)}
}

func (p *PipelineWithErrors) AddStage(stage Stage) {
	p.stages = append(p.stages, stage)
}

func (p *PipelineWithErrors) Execute(input interface{}) (interface{}, error) {
	result := input
	
	for i, stage := range p.stages {
		result = stage.Process(result)
		if result == nil {
			return nil, fmt.Errorf("stage %d returned nil result", i)
		}
	}
	
	return result, nil
}

func ErrorHandlingPipeline() {
	fmt.Println("\n--- Error Handling Pipeline Example ---")
	
	// Stage that might fail
	riskyStage := NewDataTransformer(func(input interface{}) interface{} {
		if num, ok := input.(int); ok {
			if num == 0 {
				return nil // Simulate failure
			}
			return 100 / num
		}
		return nil
	})
	
	pipeline := NewPipelineWithErrors()
	pipeline.AddStage(riskyStage)
	
	testValues := []int{10, 5, 2, 0, 1}
	
	for _, value := range testValues {
		result, err := pipeline.Execute(value)
		if err != nil {
			fmt.Printf("Input %d: Error - %v\n", value, err)
		} else {
			fmt.Printf("Input %d: Success - %v\n", value, result)
		}
	}
}

func main() {
	fmt.Println("=== Pipeline Pattern Demo ===")
	
	rand.Seed(time.Now().UnixNano())
	
	ChannelPipeline()
	TextProcessingPipeline()
	DataAnalysisPipeline()
	ErrorHandlingPipeline()
	
	fmt.Println("\nAll pipeline patterns demonstrated successfully!")
}
