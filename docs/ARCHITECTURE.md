# Architecture Documentation

## Overview
The system is designed to handle high-volume healthcare data transactions while maintaining FHIR R4 compliance and ensuring data security.

## System Architecture

\`\`\`
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Load Balancer │    │     API Gateway │    │   Rate Limiter  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │  HTTP Handlers  │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │ Service Layer   │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │ Repository Layer│
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │   PostgreSQL    │
                    └─────────────────┘
\`\`\`

## Directory Structure

\`\`\`
healthcare-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   ├── connection.go        # Database connection and pooling
│   │   └── migrations.go        # Migration management
│   ├── models/
│   │   ├── base.go              # Base FHIR types
│   │   ├── patient.go           # Patient FHIR resource
│   │   ├── observation.go       # Observation FHIR resource
│   │   └── errors.go            # Error types
│   ├── repository/
│   │   ├── base.go              # Base repository interface
│   │   ├── patient.go           # Patient data access
│   │   └── observation.go       # Observation data access
│   ├── service/
│   │   ├── patient.go           # Patient business logic
│   │   └── observation.go       # Observation business logic
│   ├── handlers/
│   │   ├── patient.go           # Patient HTTP handlers
│   │   └── observation.go       # Observation HTTP handlers
│   ├── middleware/
│   │   ├── auth.go              # Authentication middleware
│   │   ├── rate_limit.go        # Rate limiting
│   │   ├── security.go          # Security headers
│   │   ├── logging.go           # Request logging
│   │   ├── validation.go        # Input validation
│   │   └── audit.go             # Audit logging
│   ├── validation/
│   │   └── validator.go         # FHIR validation logic
│   ├── worker/
│   │   ├── pool.go              # Worker pool implementation
│   │   └── handlers.go          # Background job handlers
│   ├── concurrent/
│   │   ├── batch.go             # Batch processing
│   │   ├── pipeline.go          # Pipeline processing
│   │   └── cache.go             # Thread-safe caching
│   └── monitoring/
│       └── metrics.go           # Metrics collection
├── migrations/
│   ├── 001_create_patients_table.up.sql
│   ├── 001_create_patients_table.down.sql
│   ├── 002_create_observations_table.up.sql
│   ├── 002_create_observations_table.down.sql
│   ├── 003_create_audit_log_table.up.sql
│   └── 003_create_audit_log_table.down.sql
├── docs/
│   ├── API.md                   # API documentation
│   ├── SETUP.md                 # Setup instructions
│   └── ARCHITECTURE.md          # This file
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
└── .env.example
\`\`\`

## Core Components

### 1. HTTP Layer

**Location**: `internal/handlers/`

Responsible for:
- HTTP request/response handling
- Input validation and sanitization
- Error response formatting
- Content negotiation

**Key Features**:
- RESTful API design
- FHIR-compliant response formats
- Comprehensive error handling
- Request correlation IDs

### 2. Service Layer

**Location**: `internal/service/`

Responsible for:
- Business logic implementation
- Data transformation
- Transaction management
- Integration with external services

**Key Features**:
- Domain-driven design
- Transaction boundaries
- Business rule enforcement
- Event publishing

### 3. Repository Layer

**Location**: `internal/repository/`

Responsible for:
- Data persistence
- Query optimization
- Connection management
- Database transactions

**Key Features**:
- Interface-based design
- Query builders
- Connection pooling
- Audit trail generation

### 4. Middleware Stack

**Location**: `internal/middleware/`

Processing order:
1. **Security Headers**: CORS, CSP, security headers
2. **Rate Limiting**: Token bucket algorithm
3. **Authentication**: JWT token validation
4. **Authorization**: Role-based access control
5. **Logging**: Request/response logging
6. **Validation**: Input validation
7. **Audit**: Compliance logging

## Data Flow

### Request Processing Flow

\`\`\`
HTTP Request
     │
     ▼
Security Middleware
     │
     ▼
Rate Limiting
     │
     ▼
Authentication
     │
     ▼
Authorization
     │
     ▼
Input Validation
     │
     ▼
HTTP Handler
     │
     ▼
Service Layer
     │
     ▼
Repository Layer
     │
     ▼
Database
     │
     ▼
Response Processing
     │
     ▼
Audit Logging
     │
     ▼
HTTP Response
\`\`\`

### Background Processing Flow

\`\`\`
API Request
     │
     ▼
Job Queue
     │
     ▼
Worker Pool
     │
     ▼
Background Handler
     │
     ▼
Service Layer
     │
     ▼
Repository Layer
     │
     ▼
Database
     │
     ▼
Completion Notification
\`\`\`

## Concurrency Model

### Worker Pool Architecture

The system uses a worker pool pattern for handling background tasks:

- **Pool Size**: Configurable number of worker goroutines
- **Queue**: Buffered channel for job distribution
- **Job Types**: Data processing, notifications, cleanup
- **Error Handling**: Retry logic with exponential backoff

### Database Concurrency

- **Connection Pooling**: Optimized for high concurrency
- **Transaction Isolation**: Read committed level
- **Lock Management**: Row-level locking for updates
- **Deadlock Prevention**: Consistent lock ordering

### Cache Concurrency

- **Thread-Safe Operations**: Mutex-protected cache operations
- **TTL Management**: Background cleanup goroutines
- **Memory Management**: LRU eviction policy

## Security Architecture

### Authentication & Authorization

\`\`\`
Client Request
     │
     ▼
JWT Token Validation
     │
     ▼
User Context Extraction
     │
     ▼
Role-Based Access Control
     │
     ▼
Resource Permission Check
     │
     ▼
Request Processing
\`\`\`

### Data Protection

- **Encryption at Rest**: Database-level encryption
- **Encryption in Transit**: TLS 1.3 for all communications
- **Data Masking**: Sensitive data redaction in logs
- **Audit Trail**: Comprehensive activity logging

### Input Validation

- **Schema Validation**: FHIR resource schema compliance
- **Data Sanitization**: XSS and injection prevention
- **Business Rule Validation**: Healthcare-specific rules
- **Rate Limiting**: DDoS protection

## Database Design

### Schema Overview

\`\`\`sql
-- Core tables
patients
observations
audit_log

-- Indexes for performance
idx_patients_identifier
idx_patients_name
idx_observations_patient_id
idx_observations_code
idx_audit_log_timestamp
\`\`\`

### Data Integrity

- **Foreign Key Constraints**: Referential integrity
- **Check Constraints**: Data validation at DB level
- **Unique Constraints**: Prevent duplicates
- **Not Null Constraints**: Required field enforcement

### Performance Optimization

- **Indexing Strategy**: Query-optimized indexes
- **Partitioning**: Time-based partitioning for audit logs
- **Connection Pooling**: Optimized connection management
- **Query Optimization**: Prepared statements and query plans

## Monitoring & Observability

### Metrics Collection

- **HTTP Metrics**: Request count, duration, status codes
- **Database Metrics**: Connection pool, query performance
- **Worker Pool Metrics**: Queue size, processing time
- **Cache Metrics**: Hit ratio, eviction rate

### Health Checks

- **Liveness Probe**: Application health
- **Readiness Probe**: Dependency health
- **Database Health**: Connection and query tests
- **External Service Health**: Integration status

### Logging Strategy

- **Structured Logging**: JSON format for parsing
- **Log Levels**: Debug, Info, Warn, Error
- **Correlation IDs**: Request tracing
- **Audit Logs**: Compliance and security

## Scalability Considerations

### Horizontal Scaling

- **Stateless Design**: No server-side session state
- **Load Balancing**: Round-robin or least connections
- **Database Sharding**: Patient-based partitioning
- **Cache Distribution**: Redis cluster support

### Vertical Scaling

- **Resource Optimization**: Memory and CPU tuning
- **Connection Pooling**: Optimal pool sizing
- **Worker Pool Tuning**: Based on workload patterns
- **Database Optimization**: Query and index tuning

### Performance Targets

- **Response Time**: < 100ms for simple queries
- **Throughput**: > 1000 requests/second
- **Availability**: 99.9% uptime
- **Data Consistency**: Strong consistency for critical data

## Deployment Architecture

### Container Strategy

- **Multi-stage Builds**: Optimized image size
- **Security Scanning**: Vulnerability assessment
- **Resource Limits**: CPU and memory constraints
- **Health Checks**: Container health monitoring

### Infrastructure

- **Container Orchestration**: Kubernetes or Docker Swarm
- **Service Discovery**: DNS-based service resolution
- **Configuration Management**: Environment-based config
- **Secret Management**: Encrypted secret storage

## Compliance & Standards

### FHIR R4 Compliance

- **Resource Definitions**: Complete FHIR resource support
- **Search Parameters**: FHIR search specification
- **Bundle Support**: Transaction and batch bundles
- **Validation**: FHIR profile validation

### Healthcare Standards

- **HL7 Integration**: Message format support
- **HIPAA Compliance**: Privacy and security rules
- **Audit Requirements**: Comprehensive audit trails
- **Data Retention**: Configurable retention policies

## Future Enhancements

### Planned Features

- **GraphQL API**: Alternative query interface
- **Real-time Subscriptions**: WebSocket support
- **Advanced Analytics**: Data warehouse integration
- **Machine Learning**: Predictive analytics

### Scalability Improvements

- **Microservices**: Service decomposition
- **Event Sourcing**: Event-driven architecture
- **CQRS**: Command Query Responsibility Segregation
- **Distributed Caching**: Multi-region cache
