package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// API Gateway Pattern

// Service interface
type Service interface {
	HandleRequest(request *Request) *Response
	GetName() string
}

// Request structure
type Request struct {
	Path    string
	Method  string
	Headers map[string]string
	Body    string
	Params  map[string]string
}

// Response structure
type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

// Route configuration
type Route struct {
	Path       string
	Method     string
	Service    string
	Auth       bool
	RateLimit  int
}

// Microservice implementations
type UserService struct{}

func (us *UserService) HandleRequest(request *Request) *Response {
	switch request.Path {
	case "/users":
		return &Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"users": [{"id": 1, "name": "John"}, {"id": 2, "name": "Jane"}]}`,
		}
	case "/users/" + request.Params["id"]:
		return &Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       fmt.Sprintf(`{"user": {"id": %s, "name": "User %s"}}`, request.Params["id"], request.Params["id"]),
		}
	default:
		return &Response{
			StatusCode: 404,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "User not found"}`,
		}
	}
}

func (us *UserService) GetName() string {
	return "UserService"
}

type OrderService struct{}

func (os *OrderService) HandleRequest(request *Request) *Response {
	switch request.Path {
	case "/orders":
		return &Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"orders": [{"id": 1, "product": "Laptop"}, {"id": 2, "product": "Phone"}]}`,
		}
	case "/orders/" + request.Params["id"]:
		return &Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       fmt.Sprintf(`{"order": {"id": %s, "product": "Product %s"}}`, request.Params["id"], request.Params["id"]),
		}
	default:
		return &Response{
			StatusCode: 404,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "Order not found"}`,
		}
	}
}

func (os *OrderService) GetName() string {
	return "OrderService"
}

type ProductService struct{}

func (ps *ProductService) HandleRequest(request *Request) *Response {
	switch request.Path {
	case "/products":
		return &Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"products": [{"id": 1, "name": "Laptop", "price": 999.99}, {"id": 2, "name": "Phone", "price": 699.99}]}`,
		}
	case "/products/" + request.Params["id"]:
		return &Response{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       fmt.Sprintf(`{"product": {"id": %s, "name": "Product %s", "price": 99.99}}`, request.Params["id"], request.Params["id"]),
		}
	default:
		return &Response{
			StatusCode: 404,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "Product not found"}`,
		}
	}
}

func (ps *ProductService) GetName() string {
	return "ProductService"
}

// Authentication Service
type AuthService struct {
	validTokens map[string]string
}

func NewAuthService() *AuthService {
	return &AuthService{
		validTokens: map[string]string{
			"token123": "user1",
			"token456": "user2",
		},
	}
}

func (as *AuthService) ValidateToken(token string) bool {
	_, exists := as.validTokens[token]
	return exists
}

func (as *AuthService) GetUser(token string) string {
	return as.validTokens[token]
}

// Rate Limiter
type RateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(clientID string) bool {
	now := time.Now()
	
	// Clean old requests
	if requests, exists := rl.requests[clientID]; exists {
		validRequests := make([]time.Time, 0)
		for _, req := range requests {
			if now.Sub(req) <= rl.window {
				validRequests = append(validRequests, req)
			}
		}
		rl.requests[clientID] = validRequests
	}
	
	// Check limit
	if len(rl.requests[clientID]) >= rl.limit {
		return false
	}
	
	// Add new request
	rl.requests[clientID] = append(rl.requests[clientID], now)
	return true
}

// Load Balancer
type LoadBalancer struct {
	services []Service
	current  int
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		services: make([]Service, 0),
		current:  0,
	}
}

func (lb *LoadBalancer) AddService(service Service) {
	lb.services = append(lb.services, service)
}

func (lb *LoadBalancer) GetNextService() Service {
	if len(lb.services) == 0 {
		return nil
	}
	
	service := lb.services[lb.current]
	lb.current = (lb.current + 1) % len(lb.services)
	return service
}

// API Gateway
type APIGateway struct {
	services    map[string]Service
	routes      []Route
	authService *AuthService
	rateLimiter *RateBalancer
	loadBalancer *LoadBalancer
	middleware  []Middleware
}

type Middleware interface {
	Process(request *Request, response *Response, next func()) bool
}

// Logging Middleware
type LoggingMiddleware struct{}

func (lm *LoggingMiddleware) Process(request *Request, response *Response, next func()) bool {
	fmt.Printf("Gateway: %s %s\n", request.Method, request.Path)
	next()
	fmt.Printf("Gateway: Response %d\n", response.StatusCode)
	return true
}

// CORS Middleware
type CORSMiddleware struct{}

func (cm *CORSMiddleware) Process(request *Request, response *Response, next func()) bool {
	response.Headers["Access-Control-Allow-Origin"] = "*"
	response.Headers["Access-Control-Allow-Methods"] = "GET, POST, PUT, DELETE"
	response.Headers["Access-Control-Allow-Headers"] = "Content-Type, Authorization"
	next()
	return true
}

// Metrics Middleware
type MetricsMiddleware struct {
	requestCount map[string]int
}

