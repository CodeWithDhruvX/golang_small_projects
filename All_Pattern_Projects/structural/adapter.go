package main

import "fmt"

// Adapter Pattern

// Target interface
type Target interface {
	Request() string
}

// Adaptee (the class that needs adaptation)
type Adaptee struct{}

func (a *Adaptee) SpecificRequest() string {
	return "Specific request from Adaptee"
}

// Adapter class
type Adapter struct {
	adaptee *Adaptee
}

func NewAdapter(adaptee *Adaptee) *Adapter {
	return &Adapter{adaptee: adaptee}
}

func (a *Adapter) Request() string {
	return "Adapter: " + a.adaptee.SpecificRequest()
}

// Client code
func clientCode(target Target) {
	fmt.Println(target.Request())
}

// Audio Player Example
type AudioPlayer interface {
	Play(audioType string, fileName string)
}

type MP3Player struct{}

func (mp3 *MP3Player) Play(audioType, fileName string) {
	if audioType == "mp3" {
		fmt.Printf("Playing MP3 file: %s\n", fileName)
	} else {
		fmt.Printf("Unsupported audio type: %s\n", audioType)
	}
}

// Advanced Media Player (Adaptee)
type AdvancedMediaPlayer interface {
	PlayVLC(fileName string)
	PlayMP4(fileName string)
}

type VLCPlayer struct{}

func (vlc *VLCPlayer) PlayVLC(fileName string) {
	fmt.Printf("Playing VLC file: %s\n", fileName)
}

func (vlc *VLCPlayer) PlayMP4(fileName string) {
	// VLC player doesn't support MP4
}

type MP4Player struct{}

func (mp4 *MP4Player) PlayVLC(fileName string) {
	// MP4 player doesn't support VLC
}

func (mp4 *MP4Player) PlayMP4(fileName string) {
	fmt.Printf("Playing MP4 file: %s\n", fileName)
}

// Media Adapter
type MediaAdapter struct {
	vlcPlayer AdvancedMediaPlayer
	mp4Player AdvancedMediaPlayer
}

func NewMediaAdapter(audioType string) *MediaAdapter {
	adapter := &MediaAdapter{}
	
	if audioType == "vlc" {
		adapter.vlcPlayer = &VLCPlayer{}
	} else if audioType == "mp4" {
		adapter.mp4Player = &MP4Player{}
	}
	
	return adapter
}

func (ma *MediaAdapter) Play(audioType, fileName string) {
	if audioType == "vlc" {
		if ma.vlcPlayer != nil {
			ma.vlcPlayer.PlayVLC(fileName)
		}
	} else if audioType == "mp4" {
		if ma.mp4Player != nil {
			ma.mp4Player.PlayMP4(fileName)
		}
	}
}

// Audio Player implementation
type AudioPlayerImpl struct {
	mp3Player   *MP3Player
	mediaAdapter *MediaAdapter
}

func NewAudioPlayerImpl() *AudioPlayerImpl {
	return &AudioPlayerImpl{
		mp3Player: &MP3Player{},
	}
}

func (api *AudioPlayerImpl) Play(audioType, fileName string) {
	if audioType == "mp3" {
		api.mp3Player.Play(audioType, fileName)
	} else if audioType == "vlc" || audioType == "mp4" {
		api.mediaAdapter = NewMediaAdapter(audioType)
		api.mediaAdapter.Play(audioType, fileName)
	} else {
		fmt.Printf("Invalid media. %s format not supported\n", audioType)
	}
}

// Payment Gateway Example
type ModernPaymentGateway interface {
	ProcessPayment(amount float64) string
}

type StripeGateway struct{}

func (sg *StripeGateway) ProcessPayment(amount float64) string {
	return fmt.Sprintf("Stripe: Processed payment of $%.2f", amount)
}

// Legacy Payment System (Adaptee)
type LegacyPaymentSystem struct{}

func (lps *LegacyPaymentSystem) MakePayment(cents int) string {
	return fmt.Sprintf("Legacy: Payment of %d cents processed", cents)
}

func (lps *LegacyPaymentSystem) ValidateAccount(accountNumber string) bool {
	return len(accountNumber) == 10
}

// Payment Adapter
type PaymentAdapter struct {
	legacySystem *LegacyPaymentSystem
}

func NewPaymentAdapter(legacySystem *LegacyPaymentSystem) *PaymentAdapter {
	return &PaymentAdapter{legacySystem: legacySystem}
}

func (pa *PaymentAdapter) ProcessPayment(amount float64) string {
	cents := int(amount * 100)
	return pa.legacySystem.MakePayment(cents)
}

