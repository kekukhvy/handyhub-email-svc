# HandyHub Email Service

A microservice for asynchronous email message processing from RabbitMQ queues in the HandyHub ecosystem.

## 📋 Description

HandyHub Email Service is a **MQ Consumer service** for processing email messages with minimal REST API for monitoring. Built with Go using:
- **RabbitMQ Consumer** for reading email message queues
- **Minimal REST API** (Gin) only for health-check and testing
- **SMTP Client** for sending emails

Core functionality:

- **Reading messages** from RabbitMQ queues
- **Asynchronous processing** of email messages
- Sending emails via SMTP
- Logging all operations with flexible storage strategies
- Minimal REST API for service status monitoring

## 🏗️ Architecture

### Core Components:

1. **RabbitMQ Consumer** - main component for reading email queues
2. **SMTP Client** - sending emails via SMTP server
3. **Storage Layer** - abstraction for log storage (console/file/MongoDB)
4. **Minimal REST API** (Gin) - only for health-check and testing
5. **Configuration Management** (Viper) - configuration management
6. **Database Layer** (MongoDB) - optional persistent storage
7. **Logging** (Logrus) - structured logging

### Project Structure:

```
handyhub-email-svc/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/                 # Configuration
│   │   ├── config.go
│   │   └── cfg.yml
│   ├── database/               # MongoDB connection
│   │   └── mongodb.go
│   ├── logger/                 # Custom logging
│   │   ├── formatter.go
│   │   └── init.go
│   ├── models/                 # Data models
│   │   └── email_log.go
│   ├── server/                 # HTTP server and routes
│   │   ├── routes.go
│   │   └── server.go
│   └── storage/                # Storage abstraction
│       ├── console.go
│       ├── database.go
│       ├── factory.go
│       ├── file.go
│       └── interface.go
├── logs/                       # Logs directory
├── go.mod
├── go.sum
└── README.md
```

## 🔧 Integrations and Dependencies

### Main Dependencies:

- **RabbitMQ (AMQP)** - main dependency for reading message queues
- **Gin** - minimal HTTP web framework (monitoring only)
- **MongoDB Driver** - optional for database storage
- **Viper** - configuration management
- **Logrus** - structured logging
- **Gomail** (gopkg.in/gomail.v2) - SMTP client for email sending

### Supported Storage Types:

1. **Console Storage** - output logs to console
2. **File Storage** - save to JSON files
3. **Database Storage** - save to MongoDB

### Queue Message Data Model:

```go
type QueueMessage struct {
    Email     EmailMessage      `json:"email"`
    Priority  string            `json:"priority"`
    Metadata  map[string]string `json:"metadata,omitempty"`
    Timestamp time.Time         `json:"timestamp" bson:"timestamp"`
}

type EmailMessage struct {
    To      []string `json:"to"`
    Subject string   `json:"subject"`
    Body    string   `json:"body"`
    // additional fields...
}
```

### Email Log Data Model:

```go
type EmailLog struct {
    ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    To       []string           `json:"to" bson:"to"`
    Subject  string             `json:"subject" bson:"subject"`
    Status   string             `json:"status" bson:"status"`
    Provider string             `json:"provider" bson:"provider"`
    Attempts int                `json:"attempts" bson:"attempts"`
    SentAt   time.Time          `json:"sent_at" bson:"sent_at"`
    ErrorMsg string             `json:"error_msg,omitempty" bson:"error_msg,omitempty"`
}
```

## 🚀 Local Setup

### Prerequisites:

1. **Go 1.25+**
2. **RabbitMQ** (main dependency for queues)
3. **SMTP server** (e.g., MailHog for testing)
4. **MongoDB** (optional, only for database storage)

### 1. Clone and install dependencies:

```bash
git clone <repository-url>
cd handyhub-email-svc
go mod download
```

### 2. Environment setup:

Create `.env` file or set environment variables:

```bash
export MONGODB_URL="mongodb://localhost:27017"
export DB_NAME="handyhub"
export RABBITMQ_URL="amqp://admin:admin123@localhost:5672/"
```

