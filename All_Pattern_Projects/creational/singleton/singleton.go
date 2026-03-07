package main

import (
	"fmt"
	"sync"
	"time"
)

// Singleton Pattern

// Basic Singleton (not thread-safe)
type BasicSingleton struct {
	data string
}

var basicInstance *BasicSingleton

func GetBasicSingleton() *BasicSingleton {
	if basicInstance == nil {
		basicInstance = &BasicSingleton{data: "Basic Singleton Data"}
	}
	return basicInstance
}

func (bs *BasicSingleton) GetData() string {
	return bs.data
}

func (bs *BasicSingleton) SetData(data string) {
	bs.data = data
}

// Thread-safe Singleton using sync.Once
type ThreadSafeSingleton struct {
	data string
}

var (
	threadSafeInstance *ThreadSafeSingleton
	threadSafeOnce     sync.Once
)

func GetThreadSafeSingleton() *ThreadSafeSingleton {
	threadSafeOnce.Do(func() {
		threadSafeInstance = &ThreadSafeSingleton{data: "Thread-Safe Singleton Data"}
	})
	return threadSafeInstance
}

func (tss *ThreadSafeSingleton) GetData() string {
	return tss.data
}

func (tss *ThreadSafeSingleton) SetData(data string) {
	tss.data = data
}

// Logger Singleton
type Logger struct {
	logs []string
	mu   sync.Mutex
}

var (
	loggerInstance *Logger
	loggerOnce     sync.Once
)

func GetLogger() *Logger {
	loggerOnce.Do(func() {
		loggerInstance = &Logger{logs: make([]string, 0)}
	})
	return loggerInstance
}

func (l *Logger) Log(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s", timestamp, message)
	l.logs = append(l.logs, logEntry)
	fmt.Printf("LOG: %s\n", message)
}

func (l *Logger) GetLogs() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	logsCopy := make([]string, len(l.logs))
	copy(logsCopy, l.logs)
	return logsCopy
}

func (l *Logger) ClearLogs() {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	l.logs = l.logs[:0]
}

// Database Connection Singleton
type DatabaseConnection struct {
	connectionString string
	connected        bool
	mu               sync.RWMutex
}

var (
	dbInstance *DatabaseConnection
	dbOnce     sync.Once
)

func GetDatabaseConnection() *DatabaseConnection {
	dbOnce.Do(func() {
		dbInstance = &DatabaseConnection{
			connectionString: "localhost:5432/mydb",
			connected:        false,
		}
	})
	return dbInstance
}

func (dc *DatabaseConnection) Connect() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	
	if dc.connected {
		return fmt.Errorf("already connected to database")
	}
	
	// Simulate connection
	time.Sleep(100 * time.Millisecond)
	dc.connected = true
	fmt.Printf("Connected to database: %s\n", dc.connectionString)
	return nil
}

func (dc *DatabaseConnection) Disconnect() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	
	if !dc.connected {
		return fmt.Errorf("not connected to database")
	}
	
	// Simulate disconnection
	time.Sleep(50 * time.Millisecond)
	dc.connected = false
	fmt.Println("Disconnected from database")
	return nil
}

func (dc *DatabaseConnection) IsConnected() bool {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return dc.connected
}

func (dc *DatabaseConnection) ExecuteQuery(query string) error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	
	if !dc.connected {
		return fmt.Errorf("not connected to database")
	}
	
	fmt.Printf("Executing query: %s\n", query)
	return nil
}

// Configuration Manager Singleton
type ConfigManager struct {
	config map[string]interface{}
	mu     sync.RWMutex
}

var (
	configInstance *ConfigManager
	configOnce     sync.Once
)

func GetConfigManager() *ConfigManager {
	configOnce.Do(func() {
		configInstance = &ConfigManager{
			config: make(map[string]interface{}),
		}
		// Initialize with default values
		configInstance.config["app_name"] = "MyApplication"
		configInstance.config["version"] = "1.0.0"
		configInstance.config["debug"] = true
		configInstance.config["max_connections"] = 100
	})
	return configInstance
}

func (cm *ConfigManager) Get(key string) (interface{}, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	value, exists := cm.config[key]
	return value, exists
}

func (cm *ConfigManager) Set(key string, value interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.config[key] = value
	fmt.Printf("Config updated: %s = %v\n", key, value)
}

func (cm *ConfigManager) GetAll() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	configCopy := make(map[string]interface{})
	for k, v := range cm.config {
		configCopy[k] = v
	}
	return configCopy
}

// Cache Singleton
type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

var (
	cacheInstance *Cache
	cacheOnce     sync.Once
)

