# Golang WhatsApp Clone - Project Structure

## Overview
This is a WhatsApp clone application built with Go, featuring GraphQL API, PostgreSQL database, authentication, and real-time messaging capabilities.

## Project Structure

```
golang-whatsapp-clone/
├── api/                           # API layer
│   └── index.go                   # API entry point
├── auth/                          # Authentication system
│   ├── handlers.go.txt            # Authentication handlers
│   ├── jwt.go                     # JWT token management
│   ├── middleware.go.txt          # Authentication middleware
│   ├── oauth.go                   # OAuth implementation
│   ├── types.go                   # Authentication type definitions
│   └── utils.go                   # Authentication utilities
├── cmd/                           # Command-line applications
│   └── seed/                      # Database seeding
│       └── main.go                # Seed data main function
├── config/                        # Configuration management
│   └── config.go                  # Application configuration
├── database/                      # Database layer
│   ├── database.go                # Database connection and setup
│   ├── gen/                       # Generated SQL code (sqlc)
│   │   ├── conversation_participants.sql.go
│   │   ├── conversations.sql.go
│   │   ├── copyfrom.go
│   │   ├── db.go
│   │   ├── messages.sql.go
│   │   ├── models.go
│   │   └── users.sql.go
│   ├── migrations/                # Database migrations
│   │   ├── 000001_enable_uuid_extension.down.sql
│   │   ├── 000001_enable_uuid_extension.up.sql
│   │   ├── 000002_create_users_table.down.sql
│   │   ├── 000002_create_users_table.up.sql
│   │   ├── 000003_create_conversations_table.down.sql
│   │   ├── 000003_create_conversations_table.up.sql
│   │   ├── 000004_create_conversation_participants.down.sql
│   │   ├── 000004_create_conversation_participants.up.sql
│   │   ├── 000005_create_messages.down.sql
│   │   └── 000005_create_messages.up.sql
│   └── queries/                   # SQL query definitions
│       ├── conversation_participants.sql
│       ├── conversations.sql
│       ├── messages.sql
│       └── users.sql
├── docker-compose.yml             # Docker container orchestration
├── docs/                          # Documentation
│   └── project-structure.md       # This file
├── errors/                        # Error handling
│   └── errors.go                  # Custom error definitions
├── go.mod                         # Go module definition
├── go.sum                         # Go module checksums
├── gqlgen.yml                     # GraphQL code generation config
├── graph/                         # GraphQL implementation
│   ├── generated.go               # Auto-generated GraphQL code
│   ├── messages.graphqls          # Messages GraphQL schema
│   ├── messages.resolvers.go      # Messages resolver implementation
│   ├── model/                     # GraphQL models
│   │   └── models_gen.go          # Generated GraphQL models
│   ├── pagination.graphqls        # Pagination GraphQL schema
│   ├── resolver.go                # Main resolver interface
│   ├── response.graphqls          # Response GraphQL schema
│   ├── schema.graphqls            # Main GraphQL schema
│   ├── schema.resolvers.go        # Main schema resolver
│   ├── user.graphqls              # User GraphQL schema
│   └── user.resolvers.go          # User resolver implementation
├── handler/                       # HTTP handlers
│   ├── auth.go                    # Authentication handlers
│   ├── gql.go                     # GraphQL handler
│   ├── handler.go                 # Main HTTP handler
│   └── healthcheck.go             # Health check endpoint
├── logger/                        # Logging system
│   └── logger.go                  # Logger configuration
├── main.go                        # Application entry point
├── nginx.conf                     # Nginx configuration
├── postgres_data/                 # PostgreSQL data directory
├── README.md                      # Project readme
├── repository/                    # Data access layer
│   ├── constants.go               # Repository constants
│   ├── conversation_repository.go # Conversation data access
│   ├── message_repository.go      # Message data access
│   ├── participant_repository.go  # Participant data access
│   └── utils.go                   # Repository utilities
├── server/                        # Server implementation
│   └── server.go                  # HTTP server setup
├── service/                       # Business logic layer
│   ├── conversation_service.go    # Conversation business logic
│   └── message_service.go         # Message business logic
├── sqlc.yaml                      # SQL code generation config
├── Taskfile.yml                   # Task runner configuration
├── tmp/                           # Temporary files
├── tools.go                       # Go tools configuration
├── vercel.json                    # Vercel deployment config
└── views/                         # Frontend views
    └── chats.html                 # Chat interface HTML
```

