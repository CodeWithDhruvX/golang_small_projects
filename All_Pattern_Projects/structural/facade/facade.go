package main

import "fmt"

// Facade Pattern

// Subsystem 1
type SubsystemA struct{}

func (sa *SubsystemA) OperationA1() string {
	return "Subsystem A: Operation A1"
}

func (sa *SubsystemA) OperationA2() string {
	return "Subsystem A: Operation A2"
}

// Subsystem 2
type SubsystemB struct{}

func (sb *SubsystemB) OperationB1() string {
	return "Subsystem B: Operation B1"
}

func (sb *SubsystemB) OperationB2() string {
	return "Subsystem B: Operation B2"
}

// Subsystem 3
type SubsystemC struct{}

func (sc *SubsystemC) OperationC1() string {
	return "Subsystem C: Operation C1"
}

func (sc *SubsystemC) OperationC2() string {
	return "Subsystem C: Operation C2"
}

// Facade
type Facade struct {
	subsystemA *SubsystemA
	subsystemB *SubsystemB
	subsystemC *SubsystemC
}

func NewFacade() *Facade {
	return &Facade{
		subsystemA: &SubsystemA{},
		subsystemB: &SubsystemB{},
		subsystemC: &SubsystemC{},
	}
}

func (f *Facade) OperationWrapper1() string {
	result := "Facade initializes subsystems:\n"
	result += f.subsystemA.OperationA1() + "\n"
	result += f.subsystemB.OperationB1() + "\n"
	result += f.subsystemC.OperationC1()
	return result
}

func (f *Facade) OperationWrapper2() string {
	result := "Facade orders subsystems to perform the action:\n"
	result += f.subsystemA.OperationA2() + "\n"
	result += f.subsystemB.OperationB2() + "\n"
	result += f.subsystemC.OperationC2()
	return result
}

// Home Theater System Example
type Amplifier struct{}

func (a *Amplifier) On() {
	fmt.Println("Amplifier: On")
}

func (a *Amplifier) Off() {
	fmt.Println("Amplifier: Off")
}

func (a *Amplifier) SetSurroundSound() {
	fmt.Println("Amplifier: Surround sound set")
}

func (a *Amplifier) SetVolume(level int) {
	fmt.Printf("Amplifier: Volume set to %d\n", level)
}

type Tuner struct{}

func (t *Tuner) On() {
	fmt.Println("Tuner: On")
}

func (t *Tuner) Off() {
	fmt.Println("Tuner: Off")
}

func (t *Tuner) SetFrequency(frequency string) {
	fmt.Printf("Tuner: Frequency set to %s\n", frequency)
}

type DvdPlayer struct {
	amplifier *Amplifier
}

func NewDvdPlayer(amplifier *Amplifier) *DvdPlayer {
	return &DvdPlayer{amplifier: amplifier}
}

func (dp *DvdPlayer) On() {
	fmt.Println("DVD Player: On")
}

func (dp *DvdPlayer) Off() {
	fmt.Println("DVD Player: Off")
}

func (dp *DvdPlayer) Play(movie string) {
	fmt.Printf("DVD Player: Playing \"%s\"\n", movie)
}

func (dp *DvdPlayer) Stop() {
	fmt.Println("DVD Player: Stop")
}

func (dp *DvdPlayer) Eject() {
	fmt.Println("DVD Player: Eject")
}

type Projector struct {
	dvdPlayer *DvdPlayer
}

func NewProjector(dvdPlayer *DvdPlayer) *Projector {
	return &Projector{dvdPlayer: dvdPlayer}
}

func (p *Projector) On() {
	fmt.Println("Projector: On")
}

func (p *Projector) Off() {
	fmt.Println("Projector: Off")
}

func (p *Projector) TvMode() {
	fmt.Println("Projector: TV mode")
}

func (p *Projector) WideScreenMode() {
	fmt.Println("Projector: Wide screen mode")
}

type TheaterLights struct{}

func (tl *TheaterLights) On() {
	fmt.Println("Theater Lights: On")
}

func (tl *TheaterLights) Off() {
	fmt.Println("Theater Lights: Off")
}

func (tl *TheaterLights) Dim(level int) {
	fmt.Printf("Theater Lights: Dimming to %d%%\n", level)
}

type Screen struct{}

func (s *Screen) Up() {
	fmt.Println("Screen: Up")
}

func (s *Screen) Down() {
	fmt.Println("Screen: Down")
}

