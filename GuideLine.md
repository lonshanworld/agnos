# Agnos Hospital Management System - Complete Guide

## Table of Contents
1. [Project Overview](#project-overview)
2. [Architecture & Design Patterns](#architecture--design-patterns)
3. [Project Structure](#project-structure)
4. [Database Structure](#database-structure)
5. [Local Development Setup](#local-development-setup)
6. [Docker Setup](#docker-setup)
7. [Deployment Setup (Production VPS)](#deployment-setup-production-vps)
8. [Updating Application Code](#updating-application-code)
9. [Running Tests](#running-tests)
10. [Swagger Documentation](#swagger-documentation)
11. [Starting the Project](#starting-the-project)
12. [API Endpoints](#api-endpoints)
13. [Troubleshooting](#troubleshooting)
14. [Security Checklist](#security-checklist)

---

## Project Overview

A hospital management middleware system that allows staff to authenticate and search for patient information across multiple Hospital Information Systems (HIS).

**Tech Stack:**
- **Backend**: Go 1.24, Gin Web Framework
- **Database**: PostgreSQL 15
- **Authentication**: JWT (JSON Web Tokens)
- **Containerization**: Docker, Docker Compose
- **Reverse Proxy**: Nginx
- **Documentation**: Swagger/OpenAPI
- **Testing**: Go testing with testify/mock
- **Deployment**: Debian VPS with Cloudflare

---

## Architecture & Design Patterns

### Architecture: Clean Architecture (Layered Architecture)

This project follows **Clean Architecture** principles, also known as **Layered Architecture** or **Hexagonal Architecture**. This approach separates concerns into distinct layers, making the code:

- **Maintainable**: Changes in one layer don't affect others
- **Testable**: Each layer can be tested independently
- **Scalable**: Easy to add new features without breaking existing code
- **Database-agnostic**: Business logic is independent of data storage

### Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│                     Presentation Layer                   │
│              (Handlers - HTTP/REST API)                  │
│  • Receives HTTP requests                                │
│  • Validates input                                       │
│  • Returns HTTP responses                                │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                    Business Logic Layer                  │
│                      (Services)                          │
│  • Contains core business rules                          │
│  • Orchestrates data flow                                │
│  • Independent of frameworks                             │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                   Data Access Layer                      │
│                   (Repositories)                         │
│  • Database operations (CRUD)                            │
│  • Query building                                        │
│  • Data persistence                                      │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                    Database Layer                        │
│                  (PostgreSQL + GORM)                     │
│  • Data storage                                          │
│  • Schema management                                     │
└─────────────────────────────────────────────────────────┘
```

### Design Patterns Used

#### 1. **Repository Pattern**

**Why:** Abstracts data access logic from business logic.

**Implementation:**
- `repositories/` contains interfaces and implementations
- Each entity (Hospital, Staff, Patient) has its own repository
- Business logic (services) depends on repository interfaces, not concrete implementations

**Benefits:**
- Easy to mock repositories for testing
- Can swap database implementations without changing business logic
- Centralizes database query logic

**Example:**
```go
type HospitalRepositoryInterface interface {
    Create(hospital *models.Hospital) error
    FindByName(name string) (*models.Hospital, error)
}
```

#### 2. **Dependency Injection**

**Why:** Makes components loosely coupled and easier to test.

**Implementation:**
- Dependencies are passed through constructors ("constructor injection")
- Handlers depend on services
- Services depend on repositories
- No component creates its own dependencies

**Example:**
```go
func NewHospitalHandler(repo repositories.HospitalRepositoryInterface) *HospitalHandler {
    return &HospitalHandler{Repo: repo}
}
```

**Benefits:**
- Components can be tested in isolation with mocks
- Easy to replace implementations
- Clear dependency graph

#### 3. **Interface Segregation**

**Why:** Clients should not depend on interfaces they don't use.

**Implementation:**
- `repositories/interfaces.go` defines repository contracts
- `services/interfaces.go` defines service contracts
- Each interface is focused and minimal

**Benefits:**
- Easier to mock for testing
- Clear contracts between layers
- Prevents bloated interfaces

#### 4. **Middleware Pattern**

**Why:** Separates cross-cutting concerns (authentication, logging) from business logic.

**Implementation:**
- `middleware/auth.go` - JWT authentication
- `middleware/hospital_check.go` - Hospital validation
- Middleware functions wrap HTTP handlers

**Example:**
```go
api.GET("/patient/search", authMiddleWare, patientHandler.Search)
```

**Benefits:**
- Reusable across multiple endpoints
- Easy to add/remove middleware
- Keeps handlers focused on business logic

#### 5. **MVC-like Pattern (Handler-Service-Repository)**

**Why:** Separates concerns into distinct responsibilities.

**Layers:**
- **Handlers (Controllers)**: Handle HTTP requests/responses
- **Services**: Business logic and orchestration
- **Repositories (Models)**: Data access and persistence

**Benefits:**
- Clear separation of concerns
- Each layer has a single responsibility
- Easy to locate and modify code

### Project Organization Principles

#### 1. **Package by Feature/Layer**

**Why:** Makes it easy to find related code.

**Structure:**
```
handlers/    - HTTP request handling
services/    - Business logic
repositories/ - Data access
models/      - Data structures
middleware/  - Cross-cutting concerns
```

#### 2. **Configuration Management**

**Why:** Centralized configuration makes deployment easier.

**Implementation:**
- `config/config.go` loads environment variables
- `.env` file for local development
- Environment variables in Docker/production

#### 3. **Database Migration via Code**

**Why:** Schema is version-controlled and auto-applied.

**Implementation:**
- GORM's `AutoMigrate` in `database/postgres.go`
- Models define schema structure
- No separate migration files needed

### Why This Architecture?

#### Advantages for This Project:

1. **Testability**: Unit tests can mock repositories and services
2. **Maintainability**: Clear boundaries between layers
3. **Scalability**: Easy to add new entities (e.g., doctors, appointments)
4. **Team Collaboration**: Different developers can work on different layers
5. **Technology Independence**: Can swap Gin for another framework, or Postgres for MySQL
6. **API Documentation**: Swagger annotations live with handlers (presentation layer)
7. **Security**: Middleware handles authentication separately from business logic

#### Trade-offs:

1. **More Boilerplate**: More files and interfaces than a simple monolithic approach
2. **Learning Curve**: Team needs to understand the architecture
3. **Overkill for Small Projects**: For a 3-endpoint API, this might be over-engineered

**Justification for This Project:**

Given that this is a hospital management system that will likely grow (more entities, more complex queries, integration with external HIS systems), the upfront investment in clean architecture pays off through:
- Easy addition of new features
- Comprehensive testing
- Clear code organization
- Production-ready structure

---

## Project Structure

```
agnos_candidate_assignment/
├── main.go                      # Application entry point with Swagger config
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
├── Dockerfile                   # Multi-stage Docker build
├── docker-compose.yml           # Docker services orchestration
├── nginx.conf                   # Nginx reverse proxy config
├── .env                         # Environment variables (not in git)
├── .gitignore                   # Git ignore rules
├── README.md                    # Project README
├── GuideLine.md                 # This comprehensive guide
│
├── config/
│   └── config.go                # Configuration loader
│
├── database/
│   └── postgres.go              # Database connection & migrations
│
├── models/
│   ├── hospital.go              # Hospital model
│   ├── staff.go                 # Staff model
│   ├── patient.go               # Patient model
│   └── gender.go                # Gender enum
│
├── handlers/
│   ├── hospital_handler.go      # Hospital HTTP handlers
│   ├── staff_handler.go         # Staff authentication handlers
│   └── patient_handler.go       # Patient search handlers
│
├── services/
│   ├── interfaces.go            # Service interfaces
│   ├── auth_service.go          # Authentication business logic
│   └── patient_service.go       # Patient business logic
│
├── repositories/
│   ├── interfaces.go            # Repository interfaces
│   ├── hospital_repository.go   # Hospital data access
│   ├── staff_repository.go      # Staff data access
│   └── patient_repository.go    # Patient data access
│
├── middleware/
│   ├── auth.go                  # JWT authentication middleware
│   └── hospital_check.go        # Hospital validation middleware
│
├── scripts/
│   └── seed.go                  # Database seeding script
│
├── tests/
│   ├── hospital_handler_unit_test.go
│   ├── staff_handler_unit_test.go
│   └── patient_handler_unit_test.go
│
├── docs/                        # Auto-generated Swagger docs
│   ├── docs.go                  # (generated by swag init)
│   ├── swagger.json             # (generated by swag init)
│   └── swagger.yaml             # (generated by swag init)
│
└── utils/                       # Utility functions
```

---

## Database Structure

### Tables and Relationships

**ER Diagram:**
```
hospitals (1) ──< (N) staff
hospitals (1) ──< (N) patients
```

### 1. `hospitals` Table
| Column      | Type      | Constraints                  | Description             |
|-------------|-----------|------------------------------|-------------------------|
| id          | SERIAL    | PRIMARY KEY                  | Auto-increment ID       |
| name        | VARCHAR   | UNIQUE, NOT NULL             | Hospital name           |
| created_at  | TIMESTAMP | DEFAULT NOW()                | Record creation time    |
| updated_at  | TIMESTAMP | DEFAULT NOW()                | Last update time        |

### 2. `staffs` Table
| Column       | Type      | Constraints                  | Description             |
|--------------|-----------|------------------------------|-------------------------|
| id           | SERIAL    | PRIMARY KEY                  | Auto-increment ID       |
| hospital_id  | INTEGER   | FOREIGN KEY, NOT NULL, INDEX | Reference to hospitals  |
| user_name    | VARCHAR   | NOT NULL                     | Staff username          |
| password     | VARCHAR   | NOT NULL                     | Bcrypt hashed password  |
| name         | VARCHAR   |                              | Staff full name         |
| patient_hn   | VARCHAR   | UNIQUE                       | Staff HN number         |
| email        | VARCHAR   |                              | Staff email             |
| created_at   | TIMESTAMP | DEFAULT NOW()                | Record creation time    |
| updated_at   | TIMESTAMP | DEFAULT NOW()                | Last update time        |

### 3. `patients` Table
| Column         | Type      | Constraints                  | Description                |
|----------------|-----------|------------------------------|----------------------------|
| id             | SERIAL    | PRIMARY KEY                  | Auto-increment ID          |
| hospital_id    | INTEGER   | FOREIGN KEY, NOT NULL, INDEX | Reference to hospitals     |
| patient_hn     | VARCHAR   | UNIQUE, NOT NULL             | Hospital number            |
| national_id    | VARCHAR   | UNIQUE                       | National ID                |
| passport_id    | VARCHAR   | UNIQUE                       | Passport ID                |
| first_name_th  | VARCHAR   | NULLABLE                     | Thai first name            |
| middle_name_th | VARCHAR   | NULLABLE                     | Thai middle name           |
| last_name_th   | VARCHAR   | NULLABLE                     | Thai last name             |
| first_name_en  | VARCHAR   | NULLABLE                     | English first name         |
| middle_name_en | VARCHAR   | NULLABLE                     | English middle name        |
| last_name_en   | VARCHAR   | NULLABLE                     | English last name          |
| date_of_birth  | DATE      | NOT NULL                     | Date of birth              |
| phone_number   | VARCHAR   |                              | Contact phone              |
| email          | VARCHAR   |                              | Contact email              |
| gender         | CHAR(1)   | NOT NULL                     | M/F/O (Male/Female/Other)  |
| created_at     | TIMESTAMP | DEFAULT NOW()                | Record creation time       |
| updated_at     | TIMESTAMP | DEFAULT NOW()                | Last update time           |

**Note:** GORM automatically handles migrations. The database schema is defined in the `models/` directory.

---

## Local Development Setup

### Prerequisites
- Go 1.24 or higher
- PostgreSQL 15
- Git

### Steps

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd agnos_candidate_assignment
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Set up environment variables**
   
   Create a `.env` file in the project root:
   ```bash
   SERVER_PORT=8080
   GIN_MODE=debug
   JWT_SECRET=your-secret-key-change-this-in-production
   DATABASE_URL=postgresql://postgres:postgres@localhost:5432/agnosdb?sslmode=disable
   ```

4. **Set up PostgreSQL database**
   
   Create the database:
   ```bash
   psql -U postgres
   CREATE DATABASE agnosdb;
   \q
   ```

5. **Install Swagger CLI**
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

6. **Install Swagger dependencies**
   ```bash
   go get -u github.com/swaggo/gin-swagger
   go get -u github.com/swaggo/files
   go mod tidy
   ```

7. **Generate Swagger documentation**
   ```bash
   swag init
   ```

8. **Run the application**
   ```bash
   go run main.go
   ```

9. **Seed the database (optional)**
   ```bash
   go run ./scripts/seed.go
   ```

The server will start on `http://localhost:8080`.

---

## Docker Setup

### Prerequisites
- Docker Engine
- Docker Compose plugin

### Docker Architecture

The project uses a multi-stage Dockerfile and Docker Compose with three services:

1. **db** (PostgreSQL 15)
2. **app** (Go application)
3. **nginx** (Reverse proxy)

### Build and Run with Docker Compose

1. **Build images**
   ```bash
   docker compose build
   ```

2. **Start all services**
   ```bash
   docker compose up -d
   ```

3. **Check running services**
   ```bash
   docker compose ps
   ```

4. **View logs**
   ```bash
   docker compose logs -f app
   docker compose logs -f db
   docker compose logs -f nginx
   ```

5. **Stop services**
   ```bash
   docker compose down
   ```

6. **Stop and remove volumes (destructive)**
   ```bash
   docker compose down --volumes
   ```

### Running the Seeder in Docker

**Option 1: Using builder image (recommended)**

1. Build the builder stage:
   ```bash
   docker build --target builder -t agnos-builder .
   ```

2. Run the seeder:
   ```bash
   docker run --rm --network agnos_agnos-net --env-file .env agnos-builder sh -c "go run ./scripts/seed.go"
   ```

**Option 2: Using docker compose run**
```bash
docker compose run --rm --no-deps -e DATABASE_URL -e JWT_SECRET app sh -c "go run ./scripts/seed.go"
```

---

## Deployment Setup (Production VPS)

### Complete Step-by-Step Deployment on Debian VPS

#### Prerequisites
- Clean Debian 11+ VPS
- Root or sudo access
- Domain name (e.g., agnos.lonshan.com)
- Cloudflare account (for proxy and SSL)

### Step 1: Install Docker and Docker Compose

1. **Update packages**
   ```bash
   apt update
   apt upgrade -y
   ```

2. **Install prerequisites**
   ```bash
   apt install -y ca-certificates curl gnupg lsb-release
   ```

3. **Add Docker GPG key**
   ```bash
   mkdir -p /etc/apt/keyrings
   curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
   ```

4. **Add Docker repository**
   ```bash
   echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list
   ```

5. **Install Docker**
   ```bash
   apt update
   apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
   ```

6. **Verify installation**
   ```bash
   docker --version
   docker compose version
   ```

7. **Enable Docker service**
   ```bash
   systemctl enable --now docker
   ```

### Step 2: Transfer Project to VPS

**Option A: Using Git**
```bash
cd ~
git clone <your-repository-url> agnos
cd agnos
```

**Option B: Using SCP (from local machine)**
```bash
scp -r /local/path/to/agnos_candidate_assignment root@your-vps-ip:/root/agnos
```

### Step 3: Configure Environment Variables

Create `.env` file:
```bash
cd ~/agnos
cat > .env <<EOF
SERVER_PORT=8080
GIN_MODE=release
JWT_SECRET=$(openssl rand -hex 32)
DATABASE_URL=postgresql://postgres:postgres@db:5432/agnosdb?sslmode=disable
SILENCE_LOGS=true
EOF
```

### Step 4: Build Docker Images

1. **Build the builder image (for seeding)**
   ```bash
   docker build --target builder -t agnos-builder .
   ```

2. **Build and start services**
   ```bash
   docker compose up -d
   ```

3. **Check status**
   ```bash
   docker compose ps
   ```

### Step 5: Seed the Database

Run the seeder using the builder image:
```bash
docker run --rm --network agnos_agnos-net --env-file .env agnos-builder sh -c "go run ./scripts/seed.go"
```

Verify seeded data:
```bash
docker compose exec db psql -U postgres -d agnosdb -c "SELECT count(*) FROM hospitals;"
docker compose exec db psql -U postgres -d agnosdb -c "SELECT count(*) FROM staffs;"
docker compose exec db psql -U postgres -d agnosdb -c "SELECT count(*) FROM patients;"
```

Expected output: 3 hospitals, 12 staff, 100 patients.

### Step 6: Configure Cloudflare

1. **DNS Setup**
   - Go to Cloudflare Dashboard → DNS
   - Add an A record:
     - Name: `agnos` (or `@` for root domain)
     - Content: Your VPS IP address
     - Proxy status: **Proxied (orange cloud)**

2. **SSL/TLS Configuration**
   - Go to SSL/TLS → Overview
   - Set encryption mode to **Flexible**
   - (For production security, use Cloudflare Origin Certificate + Full (strict))

3. **Edge Certificates**
   - Go to SSL/TLS → Edge Certificates
   - Enable **Always Use HTTPS**
   - Enable **Automatic HTTPS Rewrites** (optional)

### Step 7: Verify Deployment

1. **Check services**
   ```bash
   docker compose ps
   ```
   All services should show "Up" status.

2. **Test local HTTP**
   ```bash
   curl http://localhost/
   curl http://localhost/api/health
   ```

3. **Test via domain**
   ```bash
   curl https://agnos.lonshan.com/
   curl https://agnos.lonshan.com/api/health
   ```

4. **Check Swagger UI**
   - Browser: https://agnos.lonshan.com/swagger/index.html

### Step 8: Ongoing Maintenance

**View logs**
```bash
docker compose logs -f --tail=200 app
```

**Restart services**
```bash
docker compose restart app
docker compose restart nginx
```

**Update application**
```bash
git pull origin main
docker compose build --no-cache app
docker compose up -d --no-deps --build app
```

**Cleanup Docker resources**
```bash
docker builder prune --all --force
docker image prune -f
docker container prune -f
docker network prune -f
```

**Database backup**
```bash
docker compose exec db pg_dump -U postgres agnosdb > backup_$(date +%Y%m%d).sql
```

---

## Updating Application Code

### When to Update

Update the application when:
- Code changes are pushed to the repository
- Dependencies are updated in `go.mod`
- Configuration changes are made
- Bug fixes or new features are deployed

### Update Process (Production/Docker)

#### Method 1: Full Rebuild with No Cache (Recommended for major changes)

**When to use:** Major code changes, dependency updates, or when you want to ensure a clean build.

```bash
cd ~/agnos

# Pull latest code (if using Git)
git pull origin main

# Stop the app service
docker compose stop app

# Rebuild the app image without cache
docker compose build --no-cache app

# Start the updated app
docker compose up -d app

# Verify it's running
docker compose ps
docker compose logs -f --tail=200 app
```

#### Method 2: Quick Update with Cache (For minor changes)

**When to use:** Small code changes where you want faster builds.

```bash
cd ~/agnos

# Pull latest code
git pull origin main

# Rebuild and restart app
docker compose up -d --build app

# Or separate steps:
docker compose build app
docker compose up -d app
```

#### Method 3: Update All Services

**When to use:** Changes to docker-compose.yml, nginx.conf, or multiple services.

```bash
cd ~/agnos

# Pull latest code
git pull origin main

# Rebuild all services
docker compose build --no-cache

# Restart all services
docker compose up -d

# Verify
docker compose ps
```

### Update Checklist

- [ ] **Backup database** (if schema changes are involved)
  ```bash
  docker compose exec db pg_dump -U postgres agnosdb > backup_pre_update_$(date +%Y%m%d_%H%M%S).sql
  ```

- [ ] **Pull latest code**
  ```bash
  git pull origin main
  ```

- [ ] **Update environment variables** (if needed)
  ```bash
  nano .env
  ```

- [ ] **Regenerate Swagger docs** (if API changes)
  ```bash
  # Run locally before committing, or in builder container
  swag init
  ```

- [ ] **Rebuild images**
  ```bash
  docker compose build --no-cache app
  ```

- [ ] **Stop old container**
  ```bash
  docker compose stop app
  ```

- [ ] **Start new container**
  ```bash
  docker compose up -d app
  ```

- [ ] **Check logs for errors**
  ```bash
  docker compose logs -f --tail=200 app
  ```

- [ ] **Test endpoints**
  ```bash
  curl https://agnos.lonshan.com/api/health
  curl https://agnos.lonshan.com/swagger/index.html
  ```

- [ ] **Monitor for issues**
  ```bash
  docker compose logs -f app
  ```

### Zero-Downtime Updates (Advanced)

For production systems requiring zero downtime:

1. **Build new image with different tag**
   ```bash
   docker build -t agnos-app:v2 .
   ```

2. **Start new container alongside old one** (on different port)
   ```bash
   docker run -d --name agnos-app-v2 --network agnos_agnos-net -p 8081:8080 --env-file .env agnos-app:v2
   ```

3. **Test new version**
   ```bash
   curl http://localhost:8081/api/health
   ```

4. **Update nginx to point to new container** (update upstream in nginx.conf)

5. **Reload nginx**
   ```bash
   docker compose exec nginx nginx -s reload
   ```

6. **Remove old container**
   ```bash
   docker stop agnos-app-1
   docker rm agnos-app-1
   ```

### Rollback Procedure

If the update causes issues:

**Quick rollback:**
```bash
# Stop current version
docker compose stop app

# Revert code
git reset --hard HEAD~1  # or git checkout <previous-commit>

# Rebuild previous version
docker compose build --no-cache app

# Start previous version
docker compose up -d app
```

**Restore from backup (if database changes):**
```bash
# Stop app to prevent new writes
docker compose stop app

# Restore database
cat backup_pre_update_<timestamp>.sql | docker compose exec -T db psql -U postgres agnosdb

# Start app
docker compose up -d app
```

### Common Update Scenarios

#### Scenario 1: Code Changes Only
```bash
git pull origin main
docker compose build --no-cache app
docker compose up -d app
```

#### Scenario 2: Dependency Updates (go.mod)
```bash
git pull origin main
# Rebuild forces new go mod download
docker compose build --no-cache app
docker compose up -d app
```

#### Scenario 3: Database Schema Changes
```bash
# Backup first!
docker compose exec db pg_dump -U postgres agnosdb > backup.sql

git pull origin main
docker compose build --no-cache app
docker compose up -d app
# GORM AutoMigrate will update schema automatically
```

#### Scenario 4: Environment Variable Changes
```bash
# Edit .env
nano .env

# Recreate container with new env vars
docker compose up -d --force-recreate app
```

#### Scenario 5: Nginx Configuration Changes
```bash
git pull origin main
# nginx.conf is mounted as volume, just reload
docker compose exec nginx nginx -s reload

# Or restart nginx service
docker compose restart nginx
```

### Monitoring After Update

**Watch logs in real-time:**
```bash
docker compose logs -f app
```

**Check resource usage:**
```bash
docker stats
```

**Test all endpoints:**
```bash
curl -I https://agnos.lonshan.com/
curl -I https://agnos.lonshan.com/api/health
curl https://agnos.lonshan.com/swagger/index.html
```

**Check database connectivity:**
```bash
docker compose exec db psql -U postgres agnosdb -c "SELECT COUNT(*) FROM patients;"
```

### Cleanup After Update

**Remove old/unused images:**
```bash
docker image prune -f
```

**Remove dangling images:**
```bash
docker image prune -a -f
```

**Full cleanup (careful!):**
```bash
docker system prune -a --volumes -f
```

---

## Running Tests

### Unit Tests

The project includes unit tests for all handlers using mock repositories and services.

**Run all tests:**
```bash
go test ./tests -v
```

**Run specific test file:**
```bash
go test ./tests/hospital_handler_unit_test.go -v
go test ./tests/staff_handler_unit_test.go -v
go test ./tests/patient_handler_unit_test.go -v
```

**Run tests with coverage:**
```bash
go test ./... -cover
```

**Generate coverage report:**
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

**Test output example:**
```
=== RUN   TestHospitalCreate
--- PASS: TestHospitalCreate (0.00s)
=== RUN   TestHospitalCreate_BadRequestBinding
--- PASS: TestHospitalCreate_BadRequestBinding (0.00s)
=== RUN   TestStaffRegister
--- PASS: TestStaffRegister (0.00s)
PASS
ok      agnos_candidate_assignment/tests    0.123s
```

---

## Swagger Documentation

### Overview

Swagger/OpenAPI provides interactive API documentation accessible via web browser.

### Installation

1. **Install swag CLI**
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. **Install dependencies**
   ```bash
   go get -u github.com/swaggo/gin-swagger
   go get -u github.com/swaggo/files
   go mod tidy
   ```

### Generate Documentation

Run this command in the project root:
```bash
swag init
```

This generates:
- `docs/docs.go`
- `docs/swagger.json`
- `docs/swagger.yaml`

### Access Swagger UI

**Local development:**
```
http://localhost:8080/swagger/index.html
```

**Production:**
```
https://agnos.lonshan.com/swagger/index.html
```

### Swagger Annotations

Annotations are added as comments above handlers and in `main.go`:

**General API info (main.go):**
```go
// @title           Agnos Hospital API
// @version         1.0
// @description     Hospital management system API
// @host            agnos.lonshan.com
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

**Endpoint example (handler):**
```go
// Create godoc
// @Summary      Create a new hospital
// @Description  Register a new hospital in the system
// @Tags         hospitals
// @Accept       json
// @Produce      json
// @Param        request body createHospitalRequest true "Hospital creation request"
// @Success      201  {object}  models.Hospital
// @Failure      400  {object}  map[string]string
// @Router       /hospital [post]
func (h *HospitalHandler) Create(c *gin.Context) { ... }
```

### Regenerate After Changes

After adding or modifying endpoints:
```bash
swag init
```

---

## Starting the Project

### Local Development

**With live reload (using air):**
```bash
go install github.com/cosmtrek/air@latest
air
```

**Standard run:**
```bash
go run main.go
```

### Docker

**Start all services:**
```bash
docker compose up -d
```

**Watch logs:**
```bash
docker compose logs -f app
```

### Production (VPS)

**Start services:**
```bash
cd ~/agnos
docker compose up -d
```

**Verify:**
```bash
docker compose ps
curl https://agnos.lonshan.com/api/health
```

---

## API Endpoints

### Base URL
- **Local**: `http://localhost:8080/api`
- **Production**: `https://agnos.lonshan.com/api`

### Public Endpoints

#### 1. Create Hospital
```http
POST /api/hospital
Content-Type: application/json

{
  "name": "General Hospital"
}
```

**Response (201):**
```json
{
  "id": 1,
  "name": "General Hospital",
  "created_at": "2026-01-05T10:00:00Z",
  "updated_at": "2026-01-05T10:00:00Z"
}
```

#### 2. Staff Registration
```http
POST /api/:hospital/staff/create
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

**Response (201):**
```json
{
  "staff_id": 1,
  "username": "admin"
}
```

#### 3. Staff Login
```http
POST /api/:hospital/staff/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "staff_id": 1,
  "username": "admin",
  "hospital_id": 1
}
```

#### 4. Get Patient by ID (Public)
```http
GET /api/:hospital/patient/search/:id
```

**Response (200):**
```json
{
  "id": 1,
  "hospital_id": 1,
  "patient_hn": "HN00001",
  "national_id": "1234567890123",
  "first_name_th": "สมชาย",
  "last_name_th": "ใจดี",
  "first_name_en": "Somchai",
  "last_name_en": "Jaidee",
  "date_of_birth": "1990-01-01",
  "gender": "M",
  ...
}
```

### Protected Endpoints (Require JWT)

#### 5. Search Patients
```http
GET /api/patient/search?first_name=John&last_name=Doe
Authorization: Bearer <JWT_TOKEN>
```

**Query Parameters:**
- `national_id` - National ID
- `passport_id` - Passport ID
- `first_name` - First name (any language)
- `middle_name` - Middle name (any language)
- `last_name` - Last name (any language)
- `first_name_th` - First name (Thai)
- `middle_name_th` - Middle name (Thai)
- `last_name_th` - Last name (Thai)
- `first_name_en` - First name (English)
- `middle_name_en` - Middle name (English)
- `last_name_en` - Last name (English)
- `date_of_birth` - Date of birth
- `phone_number` - Phone number
- `email` - Email

**Response (200):**
```json
{
  "patients": [
    {
      "id": 1,
      "hospital_id": 1,
      "patient_hn": "HN00001",
      "first_name_en": "John",
      "last_name_en": "Doe",
      ...
    }
  ]
}
```

### Authentication

Protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

To obtain a token, use the `/api/:hospital/staff/login` endpoint.

---

## Troubleshooting

### Common Issues

**1. Port 8080 already in use**
```bash
# Find process using port 8080
lsof -i :8080
# Kill the process
kill -9 <PID>
```

**2. Docker permission denied**
```bash
sudo usermod -aG docker $USER
newgrp docker
```

**3. Database connection refused**
- Check if PostgreSQL is running: `docker compose ps`
- Verify DATABASE_URL in `.env`
- Check firewall rules

**4. Cloudflare 521 error**
- Ensure services are running: `docker compose ps`
- Set Cloudflare SSL mode to Flexible
- Check nginx configuration
- Verify port 80 is accessible

**5. Swagger not loading**
- Run `swag init` to generate docs
- Check `docs/` folder exists
- Verify import in main.go: `_ "agnos_candidate_assignment/docs"`

---

## Security Notes

### Production Security Checklist

- [ ] Change `JWT_SECRET` to a strong random value
- [ ] Use strong database passwords
- [ ] Enable Cloudflare Origin Certificate (upgrade from Flexible to Full (strict))
- [ ] Set `GIN_MODE=release` in production
- [ ] Enable `SILENCE_LOGS=true` to reduce log verbosity
- [ ] Configure firewall to only allow ports 80, 443, and SSH
- [ ] Regularly update Docker images and dependencies
- [ ] Set up database backups
- [ ] Monitor logs for suspicious activity
- [ ] Use HTTPS for all communications
- [ ] Implement rate limiting (add middleware if needed)

---

## Support and Contact

For issues or questions, please refer to:
- Project README.md
- Swagger documentation at `/swagger/index.html`
- Repository issues page

---

**Last Updated:** January 5, 2026
**Version:** 1.0.0
