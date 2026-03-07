package main

import (
	"fmt"
	"sync"
	"time"
)

// Saga Pattern

// Saga interface
type Saga interface {
	Execute() error
	Compensate()
	GetStatus() SagaStatus
	GetSteps() []SagaStep
}

type SagaStatus int

const (
	Pending SagaStatus = iota
	Running
	Completed
	Failed
	Compensating
	Compensated
)

// SagaStep interface
type SagaStep interface {
	Execute() error
	Compensate() error
	GetName() string
	GetStatus() StepStatus
}

type StepStatus int

const (
	StepPending StepStatus = iota
	StepCompleted
	StepFailed
	StepCompensated
)

// Basic Saga Step
type BasicSagaStep struct {
	name         string
	action       func() error
	compensation func() error
	status       StepStatus
}

func NewBasicSagaStep(name string, action func() error, compensation func() error) *BasicSagaStep {
	return &BasicSagaStep{
		name:         name,
		action:       action,
		compensation: compensation,
		status:       StepPending,
	}
}

func (bss *BasicSagaStep) Execute() error {
	err := bss.action()
	if err != nil {
		bss.status = StepFailed
		return err
	}
	bss.status = StepCompleted
	return nil
}

func (bss *BasicSagaStep) Compensate() error {
	err := bss.compensation()
	if err != nil {
		fmt.Printf("Warning: Compensation failed for step %s: %v\n", bss.name, err)
	}
	bss.status = StepCompensated
	return err
}

func (bss *BasicSagaStep) GetName() string {
	return bss.name
}

func (bss *BasicSagaStep) GetStatus() StepStatus {
	return bss.status
}

// Orchestrator Saga
type OrchestratorSaga struct {
	name         string
	steps        []SagaStep
	currentStep  int
	status       SagaStatus
	compensated  []bool
	mu           sync.Mutex
	onStepStart  func(step SagaStep)
	onStepComplete func(step SagaStep)
	onStepFailed func(step SagaStep, error error)
	onSagaComplete func()
	onSagaFailed func(error error)
}

func NewOrchestratorSaga(name string) *OrchestratorSaga {
	return &OrchestratorSaga{
		name:        name,
		steps:       make([]SagaStep, 0),
		status:      Pending,
		compensated: make([]bool, 0),
	}
}

func (os *OrchestratorSaga) AddStep(step SagaStep) {
	os.steps = append(os.steps, step)
	os.compensated = append(os.compensated, false)
}

func (os *OrchestratorSaga) Execute() error {
	os.mu.Lock()
	defer os.mu.Unlock()
	
	if os.status != Pending {
		return fmt.Errorf("saga %s: cannot execute, current status is %v", os.name, os.status)
	}
	
	os.status = Running
	
	for i, step := range os.steps {
		os.currentStep = i
		
		if os.onStepStart != nil {
			os.onStepStart(step)
		}
		
		err := step.Execute()
		if err != nil {
			os.status = Failed
			
			if os.onStepFailed != nil {
				os.onStepFailed(step, err)
			}
			
			if os.onSagaFailed != nil {
				os.onSagaFailed(err)
			}
			
			// Start compensation
			go os.Compensate()
			
			return fmt.Errorf("saga %s: step %s failed: %v", os.name, step.GetName(), err)
		}
		
		os.compensated[i] = true
		
		if os.onStepComplete != nil {
			os.onStepComplete(step)
		}
	}
	
	os.status = Completed
	
	if os.onSagaComplete != nil {
		os.onSagaComplete()
	}
	
	return nil
}

func (os *OrchestratorSaga) Compensate() {
	os.mu.Lock()
	if os.status == Completed {
		os.mu.Unlock()
		return // No need to compensate
	}
	
	os.status = Compensating
	os.mu.Unlock()
	
	// Compensate in reverse order
	for i := len(os.steps) - 1; i >= 0; i-- {
		if os.compensated[i] {
			step := os.steps[i]
			
			if step.GetStatus() == StepCompleted {
				fmt.Printf("Compensating step: %s\n", step.GetName())
				step.Compensate()
			}
		}
	}
	
	os.mu.Lock()
	os.status = Compensated
	os.mu.Unlock()
}

