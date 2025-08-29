# Clean Architecture Go Server

A robust Go server implementation following Clean Architecture principles, ensuring separation of concerns, testability, and maintainability.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Project Structure](#project-structure)
- [Dependencies Flow](#dependencies-flow)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Development](#development)
- [Testing](#testing)
- [Contributing](#contributing)

## Architecture Overview

This project implements Uncle Bob's Clean Architecture pattern, organizing code into distinct layers with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────┐
│                    Frameworks & Drivers                 │
│              (Web, DB, External APIs)                   │
├─────────────────────────────────────────────────────────┤
│                Interface Adapters                       │
│            (Controllers, Presenters, Gateways)         │
├─────────────────────────────────────────────────────────┤
│                   Use Cases                             │
│              (Business Rules)                           │
├─────────────────────────────────────────────────────────┤
│                    Entities                             │
│              (Enterprise Business Rules)               │
└─────────────────────────────────────────────────────────┘
```

### Core Principles

- **Dependency Inversion**: Inner layers don't depend on outer layers
- **Independence**: Business logic is independent of frameworks, UI, and databases
- **Testability**: Business rules can be tested without external dependencies
- **Flexibility**: Easy to change databases, web frameworks, or external services

## Project Structure

```
project/
├── cmd/                    # Application entry points
│   └── server/
│       └── main.go        # Main application bootstrap
├── internal/              # Private application code
│   ├── entity/           # Enterprise business rules
│   │   └── user.go       # User entity definition
│   ├── usecase/          # Application business rules
│   │   ├── interfaces/   # Use case interfaces/contracts
│   │   │   └── user_repository.go
│   │   └── user_usecase.go
│   ├── adapter/          # Interface adapters
│   │   ├── controller/   # HTTP handlers
│   │   │   └── user_controller.go
│   │   └── repository/   # Data access implementations
│   │       └── user_repository.go
│   └── infrastructure/   # External frameworks & tools
│       ├── database/     # Database connections
│       │   └── connection.go
│       └── server/       # HTTP server setup
│           └── server.go
├── pkg/                  # Public packages (reusable)
│   └── response/         # Standard API response format
│       └── response.go
├── go.mod               # Go module definition
└── go.sum               # Go module checksums
```

### Directory Explanations

#### `/cmd`
- Contains application entry points
- Each subdirectory represents a different executable
- Minimal logic, mainly calls into `/internal`

#### `/internal`
- Private application code that cannot be imported by other applications
- Core of the Clean Architecture implementation

##### `/internal/entity`
- **Layer**: Entities (innermost layer)
- Contains enterprise business rules
- Pure business objects with no external dependencies
- Defines core data structures and business invariants

##### `/internal/usecase`
- **Layer**: Use Cases
- Application-specific business rules
- Orchestrates entities and coordinates data flow
- Contains interfaces that outer layers must implement

##### `/internal/adapter`
- **Layer**: Interface Adapters
- Converts data between use cases and external systems
- Controllers handle HTTP requests/responses
- Repositories implement data persistence interfaces

##### `/internal/infrastructure`
- **Layer**: Frameworks & Drivers (outermost layer)
- External tools, frameworks, and drivers
- Database connections, HTTP servers, external API clients

#### `/pkg`
- Public packages that can be imported by external applications
- Reusable utilities and common functionality

## Dependencies Flow

The dependency flow follows the Clean Architecture dependency rule:

```
main.go → infrastructure → adapter → usecase → entity
    ↑           ↑            ↑         ↑        ↑
    └───────────┴────────────┴─────────┴────────┘
         Dependency Direction (inward only)
```

- **Entities** have no dependencies
- **Use Cases** depend only on entities and interfaces
- **Adapters** depend on use cases and entities
- **Infrastructure** depends on adapters, use cases, and entities
- **Main** orchestrates all layers through dependency injection

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL (or your preferred database)
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd clean-architecture-go-server
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run database migrations:
```bash
# Add your migration commands here
```

5. Start the server:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080` (or your configured port).

## API Endpoints

### User Management

| Method | Endpoint    | Description      | Request Body |
|--------|-------------|------------------|--------------|
| GET    | `/users`    | List all users   | -            |
| GET    | `/users/:id`| Get user by ID   | -            |
| POST   | `/users`    | Create new user  | User JSON    |
| PUT    | `/users/:id`| Update user      | User JSON    |
| DELETE | `/users/:id`| Delete user      | -            |

### Example Request/Response

**POST /users**
```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Response**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

## Development

### Adding New Features

1. **Define Entity** (if needed):
   ```go
   // internal/entity/product.go
   type Product struct {
       ID    int
       Name  string
       Price float64
   }
   ```

2. **Create Repository Interface**:
   ```go
   // internal/usecase/interfaces/product_repository.go
   type ProductRepository interface {
       Create(product *entity.Product) error
       GetByID(id int) (*entity.Product, error)
   }
   ```

3. **Implement Use Case**:
   ```go
   // internal/usecase/product_usecase.go
   type ProductUseCase struct {
       repo interfaces.ProductRepository
   }
   ```

4. **Implement Repository**:
   ```go
   // internal/adapter/repository/product_repository.go
   type productRepository struct {
       db *sql.DB
   }
   ```

5. **Create Controller**:
   ```go
   // internal/adapter/controller/product_controller.go
   type ProductController struct {
       usecase *usecase.ProductUseCase
   }
   ```

### Code Style Guidelines

- Follow Go conventions and `gofmt` formatting
- Use meaningful variable and function names
- Keep functions small and focused
- Write comprehensive tests for each layer
- Document public interfaces and complex logic

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/usecase/...
```

### Testing Strategy

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test interactions between layers
- **Mock Dependencies**: Use interfaces to mock external dependencies

Example test structure:
```go
func TestUserUseCase_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := &mockUserRepository{}
    usecase := NewUserUseCase(mockRepo)
    
    // Act
    err := usecase.CreateUser(validUser)
    
    // Assert
    assert.NoError(t, err)
}
```

## Environment Configuration

Create a `.env` file in the project root:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=cleanarch
DB_USER=username
DB_PASSWORD=password

# Server
SERVER_PORT=8080
SERVER_HOST=localhost

# Other configurations...
```

## Deployment

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Build Commands

```bash
# Build for current platform
go build -o bin/server cmd/server/main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/main.go
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Follow the project structure and conventions
4. Write tests for your changes
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Pull Request Guidelines

- Ensure all tests pass
- Follow the existing code style
- Update documentation if needed
- Add tests for new functionality
- Keep commits atomic and well-described

## Architecture Benefits

- **Maintainability**: Clear separation of concerns makes code easier to maintain
- **Testability**: Business logic can be tested independently of external systems
- **Flexibility**: Easy to swap out databases, frameworks, or external services
- **Scalability**: Well-organized code structure supports growth
- **Team Collaboration**: Clear boundaries enable parallel development

## Common Patterns Used

- **Dependency Injection**: Injecting dependencies through constructors
- **Repository Pattern**: Abstracting data access logic
- **Interface Segregation**: Small, focused interfaces
- **Factory Pattern**: Creating objects with complex initialization

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you have questions or need help:

1. Check the documentation
2. Look through existing issues
3. Create a new issue with detailed information
4. Contact the maintainers

---

**Happy Coding!** 🚀