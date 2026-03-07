package main

import (
	"fmt"
	"time"
)

// Proxy Pattern

// Subject interface
type Subject interface {
	Request() string
}

// Real Subject
type RealSubject struct{}

func (rs *RealSubject) Request() string {
	return "RealSubject: Handling request"
}

// Proxy
type Proxy struct {
	realSubject *RealSubject
	initialized bool
}

func NewProxy() *Proxy {
	return &Proxy{}
}

func (p *Proxy) Request() string {
	if !p.initialized {
		fmt.Println("Proxy: Initializing real subject")
		p.realSubject = &RealSubject{}
		p.initialized = true
	}
	
	fmt.Println("Proxy: Pre-processing request")
	result := p.realSubject.Request()
	fmt.Println("Proxy: Post-processing request")
	return result
}

// Image Proxy Example
type Image interface {
	Display()
}

type RealImage struct {
	filename string
}

func NewRealImage(filename string) *RealImage {
	fmt.Printf("RealImage: Loading image from disk: %s\n", filename)
	return &RealImage{filename: filename}
}

func (ri *RealImage) Display() {
	fmt.Printf("RealImage: Displaying %s\n", ri.filename)
}

type ImageProxy struct {
	filename    string
	realImage   *RealImage
	loaded      bool
}

func NewImageProxy(filename string) *ImageProxy {
	return &ImageProxy{filename: filename}
}

func (ip *ImageProxy) Display() {
	if !ip.loaded {
		fmt.Printf("ImageProxy: Loading image %s on demand\n", ip.filename)
		ip.realImage = NewRealImage(ip.filename)
		ip.loaded = true
	}
	fmt.Printf("ImageProxy: Displaying %s\n", ip.filename)
	ip.realImage.Display()
}

// Database Connection Proxy Example
type Database interface {
	Query(sql string) string
	Connect() error
	Disconnect() error
}

type RealDatabase struct {
	connected bool
}

func NewRealDatabase() *RealDatabase {
	fmt.Println("RealDatabase: Establishing connection to database")
	return &RealDatabase{connected: false}
}

func (rd *RealDatabase) Connect() error {
	if !rd.connected {
		fmt.Println("RealDatabase: Connecting to database...")
		time.Sleep(100 * time.Millisecond) // Simulate connection time
		rd.connected = true
		fmt.Println("RealDatabase: Connected successfully")
	}
	return nil
}

func (rd *RealDatabase) Disconnect() error {
	if rd.connected {
		fmt.Println("RealDatabase: Disconnecting from database...")
		rd.connected = false
		fmt.Println("RealDatabase: Disconnected")
	}
	return nil
}

func (rd *RealDatabase) Query(sql string) string {
	if !rd.connected {
		return "Error: Not connected to database"
	}
	fmt.Printf("RealDatabase: Executing query: %s\n", sql)
	return "Query results"
}

type DatabaseProxy struct {
	realDatabase *RealDatabase
	queryCache   map[string]string
	connected    bool
}

func NewDatabaseProxy() *DatabaseProxy {
	return &DatabaseProxy{
		queryCache: make(map[string]string),
	}
}

func (dp *DatabaseProxy) Connect() error {
	if !dp.connected {
		fmt.Println("DatabaseProxy: Connecting through proxy...")
		err := dp.getRealDatabase().Connect()
		if err == nil {
			dp.connected = true
		}
		return err
	}
	return nil
}

func (dp *DatabaseProxy) Disconnect() error {
	if dp.connected {
		fmt.Println("DatabaseProxy: Disconnecting through proxy...")
		dp.realDatabase.Disconnect()
		dp.connected = false
		dp.queryCache = make(map[string]string) // Clear cache
	}
	return nil
}

func (dp *DatabaseProxy) Query(sql string) string {
	// Check cache first
	if result, exists := dp.queryCache[sql]; exists {
		fmt.Printf("DatabaseProxy: Returning cached result for: %s\n", sql)
		return result
	}
	
	// Execute query through real database
	fmt.Printf("DatabaseProxy: Executing query through proxy: %s\n", sql)
	result := dp.getRealDatabase().Query(sql)
	
	// Cache the result
	dp.queryCache[sql] = result
	fmt.Printf("DatabaseProxy: Cached result for: %s\n", sql)
	
	return result
}

func (dp *DatabaseProxy) getRealDatabase() *RealDatabase {
	if dp.realDatabase == nil {
		dp.realDatabase = NewRealDatabase()
	}
	return dp.realDatabase
}

// Internet Access Proxy Example
type InternetAccess interface {
	AccessWebsite(url string) string
}

type RealInternetAccess struct {
	userName string
}