func NewMetricsMiddleware() *MetricsMiddleware {
	return &MetricsMiddleware{
		requestCount: make(map[string]int),
	}
}

func (mm *MetricsMiddleware) Process(request *Request, response *Response, next func()) bool {
	key := fmt.Sprintf("%s %s", request.Method, request.Path)
	mm.requestCount[key]++
	next()
	fmt.Printf("Gateway: Request count for %s: %d\n", key, mm.requestCount[key])
	return true
}

func NewAPIGateway() *APIGateway {
	gateway := &APIGateway{
		services:     make(map[string]Service),
		routes:       make([]Route, 0),
		authService:  NewAuthService(),
		rateLimiter:  NewRateLimiter(100, time.Minute),
		loadBalancer: NewLoadBalancer(),
		middleware:   make([]Middleware, 0),
	}
	
	// Add middleware
	gateway.AddMiddleware(&LoggingMiddleware{})
	gateway.AddMiddleware(&CORSMiddleware{})
	gateway.AddMiddleware(NewMetricsMiddleware())
	
	return gateway
}

func (ag *APIGateway) AddService(name string, service Service) {
	ag.services[name] = service
	ag.loadBalancer.AddService(service)
}

func (ag *APIGateway) AddRoute(path, method, service string, auth bool, rateLimit int) {
	ag.routes = append(ag.routes, Route{
		Path:      path,
		Method:    method,
		Service:   service,
		Auth:      auth,
		RateLimit: rateLimit,
	})
}

func (ag *APIGateway) AddMiddleware(middleware Middleware) {
	ag.middleware = append(ag.middleware, middleware)
}

