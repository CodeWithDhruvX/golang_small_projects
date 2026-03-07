# Go Design Patterns

A comprehensive collection of design patterns implemented in Go, organized into five main categories with 32 different patterns total.

## 📋 Table of Contents

- [Overview](#overview)
- [Categories](#categories)
- [Installation](#installation)
- [Usage](#usage)
- [Pattern Categories](#pattern-categories)
  - [Behavioral Patterns (8)](#behavioral-patterns)
  - [Concurrency Patterns (6)](#concurrency-patterns)
  - [Creational Patterns (5)](#creational-patterns)
  - [Structural Patterns (7)](#structural-patterns)
  - [Microservices Patterns (6)](#microservices-patterns)
- [Project Structure](#project-structure)
- [Running Individual Patterns](#running-individual-patterns)
- [Contributing](#contributing)
- [License](#license)

## 🎯 Overview

This project provides practical implementations of classic design patterns in Go. Each pattern includes:

- **Complete working code** with proper Go idioms
- **Multiple examples** demonstrating different use cases
- **Detailed explanations** within the code comments
- **Interactive demonstrations** showing the pattern in action

The patterns are organized into logical categories and can be run individually or as a complete collection.

## 📚 Categories

| Category | Patterns | Description |
|-----------|-----------|-------------|
| **Behavioral** | 8 | Patterns that deal with communication between objects and assignment of responsibilities |
| **Concurrency** | 6 | Patterns specifically for handling concurrent programming in Go |
| **Creational** | 5 | Patterns that deal with object creation mechanisms |
| **Structural** | 7 | Patterns that deal with object composition and relationships |
| **Microservices** | 6 | Patterns for building distributed systems and microservices |

**Total: 32 design patterns**

## 🚀 Installation

### Prerequisites

- Go 1.19 or later
- Git (optional, for cloning)

### Setup

1. Clone or download the project:
```bash
git clone <repository-url>
cd golang_small_projects
```

2. Install dependencies:
```bash
go mod tidy
```

3. Verify installation:
```bash
go run main.go
```

## 🎮 Usage

### Interactive Menu

Run the main application to access an interactive menu:

```bash
go run main.go
```

This will show you:
- List of all available patterns
- Options to run specific patterns
- Category-based execution
- Statistics and information

### Running Individual Patterns

Each pattern can be run independently:

```bash
# Run a specific pattern
go run chain_of_responsibility/chain_of_responsibility.go

# Run with package path
go run ./chain_of_responsibility

# From any directory
go run ./chain_of_responsibility/chain_of_responsibility.go
```

### Running All Patterns

You can run all patterns in a category or all patterns at once:

```bash
# Run all behavioral patterns
go run main.go  # Choose option 3, then "behavioral"

# Run all patterns
go run main.go  # Choose option 4
```

## 📂 Pattern Categories

### Behavioral Patterns (8)

Patterns that characterize how objects interact and distribute responsibility.

#### 1. Chain of Responsibility
**File**: `chain_of_responsibility/chain_of_responsibility.go`

Passes a request along a chain of handlers. Each handler decides either to process the request or to pass it to the next handler in the chain.

**Key Concepts**:
- Decouples sender from receiver
- Dynamic chain configuration
- Flexible object composition

#### 2. Command
**File**: `command/command.go`

Encapsulates a request as an object, thereby letting you parameterize clients with different requests, queue or log requests, and support undoable operations.

**Key Concepts**:
- Request encapsulation
- Parameterization of objects
- Undo/redo functionality

#### 3. Mediator
**File**: `mediator/mediator.go`

Defines an object that centralizes communications between a set of objects. Prevents objects from referring to each other explicitly, promoting loose coupling.

**Key Concepts**:
- Centralized communication hub
- Loose coupling between components
- Simplified object interactions

#### 4. Memento
**File**: `memento/memento.go`

Captures and externalizes an object's internal state so that the object can be restored to this state later, without violating encapsulation.

**Key Concepts**:
- State preservation
- Undo/redo mechanisms
- Encapsulation maintenance

#### 5. Observer
**File**: `observer/observer.go`

Defines a one-to-many dependency between objects so that when one object changes state, all its dependents are notified and updated automatically.

**Key Concepts**:
- Publisher-subscriber model
- Event-driven architecture
- Loose coupling

#### 6. State
**File**: `state/state.go`

Allows an object to alter its behavior when its internal state changes. The object will appear to change its class.

**Key Concepts**:
- State-based behavior
- Dynamic behavior switching
- Clean state transitions

#### 7. Strategy
**File**: `strategy/strategy.go`

Defines a family of algorithms, encapsulates each one, and makes them interchangeable. Strategy lets the algorithm vary independently from clients that use it.

**Key Concepts**:
- Algorithm encapsulation
- Runtime algorithm switching
- Policy-based design

#### 8. Template Method
**File**: `template_method/template_method.go`

Defines the skeleton of an algorithm in an operation, deferring some steps to subclasses. Template Method lets subclasses redefine certain steps of an algorithm without changing the algorithm's structure.

**Key Concepts**:
- Algorithm skeleton
- Step delegation
- Inversion of control

### Concurrency Patterns (6)

Patterns specifically designed for handling concurrent programming challenges in Go.

#### 1. Barrier
**File**: `barrier/barrier.go`

A synchronization primitive that blocks until a certain number of threads have reached the barrier, then releases all threads simultaneously.

**Key Concepts**:
- Thread synchronization
- Coordination points
- Collective operations

#### 2. Fan-in/Fan-out
**File**: `fan_in_fan_out/fan_in_fan_out.go`

Concurrency pattern where multiple goroutines (fan-out) process work in parallel, and their results are collected and combined (fan-in).

**Key Concepts**:
- Parallel processing
- Result aggregation
- Work distribution

#### 3. Generator
**File**: `generator/generator.go`

A function that returns a channel that produces a sequence of values, allowing for lazy evaluation and memory-efficient iteration.

**Key Concepts**:
- Lazy evaluation
- Stream processing
- Memory efficiency

#### 4. Pipeline
**File**: `pipeline/pipeline.go`

A chain of processing stages connected by channels, where the output of one stage becomes the input to the next.

**Key Concepts**:
- Data processing pipeline
- Stage composition
- Stream transformation

#### 5. Semaphore
**File**: `semaphore/semaphore.go`

A synchronization primitive that controls access to a common resource by multiple concurrent processes.

**Key Concepts**:
- Resource management
- Access control
- Concurrency limiting

#### 6. Worker Pool
**File**: `worker_pool/worker_pool.go`

A collection of goroutines waiting for tasks to be assigned, providing efficient task processing and resource utilization.

**Key Concepts**:
- Task queue management
- Resource pooling
- Concurrent processing

### Creational Patterns (5)

Patterns that deal with object creation mechanisms, trying to create objects in a manner suitable to the situation.

#### 1. Abstract Factory
**File**: `abstract_factory/abstract_factory.go`

Provides an interface for creating families of related or dependent objects without specifying their concrete classes.

**Key Concepts**:
- Family of objects
- Platform independence
- Object composition

#### 2. Builder
**File**: `builder/builder.go`

Separates the construction of a complex object from its representation, allowing the same construction process to create different representations.

**Key Concepts**:
- Step-by-step construction
- Complex object creation
- Fluent interfaces

#### 3. Factory Method
**File**: `factory_method/factory_method.go`

Defines an interface for creating an object but lets subclasses decide which class to instantiate. Factory Method lets a class defer instantiation to subclasses.

**Key Concepts**:
- Object creation delegation
- Subclass instantiation
- Extensibility

#### 4. Prototype
**File**: `prototype/prototype.go`

Creates new objects by copying an existing object, known as the prototype. This pattern is useful when creating objects is expensive.

**Key Concepts**:
- Object cloning
- Performance optimization
- Dynamic object creation

#### 5. Singleton
**File**: `singleton/singleton.go`

Ensures a class has only one instance and provides a global point of access to it.

**Key Concepts**:
- Single instance guarantee
- Global access point
- Resource management

### Structural Patterns (7)

Patterns that deal with object composition and relationships, forming larger structures from individual objects.

#### 1. Adapter
**File**: `adapter/adapter.go`

Allows incompatible interfaces to work together. It acts as a bridge between two incompatible interfaces.

**Key Concepts**:
- Interface compatibility
- Legacy system integration
- Protocol conversion

#### 2. Bridge
**File**: `bridge/bridge.go`

Decouples an abstraction from its implementation so that the two can vary independently.

**Key Concepts**:
- Abstraction-implementation separation
- Platform independence
- Runtime binding

#### 3. Composite
**File**: `composite/composite.go`

Composes objects into tree structures to represent part-whole hierarchies. It lets clients treat individual objects and compositions of objects uniformly.

**Key Concepts**:
- Tree structures
- Uniform treatment
- Recursive composition

#### 4. Decorator
**File**: `decorator/decorator.go`

Attaches additional responsibilities to an object dynamically. Decorators provide a flexible alternative to subclassing for extending functionality.

**Key Concepts**:
- Dynamic functionality addition
- Object wrapping
- Responsibility chaining

#### 5. Facade
**File**: `facade/facade.go`

Provides a simplified interface to a complex subsystem. It defines a higher-level interface that makes the subsystem easier to use.

**Key Concepts**:
- Simplified interface
- Subsystem complexity hiding
- Layered architecture

#### 6. Flyweight
**File**: `flyweight/flyweight.go`

Reduces memory usage by sharing as much data as possible with similar objects. It's useful when a large number of similar objects need to be created.

**Key Concepts**:
- Memory optimization
- Object sharing
- Intrinsic/extrinsic state

#### 7. Proxy
**File**: `proxy/proxy.go`

Provides a surrogate or placeholder for another object to control access to it. It can add functionality before or after the request reaches the actual object.

**Key Concepts**-
- Access control
- Lazy initialization
- Protection proxy

### Microservices Patterns (6)

Patterns specifically for building distributed systems and microservices architectures.

#### 1. API Gateway
**File**: `api_gateway/api_gateway.go`

Acts as a single entry point for all requests, routing them to the appropriate microservice and handling cross-cutting concerns.

**Key Concepts**:
- Request routing
- Service aggregation
- Cross-cutting concerns

#### 2. Bulkhead
**File**: `bulkhead/bulkhead.go`

Isolates different parts of the system to prevent cascading failures. Each part has its own resources and failure isolation.

**Key Concepts**:
- Failure isolation
- Resource partitioning
- Resilience patterns

#### 3. Circuit Breaker
**File**: `circuit_breaker/circuit_breaker.go`

Detects failures and encapsulates logic to prevent them from recurring. It wraps calls to potentially failing services.

**Key Concepts**:
- Failure detection
- Automatic recovery
- Fault tolerance

#### 4. Rate Limiting
**File**: `rate_limiting/rate_limiting.go`

Controls the rate of incoming requests to prevent service overload and ensure fair resource usage.

**Key Concepts**:
- Request throttling
- Resource protection
- Fair usage

#### 5. Saga
**File**: `saga/saga.go`

Manages distributed transactions using a sequence of local transactions. If one transaction fails, compensating transactions are executed.

**Key Concepts**:
- Distributed transactions
- Compensation patterns
- Eventual consistency

#### 6. Sidecar
**File**: `sidecar/sidecar.go`

Deploys helper services alongside the main application to enhance functionality without modifying the main application code.

**Key Concepts**:
- Service augmentation
- Infrastructure concerns
- Deployment patterns

## 🏗️ Project Structure

```
golang_small_projects/
├── main.go                          # Interactive main application
├── go.mod                           # Go module file
├── README.md                        # This file
├── behavioral/                      # Behavioral patterns
│   ├── chain_of_responsibility/
│   ├── command/
│   ├── mediator/
│   ├── memento/
│   ├── observer/
│   ├── state/
│   ├── strategy/
│   └── template_method/
├── concurrency/                     # Concurrency patterns
│   ├── barrier/
│   ├── fan_in_fan_out/
│   ├── generator/
│   ├── pipeline/
│   ├── semaphore/
│   └── worker_pool/
├── creational/                      # Creational patterns
│   ├── abstract_factory/
│   ├── builder/
│   ├── factory_method/
│   ├── prototype/
│   └── singleton/
├── structural/                      # Structural patterns
│   ├── adapter/
│   ├── bridge/
│   ├── composite/
│   ├── decorator/
│   ├── facade/
│   ├── flyweight/
│   └── proxy/
└── microservices/                   # Microservices patterns
    ├── api_gateway/
    ├── bulkhead/
    ├── circuit_breaker/
    ├── rate_limiting/
    ├── saga/
    └── sidecar/
```

## 🏃 Running Individual Patterns

Each pattern can be executed independently. Here are some examples:

### Behavioral Pattern Example
```bash
# Run the Observer pattern
go run ./observer/observer.go

# Run the Strategy pattern
go run ./strategy/strategy.go
```

### Concurrency Pattern Example
```bash
# Run the Worker Pool pattern
go run ./worker_pool/worker_pool.go

# Run the Pipeline pattern
go run ./pipeline/pipeline.go
```

### Creational Pattern Example
```bash
# Run the Factory Method pattern
go run ./factory_method/factory_method.go

# Run the Singleton pattern
go run ./singleton/singleton.go
```

### Structural Pattern Example
```bash
# Run the Decorator pattern
go run ./decorator/decorator.go

# Run the Facade pattern
go run ./facade/facade.go
```

### Microservices Pattern Example
```bash
# Run the Circuit Breaker pattern
go run ./circuit_breaker/circuit_breaker.go

# Run the API Gateway pattern
go run ./api_gateway/api_gateway.go
```

## 🎓 Learning Path

For beginners, we recommend the following learning order:

1. **Start with Creational Patterns** - Understand how objects are created
   - Singleton → Factory Method → Builder → Abstract Factory → Prototype

2. **Move to Structural Patterns** - Learn how objects are composed
   - Adapter → Decorator → Facade → Proxy → Composite → Bridge → Flyweight

3. **Explore Behavioral Patterns** - Study object interactions
   - Observer → Strategy → Command → Template Method → State → Chain of Responsibility → Mediator → Memento

4. **Master Concurrency Patterns** - Essential for Go programming
   - Worker Pool → Fan-in/Fan-out → Pipeline → Semaphore → Generator → Barrier

5. **Apply Microservices Patterns** - For distributed systems
   - Circuit Breaker → API Gateway → Rate Limiting → Bulkhead → Saga → Sidecar

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Guidelines

1. Follow Go best practices and idioms
2. Add comprehensive comments explaining the pattern
3. Include multiple examples if applicable
4. Ensure the code runs without errors
5. Update documentation as needed

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Design Patterns: Elements of Reusable Object-Oriented Software (Gang of Four)
- Patterns of Enterprise Application Architecture (Martin Fowler)
- Go Concurrency Patterns (William Kennedy)
- Microservices Patterns (Chris Richardson)

## 📞 Support

If you have any questions or suggestions, please:
- Open an issue on GitHub
- Check the existing documentation
- Run the interactive menu for guidance

---

**Happy coding with Go design patterns! 🚀**