func NewRealInternetAccess(userName string) *RealInternetAccess {
	return &RealInternetAccess{userName: userName}
}

func (ria *RealInternetAccess) AccessWebsite(url string) string {
	fmt.Printf("RealInternetAccess: %s accessing %s\n", ria.userName, url)
	return fmt.Sprintf("Content of %s", url)
}

type InternetAccessProxy struct {
	realAccess   *RealInternetAccess
	userName     string
	blockedSites []string
}

func NewInternetAccessProxy(userName string) *InternetAccessProxy {
	return &InternetAccessProxy{
		userName:     userName,
		blockedSites: []string{"facebook.com", "twitter.com", "instagram.com"},
	}
}

func (iap *InternetAccessProxy) AccessWebsite(url string) string {
	// Check if site is blocked
	for _, blocked := range iap.blockedSites {
		if url == blocked {
			return fmt.Sprintf("Access to %s is blocked for user %s", url, iap.userName)
		}
	}
	
	// Check user permissions (simplified)
	if iap.userName == "admin" {
		fmt.Printf("InternetAccessProxy: Admin access granted to %s\n", url)
		return iap.getRealAccess().AccessWebsite(url)
	}
	
	// Regular user access
	fmt.Printf("InternetAccessProxy: User %s accessing %s\n", iap.userName, url)
	return iap.getRealAccess().AccessWebsite(url)
}

func (iap *InternetAccessProxy) getRealAccess() *RealInternetAccess {
	if iap.realAccess == nil {
		iap.realAccess = NewRealInternetAccess(iap.userName)
	}
	return iap.realAccess
}

// File System Proxy Example
type FileSystem interface {
	ReadFile(path string) string
	WriteFile(path, content string) error
	DeleteFile(path string) error
}

type RealFileSystem struct{}

func (rfs *RealFileSystem) ReadFile(path string) string {
	fmt.Printf("RealFileSystem: Reading file %s from disk\n", path)
	return fmt.Sprintf("Content of %s", path)
}

func (rfs *RealFileSystem) WriteFile(path, content string) error {
	fmt.Printf("RealFileSystem: Writing to file %s\n", path)
	return nil
}

func (rfs *RealFileSystem) DeleteFile(path string) error {
	fmt.Printf("RealFileSystem: Deleting file %s\n", path)
	return nil
}

type FileSystemProxy struct {
	realFileSystem *RealFileSystem
	accessControl  map[string][]string // user -> permissions
	currentUser    string
}

func NewFileSystemProxy(currentUser string) *FileSystemProxy {
	return &FileSystemProxy{
		accessControl: map[string][]string{
			"admin": {"read", "write", "delete"},
			"user":  {"read", "write"},
			"guest": {"read"},
		},
		currentUser: currentUser,
	}
}

func (fsp *FileSystemProxy) ReadFile(path string) string {
	if !fsp.hasPermission("read") {
		return fmt.Sprintf("Access denied: %s cannot read files", fsp.currentUser)
	}
	
	fmt.Printf("FileSystemProxy: %s reading file %s\n", fsp.currentUser, path)
	return fsp.getRealFileSystem().ReadFile(path)
}

func (fsp *FileSystemProxy) WriteFile(path, content string) error {
	if !fsp.hasPermission("write") {
		return fmt.Errorf("access denied: %s cannot write files", fsp.currentUser)
	}
	
	fmt.Printf("FileSystemProxy: %s writing to file %s\n", fsp.currentUser, path)
	return fsp.getRealFileSystem().WriteFile(path, content)
}

func (fsp *FileSystemProxy) DeleteFile(path string) error {
	if !fsp.hasPermission("delete") {
		return fmt.Errorf("access denied: %s cannot delete files", fsp.currentUser)
	}
	
	fmt.Printf("FileSystemProxy: %s deleting file %s\n", fsp.currentUser, path)
	return fsp.getRealFileSystem().DeleteFile(path)
}

