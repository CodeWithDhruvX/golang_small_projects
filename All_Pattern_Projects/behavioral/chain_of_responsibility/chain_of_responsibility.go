package main

import "fmt"

// Handler defines the interface for handling requests
type Handler interface {
	SetNext(handler Handler)
	Handle(request string)
}

// BaseHandler provides default implementation for the handler chain
type BaseHandler struct {
	next Handler
}

func (h *BaseHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *BaseHandler) Handle(request string) {
	if h.next != nil {
		h.next.Handle(request)
	}
}

// ConcreteHandler1 handles specific types of requests
type ConcreteHandler1 struct {
	BaseHandler
}

func (h *ConcreteHandler1) Handle(request string) {
	if request == "request1" {
		fmt.Println("ConcreteHandler1: Handling request1")
		return
	}
	fmt.Println("ConcreteHandler1: Passing to next handler")
	h.BaseHandler.Handle(request)
}

// ConcreteHandler2 handles another type of request
type ConcreteHandler2 struct {
	BaseHandler
}

func (h *ConcreteHandler2) Handle(request string) {
	if request == "request2" {
		fmt.Println("ConcreteHandler2: Handling request2")
		return
	}
	fmt.Println("ConcreteHandler2: Passing to next handler")
	h.BaseHandler.Handle(request)
}

// ConcreteHandler3 handles all remaining requests
type ConcreteHandler3 struct {
	BaseHandler
}

func (h *ConcreteHandler3) Handle(request string) {
	fmt.Printf("ConcreteHandler3: Handling default request: %s\n", request)
}

func main() {
	// Create handlers
	handler1 := &ConcreteHandler1{}
	handler2 := &ConcreteHandler2{}
	handler3 := &ConcreteHandler3{}

	// Set up the chain of responsibility
	handler1.SetNext(handler2)
	handler2.SetNext(handler3)

	// Process requests
	fmt.Println("=== Chain of Responsibility Pattern Demo ===")
	
	requests := []string{"request1", "request2", "request3", "unknown"}
	
	for _, request := range requests {
		fmt.Printf("\nProcessing: %s\n", request)
		handler1.Handle(request)
	}
}
