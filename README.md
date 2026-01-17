# Karir Nusantara API

A scalable, maintainable backend API for the Karir Nusantara job portal platform.

## 🏗️ Architecture

This project follows **Clean Architecture** with a modular monolith approach, ensuring:
- Clear separation of concerns
- Easy testing and maintenance
- Future-ready for microservices migration

### Why This Architecture?

1. **Modularity**: Each domain (auth, jobs, applications) is self-contained
2. **Testability**: Business logic is isolated from infrastructure
3. **Flexibility**: Easy to swap databases, frameworks, or external services
4. **Scalability**: Modules can be extracted to microservices when needed

## 📁 Project Structure

```
karir-nusantara-api/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Environment configuration
│   ├── database/
│   │   └── mysql.go             # Database connection
│   ├── middleware/
│   │   ├── auth.go              # JWT authentication
│   │   ├── cors.go              # CORS handling
│   │   └── logging.go           # Request logging
│   ├── modules/
│   │   ├── auth/
│   │   │   ├── entity.go        # User entity
│   │   │   ├── repository.go    # Data access
│   │   │   ├── service.go       # Business logic
│   │   │   ├── handler.go       # HTTP handlers
│   │   │   └── routes.go        # Route definitions
│   │   ├── jobs/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   ├── handler.go
│   │   │   └── routes.go
│   │   ├── cvs/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   ├── handler.go
│   │   │   └── routes.go
│   │   ├── applications/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   ├── handler.go
│   │   │   └── routes.go
│   │   └── timelines/
│   │       ├── entity.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       ├── handler.go
│   │       └── routes.go
│   └── shared/
│       ├── response/
│       │   └── response.go      # Standard API responses
│       ├── validator/
│       │   └── validator.go     # Input validation
│       └── errors/
│           └── errors.go        # Custom error types
├── migrations/
│   └── 001_initial_schema.sql   # Database migrations
├── docs/
│   └── api.md                   # API documentation
├── .env.example
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- MySQL 8.0+
- Make (optional)

### Setup

1. Clone the repository
2. Copy environment file:
   ```bash
   cp .env.example .env
   ```
3. Configure your `.env` file
4. Run migrations:
   ```bash
   make migrate-up
   ```
5. Start the server:
   ```bash
   make run
   ```

## 🔑 Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `APP_PORT` | Server port | `8080` |
| `APP_ENV` | Environment | `development` |
| `DB_HOST` | MySQL host | `localhost` |
| `DB_PORT` | MySQL port | `3306` |
| `DB_USER` | MySQL user | `root` |
| `DB_PASSWORD` | MySQL password | `password` |
| `DB_NAME` | Database name | `karir_nusantara` |
| `JWT_SECRET` | JWT signing key | `your-secret-key` |
| `JWT_EXPIRY` | Token expiry | `24h` |

## 📚 API Documentation

See [docs/api.md](docs/api.md) for complete API documentation.

## 🧪 Running Tests

```bash
make test
```

## 📦 MVP Implementation Order

1. **Phase 1**: Auth module (register, login, JWT)
2. **Phase 2**: Jobs module (CRUD, search, filter)
3. **Phase 3**: CV module (create, update, snapshot)
4. **Phase 4**: Applications module (apply, list)
5. **Phase 5**: Timeline module (events, status updates)

## 📄 License

MIT License
# Karir-Nusantara-API
