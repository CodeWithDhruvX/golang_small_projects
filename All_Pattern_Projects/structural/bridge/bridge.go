package main

import "fmt"

// Bridge Pattern

// Abstraction interface
type Abstraction interface {
	Operation() string
}

// Refined Abstraction
type RefinedAbstraction struct {
	implementation Implementation
}

func (ra *RefinedAbstraction) Operation() string {
	return "RefinedAbstraction: " + ra.implementation.OperationImplementation()
}

// Implementation interface
type Implementation interface {
	OperationImplementation() string
}

// Concrete Implementation A
type ConcreteImplementationA struct{}

func (cia *ConcreteImplementationA) OperationImplementation() string {
	return "ConcreteImplementationA: Here's the result on the platform A."
}

// Concrete Implementation B
type ConcreteImplementationB struct{}

func (cib *ConcreteImplementationB) OperationImplementation() string {
	return "ConcreteImplementationB: Here's the result on the platform B."
}

// Remote Control Example
type RemoteControl interface {
	TurnOn()
	TurnOff()
	SetChannel(channel int)
}

type TV interface {
	TurnOn()
	TurnOff()
	ChangeChannel(channel int)
	GetCurrentChannel() int
}

type BasicRemote struct {
	device TV
}

func NewBasicRemote(device TV) *BasicRemote {
	return &BasicRemote{device: device}
}

func (br *BasicRemote) TurnOn() {
	br.device.TurnOn()
}

func (br *BasicRemote) TurnOff() {
	br.device.TurnOff()
}

func (br *BasicRemote) SetChannel(channel int) {
	br.device.ChangeChannel(channel)
}

type AdvancedRemote struct {
	*BasicRemote
}

func NewAdvancedRemote(device TV) *AdvancedRemote {
	return &AdvancedRemote{BasicRemote: NewBasicRemote(device)}
}

func (ar *AdvancedRemote) TurnOn() {
	fmt.Println("Advanced Remote: Powering on with smart features...")
	ar.device.TurnOn()
}

func (ar *AdvancedRemote) TurnOff() {
	fmt.Println("Advanced Remote: Powering off with energy saving...")
	ar.device.TurnOff()
}

func (ar *AdvancedRemote) SetChannel(channel int) {
	fmt.Printf("Advanced Remote: Setting channel %d with voice control...\n", channel)
	ar.device.ChangeChannel(channel)
}

func (ar *AdvancedRemote) VoiceCommand(command string) {
	fmt.Printf("Advanced Remote: Executing voice command: %s\n", command)
	switch command {
	case "volume up":
		fmt.Println("Advanced Remote: Volume increased")
	case "volume down":
		fmt.Println("Advanced Remote: Volume decreased")
	case "mute":
		fmt.Println("Advanced Remote: Muted")
	default:
		fmt.Printf("Advanced Remote: Unknown command: %s\n", command)
	}
}

type SonyTV struct {
	currentChannel int
	isOn           bool
}

func NewSonyTV() *SonyTV {
	return &SonyTV{currentChannel: 1, isOn: false}
}

func (stv *SonyTV) TurnOn() {
	stv.isOn = true
	fmt.Println("Sony TV: Turned on")
}

func (stv *SonyTV) TurnOff() {
	stv.isOn = false
	fmt.Println("Sony TV: Turned off")
}

func (stv *SonyTV) ChangeChannel(channel int) {
	if stv.isOn {
		stv.currentChannel = channel
		fmt.Printf("Sony TV: Changed to channel %d\n", channel)
	} else {
		fmt.Println("Sony TV: Cannot change channel when off")
	}
}

func (stv *SonyTV) GetCurrentChannel() int {
	return stv.currentChannel
}

type SamsungTV struct {
	currentChannel int
	isOn           bool
}

func NewSamsungTV() *SamsungTV {
	return &SamsungTV{currentChannel: 1, isOn: false}
}

func (sstv *SamsungTV) TurnOn() {
	sstv.isOn = true
	fmt.Println("Samsung TV: Turned on with smart features")
}

func (sstv *SamsungTV) TurnOff() {
	sstv.isOn = false
	fmt.Println("Samsung TV: Turned off")
}

func (sstv *SamsungTV) ChangeChannel(channel int) {
	if sstv.isOn {
		sstv.currentChannel = channel
		fmt.Printf("Samsung TV: Changed to channel %d with HD quality\n", channel)
	} else {
		fmt.Println("Samsung TV: Cannot change channel when off")
	}
}

func (sstv *SamsungTV) GetCurrentChannel() int {
	return sstv.currentChannel
}

// Message System Example
type MessageSender interface {
	SendMessage(message string) string
}

type EmailSender struct{}

func (es *EmailSender) SendMessage(message string) string {
	return fmt.Sprintf("Email: %s", message)
}

type SMSSender struct{}

func (ss *SMSSender) SendMessage(message string) string {
	return fmt.Sprintf("SMS: %s", message)
}

type PushNotificationSender struct{}

