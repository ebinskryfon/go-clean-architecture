# Clean Architecture Go Server

A comprehensive guide to implementing Clean Architecture in Go, following Uncle Bob's principles for maintainable, testable, and scalable applications.

## 📋 Table of Contents

- [What is Clean Architecture?](#what-is-clean-architecture)
- [Project Structure](#project-structure)
- [Layer Responsibilities](#layer-responsibilities)
- [Dependency Flow](#dependency-flow)
- [Getting Started](#getting-started)
- [Implementation Examples](#implementation-examples)
- [Best Practices](#best-practices)
- [Testing Strategy](#testing-strategy)

## 🏗️ What is Clean Architecture?

Clean Architecture is a software design philosophy that emphasizes:

- **Independence of Frameworks**: Your business logic doesn't depend on external libraries
- **Testability**: Business rules can be tested without UI, database, or external services
- **Independence of UI**: The UI can change without affecting business rules
- **Independence of Database**: Business rules don't know about the database
- **Independence of External Agencies**: Business rules don't know about the outside world

### The Dependency Rule

The fundamental rule is that **dependencies point inward**. Inner layers should not know about outer layers.

```
Outer Layer → Inner Layer ✅
Inner Layer → Outer Layer ❌
```

## 📁 Project Structure

```
project/
├── cmd/                    # Application entry points
│   └── server/
│       └── main.go        # Main application bootstrap
├── internal/              # Private application code
│   ├── entity/           # Enterprise Business Rules (Innermost layer)
│   │   └── user.go       # Core business entities
│   ├── usecase/          # Application Business Rules
│   │   ├── interfaces/   # Use case interfaces/contracts
│   │   │   └── user_repository.go
│   │   └── user_usecase.go  # Business logic implementation
│   ├── adapter/          # Interface Adapters
│   │   ├── controller/   # HTTP handlers/controllers
│   │   │   └── user_controller.go
│   │   └── repository/   # Data access implementations
│   │       └── user_repository.go
│   └── infrastructure/   # Frameworks & Drivers (Outermost layer)
│       ├── database/     # Database connections
│       │   └── connection.go
│       └── server/       # HTTP server setup
│           └── server.go
├── pkg/                  # Public packages (can be imported by other projects)
│   └── response/         # Common response utilities
│       └── response.go
├── go.mod               # Go module definition
└── go.sum              # Dependency checksums
```

## 🎯 Layer Responsibilities

### 1. Entity Layer (Core/Domain)
**Location**: `internal/entity/`

- Contains enterprise-wide business rules
- Defines core business objects and their behaviors
- No dependencies on external layers
- Pure business logic, no infrastructure concerns

```go
// Example: user.go
type User struct {
    ID       string
    Email    string
    Password string
    Name     string
}

func (u *User) ValidateEmail() error {
    // Pure business logic
}
```

### 2. Use Case Layer (Application)
**Location**: `internal/usecase/`

- Contains application-specific business rules
- Orchestrates data flow to/from entities
- Defines interfaces for data access (Repository pattern)
- Independent of UI, database, or external services

```go
// Example: user_usecase.go
type UserUseCase struct {
    userRepo UserRepository // Interface, not concrete implementation
}

func (uc *UserUseCase) CreateUser(email, password, name string) (*User, error) {
    // Application-specific business logic
}
```

### 3. Adapter Layer (Interface Adapters)
**Location**: `internal/adapter/`

- Converts data between use cases and external world
- **Controllers**: Handle HTTP requests/responses
- **Repositories**: Implement data access interfaces
- **Presenters**: Format data for specific output needs

```go
// Example: user_controller.go
type UserController struct {
    userUseCase UserUseCase
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
    // Handle HTTP-specific concerns
    // Convert HTTP request to use case input
    // Call use case
    // Convert use case output to HTTP response
}
```

### 4. Infrastructure Layer (Frameworks & Drivers)
**Location**: `internal/infrastructure/`

- Outermost layer containing frameworks and tools
- Database connections, web servers, external APIs
- Implementation details that can be easily swapped
- Depends on inner layers, but inner layers don't depend on this

```go
// Example: connection.go
func NewDatabaseConnection(dsn string) (*sql.DB, error) {
    // Database-specific connection logic
}
```

### 5. Public Packages
**Location**: `pkg/`

- Reusable utilities that can be imported by other projects
- Common response structures, error handling, utilities
- Should be framework-agnostic

## 🔄 Dependency Flow

```
Infrastructure → Adapters → Use Cases → Entities

┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Infrastructure │───▶│    Adapters     │───▶│    Use Cases    │───▶│    Entities     │
│  (Database,     │    │  (Controllers,  │    │  (Business      │    │  (Core Business │
│   Web Server)   │    │   Repositories) │    │   Logic)        │    │   Rules)        │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
```

**Key Point**: Dependencies only flow inward. Entities know nothing about databases or HTTP.

## 🚀 Getting Started

### 1. Initialize the Project

```bash
mkdir clean-architecture-go
cd clean-architecture-go
go mod init your-project-name
```

### 2. Create the Directory Structure

```bash
mkdir -p cmd/server
mkdir -p internal/{entity,usecase/interfaces,adapter/{controller,repository},infrastructure/{database,server}}
mkdir -p pkg/response
```

### 3. Install Dependencies

```bash
go get github.com/gorilla/mux  # or your preferred router
go get github.com/lib/pq       # or your preferred database driver
```

### 4. Start with Entities

Begin by defining your core business entities in `internal/entity/`.

### 5. Define Use Case Interfaces

Create interfaces in `internal/usecase/interfaces/` that define what your use cases need.

### 6. Implement Use Cases

Implement your business logic in `internal/usecase/`.

### 7. Create Adapters

Implement controllers and repositories in `internal/adapter/`.

### 8. Set Up Infrastructure

Configure database connections and servers in `internal/infrastructure/`.

## 💡 Implementation Examples

### Entity Example
```go
// internal/entity/user.go
package entity

import (
    "errors"
    "regexp"
    "time"
)

type User struct {
    ID        string
    Email     string
    Password  string
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewUser(email, password, name string) (*User, error) {
    user := &User{
        Email:     email,
        Password:  password,
        Name:      name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := user.ValidateEmail(); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (u *User) ValidateEmail() error {
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(u.Email) {
        return errors.New("invalid email format")
    }
    return nil
}
```

### Repository Interface Example
```go
// internal/usecase/interfaces/user_repository.go
package interfaces

import "your-project/internal/entity"

type UserRepository interface {
    Create(user *entity.User) error
    GetByID(id string) (*entity.User, error)
    GetByEmail(email string) (*entity.User, error)
    Update(user *entity.User) error
    Delete(id string) error
}
```

### Use Case Example
```go
// internal/usecase/user_usecase.go
package usecase

import (
    "your-project/internal/entity"
    "your-project/internal/usecase/interfaces"
)

type UserUseCase struct {
    userRepo interfaces.UserRepository
}

func NewUserUseCase(userRepo interfaces.UserRepository) *UserUseCase {
    return &UserUseCase{
        userRepo: userRepo,
    }
}

func (uc *UserUseCase) CreateUser(email, password, name string) (*entity.User, error) {
    // Check if user already exists
    existingUser, _ := uc.userRepo.GetByEmail(email)
    if existingUser != nil {
        return nil, errors.New("user already exists")
    }
    
    // Create new user entity
    user, err := entity.NewUser(email, password, name)
    if err != nil {
        return nil, err
    }
    
    // Generate ID (this could be done in repository or here)
    user.ID = generateUUID() // Implement this function
    
    // Save to repository
    if err := uc.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

## ✅ Best Practices

### 1. Dependency Injection
- Use constructor injection to provide dependencies
- Inject interfaces, not concrete types
- Consider using a DI container for complex applications

### 2. Error Handling
- Define custom error types for your domain
- Don't expose internal errors to outer layers
- Use error wrapping to maintain context

### 3. Interface Design
- Keep interfaces small and focused (Interface Segregation Principle)
- Define interfaces in the package that uses them
- Use meaningful names that describe behavior

### 4. Package Organization
- Group related functionality together
- Avoid circular dependencies
- Use internal/ for private packages

### 5. Configuration
- Separate configuration from business logic
- Use environment variables or config files
- Validate configuration at startup

## 🧪 Testing Strategy

### Unit Testing
- Test entities and use cases in isolation
- Mock dependencies using interfaces
- Focus on business logic, not implementation details

```go
func TestUserUseCase_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    useCase := NewUserUseCase(mockRepo)
    
    // Act
    user, err := useCase.CreateUser("test@example.com", "password", "Test User")
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
}
```

### Integration Testing
- Test adapter layers with real implementations
- Use test databases or in-memory implementations
- Test the interaction between layers

### End-to-End Testing
- Test complete user journeys
- Use real HTTP requests
- Test with production-like environments

## 🎓 Learning Path

1. **Start Simple**: Implement a basic CRUD operation following the structure
2. **Add Complexity**: Introduce business rules and validations
3. **Test Everything**: Write comprehensive tests for each layer
4. **Refactor**: Improve your design as you learn
5. **Scale**: Add more features while maintaining the architecture

## 📚 Additional Resources

- [Clean Architecture Book by Robert C. Martin](https://www.amazon.com/Clean-Architecture-Craftsmans-Software-Structure/dp/0134494164)
- [Go Clean Architecture Examples](https://github.com/bxcodec/go-clean-arch)
- [Hexagonal Architecture in Go](https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3)

Remember: Clean Architecture is about **separation of concerns** and **dependency management**. Start simple and evolve your understanding as you build more complex applications.