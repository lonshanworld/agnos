
# Hospital Middleware System

A small middleware service to search and display patient information from multiple Hospital Information Systems (HIS). Built with Go and Gin, using PostgreSQL for storage and Docker for local development.

## Project Overview

- Purpose: Allow hospital staff to authenticate and search for patients belonging to the staff's hospital. Patient data may be fetched from external HIS APIs and stored or cached in our database.
- Tech: Go, Gin, PostgreSQL, Docker, Nginx, JWT for auth, Swagger/OpenAPI documentation

## API Documentation

**Swagger UI**: Available at `/swagger/index.html` when the server is running.

Example: `http://localhost:8080/swagger/index.html` or `https://agnos.lonshan.com/swagger/index.html`

## Project Structure

```
agnos_candidate_assignment/
├── main.go                  # Application entry point
├── models/                  # Database models (Patient, Staff, Hospital)
├── handlers/                # HTTP handlers (controllers)
├── services/                # Business logic
├── repositories/            # Database operations
├── middleware/              # Authentication middleware
├── config/                  # Configuration helpers
├── database/                # DB connection and migrations
├── utils/                   # Utility helpers
├── tests/                   # Unit tests
├── docker-compose.yml
├── Dockerfile
├── nginx.conf
├── .env.example
└── README.md
```

## API Endpoints

Base path: `/api/v1`

- `POST /api/v1/staff/create` — create staff user
- `POST /api/v1/staff/login` — staff login (returns JWT)
- `GET /api/v1/patient/search` — protected; search patients by query params

### `/api/v1/staff/create`
- Input JSON: `{ "username": "u", "password": "p", "hospital": "Hospital A" }`
- Response: success message and created staff id

### `/api/v1/staff/login`
- Input JSON: `{ "username": "u", "password": "p", "hospital": "Hospital A" }`
- Response: `{ "token": "<jwt>", "staff_id": 1, "hospital": "Hospital A" }`

### `/api/v1/patient/search`
- Requires `Authorization: Bearer <token>` header
- Query params (all optional): `national_id`, `passport_id`, `first_name`, `middle_name`, `last_name`, `date_of_birth`, `phone_number`, `email`
- Response: list of patients belonging to the staff's hospital matching criteria

## Database Model (high level)

- `hospitals` : id, name, api_url, created_at, updated_at
- `staff` : id, username, password_hash, hospital_id, created_at
- `patients` : id, hospital_id, patient_hn, national_id, passport_id, first_name_th, middle_name_th, last_name_th, first_name_en, middle_name_en, last_name_en, date_of_birth, phone_number, email, gender, created_at

ER note: `hospitals` 1 - N `staff`; `hospitals` 1 - N `patients`.

## Setup (local)

1. Copy `.env.example` to `.env` and adjust values.
2. Start services with Docker Compose:

```bash
docker-compose up --build
```

3. Or run locally (requires Postgres running):

```bash
go mod download
go run main.go
```

## Swagger Documentation Setup

### Install swag CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Install dependencies

```bash
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
go mod tidy
```

### Generate Swagger docs

Run this command in the project root to auto-generate Swagger documentation from code annotations:

```bash
swag init
```

This creates/updates the `docs/` folder with `docs.go`, `swagger.json`, and `swagger.yaml`.

### Access Swagger UI

Start the server and navigate to:
- Local: `http://localhost:8080/swagger/index.html`
- Production: `https://agnos.lonshan.com/swagger/index.html`

### Regenerate docs after changes

After adding or modifying API endpoints, run `swag init` again to update the documentation.

## Testing

Run unit tests:

```bash
go test ./... -v
```

## Next steps for implementation

1. Add `.env.example` and configuration loader in `config/`.
2. Implement DB connection and migrations in `database/`.
3. Create models in `models/` and repositories.
4. Implement JWT auth middleware and handlers for staff and patient APIs.
5. Add unit tests under `tests/`.


