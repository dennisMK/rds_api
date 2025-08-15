# RDS Healthcare API
Its now compliant with FHIR (Fast Healthcare Interoperability Resources) R4 standards.

## Features

- **FHIR R4 Compliance**: Implements FHIR Patient and Observation resources with full validation
- **High Performance**: Multithreaded architecture with connection pooling for high-volume transactions
- **Enterprise Security**: JWT authentication, role-based access control, rate limiting, and audit logging
- **Scalability**: Worker pools, batch processing, and concurrent request handling
- **Healthcare Standards**: Designed for healthcare interoperability with HL7 and FHIR compliance
- **Production Ready**: Comprehensive logging, monitoring, error handling, and graceful shutdown

## Architecture

\`\`\`
├── cmd/
│   └── server/          # Application entrypoint
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and migrations
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware (auth, logging, security)
│   ├── models/          # FHIR resource models
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic layer
│   ├── validation/      # Input validation
│   ├── worker/          # Background job processing
│   ├── concurrent/      # Concurrency utilities
│   └── monitoring/      # Performance monitoring
├── migrations/          # Database migrations
└── docs/               # Documentation
\`\`\`

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Docker (optional)

### Local Development

1. **Clone and setup**
   \`\`\`bash
   git clone <repository-url>
   cd healthcare-api
   cp .env.example .env
   \`\`\`

2. **Configure environment**
   Edit `.env` with your database credentials and JWT secret.

3. **Start PostgreSQL**
   \`\`\`bash
   # Using Docker
   docker run --name postgres -e POSTGRES_DB=rds -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:15-alpine
   
   # Or use your local PostgreSQL instance
   \`\`\`

4. **Run migrations**
   \`\`\`bash
   make migrate-up
   \`\`\`

5. **Start the server**
   \`\`\`bash
   make run
   \`\`\`

The API will be available at `http://localhost:8080`

### Docker Development

\`\`\`bash
docker-compose up --build
\`\`\`

## API Documentation

### Base URL
\`\`\`
http://localhost:8080/api/v1
\`\`\`

### Authentication

All endpoints require JWT authentication via the `Authorization` header:
\`\`\`
Authorization: Bearer <jwt-token>
\`\`\`

### Endpoints

#### Health Check
- `GET /health` - Service health status

#### Patients
- `POST /patients` - Create a new patient
- `GET /patients/{id}` - Get patient by ID
- `PUT /patients/{id}` - Update patient
- `DELETE /patients/{id}` - Delete patient
- `GET /patients` - List patients with pagination

#### Observations
- `POST /observations` - Create a new observation
- `GET /observations/{id}` - Get observation by ID
- `PUT /observations/{id}` - Update observation
- `DELETE /observations/{id}` - Delete observation
- `GET /observations` - List observations with pagination

### Request/Response Examples

#### Create Patient
\`\`\`bash
curl -X POST http://localhost:8080/api/v1/patients \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": [{
      "use": "official",
      "family": "Doe",
      "given": ["John"]
    }],
    "gender": "male",
    "birthDate": "1990-01-01T00:00:00Z"
  }'
\`\`\`

#### Create Observation
\`\`\`bash
curl -X POST http://localhost:8080/api/v1/observations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "final",
    "code": {
      "coding": [{
        "system": "http://loinc.org",
        "code": "8867-4",
        "display": "Heart rate"
      }]
    },
    "subject": {
      "reference": "Patient/{patient-id}"
    },
    "valueQuantity": {
      "value": 72,
      "unit": "beats/min",
      "system": "http://unitsofmeasure.org",
      "code": "/min"
    }
  }'
\`\`\`

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Application environment | `development` |
| `SERVER_PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | `rds` |
| `JWT_SECRET` | JWT signing secret | - |
| `LOG_LEVEL` | Log level (1-6) | `4` |

### Database Configuration

The API uses PostgreSQL with optimized connection pooling:
- Max Open Connections: 200
- Max Idle Connections: 50
- Connection Max Lifetime: 10 minutes
- Connection Max Idle Time: 2 minutes

## Security

### Authentication & Authorization