type PopcornPopper struct{}

func (pp *PopcornPopper) On() {
	fmt.Println("Popcorn Popper: On")
}

func (pp *PopcornPopper) Off() {
	fmt.Println("Popcorn Popper: Off")
}

func (pp *PopcornPopper) Pop() {
	fmt.Println("Popcorn Popper: Popping popcorn!")
}

// Home Theater Facade
type HomeTheaterFacade struct {
	amp          *Amplifier
	tuner        *Tuner
	dvd          *DvdPlayer
	projector    *Projector
	screen       *Screen
	lights       *TheaterLights
	popper       *PopcornPopper
}

func NewHomeTheaterFacade(amp *Amplifier, tuner *Tuner, dvd *DvdPlayer, projector *Projector, 
	screen *Screen, lights *TheaterLights, popper *PopcornPopper) *HomeTheaterFacade {
	return &HomeTheaterFacade{
		amp:       amp,
		tuner:     tuner,
		dvd:       dvd,
		projector: projector,
		screen:    screen,
		lights:    lights,
		popper:    popper,
	}
}

func (htf *HomeTheaterFacade) WatchMovie(movie string) {
	fmt.Println("Get ready to watch a movie...")
	htf.popper.On()
	htf.popper.Pop()
	htf.lights.On()
	htf.lights.Dim(10)
	htf.screen.Down()
	htf.projector.On()
	htf.projector.WideScreenMode()
	htf.amp.On()
	htf.amp.SetSurroundSound()
	htf.amp.SetVolume(5)
	htf.dvd.On()
	htf.dvd.Play(movie)
}

func (htf *HomeTheaterFacade) EndMovie() {
	fmt.Println("\nShutting movie theater down...")
	htf.popper.Off()
	htf.lights.On()
	htf.screen.Up()
	htf.projector.Off()
	htf.amp.Off()
	htf.dvd.Stop()
	htf.dvd.Eject()
	htf.dvd.Off()
}

// Computer System Example
type CPU struct{}

func (cpu *CPU) Freeze() {
	fmt.Println("CPU: Freezing processor")
}

func (cpu *CPU) Jump(position string) {
	fmt.Printf("CPU: Jumping to %s\n", position)
}

func (cpu *CPU) Execute() {
	fmt.Println("CPU: Executing")
}

type Memory struct{}

func (m *Memory) Load(position string, data string) {
	fmt.Printf("Memory: Loading %s at %s\n", data, position)
}

type HardDrive struct{}

func (hd *HardDrive) Read(lba, size string) string {
	fmt.Printf("Hard Drive: Reading sector %s with size %s\n", lba, size)
	return "data"
}

type GPU struct{}

func (gpu *GPU) Render() {
	fmt.Println("GPU: Rendering graphics")
}

type SoundCard struct{}

func (sc *SoundCard) PlaySound() {
	fmt.Println("Sound Card: Playing sound")
}

type NetworkCard struct{}

func (nc *NetworkCard) Connect() {
	fmt.Println("Network Card: Connecting to network")
}

// Computer Facade
type ComputerFacade struct {
	cpu        *CPU
	memory     *Memory
	hardDrive  *HardDrive
	gpu        *GPU
	soundCard  *SoundCard
	networkCard *NetworkCard
}

func NewComputerFacade() *ComputerFacade {
	return &ComputerFacade{
		cpu:         &CPU{},
		memory:      &Memory{},
		hardDrive:   &HardDrive{},
		gpu:         &GPU{},
		soundCard:   &SoundCard{},
		networkCard: &NetworkCard{},
	}
}

func (cf *ComputerFacade) Start() {
	fmt.Println("Starting computer...")
	cf.cpu.Freeze()
	cf.memory.Load("0x00", "bootloader")
	cf.cpu.Jump("0x00")
	cf.cpu.Execute()
	cf.gpu.Render()
	cf.soundCard.PlaySound()
	cf.networkCard.Connect()
	fmt.Println("Computer started successfully!")
}

func (cf *ComputerFacade) Shutdown() {
	fmt.Println("Shutting down computer...")
	cf.networkCard.Connect()
	cf.soundCard.PlaySound()
	cf.gpu.Render()
	fmt.Println("Computer shut down successfully!")
}

// Banking System Example
type AccountService struct{}

func (as *AccountService) GetAccountBalance(accountNumber string) float64 {
	fmt.Printf("Account Service: Getting balance for account %s\n", accountNumber)
	return 1000.0
}