func (fsp *FileSystemProxy) hasPermission(permission string) bool {
	permissions, exists := fsp.accessControl[fsp.currentUser]
	if !exists {
		return false
	}
	
	for _, perm := range permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

func (fsp *FileSystemProxy) getRealFileSystem() *RealFileSystem {
	if fsp.realFileSystem == nil {
		fsp.realFileSystem = &RealFileSystem{}
	}
	return fsp.realFileSystem
}

// Payment Service Proxy Example
type PaymentService interface {
	ProcessPayment(amount float64, cardNumber string) string
	RefundPayment(transactionID string, amount float64) string
}

type RealPaymentService struct{}

func (rps *RealPaymentService) ProcessPayment(amount float64, cardNumber string) string {
	fmt.Printf("RealPaymentService: Processing payment of $%.2f for card %s\n", amount, maskCard(cardNumber))
	time.Sleep(200 * time.Millisecond) // Simulate network latency
	return fmt.Sprintf("TXN_%d", time.Now().Unix())
}

func (rps *RealPaymentService) RefundPayment(transactionID string, amount float64) string {
	fmt.Printf("RealPaymentService: Refunding $%.2f for transaction %s\n", amount, transactionID)
	time.Sleep(150 * time.Millisecond) // Simulate network latency
	return fmt.Sprintf("REFUND_%d", time.Now().Unix())
}

type PaymentServiceProxy struct {
	realService    *RealPaymentService
	transactionLog []Transaction
	rateLimiter    *RateLimiter
}

type Transaction struct {
	ID        string
	Amount    float64
	CardNumber string
	Timestamp time.Time
	Status    string
}

type RateLimiter struct {
	requests    []time.Time
	maxRequests int
	timeWindow  time.Duration
}

func NewRateLimiter(maxRequests int, timeWindow time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		timeWindow:  timeWindow,
	}
}

func (rl *RateLimiter) AllowRequest() bool {
	now := time.Now()
	
	// Remove old requests outside the time window
	validRequests := make([]time.Time, 0)
	for _, req := range rl.requests {
		if now.Sub(req) <= rl.timeWindow {
			validRequests = append(validRequests, req)
		}
	}
	
	rl.requests = validRequests
	
	// Check if we can allow this request
	if len(rl.requests) < rl.maxRequests {
		rl.requests = append(rl.requests, now)
		return true
	}
	
	return false
}

func NewPaymentServiceProxy() *PaymentServiceProxy {
	return &PaymentServiceProxy{
		transactionLog: make([]Transaction, 0),
		rateLimiter:    NewRateLimiter(5, time.Minute), // 5 requests per minute
	}
}

func (psp *PaymentServiceProxy) ProcessPayment(amount float64, cardNumber string) string {
	// Rate limiting
	if !psp.rateLimiter.AllowRequest() {
		return "Error: Rate limit exceeded"
	}
	
	// Input validation
	if amount <= 0 {
		return "Error: Invalid amount"
	}
	
	if !isValidCard(cardNumber) {
		return "Error: Invalid card number"
	}
	
	// Log transaction attempt
	transaction := Transaction{
		Amount:     amount,
		CardNumber: maskCard(cardNumber),
		Timestamp:  time.Now(),
		Status:     "PROCESSING",
	}
	
	fmt.Printf("PaymentServiceProxy: Processing payment through proxy\n")
	transactionID := psp.getRealService().ProcessPayment(amount, cardNumber)
	
	transaction.ID = transactionID
	transaction.Status = "COMPLETED"
	psp.transactionLog = append(psp.transactionLog, transaction)
	
	fmt.Printf("PaymentServiceProxy: Transaction %s completed\n", transactionID)
	return transactionID
}

func (psp *PaymentServiceProxy) RefundPayment(transactionID string, amount float64) string {
	// Find transaction
	var foundTransaction *Transaction
	for _, txn := range psp.transactionLog {
		if txn.ID == transactionID {
			foundTransaction = &txn
			break
		}
	}
	
	if foundTransaction == nil {
		return "Error: Transaction not found"
	}
	
	if foundTransaction.Status != "COMPLETED" {
		return "Error: Transaction already refunded or failed"
	}
	
	fmt.Printf("PaymentServiceProxy: Processing refund through proxy\n")
	refundID := psp.getRealService().RefundPayment(transactionID, amount)
	
	// Update transaction status
	for i, txn := range psp.transactionLog {
		if txn.ID == transactionID {
			psp.transactionLog[i].Status = "REFUNDED"
			break
		}
	}
	
	fmt.Printf("PaymentServiceProxy: Refund %s completed\n", refundID)
	return refundID
}

func (psp *PaymentServiceProxy) getRealService() *RealPaymentService {
	if psp.realService == nil {
		psp.realService = &RealPaymentService{}
	}
	return psp.realService
}

func (psp *PaymentServiceProxy) GetTransactionHistory() []Transaction {
	return psp.transactionLog
}

func maskCard(cardNumber string) string {
	if len(cardNumber) <= 4 {
		return cardNumber
	}
	return "****-****-****-" + cardNumber[len(cardNumber)-4:]
}

func isValidCard(cardNumber string) bool {
	// Simple validation - check if it's all digits and has reasonable length
	for _, char := range cardNumber {
		if char < '0' || char > '9' {
			return false
		}
	}
	return len(cardNumber) >= 13 && len(cardNumber) <= 19
}