### 3. Infrastructure setup:

#### RabbitMQ:
```bash
# Using Docker
docker run -d --name rabbitmq \
  -p 5672:5672 \
  -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=admin \
  -e RABBITMQ_DEFAULT_PASS=admin123 \
  rabbitmq:3-management

# Web UI available at http://localhost:15672
# Login: admin, Password: admin123
```

#### MongoDB (only for database storage):
```bash
# Using Docker
docker run -d --name mongodb -p 27017:27017 mongo:latest

# Or install MongoDB locally
# https://docs.mongodb.com/manual/installation/
```

### 4. Start SMTP server for testing (MailHog):

```bash
# Using Docker
docker run -d --name mailhog -p 1025:1025 -p 8025:8025 mailhog/mailhog

# Web UI available at http://localhost:8025
```

### 5. Create logs directory:

```bash
mkdir -p logs
```

### 6. Start the service:

```bash
go run cmd/main.go
```

Service will start on port `:8008`

## 🧪 Testing

### 1. Infrastructure check:

#### RabbitMQ Management UI:
```bash
# Open in browser
http://localhost:15672
# Login: admin, Password: admin123
```

### 2. Service health check (minimal REST API):

```bash
curl http://localhost:8008/health
```

Expected response:
```json
{
  "status": "ok",
  "service": "email-service"
}
```

### 3. API status check:

```bash
curl http://localhost:8008/api/v1/status
```

Expected response:
```json
{
  "api_version": "v1",
  "status": "operational"
}
```

### 4. RabbitMQ message testing:

#### Create exchange and send message:

```bash
# Using RabbitMQ Management UI (http://localhost:15672):
# 1. Go to "Exchanges" section
# 2. Create exchange named "email.exchange" (type: direct)
# 3. In "Publish message" section send JSON:
```

**Example message to send to exchange:**
```json
{
  "email": {
    "to": ["recipient@example.com"],
    "subject": "Test Email from RabbitMQ",
    "body": "This is a test email sent through RabbitMQ queue"
  },
  "priority": "high",
  "metadata": {
    "template": "default",
    "user_id": "12345"
  },
  "timestamp": "2025-09-04T10:30:00Z"
}
```

#### Using RabbitMQ Management Web UI:
```bash
# 1. Open http://localhost:15672
# 2. Login (admin/admin123)
# 3. Go to "Exchanges" section
# 4. Find exchange "email.exchange" (create if not exists)
# 5. Click on exchange name
# 6. In "Publish message" section specify routing key: "email.send"
# 7. Paste JSON above into Message payload field
# 8. Click "Publish message"
# 9. Check in "Queues" section that message reached the queue
```

### 5. Minimal REST API testing (monitoring only):

```bash
curl -X POST http://localhost:8008/api/v1/test-email-log \
  -H "Content-Type: application/json"
```

Expected response:
```json
{
  "message": "Test email log stored",
  "storage_type": "database",
  "log": {
    "id": "...",
    "to": ["test@example.com"],
    "subject": "Test Email Log",
    "status": "success",
    "provider": "test",
    "attempts": 1,
    "sent_at": "2025-09-04T10:30:00Z",
    "error_msg": ""
  }
}
```

### 6. Logs and storage verification:

Depending on configured storage type:

- **Console**: Logs will be output to console
- **File**: Check file `logs/emails`
- **Database**: Connect to MongoDB and check `emails` collection

```bash
# For MongoDB
mongo
use handyhub
db.emails.find().pretty()
```

### 7. Full RabbitMQ → Email Service integration check:

#### Step 1: Verify service is running and connected to RabbitMQ
```bash
# Service logs should contain messages like:
# "Connected to RabbitMQ"
# "Listening for messages on queue: email.send"
```

#### Step 2: Send test message to exchange
```bash
# Through RabbitMQ Management UI (http://localhost:15672):
# 1. Exchanges → email.exchange → Publish message
# 2. Routing key: email.send
# 3. Payload:
{
  "email": {
    "to": ["test@example.com"],
    "subject": "Integration Test",
    "body": "This email tests the full RabbitMQ → Email Service → SMTP integration"
  },
  "priority": "high",
  "metadata": {
    "test_type": "integration"
  },
  "timestamp": "2025-09-04T10:30:00Z"
}
```