func (as *AccountService) UpdateAccountBalance(accountNumber string, amount float64) {
	fmt.Printf("Account Service: Updating balance for account %s to $%.2f\n", accountNumber, amount)
}

type TransferService struct{}

func (ts *TransferService) Transfer(fromAccount, toAccount string, amount float64) bool {
	fmt.Printf("Transfer Service: Transferring $%.2f from %s to %s\n", amount, fromAccount, toAccount)
	return true
}

type LoanService struct{}

func (ls *LoanService) CheckLoanEligibility(accountNumber string) bool {
	fmt.Printf("Loan Service: Checking loan eligibility for account %s\n", accountNumber)
	return true
}

func (ls *LoanService) ApproveLoan(accountNumber string, amount float64) {
	fmt.Printf("Loan Service: Loan of $%.2f approved for account %s\n", amount, accountNumber)
}

type NotificationService struct{}

func (ns *NotificationService) SendEmail(to, subject, body string) {
	fmt.Printf("Notification Service: Email sent to %s - %s\n", to, subject)
}

func (ns *NotificationService) SendSMS(to, message string) {
	fmt.Printf("Notification Service: SMS sent to %s - %s\n", to, message)
}

// Banking Facade
type BankingFacade struct {
	accountService    *AccountService
	transferService   *TransferService
	loanService       *LoanService
	notificationService *NotificationService
}

func NewBankingFacade() *BankingFacade {
	return &BankingFacade{
		accountService:     &AccountService{},
		transferService:    &TransferService{},
		loanService:        &LoanService{},
		notificationService: &NotificationService{},
	}
}

func (bf *BankingFacade) TransferMoney(fromAccount, toAccount string, amount float64) bool {
	fmt.Printf("Initiating transfer of $%.2f from %s to %s\n", amount, fromAccount, toAccount)
	
	balance := bf.accountService.GetAccountBalance(fromAccount)
	if balance < amount {
		fmt.Println("Insufficient funds")
		return false
	}
	
	success := bf.transferService.Transfer(fromAccount, toAccount, amount)
	if success {
		bf.accountService.UpdateAccountBalance(fromAccount, balance-amount)
		bf.accountService.UpdateAccountBalance(toAccount, bf.accountService.GetAccountBalance(toAccount)+amount)
		
		bf.notificationService.SendEmail(fromAccount, "Transfer Completed", 
			fmt.Sprintf("You transferred $%.2f to %s", amount, toAccount))
		bf.notificationService.SendEmail(toAccount, "Money Received", 
			fmt.Sprintf("You received $%.2f from %s", amount, fromAccount))
		
		fmt.Println("Transfer completed successfully!")
		return true
	}
	
	return false
}

func (bf *BankingFacade) ApplyForLoan(accountNumber string, amount float64) bool {
	fmt.Printf("Processing loan application of $%.2f for account %s\n", amount, accountNumber)
	
	eligible := bf.loanService.CheckLoanEligibility(accountNumber)
	if eligible {
		bf.loanService.ApproveLoan(accountNumber, amount)
		
		bf.notificationService.SendEmail(accountNumber, "Loan Approved", 
			fmt.Sprintf("Your loan of $%.2f has been approved", amount))
		bf.notificationService.SendSMS(accountNumber, "Loan approved! Check your email for details.")
		
		fmt.Println("Loan approved successfully!")
		return true
	}
	
	fmt.Println("Loan application rejected")
	return false
}

// E-commerce System Example
type InventoryService struct{}

func (is *InventoryService) CheckStock(productID string) int {
	fmt.Printf("Inventory Service: Checking stock for product %s\n", productID)
	return 10
}

func (is *InventoryService) ReduceStock(productID string, quantity int) {
	fmt.Printf("Inventory Service: Reducing stock for product %s by %d\n", productID, quantity)
}

type PaymentService struct{}

func (ps *PaymentService) ProcessPayment(amount float64, paymentMethod string) bool {
	fmt.Printf("Payment Service: Processing payment of $%.2f via %s\n", amount, paymentMethod)
	return true
}

type ShippingService struct{}

func (ss *ShippingService) CalculateShipping(productID string, address string) float64 {
	fmt.Printf("Shipping Service: Calculating shipping for product %s to %s\n", productID, address)
	return 9.99
}

func (ss *ShippingService) ScheduleShipment(productID string, address string) {
	fmt.Printf("Shipping Service: Scheduling shipment of product %s to %s\n", productID, address)
}