func (os *OrchestratorSaga) GetStatus() SagaStatus {
	os.mu.Lock()
	defer os.mu.Unlock()
	return os.status
}

func (os *OrchestratorSaga) GetSteps() []SagaStep {
	return os.steps
}

func (os *OrchestratorSaga) SetStepStartCallback(callback func(step SagaStep)) {
	os.onStepStart = callback
}

func (os *OrchestratorSaga) SetStepCompleteCallback(callback func(step SagaStep)) {
	os.onStepComplete = callback
}

func (os *OrchestratorSaga) SetStepFailedCallback(callback func(step SagaStep, error error)) {
	os.onStepFailed = callback
}

func (os *OrchestratorSaga) SetSagaCompleteCallback(callback func()) {
	os.onSagaComplete = callback
}

func (os *OrchestratorSaga) SetSagaFailedCallback(callback func(error error)) {
	os.onSagaFailed = callback
}

// Choreography Saga (Event-driven)
type ChoreographySaga struct {
	name    string
	steps   []ChoreographyStep
	status  SagaStatus
	current int
	mu      sync.Mutex
}

type ChoreographyStep struct {
	name         string
	action       func() error
	compensation func() error
	nextStep     *ChoreographyStep
	status       StepStatus
}

func NewChoreographySaga(name string) *ChoreographySaga {
	return &ChoreographySaga{
		name:   name,
		steps:  make([]ChoreographyStep, 0),
		status: Pending,
	}
}

func (cs *ChoreographySaga) AddStep(name string, action func() error, compensation func() error) *ChoreographyStep {
	step := ChoreographyStep{
		name:         name,
		action:       action,
		compensation: compensation,
		status:       StepPending,
	}
	
	if len(cs.steps) > 0 {
		cs.steps[len(cs.steps)-1].nextStep = &step
	}
	
	cs.steps = append(cs.steps, step)
	return &step
}

func (cs *ChoreographySaga) Execute() error {
	cs.mu.Lock()
	if cs.status != Pending {
		cs.mu.Unlock()
		return fmt.Errorf("saga %s: cannot execute, current status is %v", cs.name, cs.status)
	}
	cs.status = Running
	cs.mu.Unlock()
	
	if len(cs.steps) == 0 {
		return nil
	}
	
	current := &cs.steps[0]
	
	for current != nil {
		fmt.Printf("Executing step: %s\n", current.name)
		
		err := current.action()
		if err != nil {
			current.status = StepFailed
			cs.mu.Lock()
			cs.status = Failed
			cs.mu.Unlock()
			
			go cs.CompensateFrom(current)
			
			return fmt.Errorf("saga %s: step %s failed: %v", cs.name, current.name, err)
		}
		
		current.status = StepCompleted
		current = current.nextStep
	}
	
	cs.mu.Lock()
	cs.status = Completed
	cs.mu.Unlock()
	
	return nil
}

func (cs *ChoreographySaga) CompensateFrom(failedStep *ChoreographyStep) {
	fmt.Printf("Starting compensation from step: %s\n", failedStep.name)
	
	// Find the failed step in the list
	var failedIndex int = -1
	for i, step := range cs.steps {
		if step.name == failedStep.name {
			failedIndex = i
			break
		}
	}
	
	// Compensate all completed steps before the failed one
	for i := failedIndex - 1; i >= 0; i-- {
		step := &cs.steps[i]
		if step.status == StepCompleted {
			fmt.Printf("Compensating step: %s\n", step.name)
			step.Compensate()
			step.status = StepCompensated
		}
	}
	
	cs.mu.Lock()
	cs.status = Compensated
	cs.mu.Unlock()
}

func (cs *ChoreographySaga) GetStatus() SagaStatus {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.status
}

