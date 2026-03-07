package main

import "fmt"

// Mediator interface defines the communication protocol
type Mediator interface {
	Send(message string, colleague Colleague)
	Register(colleague Colleague)
}

// Colleague interface for participants in the mediation
type Colleague interface {
	Receive(message string)
	SetMediator(mediator Mediator)
	GetName() string
}

// ConcreteMediator implements the mediation logic
type ChatRoom struct {
	colleagues map[string]Colleague
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		colleagues: make(map[string]Colleague),
	}
}

func (cr *ChatRoom) Register(colleague Colleague) {
	cr.colleagues[colleague.GetName()] = colleague
	colleague.SetMediator(cr)
	fmt.Printf("%s joined the chat room\n", colleague.GetName())
}

func (cr *ChatRoom) Send(message string, sender Colleague) {
	for name, colleague := range cr.colleagues {
		if name != sender.GetName() {
			colleague.Receive(message)
		}
	}
}

// ConcreteColleague represents a user in the chat room
type User struct {
	name     string
	mediator Mediator
}

func NewUser(name string) *User {
	return &User{name: name}
}

func (u *User) SetMediator(mediator Mediator) {
	u.mediator = mediator
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) Receive(message string) {
	fmt.Printf("[%s] received: %s\n", u.name, message)
}

func (u *User) Send(message string) {
	fmt.Printf("[%s] sending: %s\n", u.name, message)
	u.mediator.Send(message, u)
}

// Another example with different types of colleagues
type Airplane interface {
	Colleague
	RequestLanding()
	RequestTakeoff()
	NotifyClearance()
}

type AirTrafficControlTower struct {
	mediator Mediator
	airplanes map[string]Airplane
}

func NewAirTrafficControlTower() *AirTrafficControlTower {
	return &AirTrafficControlTower{
		airplanes: make(map[string]Airplane),
	}
}

func (atc *AirTrafficControlTower) Register(colleague Colleague) {
	if airplane, ok := colleague.(Airplane); ok {
		atc.airplanes[airplane.GetName()] = airplane
		airplane.SetMediator(atc)
		fmt.Printf("Airplane %s registered with ATC\n", airplane.GetName())
	}
}

func (atc *AirTrafficControlTower) Send(message string, sender Colleague) {
	for name, airplane := range atc.airplanes {
		if name != sender.GetName() {
			airplane.Receive(message)
		}
	}
}

type CommercialAirplane struct {
	name     string
	mediator Mediator
}

func NewCommercialAirplane(name string) *CommercialAirplane {
	return &CommercialAirplane{name: name}
}

func (ca *CommercialAirplane) SetMediator(mediator Mediator) {
	ca.mediator = mediator
}

func (ca *CommercialAirplane) GetName() string {
	return ca.name
}

func (ca *CommercialAirplane) Receive(message string) {
	fmt.Printf("[Airplane %s] received: %s\n", ca.name, message)
}

func (ca *CommercialAirplane) RequestLanding() {
	fmt.Printf("[Airplane %s] requesting landing clearance\n", ca.name)
	ca.mediator.Send("Requesting landing clearance", ca)
}

func (ca *CommercialAirplane) RequestTakeoff() {
	fmt.Printf("[Airplane %s] requesting takeoff clearance\n", ca.name)
	ca.mediator.Send("Requesting takeoff clearance", ca)
}

func (ca *CommercialAirplane) NotifyClearance() {
	fmt.Printf("[Airplane %s] clearance received\n", ca.name)
}

func main() {
	fmt.Println("=== Mediator Pattern Demo ===")
	
	// Chat room example
	fmt.Println("\n--- Chat Room Example ---")
	chatRoom := NewChatRoom()
	
	alice := NewUser("Alice")
	bob := NewUser("Bob")
	charlie := NewUser("Charlie")
	
	chatRoom.Register(alice)
	chatRoom.Register(bob)
	chatRoom.Register(charlie)
	
	alice.Send("Hello everyone!")
	bob.Send("Hi Alice!")
	
	// Air traffic control example
	fmt.Println("\n--- Air Traffic Control Example ---")
	atc := NewAirTrafficControlTower()
	
	flight101 := NewCommercialAirplane("Flight101")
	flight202 := NewCommercialAirplane("Flight202")
	
	atc.Register(flight101)
	atc.Register(flight202)
	
	flight101.RequestLanding()
	flight202.RequestTakeoff()
}
