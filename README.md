# Challenge Money

A REST API service for managing customer accounts and financial transactions. 

## Tech Stack

- **Language**: Go 1.25.3
- **Web Framework**: [chi](https://github.com/go-chi/chi) v5
- **Database**: PostgreSQL 18.0
- **Database Driver**: [pgx/v5](https://github.com/jackc/pgx)
- **Validation**: [validator/v10](https://github.com/go-playground/validator)
- **Containerization**: Docker Compose

## Project Structure

```
challenge-money/
├── internal/
│   ├── account/          # Account domain logic
│   │   ├── handler.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   └── handler_test.go
│   ├── transaction/      # Transaction domain logic
│   │   ├── handler.go
│   │   ├── models.go
│   │   ├── repository.go
│   │   └── handler_test.go
│   ├── health/           # Health check handlers
│   ├── database/         # Database connection utilities
│   └── httperrors/       # Custom HTTP error handling
├── main.go               # Application entry point
├── init.sql              # Database schema and seed data
├── docker-compose.yaml   # PostgreSQL container setup
├── Makefile              # Build and development tasks
└── run                   # Quick start script
```

## Getting Started

### Prerequisites

- Go 1.25.3 or later
- Docker and Docker Compose
- (Optional) [Air](https://github.com/air-verse/air) for live reload during development

### Installation

1. Clone the repository:
```bash
git clone https://github.com/rikw22/challenge-money
cd challenge-money
```

2. Install dependencies:
```bash
go mod download
```

3. (Optional) Install Air for live reload:
```bash
go install github.com/air-verse/air@latest
```

## Running the Application

### Quick Start

Use the provided run script:
```bash
./run
```

This will:
- Start PostgreSQL container via docker-compose
- Initialize the database with schema and seed data
- Start the API server with live reload (if Air is installed)

### Alternative Methods

**Using Makefile:**
```bash
make run
```

**Manual setup:**
```bash
# Start PostgreSQL
docker compose up -d

# Run the application
go run main.go
```

The server will start on `http://localhost:8080` by default.

## API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```
Returns the health status of the service.

### Create Account
```bash
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "document_number": "12345678901"
  }'
```

**Response** (201 Created):
```json
{
  "account_id": 1,
  "document_number": "12345678901"
}
```

### Get Account
```bash
curl http://localhost:8080/accounts/1
```

**Response** (200 OK):
```json
{
  "account_id": 1,
  "document_number": "12345678901"
}
```

### Create Transaction
```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": 1,
    "operation_type_id": 1,
    "amount": 123.45
  }'
```

**Operation Types:**
- `1` - Normal Purchase (negative amount)
- `2` - Purchase with installments (negative amount)
- `3` - Withdrawal (negative amount)
- `4` - Credit Voucher (positive amount)

**Response** (201 Created):
```json
{
  "transaction_id": "019a096b-ad9f-7f0e-88a4-9c93a754b029",
  "account_id": 1,
  "operation_type_id": 1,
  "amount": -123.45
}
```

## Database Schema

### Connection Details

When using docker compose:
- **Host**: localhost
- **Port**: 5432
- **Database**: postgres
- **User**: postgres
- **Password**: postgres

## Testing

Run all tests:
```bash
make test
# or
go test ./...
```

Run tests with coverage:
```bash
make test-coverage
# or
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Development

### Building the Application

```bash
make build
# or
go build
```

### Docker Commands

```bash
# Start PostgreSQL
make docker-up

# Stop PostgreSQL
make docker-down
```

### Environment Variables

| Variable       | Description                  | Default                                                                  |
|----------------|------------------------------|--------------------------------------------------------------------------|
| `PORT`         | Server port                  | `8080`                                                                   |
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable` |


## Future Improvements

- [ ] Graceful Shutdown - https://github.com/go-chi/chi/blob/master/_examples/graceful/main.go
- [ ] Database migrations (e.g., golang-migrate, goose)
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Rate limiting
- [ ] Authentication and authorization
- [ ] Audit logging
- [ ] Transaction rollback mechanisms
- [ ] Balance calculation and tracking