func GetCache() *Cache {
	cacheOnce.Do(func() {
		cacheInstance = &Cache{
			data: make(map[string]interface{}),
		}
	})
	return cacheInstance
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.data[key] = value
	fmt.Printf("Cache: Set %s\n", key)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	value, exists := c.data[key]
	if exists {
		fmt.Printf("Cache: Hit %s\n", key)
	} else {
		fmt.Printf("Cache: Miss %s\n", key)
	}
	return value, exists
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.data, key)
	fmt.Printf("Cache: Deleted %s\n", key)
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.data = make(map[string]interface{})
	fmt.Println("Cache: Cleared")
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return len(c.data)
}

// Payment Gateway Singleton
type PaymentGateway struct {
	name        string
	apiKey      string
	transactions []Transaction
	mu          sync.RWMutex
}

type Transaction struct {
	ID        string
	Amount    float64
	Currency  string
	Status    string
	Timestamp time.Time
}

var (
	paymentInstance *PaymentGateway
	paymentOnce     sync.Once
)

func GetPaymentGateway() *PaymentGateway {
	paymentOnce.Do(func() {
		paymentInstance = &PaymentGateway{
			name:         "StripeGateway",
			apiKey:       "sk_test_123456789",
			transactions: make([]Transaction, 0),
		}
	})
	return paymentInstance
}

func (pg *PaymentGateway) ProcessPayment(amount float64, currency string) (string, error) {
	pg.mu.Lock()
	defer pg.mu.Unlock()
	
	// Simulate payment processing
	time.Sleep(200 * time.Millisecond)
	
	transactionID := fmt.Sprintf("txn_%d", time.Now().Unix())
	transaction := Transaction{
		ID:        transactionID,
		Amount:    amount,
		Currency:  currency,
		Status:    "completed",
		Timestamp: time.Now(),
	}
	
	pg.transactions = append(pg.transactions, transaction)
	
	fmt.Printf("Payment processed: $%.2f %s (ID: %s)\n", amount, currency, transactionID)
	return transactionID, nil
}

func (pg *PaymentGateway) GetTransaction(transactionID string) (*Transaction, error) {
	pg.mu.RLock()
	defer pg.mu.RUnlock()
	
	for _, txn := range pg.transactions {
		if txn.ID == transactionID {
			return &txn, nil
		}
	}
	return nil, fmt.Errorf("transaction not found: %s", transactionID)
}

func (pg *PaymentGateway) GetAllTransactions() []Transaction {
	pg.mu.RLock()
	defer pg.mu.RUnlock()
	
	transactionsCopy := make([]Transaction, len(pg.transactions))
	copy(transactionsCopy, pg.transactions)
	return transactionsCopy
}

// Email Service Singleton
type EmailService struct {
	provider string
	apiKey   string
	sent     []Email
	mu       sync.RWMutex
}

type Email struct {
	ID      string
	To      string
	Subject string
	Body    string
	Sent    bool
}

var (
	emailInstance *EmailService
	emailOnce     sync.Once
)

func GetEmailService() *EmailService {
	emailOnce.Do(func() {
		emailInstance = &EmailService{
			provider: "SendGrid",
			apiKey:   "SG.123456789",
			sent:     make([]Email, 0),
		}
	})
	return emailInstance
}

func (es *EmailService) SendEmail(to, subject, body string) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	
	// Simulate email sending
	time.Sleep(150 * time.Millisecond)
	
	emailID := fmt.Sprintf("email_%d", time.Now().Unix())
	email := Email{
		ID:      emailID,
		To:      to,
		Subject: subject,
		Body:    body,
		Sent:    true,
	}
	
	es.sent = append(es.sent, email)
	
	fmt.Printf("Email sent to %s: %s (ID: %s)\n", to, subject, emailID)
	return nil
}

func (es *EmailService) GetSentEmails() []Email {
	es.mu.RLock()
	defer es.mu.RUnlock()
	
	emailsCopy := make([]Email, len(es.sent))
	copy(emailsCopy, es.sent)
	return emailsCopy
}

func demonstrateSingleton() {
	fmt.Println("--- Basic Singleton ---")
	singleton1 := GetBasicSingleton()
	singleton2 := GetBasicSingleton()
	
	fmt.Printf("Instance 1 data: %s\n", singleton1.GetData())
	fmt.Printf("Instance 2 data: %s\n", singleton2.GetData())
	fmt.Printf("Same instance? %t\n\n", singleton1 == singleton2)
	
	singleton1.SetData("Modified Data")
	fmt.Printf("Instance 1 data: %s\n", singleton1.GetData())
	fmt.Printf("Instance 2 data: %s\n", singleton2.GetData())
}

