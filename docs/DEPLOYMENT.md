# Deployment Guide

## Overview

This guide covers deploying the RDS Healthcare API in various environments, from development to production.

## Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Docker (optional)
- Kubernetes (for container orchestration)

## Environment Configuration

### Environment Variables

Create a `.env` file or set environment variables:

\`\`\`bash
# Server Configuration
ENVIRONMENT=production
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30
SERVER_WRITE_TIMEOUT=30
SERVER_IDLE_TIMEOUT=120

# Database Configuration
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=healthcare_api
DB_PASSWORD=secure-password
DB_NAME=healthcare_db
DB_SSL_MODE=require

# Security Configuration
JWT_SECRET=your-256-bit-secret-key
JWT_EXPIRATION=3600

# Logging
LOG_LEVEL=4
\`\`\`

### Security Considerations

1. **JWT Secret**: Use a cryptographically secure random string (256 bits minimum)
2. **Database Password**: Use strong, unique passwords
3. **SSL/TLS**: Always use SSL in production (`DB_SSL_MODE=require`)
4. **Environment Isolation**: Never use development credentials in production

## Local Development

### Using Go Directly

1. **Install dependencies**
   \`\`\`bash
   go mod download
   \`\`\`

2. **Setup database**
   \`\`\`bash
   # Start PostgreSQL
   docker run --name postgres-dev \
     -e POSTGRES_DB=healthcare_db \
     -e POSTGRES_USER=healthcare_api \
     -e POSTGRES_PASSWORD=dev-password \
     -p 5432:5432 -d postgres:15-alpine
   \`\`\`

3. **Run migrations**
   \`\`\`bash
   make migrate-up
   \`\`\`

4. **Start the server**
   \`\`\`bash
   make run
   \`\`\`

### Using Docker Compose

\`\`\`bash
# Start all services
docker-compose up --build

# Run in background
docker-compose up -d --build

# View logs
docker-compose logs -f api

# Stop services
docker-compose down
\`\`\`

## Production Deployment

### Binary Deployment

1. **Build the application**
   \`\`\`bash
   # Build for Linux (if cross-compiling)
   GOOS=linux GOARCH=amd64 go build -o healthcare-api cmd/server/main.go
   
   # Or use make
   make build
   \`\`\`

2. **Prepare the server**
   \`\`\`bash
   # Create application user
   sudo useradd -r -s /bin/false healthcare-api
   
   # Create directories
   sudo mkdir -p /opt/healthcare-api/{bin,logs,migrations}
   sudo chown -R healthcare-api:healthcare-api /opt/healthcare-api
   \`\`\`

3. **Deploy files**
   \`\`\`bash
   # Copy binary
   sudo cp healthcare-api /opt/healthcare-api/bin/
   
   # Copy migrations
   sudo cp -r migrations/* /opt/healthcare-api/migrations/
   
   # Set permissions
   sudo chmod +x /opt/healthcare-api/bin/healthcare-api
   \`\`\`

4. **Create systemd service**
   \`\`\`bash
   sudo tee /etc/systemd/system/healthcare-api.service > /dev/null <<EOF
   [Unit]
   Description=Healthcare API
   After=network.target postgresql.service
   
   [Service]
   Type=simple
   User=healthcare-api
   Group=healthcare-api
   WorkingDirectory=/opt/healthcare-api
   ExecStart=/opt/healthcare-api/bin/healthcare-api
   EnvironmentFile=/opt/healthcare-api/.env
   Restart=always
   RestartSec=5
   StandardOutput=journal
   StandardError=journal
   
   [Install]
   WantedBy=multi-user.target
   EOF
   \`\`\`

5. **Start the service**
   \`\`\`bash
   sudo systemctl daemon-reload
   sudo systemctl enable healthcare-api
   sudo systemctl start healthcare-api
   sudo systemctl status healthcare-api
   \`\`\`

### Docker Deployment

1. **Build Docker image**
   \`\`\`bash
   docker build -t healthcare-api:latest .
   \`\`\`

2. **Run container**
   \`\`\`bash
   docker run -d \
     --name healthcare-api \
     --restart unless-stopped \
     -p 8080:8080 \
     --env-file .env \
     healthcare-api:latest
   \`\`\`

3. **With Docker Compose (Production)**
   \`\`\`yaml
   version: '3.8'
   
   services:
     postgres:
       image: postgres:15-alpine
       environment:
         POSTGRES_DB: healthcare_db
         POSTGRES_USER: healthcare_api
         POSTGRES_PASSWORD_FILE: /run/secrets/db_password
       volumes:
         - postgres_data:/var/lib/postgresql/data
       secrets:
         - db_password
       networks:
         - healthcare_network
   
     api:
       image: healthcare-api:latest
       depends_on:
         - postgres
       environment:
         DB_HOST: postgres
         DB_PASSWORD_FILE: /run/secrets/db_password
         JWT_SECRET_FILE: /run/secrets/jwt_secret
       secrets:
         - db_password
         - jwt_secret
       networks:
         - healthcare_network
       ports:
         - "8080:8080"
   
   secrets:
     db_password:
       file: ./secrets/db_password.txt
     jwt_secret:
       file: ./secrets/jwt_secret.txt
   
   volumes:
     postgres_data:
   
   networks:
     healthcare_network:
   \`\`\`

### Kubernetes Deployment

1. **Create namespace**
   \`\`\`yaml
   apiVersion: v1
   kind: Namespace
   metadata:
     name: healthcare-api
   \`\`\`

2. **Create secrets**
   \`\`\`yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: healthcare-api-secrets
     namespace: healthcare-api
   type: Opaque
   data:
     db-password: <base64-encoded-password>
     jwt-secret: <base64-encoded-secret>
   \`\`\`

3. **Create ConfigMap**
   \`\`\`yaml
   apiVersion: v1
   kind: ConfigMap
   metadata:
     name: healthcare-api-config
     namespace: healthcare-api
   data:
     ENVIRONMENT: "production"
     SERVER_PORT: "8080"
     DB_HOST: "postgres-service"
     DB_PORT: "5432"
     DB_USER: "healthcare_api"
     DB_NAME: "healthcare_db"
     DB_SSL_MODE: "require"
     LOG_LEVEL: "4"
   \`\`\`

4. **Create Deployment**
   \`\`\`yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: healthcare-api
     namespace: healthcare-api
   spec:
     replicas: 3
     selector:
       matchLabels:
         app: healthcare-api
     template:
       metadata:
         labels:
           app: healthcare-api
       spec:
         containers:
         - name: healthcare-api
           image: healthcare-api:latest
           ports:
           - containerPort: 8080
           envFrom:
           - configMapRef:
               name: healthcare-api-config
           env:
           - name: DB_PASSWORD
             valueFrom:
               secretKeyRef:
                 name: healthcare-api-secrets
                 key: db-password
           - name: JWT_SECRET
             valueFrom:
               secretKeyRef:
                 name: healthcare-api-secrets
                 key: jwt-secret
           livenessProbe:
             httpGet:
               path: /health
               port: 8080
             initialDelaySeconds: 30
             periodSeconds: 10
           readinessProbe:
             httpGet:
               path: /health
               port: 8080
             initialDelaySeconds: 5
             periodSeconds: 5
           resources:
             requests:
               memory: "256Mi"
               cpu: "250m"
             limits:
               memory: "512Mi"
               cpu: "500m"
   \`\`\`

5. **Create Service**
   \`\`\`yaml
   apiVersion: v1
   kind: Service
   metadata:
     name: healthcare-api-service
     namespace: healthcare-api
   spec:
     selector:
       app: healthcare-api
     ports:
     - protocol: TCP
       port: 80
       targetPort: 8080
     type: ClusterIP
   \`\`\`

6. **Create Ingress**
   \`\`\`yaml
   apiVersion: networking.k8s.io/v1
   kind: Ingress
   metadata:
     name: healthcare-api-ingress
     namespace: healthcare-api
     annotations:
       kubernetes.io/ingress.class: nginx
       cert-manager.io/cluster-issuer: letsencrypt-prod
       nginx.ingress.kubernetes.io/ssl-redirect: "true"
   spec:
     tls:
     - hosts:
       - api.healthcare.example.com
       secretName: healthcare-api-tls
     rules:
     - host: api.healthcare.example.com
       http:
         paths:
         - path: /
           pathType: Prefix
           backend:
             service:
               name: healthcare-api-service
               port:
                 number: 80
   \`\`\`

## Database Setup

### PostgreSQL Configuration

1. **Create database and user**
   \`\`\`sql
   CREATE DATABASE healthcare_db;
   CREATE USER healthcare_api WITH ENCRYPTED PASSWORD 'secure-password';
   GRANT ALL PRIVILEGES ON DATABASE healthcare_db TO healthcare_api;
   \`\`\`

2. **Configure PostgreSQL for production**
   \`\`\`bash
   # postgresql.conf
   max_connections = 200
   shared_buffers = 256MB
   effective_cache_size = 1GB
   maintenance_work_mem = 64MB
   checkpoint_completion_target = 0.9
   wal_buffers = 16MB
   default_statistics_target = 100
   random_page_cost = 1.1
   effective_io_concurrency = 200
   
   # Enable SSL
   ssl = on
   ssl_cert_file = 'server.crt'
   ssl_key_file = 'server.key'
   \`\`\`

3. **Run migrations**
   \`\`\`bash
   # Using migrate tool
   migrate -path migrations -database "postgres://user:pass@host:port/db?sslmode=require" up
   
   # Or using make
   make migrate-up
   \`\`\`

### Database Backup

1. **Automated backups**
   \`\`\`bash
   #!/bin/bash
   # backup.sh
   BACKUP_DIR="/backups"
   DATE=$(date +%Y%m%d_%H%M%S)
   
   pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME | gzip > $BACKUP_DIR/healthcare_db_$DATE.sql.gz
   
   # Keep only last 30 days
   find $BACKUP_DIR -name "healthcare_db_*.sql.gz" -mtime +30 -delete
   \`\`\`

2. **Cron job**
   \`\`\`bash
   # Daily backup at 2 AM
   0 2 * * * /opt/healthcare-api/scripts/backup.sh
   \`\`\`

## Monitoring and Logging

### Application Logs

1. **Log rotation with logrotate**
   \`\`\`bash
   # /etc/logrotate.d/healthcare-api
   /opt/healthcare-api/logs/*.log {
       daily
       missingok
       rotate 30
       compress
       delaycompress
       notifempty
       create 644 healthcare-api healthcare-api
       postrotate
           systemctl reload healthcare-api
       endscript
   }
   \`\`\`

### Health Checks

1. **Application health**
   \`\`\`bash
   curl -f http://localhost:8080/health || exit 1
   \`\`\`

2. **Database connectivity**
   \`\`\`bash
   # Check database connection
   pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER
   \`\`\`

### Metrics Collection

1. **Prometheus configuration**
   \`\`\`yaml
   # prometheus.yml
   scrape_configs:
   - job_name: 'healthcare-api'
     static_configs:
     - targets: ['localhost:8080']
     metrics_path: '/metrics'
     scrape_interval: 15s
   \`\`\`

## Security Hardening

### Network Security

1. **Firewall rules**
   \`\`\`bash
   # Allow only necessary ports
   ufw allow 22/tcp    # SSH
   ufw allow 8080/tcp  # API
   ufw allow 5432/tcp  # PostgreSQL (if external)
   ufw enable
   \`\`\`

2. **Reverse proxy with Nginx**
   ```nginx
   server {
       listen 80;
       server_name api.healthcare.example.com;
       return 301 https://$server_name$request_uri;
   }
   
   server {
       listen 443 ssl http2;
       server_name api.healthcare.example.com;
       
       ssl_certificate /path/to/cert.pem;
       ssl_certificate_key /path/to/key.pem;
       
       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }
   \`\`\`

### SSL/TLS Configuration

1. **Generate certificates**
   \`\`\`bash
   # Using Let's Encrypt
   certbot --nginx -d api.healthcare.example.com
   \`\`\`

2. **SSL configuration**
   \`\`\`bash
   # Strong SSL configuration
   ssl_protocols TLSv1.2 TLSv1.3;
   ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
   ssl_prefer_server_ciphers off;
   ssl_session_cache shared:SSL:10m;
   \`\`\`

## Troubleshooting

### Common Issues

1. **Database connection errors**
   \`\`\`bash
   # Check database status
   systemctl status postgresql
   
   # Check connectivity
   telnet $DB_HOST $DB_PORT
   
   # Check logs
   tail -f /var/log/postgresql/postgresql-*.log
   \`\`\`

2. **High memory usage**
   \`\`\`bash
   # Check memory usage
   free -h
   
   # Check application memory
   ps aux | grep healthcare-api
   
   # Adjust connection pool settings
   export DB_MAX_OPEN_CONNS=50
   export DB_MAX_IDLE_CONNS=10
   \`\`\`

3. **Performance issues**
   \`\`\`bash
   # Check database performance
   SELECT * FROM pg_stat_activity;
   
   # Check slow queries
   SELECT query, mean_time, calls 
   FROM pg_stat_statements 
   ORDER BY mean_time DESC LIMIT 10;
   \`\`\`

### Log Analysis

1. **Application logs**
   \`\`\`bash
   # View recent errors
   journalctl -u healthcare-api --since "1 hour ago" | grep ERROR
   
   # Follow logs
   journalctl -u healthcare-api -f
   \`\`\`

2. **Database logs**
   \`\`\`bash
   # PostgreSQL logs
   tail -f /var/log/postgresql/postgresql-*.log
   
   # Slow query log
   grep "duration:" /var/log/postgresql/postgresql-*.log
   \`\`\`

## Maintenance

### Regular Tasks

1. **Database maintenance**
   \`\`\`bash
   # Vacuum and analyze
   psql -c "VACUUM ANALYZE;"
   
   # Reindex
   psql -c "REINDEX DATABASE healthcare_db;"
   \`\`\`

2. **Log cleanup**
   \`\`\`bash
   # Clean old logs
   find /opt/healthcare-api/logs -name "*.log" -mtime +30 -delete
   
   # Clean journal logs
   journalctl --vacuum-time=30d
   \`\`\`

3. **Security updates**
   \`\`\`bash
   # Update system packages
   apt update && apt upgrade
   
   # Update Go dependencies
   go mod tidy && go mod download
   \`\`\`

This deployment guide provides comprehensive instructions for deploying the Healthcare API in various environments with proper security, monitoring, and maintenance procedures.
