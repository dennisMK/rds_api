# API Documentation

## Overview
All endpoints return JSON and follow REST conventions.

## Base URL

\`\`\`
https://api.healthcare.example.com/api/v1
\`\`\`

## Authentication

All API requests must include a valid JWT token in the Authorization header:

\`\`\`
Authorization: Bearer <jwt-token>
\`\`\`

### Obtaining a Token

Tokens are issued by your authentication provider and should include:
- User ID and username
- Roles (admin, clinician, patient)
- Scopes (read, write, delete)

## Error Handling

The API uses FHIR OperationOutcome resources for error responses:

\`\`\`json
{
  "resourceType": "OperationOutcome",
  "issue": [{
    "severity": "error",
    "code": "invalid",
    "diagnostics": "Invalid patient ID format"
  }]
}
\`\`\`

### HTTP Status Codes

- `200 OK` - Successful GET/PUT
- `201 Created` - Successful POST
- `204 No Content` - Successful DELETE
- `400 Bad Request` - Invalid request format
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `422 Unprocessable Entity` - Validation errors
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

## Pagination

List endpoints support pagination using query parameters:

- `limit` - Number of items per page (default: 20, max: 100)
- `offset` - Number of items to skip (default: 0)

Example:
\`\`\`
GET /api/v1/patients?limit=50&offset=100
\`\`\`

Response includes pagination metadata:
\`\`\`json
{
  "resourceType": "Bundle",
  "total": 1500,
  "entry": [...],
  "link": [
    {
      "relation": "next",
      "url": "/api/v1/patients?limit=50&offset=150"
    }
  ]
}
\`\`\`

## Patient Endpoints

### Create Patient

**POST** `/patients`

Creates a new patient record.

**Required Scopes**: `patient:write`

**Request Body**:
\`\`\`json
{
  "name": [{
    "use": "official",
    "family": "Doe",
    "given": ["John", "Michael"]
  }],
  "gender": "male",
  "birthDate": "1990-01-15T00:00:00Z",
  "telecom": [{
    "system": "phone",
    "value": "+1-555-123-4567",
    "use": "home"
  }],
  "address": [{
    "use": "home",
    "line": ["123 Main St"],
    "city": "Anytown",
    "state": "CA",
    "postalCode": "12345",
    "country": "US"
  }]
}
\`\`\`

**Response**: `201 Created`
\`\`\`json
{
  "resourceType": "Patient",
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": [{
    "use": "official",
    "family": "Doe",
    "given": ["John", "Michael"]
  }],
  "gender": "male",
  "birthDate": "1990-01-15T00:00:00Z",
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z",
  "version": 1
}
\`\`\`

### Get Patient

**GET** `/patients/{id}`

Retrieves a patient by ID.

**Required Scopes**: `patient:read`

**Response**: `200 OK`
\`\`\`json
{
  "resourceType": "Patient",
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": [{
    "use": "official",
    "family": "Doe",
    "given": ["John", "Michael"]
  }],
  "gender": "male",
  "birthDate": "1990-01-15T00:00:00Z"
}
\`\`\`

### Update Patient

**PUT** `/patients/{id}`

Updates an existing patient record.

**Required Scopes**: `patient:write`

**Request Body**: Same as Create Patient (partial updates supported)

**Response**: `200 OK` with updated patient resource

### Delete Patient

**DELETE** `/patients/{id}`

Deletes a patient record.

**Required Scopes**: `patient:delete`

**Response**: `204 No Content`

### List Patients

**GET** `/patients`

Retrieves a paginated list of patients.

**Required Scopes**: `patient:read`

**Query Parameters**:
- `limit` - Items per page (default: 20)
- `offset` - Items to skip (default: 0)

**Response**: `200 OK`
\`\`\`json
{
  "resourceType": "Bundle",
  "type": "searchset",
  "total": 150,
  "entry": [{
    "fullUrl": "/api/v1/patients/550e8400-e29b-41d4-a716-446655440000",
    "resource": {
      "resourceType": "Patient",
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": [{"family": "Doe", "given": ["John"]}]
    }
  }]
}
\`\`\`

## Observation Endpoints

### Create Observation

**POST** `/observations`

Creates a new observation record.

**Required Scopes**: `observation:write`

**Request Body**:
\`\`\`json
{
  "status": "final",
  "category": [{
    "coding": [{
      "system": "http://terminology.hl7.org/CodeSystem/observation-category",
      "code": "vital-signs",
      "display": "Vital Signs"
    }]
  }],
  "code": {
    "coding": [{
      "system": "http://loinc.org",
      "code": "8867-4",
      "display": "Heart rate"
    }]
  },
  "subject": {
    "reference": "Patient/550e8400-e29b-41d4-a716-446655440000"
  },
  "effectiveDateTime": "2024-01-15T10:30:00Z",
  "valueQuantity": {
    "value": 72,
    "unit": "beats/min",
    "system": "http://unitsofmeasure.org",
    "code": "/min"
  }
}
\`\`\`

**Response**: `201 Created` with observation resource

### Get Observation

**GET** `/observations/{id}`

Retrieves an observation by ID.

**Required Scopes**: `observation:read`

### Update Observation

**PUT** `/observations/{id}`

Updates an existing observation record.

**Required Scopes**: `observation:write`

### Delete Observation

**DELETE** `/observations/{id}`

Deletes an observation record.

**Required Scopes**: `observation:delete`

### List Observations

**GET** `/observations`

Retrieves a paginated list of observations.

**Required Scopes**: `observation:read`

## FHIR Data Types

### HumanName
\`\`\`json
{
  "use": "official|usual|temp|nickname|anonymous|old|maiden",
  "text": "Full name as displayed",
  "family": "Family name",
  "given": ["Given", "names"],
  "prefix": ["Mr.", "Dr."],
  "suffix": ["Jr.", "III"]
}
\`\`\`

### ContactPoint
\`\`\`json
{
  "system": "phone|fax|email|pager|url|sms|other",
  "value": "Contact value",
  "use": "home|work|temp|old|mobile",
  "rank": 1
}
\`\`\`

### Address
\`\`\`json
{
  "use": "home|work|temp|old|billing",
  "type": "postal|physical|both",
  "text": "Full address text",
  "line": ["Street address lines"],
  "city": "City name",
  "district": "District/County",
  "state": "State/Province",
  "postalCode": "Postal code",
  "country": "Country"
}
\`\`\`

### CodeableConcept
\`\`\`json
{
  "coding": [{
    "system": "http://terminology-system-uri",
    "version": "System version",
    "code": "Code value",
    "display": "Human readable",
    "userSelected": true
  }],
  "text": "Plain text representation"
}
\`\`\`

### Quantity
\`\`\`json
{
  "value": 72.5,
  "comparator": "<|<=|>=|>",
  "unit": "Human readable unit",
  "system": "http://unitsofmeasure.org",
  "code": "UCUM code"
}
\`\`\`

## Rate Limiting

The API implements rate limiting to ensure fair usage:

- **Default Limit**: 100 requests per minute per IP address
- **Burst Capacity**: 20 requests
- **Headers**: Rate limit information is included in response headers

\`\`\`
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 2024-01-15T10:31:00Z
\`\`\`

When rate limit is exceeded:
\`\`\`json
{
  "resourceType": "OperationOutcome",
  "issue": [{
    "severity": "error",
    "code": "throttled",
    "diagnostics": "Rate limit exceeded"
  }]
}
\`\`\`

## Security Considerations

### Data Privacy
- All patient data is encrypted at rest and in transit
- Access is logged for audit purposes
- Minimum necessary access principle applies

### Authentication
- JWT tokens must be valid and not expired
- Tokens should be refreshed before expiration
- Use HTTPS in production environments

### Authorization
- Role-based access control (RBAC)
- Scope-based permissions for fine-grained access
- Resource-level access controls

## Monitoring and Observability

### Request Tracing
Each request includes a unique `X-Request-ID` header for tracing:
\`\`\`
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
\`\`\`

### Health Checks
Monitor API health using:
\`\`\`
GET /health
\`\`\`

Response:
\`\`\`json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
\`\`\`

### Metrics
System metrics are available at `/metrics` (admin access required):
- Request counts and error rates
- Response time percentiles
- Database connection statistics
- Cache performance metrics
