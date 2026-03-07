package main

import "fmt"

// State interface defines the behavior for each state
type State interface {
	InsertCoin(vendingMachine *VendingMachine)
	EjectCoin(vendingMachine *VendingMachine)
	PressButton(vendingMachine *VendingMachine)
	Dispense(vendingMachine *VendingMachine)
}

// Context maintains the current state and delegates state-specific behavior
type VendingMachine struct {
	hasCoin       bool
	itemDispensed bool
	currentState  State
	noCoinState   State
	hasCoinState  State
	soldState     State
}

func NewVendingMachine() *VendingMachine {
	vm := &VendingMachine{}
	
	noCoinState := &NoCoinState{vendingMachine: vm}
	hasCoinState := &HasCoinState{vendingMachine: vm}
	soldState := &SoldState{vendingMachine: vm}
	
	vm.noCoinState = noCoinState
	vm.hasCoinState = hasCoinState
	vm.soldState = soldState
	vm.currentState = noCoinState
	
	return vm
}

func (vm *VendingMachine) InsertCoin() {
	vm.currentState.InsertCoin(vm)
}

func (vm *VendingMachine) EjectCoin() {
	vm.currentState.EjectCoin(vm)
}

func (vm *VendingMachine) PressButton() {
	vm.currentState.PressButton(vm)
}

func (vm *VendingMachine) Dispense() {
	vm.currentState.Dispense(vm)
}

func (vm *VendingMachine) SetState(state State) {
	vm.currentState = state
}

func (vm *VendingMachine) GetNoCoinState() State    { return vm.noCoinState }
func (vm *VendingMachine) GetHasCoinState() State   { return vm.hasCoinState }
func (vm *VendingMachine) GetSoldState() State      { return vm.soldState }

// ConcreteState1 - No coin inserted
type NoCoinState struct {
	vendingMachine *VendingMachine
}

func (ncs *NoCoinState) InsertCoin(vm *VendingMachine) {
	fmt.Println("Coin inserted. You can press the button now.")
	vm.SetState(vm.GetHasCoinState())
}

func (ncs *NoCoinState) EjectCoin(vm *VendingMachine) {
	fmt.Println("You haven't inserted a coin.")
}

func (ncs *NoCoinState) PressButton(vm *VendingMachine) {
	fmt.Println("You pressed the button, but you haven't inserted a coin.")
}

func (ncs *NoCoinState) Dispense(vm *VendingMachine) {
	fmt.Println("You need to pay first.")
}

// ConcreteState2 - Coin inserted
type HasCoinState struct {
	vendingMachine *VendingMachine
}

func (hcs *HasCoinState) InsertCoin(vm *VendingMachine) {
	fmt.Println("You can't insert another coin.")
}

func (hcs *HasCoinState) EjectCoin(vm *VendingMachine) {
	fmt.Println("Coin returned.")
	vm.SetState(vm.GetNoCoinState())
}

func (hcs *HasCoinState) PressButton(vm *VendingMachine) {
	fmt.Println("You pressed the button...")
	vm.SetState(vm.GetSoldState())
}

func (hcs *HasCoinState) Dispense(vm *VendingMachine) {
	fmt.Println("No item dispensed.")
}

// ConcreteState3 - Item sold
type SoldState struct {
	vendingMachine *VendingMachine
}

func (ss *SoldState) InsertCoin(vm *VendingMachine) {
	fmt.Println("Please wait, we're already giving you an item.")
}

func (ss *SoldState) EjectCoin(vm *VendingMachine) {
	fmt.Println("Sorry, you already pressed the button.")
}

func (ss *SoldState) PressButton(vm *VendingMachine) {
	fmt.Println("Pressing the button twice doesn't get you another item!")
}

func (ss *SoldState) Dispense(vm *VendingMachine) {
	fmt.Println("Item dispensed. Thank you for your purchase!")
	vm.SetState(vm.GetNoCoinState())
}

// Traffic Light example
type TrafficLightState interface {
	Change(trafficLight *TrafficLight)
	GetColor() string
}

type TrafficLight struct {
	currentState TrafficLightState
	redState     TrafficLightState
	greenState   TrafficLightState
	yellowState  TrafficLightState
}

func NewTrafficLight() *TrafficLight {
	tl := &TrafficLight{}
	
	redState := &RedState{trafficLight: tl}
	yellowState := &YellowState{trafficLight: tl}
	greenState := &GreenState{trafficLight: tl}
	
	tl.redState = redState
	tl.yellowState = yellowState
	tl.greenState = greenState
	tl.currentState = redState
	
	return tl
}

func (tl *TrafficLight) Change() {
	tl.currentState.Change(tl)
}

func (tl *TrafficLight) SetState(state TrafficLightState) {
	tl.currentState = state
}

func (tl *TrafficLight) GetColor() string {
	return tl.currentState.GetColor()
}

func (tl *TrafficLight) GetRedState() TrafficLightState    { return tl.redState }
func (tl *TrafficLight) GetYellowState() TrafficLightState { return tl.yellowState }
func (tl *TrafficLight) GetGreenState() TrafficLightState  { return tl.greenState }

type RedState struct {
	trafficLight *TrafficLight
}

func (rs *RedState) Change(tl *TrafficLight) {
	fmt.Println("Traffic light changing from Red to Green")
	tl.SetState(tl.GetGreenState())
}

func (rs *RedState) GetColor() string {
	return "RED"
}

type GreenState struct {
	trafficLight *TrafficLight
}

func (gs *GreenState) Change(tl *TrafficLight) {
	fmt.Println("Traffic light changing from Green to Yellow")
	tl.SetState(tl.GetYellowState())
}

func (gs *GreenState) GetColor() string {
	return "GREEN"
}

type YellowState struct {
	trafficLight *TrafficLight
}

func (ys *YellowState) Change(tl *TrafficLight) {
	fmt.Println("Traffic light changing from Yellow to Red")
	tl.SetState(tl.GetRedState())
}

func (ys *YellowState) GetColor() string {
	return "YELLOW"
}

func main() {
	fmt.Println("=== State Pattern Demo ===")
	
	// Vending Machine example
	fmt.Println("\n--- Vending Machine Example ---")
	vendingMachine := NewVendingMachine()
	
	fmt.Println("Current state: No coin")
	vendingMachine.InsertCoin()
	
	fmt.Println("\nTrying to insert another coin:")
	vendingMachine.InsertCoin()
	
	fmt.Println("\nPressing button:")
	vendingMachine.PressButton()
	
	fmt.Println("\nTrying to eject coin after pressing button:")
	vendingMachine.EjectCoin()
	
	fmt.Println("\nPressing button again:")
	vendingMachine.PressButton()
	
	fmt.Println("\n--- Traffic Light Example ---")
	trafficLight := NewTrafficLight()
	
	for i := 0; i < 6; i++ {
		fmt.Printf("Current light color: %s\n", trafficLight.GetColor())
		trafficLight.Change()
	}
}