// Temperature Sensor Example
type TemperatureSensor interface {
	GetTemperatureCelsius() float64
}

type FahrenheitSensor struct{}

func (fs *FahrenheitSensor) GetTemperatureFahrenheit() float64 {
	return 98.6 // Normal body temperature in Fahrenheit
}

// Temperature Adapter
type TemperatureAdapter struct {
	fahrenheitSensor *FahrenheitSensor
}

func NewTemperatureAdapter(fahrenheitSensor *FahrenheitSensor) *TemperatureAdapter {
	return &TemperatureAdapter{fahrenheitSensor: fahrenheitSensor}
}

func (ta *TemperatureAdapter) GetTemperatureCelsius() float64 {
	fahrenheit := ta.fahrenheitSensor.GetTemperatureFahrenheit()
	celsius := (fahrenheit - 32) * 5 / 9
	return celsius
}

// Database Adapter Example
type ModernDatabase interface {
	Query(sql string) []map[string]interface{}
	Insert(table string, data map[string]interface{}) error
	Update(table string, id string, data map[string]interface{}) error
	Delete(table string, id string) error
}

type LegacyDatabase struct{}

func (ld *LegacyDatabase) ExecuteCommand(command string) string {
	switch command {
	case "SELECT_ALL":
		return "Data1,Data2,Data3"
	case "INSERT_DATA":
		return "Data inserted"
	case "UPDATE_DATA":
		return "Data updated"
	case "DELETE_DATA":
		return "Data deleted"
	default:
		return "Unknown command"
	}
}

// Database Adapter
type DatabaseAdapter struct {
	legacyDB *LegacyDatabase
}

func NewDatabaseAdapter(legacyDB *LegacyDatabase) *DatabaseAdapter {
	return &DatabaseAdapter{legacyDB: legacyDB}
}

func (da *DatabaseAdapter) Query(sql string) []map[string]interface{} {
	result := da.legacyDB.ExecuteCommand("SELECT_ALL")
	
	// Parse legacy result into modern format
	records := []map[string]interface{}{}
	data := map[string]interface{}{
		"data": result,
	}
	records = append(records, data)
	
	return records
}

func (da *DatabaseAdapter) Insert(table string, data map[string]interface{}) error {
	da.legacyDB.ExecuteCommand("INSERT_DATA")
	return nil
}

func (da *DatabaseAdapter) Update(table string, id string, data map[string]interface{}) error {
	da.legacyDB.ExecuteCommand("UPDATE_DATA")
	return nil
}

func (da *DatabaseAdapter) Delete(table string, id string) error {
	da.legacyDB.ExecuteCommand("DELETE_DATA")
	return nil
}

// File System Adapter Example
type ModernFileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	DeleteFile(path string) error
	ListFiles(directory string) ([]string, error)
}

type LegacyFileSystem struct{}

func (lfs *LegacyFileSystem) OpenFile(filename string) string {
	return fmt.Sprintf("Content of %s", filename)
}

func (lfs *LegacyFileSystem) SaveFile(filename, content string) string {
	return fmt.Sprintf("Saved %s with content: %s", filename, content)
}

func (lfs *LegacyFileSystem) RemoveFile(filename string) string {
	return fmt.Sprintf("Removed %s", filename)
}

func (lfs *LegacyFileSystem) GetDirectoryContents(dirname string) []string {
	return []string{"file1.txt", "file2.txt", "file3.txt"}
}

// File System Adapter
type FileSystemAdapter struct {
	legacyFS *LegacyFileSystem
}

func NewFileSystemAdapter(legacyFS *LegacyFileSystem) *FileSystemAdapter {
	return &FileSystemAdapter{legacyFS: legacyFS}
}

func (fsa *FileSystemAdapter) ReadFile(path string) ([]byte, error) {
	content := fsa.legacyFS.OpenFile(path)
	return []byte(content), nil
}

func (fsa *FileSystemAdapter) WriteFile(path string, data []byte) error {
	content := string(data)
	fsa.legacyFS.SaveFile(path, content)
	return nil
}

func (fsa *FileSystemAdapter) DeleteFile(path string) error {
	fsa.legacyFS.RemoveFile(path)
	return nil
}

func (fsa *FileSystemAdapter) ListFiles(directory string) ([]string, error) {
	return fsa.legacyFS.GetDirectoryContents(directory), nil
}

// Email Service Adapter Example
type EmailService interface {
	SendEmail(to, subject, body string) error
	SendBulkEmail(recipients []string, subject, body string) error
}

type SMSService struct{}

func (sms *SMSService) SendSMS(phoneNumber, message string) string {
	return fmt.Sprintf("SMS sent to %s: %s", phoneNumber, message)
}