func (pns *PushNotificationSender) SendMessage(message string) string {
	return fmt.Sprintf("Push: %s", message)
}

type MessageSystem struct {
	sender MessageSender
}

func NewMessageSystem(sender MessageSender) *MessageSystem {
	return &MessageSystem{sender: sender}
}

func (ms *MessageSystem) SendUrgentMessage(message string) {
	urgentMsg := fmt.Sprintf("URGENT: %s", message)
	fmt.Println(ms.sender.SendMessage(urgentMsg))
}

func (ms *MessageSystem) SendNormalMessage(message string) {
	fmt.Println(ms.sender.SendMessage(message))
}

type AdvancedMessageSystem struct {
	*MessageSystem
}

func NewAdvancedMessageSystem(sender MessageSender) *AdvancedMessageSystem {
	return &AdvancedMessageSystem{MessageSystem: NewMessageSystem(sender)}
}

func (ams *AdvancedMessageSystem) SendUrgentMessage(message string) {
	urgentMsg := fmt.Sprintf("🚨 URGENT 🚨: %s", message)
	fmt.Println(ams.sender.SendMessage(urgentMsg))
}

func (ams *AdvancedMessageSystem) SendNormalMessage(message string) {
	normalMsg := fmt.Sprintf("ℹ️ INFO: %s", message)
	fmt.Println(ams.sender.SendMessage(normalMsg))
}

func (ams *AdvancedMessageSystem) SendScheduledMessage(message string, delay int) {
	fmt.Printf("Scheduling message to be sent in %d minutes: %s\n", delay, message)
	scheduledMsg := fmt.Sprintf("SCHEDULED: %s", message)
	fmt.Println(ams.sender.SendMessage(scheduledMsg))
}

// Drawing System Example
type DrawingAPI interface {
	DrawCircle(x, y, radius int)
	DrawRectangle(x, y, width, height int)
}

type WindowsDrawingAPI struct{}

func (wda *WindowsDrawingAPI) DrawCircle(x, y, radius int) {
	fmt.Printf("Windows API: Drawing circle at (%d,%d) with radius %d\n", x, y, radius)
}

func (wda *WindowsDrawingAPI) DrawRectangle(x, y, width, height int) {
	fmt.Printf("Windows API: Drawing rectangle at (%d,%d) with size %dx%d\n", x, y, width, height)
}

type MacDrawingAPI struct{}

func (mda *MacDrawingAPI) DrawCircle(x, y, radius int) {
	fmt.Printf("Mac API: Drawing circle at (%d,%d) with radius %d using Core Graphics\n", x, y, radius)
}

func (mda *MacDrawingAPI) DrawRectangle(x, y, width, height int) {
	fmt.Printf("Mac API: Drawing rectangle at (%d,%d) with size %dx%d using Core Graphics\n", x, y, width, height)
}

type Shape interface {
	Draw()
}

type Circle struct {
	x, y, radius int
	api          DrawingAPI
}

func NewCircle(x, y, radius int, api DrawingAPI) *Circle {
	return &Circle{x: x, y: y, radius: radius, api: api}
}

func (c *Circle) Draw() {
	c.api.DrawCircle(c.x, c.y, c.radius)
}

type Rectangle struct {
	x, y, width, height int
	api                 DrawingAPI
}

func NewRectangle(x, y, width, height int, api DrawingAPI) *Rectangle {
	return &Rectangle{x: x, y: y, width: width, height: height, api: api}
}

func (r *Rectangle) Draw() {
	r.api.DrawRectangle(r.x, r.y, r.width, r.height)
}

// File System Example
type FileSystem interface {
	ReadFile(path string) string
	WriteFile(path, content string)
	DeleteFile(path string)
}

type LocalFileSystem struct{}

func (lfs *LocalFileSystem) ReadFile(path string) string {
	return fmt.Sprintf("Local FS: Reading file %s", path)
}

func (lfs *LocalFileSystem) WriteFile(path, content string) {
	fmt.Printf("Local FS: Writing to file %s with content: %s\n", path, content)
}

func (lfs *LocalFileSystem) DeleteFile(path string) {
	fmt.Printf("Local FS: Deleting file %s\n", path)
}

type CloudFileSystem struct{}

func (cfs *CloudFileSystem) ReadFile(path string) string {
	return fmt.Sprintf("Cloud FS: Reading file %s from cloud storage", path)
}

func (cfs *CloudFileSystem) WriteFile(path, content string) {
	fmt.Printf("Cloud FS: Writing to file %s in cloud storage with content: %s\n", path, content)
}

func (cfs *CloudFileSystem) DeleteFile(path string) {
	fmt.Printf("Cloud FS: Deleting file %s from cloud storage\n", path)
}

type FileManager struct {
	fileSystem FileSystem
}

func NewFileManager(fileSystem FileSystem) *FileManager {
	return &FileManager{fileSystem: fileSystem}
}

func (fm *FileManager) OpenFile(path string) {
	content := fm.fileSystem.ReadFile(path)
	fmt.Println(content)
}