func (cs *ChoreographySaga) GetSteps() []SagaStep {
	// Convert to SagaStep interface
	steps := make([]SagaStep, len(cs.steps))
	for i, step := range cs.steps {
		steps[i] = &step
	}
	return steps
}

// Order Management Saga Example
type OrderService struct{}

func (os *OrderService) CreateOrder(orderID string) error {
	fmt.Printf("Order Service: Creating order %s\n", orderID)
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (os *OrderService) CancelOrder(orderID string) error {
	fmt.Printf("Order Service: Canceling order %s\n", orderID)
	return nil
}

type PaymentService struct{}

func (ps *PaymentService) ProcessPayment(orderID string, amount float64) error {
	fmt.Printf("Payment Service: Processing payment $%.2f for order %s\n", amount, orderID)
	time.Sleep(150 * time.Millisecond)
	
	// Simulate occasional failure
	if orderID == "FAIL_ORDER" {
		return fmt.Errorf("payment processing failed")
	}
	
	return nil
}

func (ps *PaymentService) RefundPayment(orderID string, amount float64) error {
	fmt.Printf("Payment Service: Refunding $%.2f for order %s\n", amount, orderID)
	return nil
}

type InventoryService struct{}

func (is *InventoryService) ReserveItems(items map[string]int) error {
	fmt.Printf("Inventory Service: Reserving items %v\n", items)
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (is *InventoryService) ReleaseItems(items map[string]int) error {
	fmt.Printf("Inventory Service: Releasing items %v\n", items)
	return nil
}

type ShippingService struct{}

func (ss *ShippingService) ScheduleDelivery(orderID string, address string) error {
	fmt.Printf("Shipping Service: Scheduling delivery for order %s to %s\n", orderID, address)
	time.Sleep(200 * time.Millisecond)
	return nil
}

func (ss *ShippingService) CancelDelivery(orderID string) error {
	fmt.Printf("Shipping Service: Canceling delivery for order %s\n", orderID)
	return nil
}

type NotificationService struct{}

func (ns *NotificationService) SendConfirmation(orderID string, email string) error {
	fmt.Printf("Notification Service: Sending confirmation for order %s to %s\n", orderID, email)
	time.Sleep(50 * time.Millisecond)
	return nil
}

func (ns *NotificationService) SendCancellation(orderID string, email string) error {
	fmt.Printf("Notification Service: Sending cancellation for order %s to %s\n", orderID, email)
	return nil
}

// Order Saga using Orchestrator pattern
func CreateOrderSagaOrchestrator(orderID string, amount float64, items map[string]int, address, email string) *OrchestratorSaga {
	orderService := &OrderService{}
	paymentService := &PaymentService{}
	inventoryService := &InventoryService{}
	shippingService := &ShippingService{}
	notificationService := &NotificationService{}
	
	saga := NewOrchestratorSaga(fmt.Sprintf("OrderSaga-%s", orderID))
	
	// Step 1: Create Order
	saga.AddStep(NewBasicSagaStep(
		"CreateOrder",
		func() error {
			return orderService.CreateOrder(orderID)
		},
		func() error {
			return orderService.CancelOrder(orderID)
		},
	))
	
	// Step 2: Process Payment
	saga.AddStep(NewBasicSagaStep(
		"ProcessPayment",
		func() error {
			return paymentService.ProcessPayment(orderID, amount)
		},
		func() error {
			return paymentService.RefundPayment(orderID, amount)
		},
	))
	
	// Step 3: Reserve Inventory
	saga.AddStep(NewBasicSagaStep(
		"ReserveInventory",
		func() error {
			return inventoryService.ReserveItems(items)
		},
		func() error {
			return inventoryService.ReleaseItems(items)
		},
	))
	
	// Step 4: Schedule Delivery
	saga.AddStep(NewBasicSagaStep(
		"ScheduleDelivery",
		func() error {
			return shippingService.ScheduleDelivery(orderID, address)
		},
		func() error {
			return shippingService.CancelDelivery(orderID)
		},
	))
	
	// Step 5: Send Confirmation
	saga.AddStep(NewBasicSagaStep(
		"SendConfirmation",
		func() error {
			return notificationService.SendConfirmation(orderID, email)
		},
		func() error {
			return notificationService.SendCancellation(orderID, email)
		},
	))
	
	// Set callbacks
	saga.SetStepStartCallback(func(step SagaStep) {
		fmt.Printf("Starting step: %s\n", step.GetName())
	})
	
	saga.SetStepCompleteCallback(func(step SagaStep) {
		fmt.Printf("Completed step: %s\n", step.GetName())
	})
	
	saga.SetStepFailedCallback(func(step SagaStep, err error) {
		fmt.Printf("Failed step: %s - %v\n", step.GetName(), err)
	})
	
	saga.SetSagaCompleteCallback(func() {
		fmt.Printf("Saga %s completed successfully!\n", saga.name)
	})
	
	saga.SetSagaFailedCallback(func(err error) {
		fmt.Printf("Saga %s failed: %v\n", saga.name, err)
	})
	
	return saga
}

// Order Saga using Choreography pattern
func CreateOrderSagaChoreography(orderID string, amount float64, items map[string]int, address, email string) *ChoreographySaga {
	orderService := &OrderService{}
	paymentService := &PaymentService{}
	inventoryService := &InventoryService{}
	shippingService := &ShippingService{}
	notificationService := &NotificationService{}
	
	saga := NewChoreographySaga(fmt.Sprintf("OrderSaga-%s", orderID))
	
	// Create steps in order
	saga.AddStep(
		"CreateOrder",
		func() error {
			return orderService.CreateOrder(orderID)
		},
		func() error {
			return orderService.CancelOrder(orderID)
		},
	)
	
	saga.AddStep(
		"ProcessPayment",
		func() error {
			return paymentService.ProcessPayment(orderID, amount)
		},
		func() error {
			return paymentService.RefundPayment(orderID, amount)
		},
	)
	
	saga.AddStep(
		"ReserveInventory",
		func() error {
			return inventoryService.ReserveItems(items)
		},
		func() error {
			return inventoryService.ReleaseItems(items)
		},
	)
	
	saga.AddStep(
		"ScheduleDelivery",
		func() error {
			return shippingService.ScheduleDelivery(orderID, address)
		},
		func() error {
			return shippingService.CancelDelivery(orderID)
		},
	)
	
	saga.AddStep(
		"SendConfirmation",
		func() error {
			return notificationService.SendConfirmation(orderID, email)
		},
		func() error {
			return notificationService.SendCancellation(orderID, email)
		},
	)
	
	return saga
}

// Saga Manager
type SagaManager struct {
	sagas map[string]Saga
	mu    sync.RWMutex
}

func NewSagaManager() *SagaManager {
	return &SagaManager{
		sagas: make(map[string]Saga),
	}
}

func (sm *SagaManager) ExecuteSaga(saga Saga) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	sagaID := fmt.Sprintf("saga_%d", time.Now().UnixNano())
	sm.sagas[sagaID] = saga
	
	return saga.Execute()
}

func (sm *SagaManager) GetSagaStatus(sagaID string) SagaStatus {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	if saga, exists := sm.sagas[sagaID]; exists {
		return saga.GetStatus()
	}
	return Pending
}

func (sm *SagaManager) GetAllStatus() map[string]SagaStatus {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	status := make(map[string]SagaStatus)
	for id, saga := range sm.sagas {
		status[id] = saga.GetStatus()
	}
	return status
}

func demonstrateOrchestratorSaga() {
	fmt.Println("--- Orchestrator Saga Demo ---")
	
	orderID := "ORDER-001"
	amount := 299.99
	items := map[string]int{"laptop": 1, "mouse": 1}
	address := "123 Main St, City, State"
	email := "customer@example.com"
	
	saga := CreateOrderSagaOrchestrator(orderID, amount, items, address, email)
	
	fmt.Printf("Executing order saga for %s\n", orderID)
	err := saga.Execute()
	
	if err != nil {
		fmt.Printf("Saga execution failed: %v\n", err)
	} else {
		fmt.Printf("Saga execution succeeded\n")
	}
	
	// Wait for compensation if needed
	time.Sleep(1 * time.Second)
	fmt.Printf("Final saga status: %v\n", saga.GetStatus())
}

func demonstrateChoreographySaga() {
	fmt.Println("\n--- Choreography Saga Demo ---")
	
	orderID := "ORDER-002"
	amount := 199.99
	items := map[string]int{"phone": 1}
	address := "456 Oak Ave, City, State"
	email := "user@example.com"
	
	saga := CreateOrderSagaChoreography(orderID, amount, items, address, email)
	
	fmt.Printf("Executing order saga for %s\n", orderID)
	err := saga.Execute()
	
	if err != nil {
		fmt.Printf("Saga execution failed: %v\n", err)
	} else {
		fmt.Printf("Saga execution succeeded\n")
	}
	
	// Wait for compensation if needed
	time.Sleep(1 * time.Second)
	fmt.Printf("Final saga status: %v\n", saga.GetStatus())
}

func demonstrateFailedSaga() {
	fmt.Println("\n--- Failed Saga Demo (with compensation) ---")
	
	orderID := "FAIL_ORDER" // This will trigger payment failure
	amount := 399.99
	items := map[string]int{"tablet": 1}
	address := "789 Pine Rd, City, State"
	email := "admin@example.com"
	
	saga := CreateOrderSagaOrchestrator(orderID, amount, items, address, email)
	
	fmt.Printf("Executing order saga for %s (will fail)\n", orderID)
	err := saga.Execute()
	
	if err != nil {
		fmt.Printf("Saga execution failed as expected: %v\n", err)
	} else {
		fmt.Printf("Saga execution succeeded unexpectedly\n")
	}
	
	// Wait for compensation
	time.Sleep(2 * time.Second)
	fmt.Printf("Final saga status: %v\n", saga.GetStatus())
	
	// Show step statuses
	fmt.Println("\nStep statuses:")
	for _, step := range saga.GetSteps() {
		fmt.Printf("  %s: %v\n", step.GetName(), step.GetStatus())
	}
}

func demonstrateSagaManager() {
	fmt.Println("\n--- Saga Manager Demo ---")
	
	manager := NewSagaManager()
	
	// Create multiple sagas
	saga1 := CreateOrderSagaOrchestrator("ORDER-003", 99.99, map[string]int{"book": 2}, "111 First St", "user1@example.com")
	saga2 := CreateOrderSagaOrchestrator("ORDER-004", 199.99, map[string]int{"headphones": 1}, "222 Second Ave", "user2@example.com")
	saga3 := CreateOrderSagaOrchestrator("FAIL_ORDER2", 299.99, map[string]int{"camera": 1}, "333 Third Blvd", "user3@example.com")
	
	// Execute sagas concurrently
	var wg sync.WaitGroup
	
	sagas := []Saga{saga1, saga2, saga3}
	for i, saga := range sagas {
		wg.Add(1)
		go func(index int, s Saga) {
			defer wg.Done()
			fmt.Printf("Starting saga %d\n", index+1)
			err := manager.ExecuteSaga(s)
			if err != nil {
				fmt.Printf("Saga %d failed: %v\n", index+1, err)
			} else {
				fmt.Printf("Saga %d succeeded\n", index+1)
			}
		}(i, saga)
	}
	
	wg.Wait()
	
	// Wait for any compensations
	time.Sleep(2 * time.Second)
	
	// Show all statuses
	fmt.Println("\nAll saga statuses:")
	for id, status := range manager.GetAllStatus() {
		fmt.Printf("  %s: %v\n", id, status)
	}
}

func main() {
	fmt.Println("=== Saga Pattern Demo ===")
	
	demonstrateOrchestratorSaga()
	demonstrateChoreographySaga()
	demonstrateFailedSaga()
	demonstrateSagaManager()
	
	fmt.Println("\nAll saga patterns demonstrated successfully!")
}