func demonstrateThreadSafeSingleton() {
	fmt.Println("--- Thread-Safe Singleton ---")
	singleton1 := GetThreadSafeSingleton()
	singleton2 := GetThreadSafeSingleton()
	
	fmt.Printf("Instance 1 data: %s\n", singleton1.GetData())
	fmt.Printf("Instance 2 data: %s\n", singleton2.GetData())
	fmt.Printf("Same instance? %t\n\n", singleton1 == singleton2)
}

func demonstrateLoggerSingleton() {
	fmt.Println("--- Logger Singleton ---")
	logger1 := GetLogger()
	logger2 := GetLogger()
	
	logger1.Log("Application started")
	logger2.Log("User logged in")
	
	fmt.Printf("Same logger instance? %t\n", logger1 == logger2)
	
	logs := logger1.GetLogs()
	fmt.Printf("Total logs: %d\n", len(logs))
}

func demonstrateDatabaseSingleton() {
	fmt.Println("\n--- Database Connection Singleton ---")
	db1 := GetDatabaseConnection()
	db2 := GetDatabaseConnection()
	
	fmt.Printf("Same database instance? %t\n", db1 == db2)
	
	db1.Connect()
	fmt.Printf("DB1 connected: %t\n", db1.IsConnected())
	fmt.Printf("DB2 connected: %t\n", db2.IsConnected())
	
	db2.ExecuteQuery("SELECT * FROM users")
	
	db2.Disconnect()
	fmt.Printf("DB1 connected: %t\n", db1.IsConnected())
	fmt.Printf("DB2 connected: %t\n", db2.IsConnected())
}

func demonstrateConfigSingleton() {
	fmt.Println("\n--- Configuration Manager Singleton ---")
	config1 := GetConfigManager()
	config2 := GetConfigManager()
	
	fmt.Printf("Same config instance? %t\n", config1 == config2)
	
	// Get initial config
	appName, _ := config1.Get("app_name")
	fmt.Printf("App name: %s\n", appName)
	
	// Update config
	config2.Set("app_name", "UpdatedApplication")
	
	// Check if both instances see the change
	appName1, _ := config1.Get("app_name")
	appName2, _ := config2.Get("app_name")
	fmt.Printf("Config1 app name: %s\n", appName1)
	fmt.Printf("Config2 app name: %s\n", appName2)
}

func demonstrateCacheSingleton() {
	fmt.Println("\n--- Cache Singleton ---")
	cache1 := GetCache()
	cache2 := GetCache()
	
	fmt.Printf("Same cache instance? %t\n", cache1 == cache2)
	
	cache1.Set("user:1", "John Doe")
	cache1.Set("user:2", "Jane Smith")
	
	value1, _ := cache2.Get("user:1")
	value2, _ := cache1.Get("user:2")
	
	fmt.Printf("User 1: %s\n", value1)
	fmt.Printf("User 2: %s\n", value2)
	fmt.Printf("Cache size: %d\n", cache2.Size())
}

func demonstratePaymentGatewaySingleton() {
	fmt.Println("\n--- Payment Gateway Singleton ---")
	gateway1 := GetPaymentGateway()
	gateway2 := GetPaymentGateway()
	
	fmt.Printf("Same gateway instance? %t\n", gateway1 == gateway2)
	
	tx1, _ := gateway1.ProcessPayment(100.50, "USD")
	tx2, _ := gateway2.ProcessPayment(75.25, "EUR")
	
	fmt.Printf("Gateway1 transactions: %d\n", len(gateway1.GetAllTransactions()))
	fmt.Printf("Gateway2 transactions: %d\n", len(gateway2.GetAllTransactions()))
	
	transaction, _ := gateway1.GetTransaction(tx1)
	fmt.Printf("Transaction %s: $%.2f %s\n", transaction.ID, transaction.Amount, transaction.Currency)
}

func demonstrateEmailServiceSingleton() {
	fmt.Println("\n--- Email Service Singleton ---")
	email1 := GetEmailService()
	email2 := GetEmailService()
	
	fmt.Printf("Same email service instance? %t\n", email1 == email2)
	
	email1.SendEmail("user@example.com", "Welcome!", "Welcome to our service!")
	email2.SendEmail("admin@example.com", "Alert", "System maintenance scheduled.")
	
	sentEmails := email1.GetSentEmails()
	fmt.Printf("Total emails sent: %d\n", len(sentEmails))
}

func main() {
	fmt.Println("=== Singleton Pattern Demo ===")
	
	demonstrateSingleton()
	demonstrateThreadSafeSingleton()
	demonstrateLoggerSingleton()
	demonstrateDatabaseSingleton()
	demonstrateConfigSingleton()
	demonstrateCacheSingleton()
	demonstratePaymentGatewaySingleton()
	demonstrateEmailServiceSingleton()
	
	fmt.Println("\nAll singleton patterns demonstrated successfully!")
}