type OrderService struct{}

func (os *OrderService) CreateOrder(customerID string, items []string) string {
	fmt.Printf("Order Service: Creating order for customer %s with items %v\n", customerID, items)
	return "ORDER12345"
}

func (os *OrderService) UpdateOrderStatus(orderID, status string) {
	fmt.Printf("Order Service: Updating order %s status to %s\n", orderID, status)
}

// E-commerce Facade
type ECommerceFacade struct {
	inventoryService *InventoryService
	paymentService   *PaymentService
	shippingService  *ShippingService
	orderService     *OrderService
}

func NewECommerceFacade() *ECommerceFacade {
	return &ECommerceFacade{
		inventoryService: &InventoryService{},
		paymentService:   &PaymentService{},
		shippingService:  &ShippingService{},
		orderService:     &OrderService{},
	}
}

func (ecf *ECommerceFacade) PlaceOrder(customerID string, productID string, paymentMethod string, shippingAddress string) bool {
	fmt.Printf("Processing order for customer %s\n", customerID)
	
	// Check inventory
	stock := ecf.inventoryService.CheckStock(productID)
	if stock <= 0 {
		fmt.Println("Product out of stock")
		return false
	}
	
	// Create order
	orderID := ecf.orderService.CreateOrder(customerID, []string{productID})
	
	// Process payment
	productPrice := 99.99
	shippingCost := ecf.shippingService.CalculateShipping(productID, shippingAddress)
	totalAmount := productPrice + shippingCost
	
	paymentSuccess := ecf.paymentService.ProcessPayment(totalAmount, paymentMethod)
	if !paymentSuccess {
		ecf.orderService.UpdateOrderStatus(orderID, "PAYMENT_FAILED")
		fmt.Println("Payment failed")
		return false
	}
	
	// Update inventory
	ecf.inventoryService.ReduceStock(productID, 1)
	
	// Schedule shipment
	ecf.shippingService.ScheduleShipment(productID, shippingAddress)
	
	// Update order status
	ecf.orderService.UpdateOrderStatus(orderID, "CONFIRMED")
	
	fmt.Printf("Order placed successfully! Order ID: %s\n", orderID)
	fmt.Printf("Total charged: $%.2f ($%.2f + $%.2f shipping)\n", totalAmount, productPrice, shippingCost)
	
	return true
}

func main() {
	fmt.Println("=== Facade Pattern Demo ===")
	
	// Basic example
	fmt.Println("\n--- Basic Facade Example ---")
	
	subsystemA := &SubsystemA{}
	subsystemB := &SubsystemB{}
	subsystemC := &SubsystemC{}
	
	facade := NewFacade()
	fmt.Println(facade.OperationWrapper1())
	fmt.Println()
	fmt.Println(facade.OperationWrapper2())
	
	// Direct access to subsystems (complex)
	fmt.Println("\nDirect access to subsystems:")
	fmt.Println(subsystemA.OperationA1())
	fmt.Println(subsystemB.OperationB1())
	fmt.Println(subsystemC.OperationC1())
	
	// Home Theater System example
	fmt.Println("\n--- Home Theater System Example ---")
	
	amp := &Amplifier{}
	tuner := &Tuner{}
	dvd := NewDvdPlayer(amp)
	projector := NewProjector(dvd)
	lights := &TheaterLights{}
	screen := &Screen{}
	popper := &PopcornPopper{}
	
	homeTheater := NewHomeTheaterFacade(amp, tuner, dvd, projector, screen, lights, popper)
	homeTheater.WatchMovie("The Matrix")
	homeTheater.EndMovie()
	
	// Computer System example
	fmt.Println("\n--- Computer System Example ---")
	
	computer := NewComputerFacade()
	computer.Start()
	computer.Shutdown()
	
	// Banking System example
	fmt.Println("\n--- Banking System Example ---")
	
	banking := NewBankingFacade()
	banking.TransferMoney("ACC123", "ACC456", 250.0)
	banking.ApplyForLoan("ACC789", 5000.0)
	
	// E-commerce System example
	fmt.Println("\n--- E-commerce System Example ---")
	
	ecommerce := NewECommerceFacade()
	ecommerce.PlaceOrder("CUST001", "PROD123", "Credit Card", "123 Main St, City, State")
	
	fmt.Println("\nAll facade patterns demonstrated successfully!")
}