func main() {
	fmt.Println("=== Proxy Pattern Demo ===")
	
	// Basic example
	fmt.Println("\n--- Basic Proxy Example ---")
	
	proxy := &Proxy{}
	fmt.Println(proxy.Request())
	fmt.Println(proxy.Request()) // Second call should use initialized real subject
	
	// Image Proxy example
	fmt.Println("\n--- Image Proxy Example ---")
	
	// Create image proxies (images are not loaded yet)
	image1 := NewImageProxy("photo1.jpg")
	image2 := NewImageProxy("photo2.jpg")
	image3 := NewImageProxy("photo1.jpg") // Same as image1
	
	fmt.Println("Displaying images (loading on demand):")
	image1.Display()
	image2.Display()
	image3.Display() // Should use already loaded image
	
	// Database Connection Proxy example
	fmt.Println("\n--- Database Connection Proxy Example ---")
	
	dbProxy := NewDatabaseProxy()
	
	// First query - should connect and cache result
	fmt.Println("First query:")
	result1 := dbProxy.Query("SELECT * FROM users")
	fmt.Printf("Result: %s\n", result1)
	
	// Second query - should connect and cache result
	fmt.Println("\nSecond query:")
	result2 := dbProxy.Query("SELECT * FROM products")
	fmt.Printf("Result: %s\n", result2)
	
	// Third query - should return cached result
	fmt.Println("\nThird query (cached):")
	result3 := dbProxy.Query("SELECT * FROM users")
	fmt.Printf("Result: %s\n", result3)
	
	dbProxy.Disconnect()
	
	// Internet Access Proxy example
	fmt.Println("\n--- Internet Access Proxy Example ---")
	
	adminAccess := NewInternetAccessProxy("admin")
	userAccess := NewInternetAccessProxy("john")
	
	fmt.Println("Admin accessing websites:")
	fmt.Println(adminAccess.AccessWebsite("google.com"))
	fmt.Println(adminAccess.AccessWebsite("facebook.com"))
	
	fmt.Println("\nUser accessing websites:")
	fmt.Println(userAccess.AccessWebsite("google.com"))
	fmt.Println(userAccess.AccessWebsite("facebook.com"))
	fmt.Println(userAccess.AccessWebsite("twitter.com"))
	
	// File System Proxy example
	fmt.Println("\n--- File System Proxy Example ---")
	
	adminFS := NewFileSystemProxy("admin")
	userFS := NewFileSystemProxy("user")
	guestFS := NewFileSystemProxy("guest")
	
	fmt.Println("Admin file operations:")
	fmt.Println(adminFS.ReadFile("document.txt"))
	adminFS.WriteFile("output.txt", "Hello World")
	adminFS.DeleteFile("temp.txt")
	
	fmt.Println("\nUser file operations:")
	fmt.Println(userFS.ReadFile("document.txt"))
	userFS.WriteFile("user_output.txt", "User content")
	fmt.Println(userFS.DeleteFile("temp.txt")) // Should fail
	
	fmt.Println("\nGuest file operations:")
	fmt.Println(guestFS.ReadFile("document.txt"))
	fmt.Println(guestFS.WriteFile("guest_output.txt", "Guest content")) // Should fail
	fmt.Println(guestFS.DeleteFile("temp.txt")) // Should fail
	
	// Payment Service Proxy example
	fmt.Println("\n--- Payment Service Proxy Example ---")
	
	paymentProxy := NewPaymentServiceProxy()
	
	// Process payments
	fmt.Println("Processing payments:")
	txn1 := paymentProxy.ProcessPayment(100.50, "1234567890123456")
	txn2 := paymentProxy.ProcessPayment(75.25, "9876543210987654")
	txn3 := paymentProxy.ProcessPayment(-50.00, "1111222233334444") // Should fail
	txn4 := paymentProxy.ProcessPayment(200.00, "invalid_card") // Should fail
	
	fmt.Println("\nTransaction history:")
	for _, txn := range paymentProxy.GetTransactionHistory() {
		fmt.Printf("ID: %s, Amount: $%.2f, Card: %s, Status: %s\n", 
			txn.ID, txn.Amount, txn.CardNumber, txn.Status)
	}
	
	fmt.Println("\nProcessing refunds:")
	refund1 := paymentProxy.RefundPayment(txn1, 50.00)
	refund2 := paymentProxy.RefundPayment("INVALID_TXN", 25.00) // Should fail
	
	fmt.Println("\nUpdated transaction history:")
	for _, txn := range paymentProxy.GetTransactionHistory() {
		fmt.Printf("ID: %s, Amount: $%.2f, Card: %s, Status: %s\n", 
			txn.ID, txn.Amount, txn.CardNumber, txn.Status)
	}
	
	fmt.Println("\nAll proxy patterns demonstrated successfully!")
}