#### Step 3: Real-time processing monitoring
```bash
# Terminal 1: Service logs
tail -f logs/sys.log

# Terminal 2: RabbitMQ queue monitoring
# Through Web UI: Queues → email.send → refresh page

# Terminal 3: (If using database storage) MongoDB
mongo handyhub --eval "db.emails.find().limit(5).sort({sent_at:-1})"
```

#### Step 4: Results verification
```bash
# 1. Service logs should show:
#    - "Message received from queue"
#    - "Email sent successfully" or "Email send failed"
#    - "Email log stored"

# 2. In MailHog UI (http://localhost:8025):
#    - Sent email should appear

# 3. RabbitMQ queue should be empty (messages = 0)

# 4. (If database storage) MongoDB should have new record:
mongo handyhub --eval "db.emails.findOne({subject: 'Integration Test'})"
```

#### Step 5: Error handling test
```bash
# Stop MailHog and send message:
docker stop mailhog

# Send message - should get SMTP error
# Through Web UI send:
{
  "email": {
    "to": ["error-test@example.com"],
    "subject": "Error Test",
    "body": "This should fail"
  },
  "priority": "high",
  "timestamp": "2025-09-04T10:30:00Z"
}

# Logs should show:
# - "SMTP connection failed"
# - "Email log stored with error status"

# Start MailHog back:
docker start mailhog
```

## ⚙️ Configuration

### Main parameters in `internal/config/cfg.yml`:

```yaml
database:
  url: "mongodb://localhost:27017"
  dbname: "handyhub"
  email-collection: "emails"
  timeout: 10

storage:
  type: database  # console | file | database
  file:
    path: "logs/emails"
    max-size: 10
    max-files: 5

rabbitmq:
  url: "amqp://guest:guest@localhost:5672/"
  exchange: "email.exchange"
  queue: "email.send"
  routing-key: "email.send"

server:
  port: ":8008"
  read-timeout: 30
  write-timeout: 30
  idle-timeout: 60

app:
  name: "handyhub-email-svc"
  timeout: 30

logs:
  level: "info"
  log-path: "logs/sys.log"

email:
  smtp-host: "localhost"
  smtp-port: 1025
  smtp-user: ""
  smtp-password: ""
  from-email: "noreply@handyhub.com"
  from-name: "HandyHub"
```

### Environment Variables:

- `MONGODB_URL` - MongoDB connection URL
- `DB_NAME` - Database name
- `RABBITMQ_URL` - RabbitMQ connection URL

## 🔍 Monitoring and Logging

Service uses structured logging with color highlighting:

- **DEBUG/TRACE**: Cyan
- **INFO**: Green
- **WARN**: Yellow
- **ERROR/FATAL/PANIC**: Red

Logs are written simultaneously to:
- Standard output (stdout)
- Log file (configurable in configuration)

## 🐳 Docker

To containerize, create `Dockerfile`:

```dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/internal/config ./internal/config

CMD ["./main"]
```

## 📝 API Documentation

### Endpoints (monitoring only):

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/health` | Service health check |
| GET    | `/api/v1/status` | API status |
| POST   | `/api/v1/test-email-log` | Test email log creation |

> **Note:** The main functionality of the service is processing messages from RabbitMQ, not REST API.

## 🤝 Development

### Adding new storage types:

1. Create new file in `internal/storage/`
2. Implement `EmailStorage` interface
3. Add new type to `factory.go`

### Adding new message handlers:

1. Create new handler in `internal/handlers/` package
2. Register handler for specific message types
3. Add routing to consumer

## 📊 Performance

- **Asynchronous processing** through RabbitMQ consumer
- **Goroutines** for parallel message processing
- Graceful shutdown with 5-second timeout
- Configurable HTTP server timeouts (monitoring only)
- MongoDB connection pool with configurable timeout
- **Reconnection** to RabbitMQ on connection loss