func (fm *FileManager) SaveFile(path, content string) {
	fm.fileSystem.WriteFile(path, content)
}

func (fm *FileManager) RemoveFile(path string) {
	fm.fileSystem.DeleteFile(path)
}

type AdvancedFileManager struct {
	*FileManager
}

func NewAdvancedFileManager(fileSystem FileSystem) *AdvancedFileManager {
	return &AdvancedFileManager{FileManager: NewFileManager(fileSystem)}
}

func (afm *AdvancedFileManager) OpenFile(path string) {
	fmt.Printf("Advanced File Manager: Opening file with enhanced features...\n")
	content := afm.fileSystem.ReadFile(path)
	fmt.Println(content)
}

func (afm *AdvancedFileManager) SaveFile(path, content string) {
	fmt.Printf("Advanced File Manager: Saving file with compression and encryption...\n")
	afm.fileSystem.WriteFile(path, content)
}

func (afm *AdvancedFileManager) RemoveFile(path string) {
	fmt.Printf("Advanced File Manager: Moving file to recycle bin before deletion...\n")
	afm.fileSystem.DeleteFile(path)
}

func (afm *AdvancedFileManager) BackupFile(path string) {
	fmt.Printf("Advanced File Manager: Creating backup of file %s\n", path)
	backupPath := path + ".backup"
	content := afm.fileSystem.ReadFile(path)
	afm.fileSystem.WriteFile(backupPath, content)
}

func main() {
	fmt.Println("=== Bridge Pattern Demo ===")
	
	// Basic example
	fmt.Println("\n--- Basic Bridge Example ---")
	implementationA := &ConcreteImplementationA{}
	abstraction := &RefinedAbstraction{implementation: implementationA}
	fmt.Println(abstraction.Operation())
	
	implementationB := &ConcreteImplementationB{}
	abstraction = &RefinedAbstraction{implementation: implementationB}
	fmt.Println(abstraction.Operation())
	
	// Remote Control example
	fmt.Println("\n--- Remote Control Example ---")
	
	fmt.Println("Basic Remote with Sony TV:")
	sonyTV := NewSonyTV()
	basicRemote := NewBasicRemote(sonyTV)
	basicRemote.TurnOn()
	basicRemote.SetChannel(5)
	basicRemote.TurnOff()
	
	fmt.Println("\nAdvanced Remote with Samsung TV:")
	samsungTV := NewSamsungTV()
	advancedRemote := NewAdvancedRemote(samsungTV)
	advancedRemote.TurnOn()
	advancedRemote.SetChannel(10)
	advancedRemote.VoiceCommand("volume up")
	advancedRemote.TurnOff()
	
	// Message System example
	fmt.Println("\n--- Message System Example ---")
	
	fmt.Println("Basic Message System with Email:")
	emailSystem := NewMessageSystem(&EmailSender{})
	emailSystem.SendNormalMessage("Hello World")
	emailSystem.SendUrgentMessage("Server is down!")
	
	fmt.Println("\nAdvanced Message System with SMS:")
	advancedSystem := NewAdvancedMessageSystem(&SMSSender{})
	advancedSystem.SendNormalMessage("Meeting at 3 PM")
	advancedSystem.SendUrgentMessage("Fire alarm activated")
	advancedSystem.SendScheduledMessage("Happy Birthday!", 60)
	
	// Drawing System example
	fmt.Println("\n--- Drawing System Example ---")
	
	fmt.Println("Drawing shapes on Windows:")
	windowsAPI := &WindowsDrawingAPI{}
	circle1 := NewCircle(10, 20, 5, windowsAPI)
	rectangle1 := NewRectangle(30, 40, 100, 50, windowsAPI)
	circle1.Draw()
	rectangle1.Draw()
	
	fmt.Println("\nDrawing shapes on Mac:")
	macAPI := &MacDrawingAPI{}
	circle2 := NewCircle(50, 60, 8, macAPI)
	rectangle2 := NewRectangle(70, 80, 120, 60, macAPI)
	circle2.Draw()
	rectangle2.Draw()
	
	// File System example
	fmt.Println("\n--- File System Example ---")
	
	fmt.Println("Basic File Manager with Local FS:")
	localFS := &LocalFileSystem{}
	localFileManager := NewFileManager(localFS)
	localFileManager.OpenFile("document.txt")
	localFileManager.SaveFile("output.txt", "Hello World")
	localFileManager.RemoveFile("temp.txt")
	
	fmt.Println("\nAdvanced File Manager with Cloud FS:")
	cloudFS := &CloudFileSystem{}
	advancedFileManager := NewAdvancedFileManager(cloudFS)
	advancedFileManager.OpenFile("cloud_document.txt")
	advancedFileManager.SaveFile("cloud_output.txt", "Hello Cloud")
	advancedFileManager.BackupFile("important_file.txt")
	advancedFileManager.RemoveFile("old_file.txt")
	
	fmt.Println("\nAll bridge patterns demonstrated successfully!")
}