- **JWT Tokens**: All API endpoints require valid JWT tokens
- **Role-Based Access**: Support for user roles (admin, clinician, patient)
- **Scope-Based Access**: Fine-grained permissions using OAuth2-style scopes

### Security Headers

- Content Security Policy (CSP)
- X-Frame-Options: DENY
- X-Content-Type-Options: nosniff
- Strict-Transport-Security (HTTPS)
- Referrer-Policy: strict-origin-when-cross-origin

### Rate Limiting

- Default: 100 requests per minute per IP
- Configurable per endpoint
- Token bucket algorithm with burst capacity

### Audit Logging

All API requests are logged with:
- Request ID for tracing
- User identification
- Resource access patterns
- Compliance with healthcare regulations

## Performance

### Concurrency Features

- **Worker Pools**: Background job processing with configurable workers
- **Batch Processing**: Efficient bulk operations with controlled parallelism
- **Pipeline Processing**: Data transformation with concurrent stages
- **Connection Pooling**: Optimized database connections for high throughput

### Monitoring

Built-in metrics collection:
- Request count and error rates
- Response times and latency percentiles
- Database connection statistics
- Cache hit rates
- Worker pool performance

Access metrics at: `GET /metrics` (requires admin role)

## Deployment

### Production Deployment

1. **Build the application**
   \`\`\`bash
   make build
   \`\`\`

2. **Set production environment variables**
   \`\`\`bash
   export ENVIRONMENT=production
   export JWT_SECRET=<secure-random-string>
   export DB_PASSWORD=<secure-password>
   \`\`\`

3. **Run migrations**
   \`\`\`bash
   make migrate-up
   \`\`\`

4. **Start the server**
   \`\`\`bash
   ./bin/server
   \`\`\`

### Docker Deployment

\`\`\`bash
docker build -t healthcare-api .
docker run -p 8080:8080 --env-file .env healthcare-api
\`\`\`

### Kubernetes Deployment

See `k8s/` directory for Kubernetes manifests.

## Development

### Running Tests

\`\`\`bash
make test
\`\`\`

### Code Formatting

\`\`\`bash
make fmt
\`\`\`

### Linting

\`\`\`bash
make lint
\`\`\`

### Database Migrations

Create new migration:
\`\`\`bash
migrate create -ext sql -dir migrations -seq migration_name
\`\`\`

Run migrations:
\`\`\`bash
make migrate-up
\`\`\`

Rollback migrations:
\`\`\`bash
make migrate-down
\`\`\`

## FHIR Compliance

This API implements FHIR R4 resources with full validation:

### Supported Resources

- **Patient**: Demographics, contact information, identifiers
- **Observation**: Clinical observations, vital signs, lab results

### FHIR Features

- Complete resource validation
- FHIR-compliant error responses (OperationOutcome)
- Bundle responses for search operations
- Proper HTTP status codes and headers
- FHIR-standard pagination

### Standards Compliance

- HL7 FHIR R4 specification
- Healthcare interoperability standards
- HIPAA compliance considerations
- Audit logging for regulatory requirements

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check PostgreSQL is running
   - Verify connection parameters in `.env`
   - Ensure database exists

2. **JWT Token Invalid**
   - Check JWT_SECRET configuration
   - Verify token hasn't expired
   - Ensure proper Authorization header format

3. **Rate Limit Exceeded**
   - Implement exponential backoff
   - Check rate limiting configuration
   - Consider upgrading to higher tier

### Logs

Application logs are structured JSON format:
\`\`\`bash
# View logs
docker logs healthcare-api

# Follow logs
docker logs -f healthcare-api
\`\`\`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run linting and tests
6. Submit a pull request

### Code Standards

- Follow Go best practices
- Write comprehensive tests
- Document public APIs
- Use structured logging
- Maintain FHIR compliance

## License

MIT License - see LICENSE file for details.

## Support

For issues and questions:
- Create an issue in the repository
- Check the troubleshooting guide
- Review the API documentation

---

**Note**: This is a healthcare API handling sensitive data. Ensure proper security measures, compliance with healthcare regulations (HIPAA, GDPR), and regular security audits in production environments.