func (ag *APIGateway) HandleRequest(request *Request) *Response {
	// Find matching route
	route := ag.findRoute(request.Path, request.Method)
	if route == nil {
		return &Response{
			StatusCode: 404,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "Route not found"}`,
		}
	}
	
	// Check authentication
	if route.Auth {
		token := request.Headers["Authorization"]
		if !ag.authService.ValidateToken(token) {
			return &Response{
				StatusCode: 401,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       `{"error": "Unauthorized"}`,
			}
		}
	}
	
	// Check rate limiting
	clientID := request.Headers["X-Client-ID"]
	if clientID == "" {
		clientID = "anonymous"
	}
	
	if !ag.rateLimiter.Allow(clientID) {
		return &Response{
			StatusCode: 429,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "Rate limit exceeded"}`,
		}
	}
	
	// Get service
	service, exists := ag.services[route.Service]
	if !exists {
		return &Response{
			StatusCode: 503,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "Service unavailable"}`,
		}
	}
	
	// Process request through middleware chain
	response := &Response{StatusCode: 200, Headers: make(map[string]string)}
	
	// Execute middleware chain
	for _, middleware := range ag.middleware {
		continueProcessing := middleware.Process(request, response, func() {
			// This would normally call the next middleware or service
		})
		if !continueProcessing {
			return response
		}
	}
	
	// Forward to service
	response = service.HandleRequest(request)
	
	return response
}

func (ag *APIGateway) findRoute(path, method string) *Route {
	for _, route := range ag.routes {
		if ag.pathMatches(route.Path, path) && route.Method == method {
			return &route
		}
	}
	return nil
}

func (ag *APIGateway) pathMatches(routePath, requestPath string) bool {
	// Simple path matching (in real implementation, use regex)
	if routePath == requestPath {
		return true
	}
	
	// Handle parameterized paths like "/users/{id}"
	routeParts := strings.Split(routePath, "/")
	requestParts := strings.Split(requestPath, "/")
	
	if len(routeParts) != len(requestParts) {
		return false
	}
	
	for i, routePart := range routeParts {
		if strings.HasPrefix(routePart, "{") && strings.HasSuffix(routePart, "}") {
			continue // Parameter
		}
		if routePart != requestParts[i] {
			return false
		}
	}
	
	return true
}

func (ag *APIGateway) extractParams(routePath, requestPath string) map[string]string {
	params := make(map[string]string)
	
	routeParts := strings.Split(routePath, "/")
	requestParts := strings.Split(requestPath, "/")
	
	for i, routePart := range routeParts {
		if strings.HasPrefix(routePart, "{") && strings.HasSuffix(routePart, "}") {
			paramName := routePart[1 : len(routePart)-1]
			if i < len(requestParts) {
				params[paramName] = requestParts[i]
			}
		}
	}
	
	return params
}

// Request Aggregator
type RequestAggregator struct {
	gateway *APIGateway
}

func NewRequestAggregator(gateway *APIGateway) *RequestAggregator {
	return &RequestAggregator{gateway: gateway}
}

func (ra *RequestAggregator) AggregateUserOrders(userID string) *Response {
	// Get user info
	userRequest := &Request{
		Path:    "/users/" + userID,
		Method:  "GET",
		Headers: map[string]string{"Authorization": "token123"},
		Params:  map[string]string{"id": userID},
	}
	
	userResponse := ra.gateway.HandleRequest(userRequest)
	
	// Get orders
	orderRequest := &Request{
		Path:    "/orders",
		Method:  "GET",
		Headers: map[string]string{"Authorization": "token123"},
		Params:  map[string]string{},
	}
	
	orderResponse := ra.gateway.HandleRequest(orderRequest)
	
	// Combine responses
	aggregatedBody := fmt.Sprintf(`{
		"user": %s,
		"orders": %s
	}`, userResponse.Body, orderResponse.Body)
	
	return &Response{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       aggregatedBody,
	}
}

// Service Discovery
type ServiceDiscovery struct {
	services map[string][]string // service -> endpoints
}

func NewServiceDiscovery() *ServiceDiscovery {
	return &ServiceDiscovery{
		services: make(map[string][]string),
	}
}

func (sd *ServiceDiscovery) RegisterService(serviceName, endpoint string) {
	if sd.services[serviceName] == nil {
		sd.services[serviceName] = make([]string, 0)
	}
	sd.services[serviceName] = append(sd.services[serviceName], endpoint)
	fmt.Printf("Service Discovery: Registered %s at %s\n", serviceName, endpoint)
}

func (sd *ServiceDiscovery) GetEndpoints(serviceName string) []string {
	return sd.services[serviceName]
}

func (sd *ServiceDiscovery) HealthCheck(serviceName string) {
	endpoints := sd.GetEndpoints(serviceName)
	fmt.Printf("Service Discovery: Health checking %s at %v\n", serviceName, endpoints)
}

func main() {
	fmt.Println("=== API Gateway Pattern Demo ===")
	
	// Create API Gateway
	gateway := NewAPIGateway()
	
	// Register services
	gateway.AddService("UserService", &UserService{})
	gateway.AddService("OrderService", &OrderService{})
	gateway.AddService("ProductService", &ProductService{})
	
	// Add routes
	gateway.AddRoute("/users", "GET", "UserService", false, 100)
	gateway.AddRoute("/users/{id}", "GET", "UserService", false, 50)
	gateway.AddRoute("/orders", "GET", "OrderService", true, 100)
	gateway.AddRoute("/orders/{id}", "GET", "OrderService", true, 50)
	gateway.AddRoute("/products", "GET", "ProductService", false, 100)
	gateway.AddRoute("/products/{id}", "GET", "ProductService", false, 50)
	
	// Service Discovery
	serviceDiscovery := NewServiceDiscovery()
	serviceDiscovery.RegisterService("UserService", "http://localhost:8001")
	serviceDiscovery.RegisterService("OrderService", "http://localhost:8002")
	serviceDiscovery.RegisterService("ProductService", "http://localhost:8003")
	
	// Test requests
	fmt.Println("\n--- Testing API Gateway ---")
	
	// Public request (no auth required)
	fmt.Println("\n1. Public request to /users:")
	request1 := &Request{
		Path:   "/users",
		Method: "GET",
		Headers: map[string]string{"X-Client-ID": "client1"},
		Params: map[string]string{},
	}
	response1 := gateway.HandleRequest(request1)
	fmt.Printf("Status: %d, Body: %s\n", response1.StatusCode, response1.Body)
	
	// Authenticated request
	fmt.Println("\n2. Authenticated request to /orders:")
	request2 := &Request{
		Path:   "/orders",
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "token123",
			"X-Client-ID":   "client1",
		},
		Params: map[string]string{},
	}
	response2 := gateway.HandleRequest(request2)
	fmt.Printf("Status: %d, Body: %s\n", response2.StatusCode, response2.Body)
	
	// Parameterized request
	fmt.Println("\n3. Parameterized request to /users/123:")
	request3 := &Request{
		Path:   "/users/123",
		Method: "GET",
		Headers: map[string]string{"X-Client-ID": "client1"},
		Params:  gateway.extractParams("/users/{id}", "/users/123"),
	}
	response3 := gateway.HandleRequest(request3)
	fmt.Printf("Status: %d, Body: %s\n", response3.StatusCode, response3.Body)
	
	// Unauthorized request
	fmt.Println("\n4. Unauthorized request to /orders:")
	request4 := &Request{
		Path:   "/orders",
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "invalid_token",
			"X-Client-ID":   "client2",
		},
		Params: map[string]string{},
	}
	response4 := gateway.HandleRequest(request4)
	fmt.Printf("Status: %d, Body: %s\n", response4.StatusCode, response4.Body)
	
	// Request Aggregation
	fmt.Println("\n--- Request Aggregation ---")
	aggregator := NewRequestAggregator(gateway)
	aggregatedResponse := aggregator.AggregateUserOrders("123")
	fmt.Printf("Aggregated Response: %s\n", aggregatedResponse.Body)
	
	// Service Discovery Health Check
	fmt.Println("\n--- Service Discovery ---")
	serviceDiscovery.HealthCheck("UserService")
	serviceDiscovery.HealthCheck("OrderService")
	serviceDiscovery.HealthCheck("ProductService")
	
	fmt.Println("\nAll API Gateway patterns demonstrated successfully!")
}
