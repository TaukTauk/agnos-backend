# Agnos Backend — Hospital Middleware System

A RESTful API middleware system built with Go, Gin, PostgreSQL, Docker, and Nginx that enables hospital staff to search patient information from Hospital Information Systems (HIS).

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL 16 |
| Auth | JWT (httpOnly cookie) + API Key (header) |
| Containerization | Docker + Docker Compose |
| Reverse Proxy | Nginx |
| Testing | Testify (assert + mock) |

---

## Project Structure

```
agnos-backend/
├── cmd/
│   └── main.go                  # Entry point, router, dependency wiring
├── internal/
│   ├── handler/                 # HTTP layer — parse request, return response
│   │   ├── staff_handler.go
│   │   └── patient_handler.go
│   ├── service/                 # Business logic layer
│   │   ├── staff_service.go
│   │   └── patient_service.go
│   ├── repository/              # Database access layer
│   │   ├── hospital_repo.go
│   │   ├── staff_repo.go
│   │   └── patient_repo.go
│   ├── model/                   # GORM models (DB schema)
│   │   ├── hospital.go
│   │   ├── staff.go
│   │   └── patient.go
│   ├── dto/                     # Request/Response shapes
│   │   ├── staff_dto.go
│   │   ├── patient_dto.go
│   │   └── common.go
│   └── middleware/              # API Key, JWT auth, Rate limiting
│       ├── auth.go
│       └── rate_limit.go
├── config/
│   └── config.go                # DB connection, migrations
├── migrations/
│   └── seed.sql                 # DB schema + seed data
├── nginx/
│   └── nginx.conf               # Reverse proxy config
├── tests/
│   ├── mocks/                   # Repository mocks for unit tests
│   └── services/                # Service layer unit tests
├── .env.example                 # Environment variable template
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

---

## Prerequisites

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Go 1.25+](https://go.dev/dl/) (only needed for running tests locally)

---

## Setup & Running

### 1. Clone the repository

```bash
git clone https://github.com/TaukTauk/agnos-backend.git
cd agnos-backend
```

### 2. Configure environment variables

```bash
cp .env.example .env
```

Edit `.env` with your values — see [Environment Variables](#environment-variables) below.

### 3. Start all services

```bash
docker compose up --build -d
```

This starts:
- **PostgreSQL** on port `5434` (mapped from internal 5432)
- **Go app** on port `8080`
- **Nginx** on port `80` (main entry point)

### 4. Fresh start (reset all data)

```bash
docker compose down -v
docker compose up --build -d
```

The `-v` flag removes the database volume, causing the seed to re-run.

---

## Environment Variables

| Variable | Description | Default |
|---|---|---|
| `APP_PORT` | Go app port | `8080` |
| `GIN_MODE` | Gin mode (`debug` or `release`) | `debug` |
| `API_KEY` | API key for all requests | — |
| `DB_HOST` | PostgreSQL host | `postgres` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL user | — |
| `DB_PASSWORD` | PostgreSQL password | — |
| `DB_NAME` | PostgreSQL database name | — |
| `JWT_SECRET` | Secret key for JWT signing | — |
| `JWT_EXPIRY_HOURS` | JWT expiry in hours | `24` |
| `CORS_ORIGIN` | Allowed CORS origins (comma-separated) | — |
| `RATE_LIMIT_LOGIN` | Max login requests per minute per IP | `5` |
| `RATE_LIMIT_CREATE_STAFF` | Max create staff requests per minute per IP | `10` |
| `RATE_LIMIT_PATIENT_SEARCH` | Max patient search requests per minute per IP | `30` |

---

## API Endpoints

All endpoints require the `X-API-Key` header.

### Public

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Health check |
| `POST` | `/staff/create` | Create a new staff member |
| `POST` | `/staff/login` | Login and receive JWT cookie |

### Protected (JWT required)

| Method | Path | Description |
|---|---|---|
| `POST` | `/staff/logout` | Logout and invalidate JWT |
| `GET` | `/patient/search` | Search patients (hospital-isolated) |

### Patient Search Query Parameters

All fields are optional:

```
national_id, passport_id, first_name, middle_name, last_name,
date_of_birth, phone_number, email, page, page_size
```

---

## Seed Data

The database is pre-seeded with:

| Entity | Data |
|---|---|
| Hospitals | Hospital A (`HOSPITAL_A`), Hospital B (`HOSPITAL_B`) |
| Staff | `staff.hospital.a` / `staff.hospital.b` — password: `password` |
| Patients | 5 patients in Hospital A, 2 patients in Hospital B |

---

## Running Tests

```bash
# Run all tests
go test ./tests/... -v

# Run with coverage
go test ./tests/... -v -cover
```

### Test Coverage

| Test | Type |
|---|---|
| `TestCreateStaff_Success` | Positive |
| `TestCreateStaff_HospitalNotFound` | Negative |
| `TestCreateStaff_DuplicateUsername` | Negative |
| `TestCreateStaff_WeakPassword` | Negative |
| `TestLogin_Success` | Positive |
| `TestLogin_HospitalNotFound` | Negative |
| `TestLogin_StaffNotFound` | Negative |
| `TestLogin_WrongPassword` | Negative |
| `TestSearchPatient_Success` | Positive |
| `TestSearchPatient_NoResults` | Negative |
| `TestSearchPatient_DefaultPagination` | Edge case |
| `TestSearchPatient_RepositoryError` | Negative |
| `TestSearchPatient_PaginationCalculation` | Edge case |

---

## Postman Collection

Import `agnos-backend.postman_collection.json` into Postman.

1. Set the `api_key` collection variable to match your `.env` `API_KEY`
2. Set `base_url` to `http://localhost` for `nginx` (better for production) or `http://localhost:8080` to go directly to `Go app`
3. Run **Login** first — JWT cookie is saved automatically
4. All protected requests will use the cookie automatically

---

## Security Features

- **Two-layer auth** — API key (app-level) + JWT (user-level)
- **httpOnly JWT cookie** — prevents XSS token theft
- **JWT blacklist** — invalidates tokens on logout
- **bcrypt password hashing** — with strong password enforcement
- **Rate limiting** — per endpoint, configurable via `.env`
- **CORS** — configurable allowed origins via `.env`
- **Hospital isolation** — staff can only search patients from their own hospital

---

## License

MIT