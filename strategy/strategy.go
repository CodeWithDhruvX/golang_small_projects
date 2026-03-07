package main

import (
	"fmt"
)

// Strategy interface defines the algorithm contract
type PaymentStrategy interface {
	Pay(amount float64)
}

// ConcreteStrategy1 - Credit Card payment
type CreditCardPayment struct {
	cardNumber string
	name       string
}

func NewCreditCardPayment(cardNumber, name string) *CreditCardPayment {
	return &CreditCardPayment{
		cardNumber: cardNumber,
		name:       name,
	}
}

func (ccp *CreditCardPayment) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using Credit Card ending in %s (Name: %s)\n", 
		amount, ccp.cardNumber[len(ccp.cardNumber)-4:], ccp.name)
}

// ConcreteStrategy2 - PayPal payment
type PayPalPayment struct {
	email string
}

func NewPayPalPayment(email string) *PayPalPayment {
	return &PayPalPayment{email: email}
}

func (pp *PayPalPayment) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using PayPal (Email: %s)\n", amount, pp.email)
}

// ConcreteStrategy3 - Bitcoin payment
type BitcoinPayment struct {
	walletAddress string
}

func NewBitcoinPayment(walletAddress string) *BitcoinPayment {
	return &BitcoinPayment{walletAddress: walletAddress}
}

func (bp *BitcoinPayment) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using Bitcoin (Wallet: %s...%s)\n", 
		amount, bp.walletAddress[:6], bp.walletAddress[len(bp.walletAddress)-4:])
}

// Context that uses the strategy
type ShoppingCart struct {
	amount   float64
	strategy PaymentStrategy
}

func NewShoppingCart(amount float64) *ShoppingCart {
	return &ShoppingCart{amount: amount}
}

func (sc *ShoppingCart) SetPaymentStrategy(strategy PaymentStrategy) {
	sc.strategy = strategy
}

func (sc *ShoppingCart) Checkout() {
	if sc.strategy == nil {
		fmt.Println("Please select a payment method")
		return
	}
	
	fmt.Printf("Checking out cart with total: $%.2f\n", sc.amount)
	sc.strategy.Pay(sc.amount)
}

// Sorting Strategy example
type SortingStrategy interface {
	Sort(data []int) []int
}

type BubbleSort struct{}

func (bs *BubbleSort) Sort(data []int) []int {
	fmt.Println("Using Bubble Sort")
	n := len(data)
	result := make([]int, n)
	copy(result, data)
	
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}

type QuickSort struct{}

func (qs *QuickSort) Sort(data []int) []int {
	fmt.Println("Using Quick Sort")
	result := make([]int, len(data))
	copy(result, data)
	return qs.quickSort(result, 0, len(result)-1)
}

func (qs *QuickSort) quickSort(arr []int, low, high int) []int {
	if low < high {
		pi := qs.partition(arr, low, high)
		qs.quickSort(arr, low, pi-1)
		qs.quickSort(arr, pi+1, high)
	}
	return arr
}

func (qs *QuickSort) partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1
	
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

type MergeSort struct{}

func (ms *MergeSort) Sort(data []int) []int {
	fmt.Println("Using Merge Sort")
	result := make([]int, len(data))
	copy(result, data)
	return ms.mergeSort(result)
}

func (ms *MergeSort) mergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	
	mid := len(arr) / 2
	left := ms.mergeSort(arr[:mid])
	right := ms.mergeSort(arr[mid:])
	
	return ms.merge(left, right)
}

func (ms *MergeSort) merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0
	
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	
	return result
}

type Sorter struct {
	strategy SortingStrategy
	data     []int
}

func NewSorter(data []int) *Sorter {
	return &Sorter{data: data}
}

func (s *Sorter) SetSortingStrategy(strategy SortingStrategy) {
	s.strategy = strategy
}

func (s *Sorter) PerformSort() []int {
	if s.strategy == nil {
		fmt.Println("Please select a sorting strategy")
		return s.data
	}
	
	fmt.Printf("Original data: %v\n", s.data)
	return s.strategy.Sort(s.data)
}

// Compression Strategy example
type CompressionStrategy interface {
	Compress(data string) string
}

type ZipCompression struct{}

func (zc *ZipCompression) Compress(data string) string {
	fmt.Println("Compressing using ZIP algorithm")
	return fmt.Sprintf("ZIP:[%s]", data)
}

type GzipCompression struct{}

func (gc *GzipCompression) Compress(data string) string {
	fmt.Println("Compressing using GZIP algorithm")
	return fmt.Sprintf("GZIP:[%s]", data)
}

type RarCompression struct{}

func (rc *RarCompression) Compress(data string) string {
	fmt.Println("Compressing using RAR algorithm")
	return fmt.Sprintf("RAR:[%s]", data)
}

type Compressor struct {
	strategy CompressionStrategy
	data     string
}

func NewCompressor(data string) *Compressor {
	return &Compressor{data: data}
}

func (c *Compressor) SetCompressionStrategy(strategy CompressionStrategy) {
	c.strategy = strategy
}

func (c *Compressor) Compress() string {
	if c.strategy == nil {
		fmt.Println("Please select a compression strategy")
		return c.data
	}
	
	return c.strategy.Compress(c.data)
}

func main() {
	fmt.Println("=== Strategy Pattern Demo ===")
	
	// Payment Strategy example
	fmt.Println("\n--- Payment Strategy Example ---")
	cart := NewShoppingCart(150.75)
	
	fmt.Println("Selecting Credit Card payment:")
	cart.SetPaymentStrategy(NewCreditCardPayment("1234567890123456", "John Doe"))
	cart.Checkout()
	
	fmt.Println("\nSelecting PayPal payment:")
	cart.SetPaymentStrategy(NewPayPalPayment("john.doe@example.com"))
	cart.Checkout()
	
	fmt.Println("\nSelecting Bitcoin payment:")
	cart.SetPaymentStrategy(NewBitcoinPayment("1A2B3C4D5E6F7G8H9I0J"))
	cart.Checkout()
	
	// Sorting Strategy example
	fmt.Println("\n--- Sorting Strategy Example ---")
	data := []int{64, 34, 25, 12, 22, 11, 90}
	sorter := NewSorter(data)
	
	sorter.SetSortingStrategy(&BubbleSort{})
	sorted := sorter.PerformSort()
	fmt.Printf("Sorted result: %v\n\n", sorted)
	
	sorter.SetSortingStrategy(&QuickSort{})
	sorted = sorter.PerformSort()
	fmt.Printf("Sorted result: %v\n\n", sorted)
	
	sorter.SetSortingStrategy(&MergeSort{})
	sorted = sorter.PerformSort()
	fmt.Printf("Sorted result: %v\n\n", sorted)
	
	// Compression Strategy example
	fmt.Println("--- Compression Strategy Example ---")
	compressor := NewCompressor("This is some important data to compress")
	
	compressor.SetCompressionStrategy(&ZipCompression{})
	compressed := compressor.Compress()
	fmt.Printf("Compressed result: %s\n\n", compressed)
	
	compressor.SetCompressionStrategy(&GzipCompression{})
	compressed = compressor.Compress()
	fmt.Printf("Compressed result: %s\n\n", compressed)
	
	compressor.SetCompressionStrategy(&RarCompression{})
	compressed = compressor.Compress()
	fmt.Printf("Compressed result: %s\n", compressed)
}