func (sms *SMSService) SendBulkSMS(phoneNumbers []string, message string) string {
	return fmt.Sprintf("Bulk SMS sent to %d numbers", len(phoneNumbers))
}

// Email to SMS Adapter
type EmailToSMSAdapter struct {
	smsService *SMSService
}

func NewEmailToSMSAdapter(smsService *SMSService) *EmailToSMSAdapter {
	return &EmailToSMSAdapter{smsService: smsService}
}

func (etsa *EmailToSMSAdapter) SendEmail(to, subject, body string) error {
	// Convert email to SMS format
	phoneNumber := to // Assuming email address is actually a phone number
	message := fmt.Sprintf("Subject: %s\nBody: %s", subject, body)
	etsa.smsService.SendSMS(phoneNumber, message)
	return nil
}

func (etsa *EmailToSMSAdapter) SendBulkEmail(recipients []string, subject, body string) error {
	message := fmt.Sprintf("Subject: %s\nBody: %s", subject, body)
	etsa.smsService.SendBulkSMS(recipients, message)
	return nil
}

func main() {
	fmt.Println("=== Adapter Pattern Demo ===")
	
	// Basic example
	fmt.Println("\n--- Basic Adapter Example ---")
	adaptee := &Adaptee{}
	adapter := NewAdapter(adaptee)
	
	fmt.Println("Client code with adapter:")
	clientCode(adapter)
	
	fmt.Println("\nDirect call to adaptee:")
	fmt.Println(adaptee.SpecificRequest())
	
	// Audio Player example
	fmt.Println("\n--- Audio Player Example ---")
	audioPlayer := NewAudioPlayerImpl()
	
	audioPlayer.Play("mp3", "song.mp3")
	audioPlayer.Play("mp4", "movie.mp4")
	audioPlayer.Play("vlc", "video.vlc")
	audioPlayer.Play("avi", "movie.avi")
	
	// Payment Gateway example
	fmt.Println("\n--- Payment Gateway Example ---")
	
	fmt.Println("Using modern payment gateway:")
	stripeGateway := &StripeGateway{}
	fmt.Println(stripeGateway.ProcessPayment(100.50))
	
	fmt.Println("\nUsing legacy payment system through adapter:")
	legacySystem := &LegacyPaymentSystem{}
	paymentAdapter := NewPaymentAdapter(legacySystem)
	fmt.Println(paymentAdapter.ProcessPayment(75.25))
	
	// Temperature Sensor example
	fmt.Println("\n--- Temperature Sensor Example ---")
	fahrenheitSensor := &FahrenheitSensor{}
	temperatureAdapter := NewTemperatureAdapter(fahrenheitSensor)
	
	celsius := temperatureAdapter.GetTemperatureCelsius()
	fmt.Printf("Temperature in Celsius: %.2f°C\n", celsius)
	
	// Database Adapter example
	fmt.Println("\n--- Database Adapter Example ---")
	legacyDB := &LegacyDatabase{}
	dbAdapter := NewDatabaseAdapter(legacyDB)
	
	fmt.Println("Querying database:")
	results := dbAdapter.Query("SELECT * FROM users")
	fmt.Printf("Results: %v\n", results)
	
	fmt.Println("\nInserting data:")
	dbAdapter.Insert("users", map[string]interface{}{"name": "John", "age": 30})
	
	// File System Adapter example
	fmt.Println("\n--- File System Adapter Example ---")
	legacyFS := &LegacyFileSystem{}
	fsAdapter := NewFileSystemAdapter(legacyFS)
	
	fmt.Println("Reading file:")
	content, _ := fsAdapter.ReadFile("test.txt")
	fmt.Printf("Content: %s\n", string(content))
	
	fmt.Println("\nWriting file:")
	fsAdapter.WriteFile("output.txt", []byte("Hello, World!"))
	
	fmt.Println("\nListing files:")
	files, _ := fsAdapter.ListFiles("/documents")
	fmt.Printf("Files: %v\n", files)
	
	// Email Service Adapter example
	fmt.Println("\n--- Email Service Adapter Example ---")
	smsService := &SMSService{}
	emailAdapter := NewEmailToSMSAdapter(smsService)
	
	fmt.Println("Sending email through SMS adapter:")
	emailAdapter.SendEmail("5551234567", "Welcome", "Welcome to our service!")
	
	fmt.Println("\nSending bulk email through SMS adapter:")
	recipients := []string{"5551234567", "5559876543", "5555555555"}
	emailAdapter.SendBulkEmail(recipients, "Newsletter", "Latest updates...")
	
	fmt.Println("\nAll adapter patterns demonstrated successfully!")
}