## Architecture Layers

### 1. **API Layer** (`/api`)
- Entry point for external API requests
- Handles routing and request/response formatting

### 2. **Authentication Layer** (`/auth`)
- JWT token management
- OAuth integration
- Authentication middleware
- User authentication handlers

### 3. **GraphQL Layer** (`/graph`)
- GraphQL schema definitions
- Resolver implementations
- Auto-generated code from gqlgen
- Supports real-time messaging

### 4. **Service Layer** (`/service`)
- Business logic implementation
- Conversation management
- Message handling
- Data validation and processing

### 5. **Repository Layer** (`/repository`)
- Data access abstraction
- Database operations
- Query execution
- Data persistence

### 6. **Database Layer** (`/database`)
- PostgreSQL database setup
- Migration management
- SQL query definitions
- Generated Go code from SQLC

### 7. **Handler Layer** (`/handler`)
- HTTP request handling
- GraphQL endpoint
- Health checks
- Authentication endpoints

## Key Technologies

- **Backend**: Go (Golang) 1.24.3
- **Database**: PostgreSQL with UUID support
- **API**: GraphQL with gqlgen
- **Authentication**: JWT + OAuth
- **Code Generation**: SQLC for database operations
- **Containerization**: Docker with docker-compose
- **Deployment**: Vercel support
- **Task Management**: Taskfile

## Core Dependencies & Packages

### **GraphQL & API**
- **gqlgen** (`github.com/99designs/gqlgen v0.17.78`) - GraphQL code generation
- **gqlparser** (`github.com/vektah/gqlparser/v2 v2.5.30`) - GraphQL parsing and validation
- **gorilla/websocket** (`github.com/gorilla/websocket v1.5.0`) - WebSocket support for real-time messaging

### **Database & ORM**
- **pgx** (`github.com/jackc/pgx/v5 v5.7.5`) - PostgreSQL driver with high performance
- **sqlc** (`github.com/sqlc-dev/sqlc v1.29.0`) - SQL code generation tool
- **google/uuid** (`github.com/google/uuid v1.6.0`) - UUID generation and handling

### **Authentication & Security**
- **jwt** (`github.com/golang-jwt/jwt/v5 v5.3.0`) - JWT token implementation
- **oauth2** (`golang.org/x/oauth2 v0.30.0`) - OAuth 2.0 client implementation

### **Configuration & Environment**
- **viper** (`github.com/spf13/viper v1.20.1`) - Configuration management
- **godotenv** (`github.com/joho/godotenv v1.5.1`) - Environment variable loading

### **Logging & Monitoring**
- **zerolog** (`github.com/rs/zerolog v1.34.0`) - Structured logging library

### **Development Tools**
- **sqlc** - Database code generation (configured in `sqlc.yaml`)
- **gqlgen** - GraphQL code generation (configured in `gqlgen.yml`)

### **Database Schema & Migrations**
The application uses a well-structured database with:
- UUID extension support
- Users table for authentication
- Conversations table for chat rooms
- Conversation participants for user management
- Messages table for chat content

## Database Schema

The application uses the following main entities:
- **Users**: User accounts and profiles
- **Conversations**: Chat conversations
- **Conversation Participants**: Users in conversations
- **Messages**: Individual chat messages

## Development Workflow

1. **Database Changes**: Update SQL files in `/database/queries/` and run migrations
2. **GraphQL Changes**: Modify `.graphqls` files and regenerate with gqlgen
3. **Business Logic**: Implement in `/service/` layer
4. **Data Access**: Use repository pattern in `/repository/` layer
5. **API Endpoints**: Handle in `/handler/` layer

## Getting Started

1. Ensure Go 1.x is installed
2. Run `docker-compose up` to start PostgreSQL
3. Execute database migrations
4. Run `go mod tidy` to install dependencies
5. Start the application with `go run main.go`

## File Naming Conventions

- **Go files**: snake_case.go
- **GraphQL schemas**: *.graphqls
- **SQL files**: snake_case.sql
- **Migrations**: numbered_sequence_description.up/down.sql
- **Configuration**: kebab-case.yml/.yaml

This structure follows clean architecture principles with clear separation of concerns between layers, making the codebase maintainable and scalable.
