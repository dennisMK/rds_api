# Setup Guide

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 13 or higher
- Docker and Docker Compose (optional)
- Make (optional, for using Makefile commands)

## Local Development Setup

### 1. Clone the Repository

\`\`\`bash
git clone <repository-url>
cd healthcare-api
\`\`\`

### 2. Install Dependencies

\`\`\`bash
go mod download
\`\`\`

### 3. Database Setup

#### Option A: Using Docker Compose (Recommended)

\`\`\`bash
# Start PostgreSQL container
docker-compose up -d postgres

# Wait for database to be ready
sleep 10

# Run migrations
make migrate-up
\`\`\`

#### Option B: Local PostgreSQL

1. Install PostgreSQL locally
2. Create database:
\`\`\`sql
CREATE DATABASE rds;
CREATE USER healthcare_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE rds TO healthcare_user;
\`\`\`

3. Update environment variables (see step 4)
4. Run migrations:
\`\`\`bash
make migrate-up
\`\`\`

### 4. Environment Configuration

Copy the example environment file:
\`\`\`bash
cp .env.example .env
\`\`\`

Update `.env` with your configuration:
\`\`\`env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=healthcare_user
DB_PASSWORD=your_password
DB_NAME=rds
DB_SSL_MODE=disable

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY=24h

# Rate Limiting
RATE_LIMIT_REQUESTS=1000
RATE_LIMIT_WINDOW=1h

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Worker Pool
WORKER_POOL_SIZE=10
WORKER_QUEUE_SIZE=1000

# Cache
CACHE_TTL=1h
CACHE_CLEANUP_INTERVAL=10m
\`\`\`

### 5. Run the Application

#### Using Make (Recommended)
\`\`\`bash
# Development mode with hot reload
make dev

# Production mode
make run
\`\`\`

#### Using Go directly
\`\`\`bash
# Development
go run cmd/server/main.go

# Build and run
go build -o bin/server cmd/server/main.go
./bin/server
\`\`\`

### 6. Verify Installation

Check if the API is running:
\`\`\`bash
curl http://localhost:8080/health
\`\`\`

Expected response:
\`\`\`json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
\`\`\`

## Docker Deployment

### Build and Run with Docker Compose

\`\`\`bash
# Build and start all services
docker-compose up --build

# Run in background
docker-compose up -d --build
\`\`\`

### Build Docker Image Only

\`\`\`bash
# Build image
docker build -t healthcare-api .

# Run container
docker run -p 8080:8080 --env-file .env healthcare-api
\`\`\`

## Database Migrations

### Create New Migration

\`\`\`bash
# Create migration files
make migrate-create name=add_new_table
\`\`\`

This creates two files:
- `migrations/XXX_add_new_table.up.sql`
- `migrations/XXX_add_new_table.down.sql`

### Run Migrations

\`\`\`bash
# Apply all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Check migration status
make migrate-status
\`\`\`

## Testing

### Run Tests

\`\`\`bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration
\`\`\`

### Test Database Setup

For integration tests, set up a separate test database:

\`\`\`env
# .env.test
DB_NAME=rds_test
# ... other config
\`\`\`

## Development Tools

### Code Quality

\`\`\`bash
# Format code
make fmt

# Run linter
make lint

# Run security check
make security
\`\`\`

### API Documentation

Generate API documentation:
\`\`\`bash
# Generate OpenAPI spec
make docs

# Serve documentation locally
make serve-docs
\`\`\`

## Monitoring and Observability

### Health Checks

The API provides several health check endpoints:

- `GET /health` - Basic health check
- `GET /health/ready` - Readiness probe (includes DB connectivity)
- `GET /health/live` - Liveness probe

### Metrics

Metrics are exposed at `/metrics` in Prometheus format:

\`\`\`bash
curl http://localhost:8080/metrics
\`\`\`

Key metrics include:
- HTTP request duration and count
- Database connection pool stats
- Worker pool utilization
- Cache hit/miss rates

### Logging

Structured logging is configured with different levels:
- `debug`: Detailed debugging information
- `info`: General operational messages
- `warn`: Warning conditions
- `error`: Error conditions

Logs include correlation IDs for request tracing.

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check PostgreSQL is running
   - Verify connection parameters in `.env`
   - Ensure database exists and user has permissions

2. **Migration Errors**
   - Check database connectivity
   - Verify migration files syntax
   - Check if migrations were partially applied

3. **Port Already in Use**
   - Change `SERVER_PORT` in `.env`
   - Kill process using the port: `lsof -ti:8080 | xargs kill`

4. **JWT Token Issues**
   - Ensure `JWT_SECRET` is set and consistent
   - Check token expiry settings
   - Verify token format in requests

### Debug Mode

Enable debug logging:
\`\`\`env
LOG_LEVEL=debug
\`\`\`

This provides detailed information about:
- Database queries
- HTTP request/response details
- Worker pool operations
- Cache operations

### Performance Tuning

For high-load scenarios, adjust these settings:

\`\`\`env
# Database connections
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# Worker pool
WORKER_POOL_SIZE=20
WORKER_QUEUE_SIZE=2000

# Rate limiting
RATE_LIMIT_REQUESTS=2000
RATE_LIMIT_WINDOW=1h
\`\`\`

## Security Considerations

1. **Change Default Secrets**: Update `JWT_SECRET` and database passwords
2. **Use HTTPS**: Configure TLS certificates for production
3. **Database Security**: Use SSL connections and restricted user permissions
4. **Rate Limiting**: Adjust limits based on expected load
5. **Audit Logging**: Ensure audit logs are properly stored and monitored
6. **Input Validation**: All inputs are validated against FHIR schemas
7. **Authentication**: Implement proper user authentication and authorization
