<div align="center">

# E-Commerce API

**Production-grade RESTful API for e-commerce platforms**

Built with Go · PostgreSQL · Clean Architecture · DevSecOps

[![CI Pipeline](https://img.shields.io/github/actions/workflow/status/akbarandriansyah22/BackendProject_and_Portofolio/ci.yml?branch=main&label=CI%20Pipeline&logo=github&logoColor=white)](https://github.com/akbarandriansyah22/BackendProject_and_Portofolio/actions)
[![CD Staging](https://img.shields.io/github/actions/workflow/status/akbarandriansyah22/BackendProject_and_Portofolio/cd.yml?branch=main&label=CD%20Staging&logo=docker&logoColor=white)](https://github.com/akbarandriansyah22/BackendProject_and_Portofolio/actions)
[![Go Version](https://img.shields.io/badge/Go-1.25.2-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![Framework](https://img.shields.io/badge/Fiber-v2.52.10-00ACD7?logo=go&logoColor=white)](https://gofiber.io/)
[![Database](https://img.shields.io/badge/PostgreSQL-12%2B-336791?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./server/LICENSE)

</div>

---

## Table of Contents

- [Overview](#overview)
- [Key Highlights](#key-highlights)
- [Technology Stack](#technology-stack)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)
- [Getting Started](#getting-started)
- [Environment Configuration](#environment-configuration)
- [API Reference](#api-reference)
- [Authentication & Authorization](#authentication--authorization)
- [CI/CD Pipeline](#cicd-pipeline)
- [Monitoring & Observability](#monitoring--observability)
- [Testing](#testing)
- [Docker](#docker)
- [Security](#security)
- [License](#license)

---

## Overview

This project is a **production-grade RESTful API** for an e-commerce platform, written in **Go (Golang)** and structured with **Clean Architecture** principles. It covers the full e-commerce lifecycle including user authentication, product catalog management, shopping cart operations, order processing, and payment tracking.

Beyond core application logic, the project implements a complete **DevSecOps pipeline** with automated secret scanning, static application security testing (SAST), container vulnerability scanning, and software bill of materials (SBOM) generation on every commit. Observability is handled by a full **Prometheus, Grafana, Loki, and Alertmanager** stack with pre-built dashboards and alert rules.

This project was built to demonstrate practical proficiency in backend engineering, software design, security engineering, and modern DevOps practices at a production level.

---

## Key Highlights

- **Clean Architecture** — strict separation between handlers, services, repositories, and domain models via port interfaces, enabling full testability without a live database
- **Custom SQL Injection Protection** — whitelist-based `SQLBuilder` and `WhereClause` builder utilities covering dynamic query construction, with dedicated unit tests and benchmarks
- **Full DevSecOps Pipeline** — Gitleaks (secrets), GoSec (SAST), Trivy (filesystem + container image), and SBOM generation automated on every push; SARIF reports surfaced in the GitHub Security tab
- **Observability Stack** — Prometheus metrics endpoint, pre-provisioned Grafana dashboards, Loki log aggregation via Promtail, and Alertmanager with configurable routing
- **Container-Ready** — multi-stage Docker build producing a minimal Alpine-based image, automatically published to GitHub Container Registry (GHCR) with semantic tags on every push
- **Role-Based Access Control** — JWT authentication with RBAC enforced at the middleware layer, supporting Admin and Customer roles with multi-role support and startup secret validation
- **Structured Logging** — `go.uber.org/zap` for high-performance, structured, leveled logging across all application layers
- **Runtime Security Hardening** — Fiber configured with `ReadTimeout`, `WriteTimeout`, `IdleTimeout`, and 2MB `BodyLimit` to prevent Slowloris and OOM DoS attacks; all containers run as non-root user (`appuser`) with `no-new-privileges`
- **Authenticated Metrics Endpoint** — `/metrics` protected with Bearer token authentication and timing-safe comparison (`secureCompare`) to prevent brute-force and timing attacks

---

## Technology Stack

### Backend

| Technology               | Version  | Purpose                           |
| ------------------------ | -------- | --------------------------------- |
| Go                       | 1.25.2   | Primary language                  |
| Fiber                    | v2.52.10 | HTTP web framework                |
| PostgreSQL               | 12+      | Relational database               |
| golang-jwt/jwt           | v5.3.0   | JWT authentication                |
| golang.org/x/crypto      | v0.47.0  | bcrypt password hashing           |
| go.uber.org/zap          | v1.27.1  | Structured logging                |
| prometheus/client_golang | v1.23.2  | Metrics collection and exposition |
| joho/godotenv            | v1.5.1   | Environment variable management   |
| lib/pq                   | v1.10.9  | PostgreSQL driver                 |
| stretchr/testify         | v1.11.1  | Test assertions and mocking       |

### DevSecOps & Infrastructure

| Tool                 | Purpose                                               |
| -------------------- | ----------------------------------------------------- |
| GitHub Actions       | CI/CD pipeline automation                             |
| Docker + GHCR        | Containerization and image registry                   |
| Gitleaks             | Secret and credential scanning                        |
| golangci-lint v2.7.2 | Static analysis and linting                           |
| GoSec v2.22.11       | Go-specific SAST with SARIF output                    |
| Trivy v0.69.3        | Filesystem and container image vulnerability scanning |
| CycloneDX            | Software Bill of Materials (SBOM) generation          |
| Prometheus           | Metrics collection and alerting                       |
| Grafana              | Metrics visualization and dashboards                  |
| Loki + Promtail      | Log aggregation and shipping                          |
| Alertmanager         | Alert routing and notification management             |

---

## Architecture

The application follows **Clean Architecture** with four distinct layers. Each layer communicates only through defined port interfaces, which decouples implementation details and makes each layer independently testable.

```
+--------------------------------------------------------------+
|                        Handler Layer                         |
|   HTTP request parsing, input validation, response           |
|   formatting. Contains no business logic.                    |
+--------------------------------------------------------------+
                             |
                      port interfaces
                             |
+--------------------------------------------------------------+
|                        Service Layer                         |
|   Business logic, domain rules, transaction orchestration.   |
|   No knowledge of HTTP, SQL, or infrastructure concerns.     |
+--------------------------------------------------------------+
                             |
                      port interfaces
                             |
+--------------------------------------------------------------+
|                      Repository Layer                        |
|   Data access, SQL queries, and database transactions.       |
|   Returns domain models only.                                |
+--------------------------------------------------------------+
                             |
+--------------------------------------------------------------+
|                      PostgreSQL Database                     |
+--------------------------------------------------------------+
```

Port interfaces are defined in `server/internal/ports/` and allow the service layer to be fully unit-tested using mock repositories without requiring a live database connection.

---

## Project Structure

```
e-commerce-api/
|
+-- .github/
|   +-- workflows/
|       +-- ci.yml                  # Go CI + DevSecOps pipeline
|       +-- cd.yml                  # Continuous deployment to staging (GHCR)
|
+-- monitoring/
|   +-- prometheus.yml              # Prometheus scrape configuration
|   +-- alert.rules.yml             # Alert rule definitions (9 rules)
|   +-- alertmanager.yml            # Alert routing and receiver configuration
|   +-- grafana/
|   |   +-- dashboards/             # Pre-built Grafana dashboard JSON files
|   |   +-- provisioning/           # Datasource and dashboard auto-provisioning
|   +-- loki/
|   |   +-- loki-config.yml         # Loki log storage configuration
|   +-- promtail/
|       +-- promtail-config.yml     # Docker log shipping to Loki
|
+-- server/
|   +-- cmd/
|   |   +-- main.go                 # Application entry point
|   |
|   +-- internal/
|       +-- config/                 # Environment configuration loader
|       +-- database/               # Connection pool, retry logic, query helpers
|       +-- handler/                # HTTP handlers (one per domain)
|       |   +-- auth_handler.go
|       |   +-- cart_handler.go
|       |   +-- category_handler.go
|       |   +-- order_handler.go
|       |   +-- Product_handler.go
|       +-- middleware/             # JWT auth, CORS, rate limiter, request logger middleware
|       +-- models/                 # Domain models, DTOs, request/response structs
|       +-- observability/          # Logger interface, Zap implementation, Prometheus metrics
|       +-- ports/                  # Repository and service interface definitions
|       +-- querybuilder/           # SQL injection-safe dynamic query construction
|       +-- repository/             # PostgreSQL data access implementations
|       +-- security/               # JWT, bcrypt, RBAC, rate limiter, password policy
|       +-- service/                # Business logic implementations
|       +-- test/                   # Service-level unit tests and mock implementations
|       +-- utils/                  # HTTP response helpers, pagination utilities
|
+-- Dockerfile                      # Multi-stage build (builder + alpine:3.20.3 runtime, non-root)
+-- go.mod
+-- go.sum
+-- README.md
```

---

## Database Schema

```
+---------+       +----------+       +--------+       +------------+
|  roles  +------>+  users   +------>+ carts  +------>+ cart_items |
+---------+       +----------+       +--------+       +-----+------+
                                                            |
+------------+       +----------+                          |
| categories +------>+ products +<-------------------------+
+------------+       +----+-----+
      |                   |
      |       +-----------------------+
      +------>+  product_categories   |
              +-----------------------+

+----------+       +--------------+       +----------+
|  orders  +------>+  order_items |       | payments |
+----+-----+       +--------------+       +----+-----+
     |                                         ^
     +-----------------------------------------+
```

| Table                | Description                                                                     |
| -------------------- | ------------------------------------------------------------------------------- |
| `roles`              | System roles: Admin (1), Customer (2)                                           |
| `users`              | User accounts with role association, hashed password, email, and account status |
| `products`           | Product catalog with name, slug, SKU, price, stock, and active flag             |
| `categories`         | Categories supporting parent-child hierarchy for nested navigation              |
| `product_categories` | Many-to-many join table between products and categories                         |
| `carts`              | One active cart per user                                                        |
| `cart_items`         | Cart line items with quantity and price snapshot at time of addition            |
| `orders`             | Customer orders with status, payment method, shipping address, and total        |
| `order_items`        | Order line items with price and quantity captured at checkout                   |
| `payments`           | Payment transactions with method, status, and external transaction reference    |

---

## Getting Started

### Prerequisites

- Go 1.25.2 or later
- PostgreSQL 12 or later
- Docker (optional, for containerized runs)
- Git

### Installation

**1. Clone the repository**

```bash
git clone https://github.com/akbarandriansyah22/BackendProject_and_Portofolio.git
cd BackendProject_and_Portofolio/e-commerce-api
```

**2. Install dependencies**

```bash
go mod tidy
go mod download
go mod verify
```

**3. Create the PostgreSQL database**

```sql
CREATE DATABASE ecommerce;
```

**4. Apply schema migrations**

Run the provided SQL schema files using `psql` or your preferred migration tool.

**5. Configure the environment**

```bash
cp .env.example .env
# Open .env and fill in your local values
# Generate secrets with: openssl rand -hex 32
```

**6. Start the application**

```bash
go run server/cmd/main.go
```

**7. Verify the server is healthy**

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

---

## Environment Configuration

All configuration is managed via environment variables loaded from a `.env` file using `godotenv`. Copy `.env.example` to `.env` and populate the values:

```env
# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ecommerce
DB_SSLMODE=disable

# JWT
JWT_SECRET=replace-with-a-strong-random-secret   # min 32 chars, generate: openssl rand -hex 32
JWT_EXPIRATION_HOURS=24

# Metrics
METRICS_TOKEN=replace-with-strong-random-token   # generate: openssl rand -hex 32

# CORS
CORS_ALLOWED_ORIGINS=*

# Application
APP_NAME=E-Commerce API
APP_VERSION=1.0.0
LOG_LEVEL=info
```

> **Security Note:** Never commit `.env` to version control. Generate `JWT_SECRET` and `METRICS_TOKEN` with `openssl rand -hex 32`. The application will refuse to start if these are missing or use known-weak values. In production, inject secrets via your platform's secrets manager and set `DB_SSLMODE=require`.

---

## API Reference

All responses use a consistent JSON envelope:

```json
{
  "success": true,
  "message": "Human-readable status message",
  "data": {}
}
```

Paginated responses include a `meta` object:

```json
{
  "success": true,
  "data": [],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 100,
    "total_pages": 10
  }
}
```

### Authentication

| Method | Endpoint                    | Description                                  | Access |
| ------ | --------------------------- | -------------------------------------------- | ------ |
| POST   | `/api/auth/register`        | Register a new user account                  | Public |
| POST   | `/api/auth/login`           | Authenticate and receive a JWT token         | Public |
| GET    | `/api/auth/profile`         | Get the current authenticated user's profile | JWT    |
| PUT    | `/api/auth/profile`         | Update the current user's profile            | JWT    |
| PUT    | `/api/auth/change-password` | Change the current user's password           | JWT    |

### Products

| Method | Endpoint                             | Description                              | Access |
| ------ | ------------------------------------ | ---------------------------------------- | ------ |
| GET    | `/api/products`                      | List all active products with pagination | Public |
| GET    | `/api/products/:id`                  | Get product detail by ID                 | Public |
| GET    | `/api/products/slug/:slug`           | Get product detail by URL slug           | Public |
| GET    | `/api/products/search?q=`            | Full-text search across products         | Public |
| GET    | `/api/products/category/:id`         | List products within a specific category | Public |
| POST   | `/api/admin/products`                | Create a new product                     | Admin  |
| PUT    | `/api/admin/products/:id`            | Update an existing product               | Admin  |
| DELETE | `/api/admin/products/:id`            | Soft-delete a product                    | Admin  |
| PUT    | `/api/admin/products/:id/activate`   | Activate a product listing               | Admin  |
| PUT    | `/api/admin/products/:id/deactivate` | Deactivate a product listing             | Admin  |
| PATCH  | `/api/admin/products/:id/stock`      | Adjust product stock quantity            | Admin  |

### Categories

| Method | Endpoint                            | Description                             | Access |
| ------ | ----------------------------------- | --------------------------------------- | ------ |
| GET    | `/api/categories`                   | List all active categories              | Public |
| GET    | `/api/categories/:id`               | Get category detail by ID               | Public |
| GET    | `/api/categories/:id/products`      | List all products in a category         | Public |
| GET    | `/api/categories/:id/subcategories` | List subcategories of a parent category | Public |
| POST   | `/api/admin/categories`             | Create a new category                   | Admin  |
| PUT    | `/api/admin/categories/:id`         | Update a category                       | Admin  |
| DELETE | `/api/admin/categories/:id`         | Delete a category                       | Admin  |
| PATCH  | `/api/admin/categories/:id/status`  | Toggle category active status           | Admin  |
| GET    | `/api/admin/categories/stats`       | Retrieve category statistics            | Admin  |

### Cart

| Method | Endpoint              | Description                               | Access |
| ------ | --------------------- | ----------------------------------------- | ------ |
| GET    | `/api/cart`           | Get the current user's cart and its items | JWT    |
| POST   | `/api/cart/items`     | Add an item to the cart                   | JWT    |
| DELETE | `/api/cart/items/:id` | Remove a specific item from the cart      | JWT    |
| DELETE | `/api/cart`           | Clear all items from the cart             | JWT    |

### Orders

| Method | Endpoint                          | Description                                      | Access |
| ------ | --------------------------------- | ------------------------------------------------ | ------ |
| POST   | `/api/orders`                     | Create an order from the current cart (checkout) | JWT    |
| GET    | `/api/orders`                     | List all orders belonging to the current user    | JWT    |
| GET    | `/api/orders/:id`                 | Get order detail by ID                           | JWT    |
| GET    | `/api/orders/number/:orderNumber` | Get order detail by order number                 | JWT    |
| POST   | `/api/orders/:id/cancel`          | Cancel a pending order                           | JWT    |
| GET    | `/api/admin/orders`               | List all orders with filtering and sorting       | Admin  |
| PUT    | `/api/admin/orders/:id/status`    | Update the status of any order                   | Admin  |
| GET    | `/api/admin/orders/stats`         | Retrieve aggregated order statistics             | Admin  |

### System

| Method | Endpoint   | Description                                         |
| ------ | ---------- | --------------------------------------------------- |
| GET    | `/health`  | Application and database health check               |
| GET    | `/metrics` | Prometheus metrics endpoint (Bearer token required) |

---

## Authentication & Authorization

Authentication is implemented using **JSON Web Tokens (JWT)** signed with HMAC-SHA256. Tokens are issued at login and must be included in the `Authorization` header of every protected request:

```
Authorization: Bearer <token>
```

### JWT Payload

```json
{
  "user_id": 1,
  "email": "user@example.com",
  "role_id": 2,
  "full_name": "John Doe",
  "iss": "e-commerce-api",
  "iat": 1700000000,
  "exp": 1700086400
}
```

### Role Definitions

| Role ID | Name     | Permissions                                                           |
| ------- | -------- | --------------------------------------------------------------------- |
| 1       | Admin    | Full access to all endpoints, including all admin management routes   |
| 2       | Customer | Access to own profile, cart, and orders; all public catalog endpoints |

Authorization is enforced by the `RequireRole()` middleware, which supports multiple allowed role IDs per route group and validates that `role_id` is a non-zero value before allowing access.

---

## CI/CD Pipeline

### Continuous Integration (`ci.yml`)

Triggered on every push to `main` and every pull request. The pipeline runs the following steps in order:

```
Checkout -> Setup Go 1.25.2 (with module cache)
         -> Dependency verification (go mod tidy / download / verify)
         -> Secret scanning          (Gitleaks)
         -> Static analysis          (golangci-lint v2.7.2)
         -> SAST                     (GoSec v2.22.11    ->  gosec.sarif)
         -> Filesystem scan          (Trivy v0.69.3     ->  trivy-fs.sarif)
         -> Build Docker image       (e-commerce-api:ci)
         -> Container image scan     (Trivy v0.69.3     ->  trivy-image.sarif)
         -> SBOM generation          (CycloneDX         ->  sbom.json artifact)
         -> Run tests                (go test -v -race -coverprofile)
         -> Build binary             (go build          ->  bin/ecommerce-api artifact, 7-day retention)
```

SARIF reports from GoSec and Trivy are automatically uploaded to the **GitHub Security tab** for centralized vulnerability tracking. Trivy is installed directly from GitHub Releases (`v0.69.3`) for reproducible, version-pinned builds.

### Continuous Deployment (`cd.yml`)

Triggered on successful CI completion or direct push to `main`. Builds and publishes a multi-platform Docker image to GitHub Container Registry (GHCR):

```
Registry: ghcr.io/akbarandriansyah22/backendproject_and_portofolio/e-commerce-api

Tags published on each run:
  main              Branch name tag
  main-<git-sha>    Commit SHA for full build traceability
  staging           Environment label
  latest            Most recent successful build
```

---

## Monitoring & Observability

### Prometheus Metrics

The API exposes application metrics at the `/metrics` endpoint (Bearer token required). Key custom metrics:

| Metric                          | Type      | Labels               | Description                                    |
| ------------------------------- | --------- | -------------------- | ---------------------------------------------- |
| `http_requests_total`           | Counter   | method, path, status | Total HTTP requests processed                  |
| `http_request_duration_seconds` | Histogram | method, path         | Request latency distribution (p50 / p95 / p99) |
| `http_errors_total`             | Counter   | method, path, status | Total HTTP errors (4xx and 5xx)                |
| `http_request_size_bytes`       | Histogram | method, path         | Incoming request payload size                  |
| `http_response_size_bytes`      | Histogram | method, path, status | Outgoing response payload size                 |
| `db_connections_active`         | Gauge     | —                    | Currently active database connections          |
| `db_connection_errors_total`    | Counter   | reason               | Database connection failures by reason         |
| `app_info`                      | Gauge     | version              | Static application version metadata            |

### Grafana Dashboard

A pre-built dashboard (`monitoring/grafana/dashboards/ecommerce-api-dashboard.json`) is automatically provisioned and includes the following panels:

- 5xx error rate with configurable threshold coloring
- Request latency percentiles (p50, p95, p99) over time
- Request rate by HTTP method (GET, POST, PUT, DELETE)
- HTTP status code distribution (2xx, 4xx, 5xx) as stacked time series
- Request volume breakdown by endpoint
- Active database connection gauge

### Alert Rules

Nine alert rules are pre-configured in `monitoring/alert.rules.yml`:

| Alert                             | Severity | Trigger Condition                                 |
| --------------------------------- | -------- | ------------------------------------------------- |
| `ServiceDown`                     | Critical | API unreachable for more than 2 minutes           |
| `HighErrorRate`                   | Warning  | 5xx error rate exceeds 5% over a 5-minute window  |
| `CriticalErrorRate`               | Critical | 5xx error rate exceeds 10% over a 2-minute window |
| `HighLatencyP95`                  | Warning  | p95 latency exceeds 1 second for 10 minutes       |
| `HighLatencyP99`                  | Critical | p99 latency exceeds 5 seconds for 5 minutes       |
| `DatabaseConnectionError`         | Critical | Any DB connection errors sustained over 2 minutes |
| `DatabaseConnectionPoolExhausted` | Warning  | Active connections exceed 20 for 5 minutes        |
| `HealthCheckFailing`              | Critical | `/health` returns non-200 for more than 2 minutes |
| `PrometheusDown`                  | Critical | Prometheus itself is unreachable                  |

Alertmanager (`monitoring/alertmanager.yml`) routes alerts by severity to configurable receivers including Slack, email, and generic webhooks.

---

## Testing

Run the full test suite with race condition detection enabled:

```bash
go test -v -race ./...
```

Generate a coverage report:

```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Test Structure

| Location                        | Contents                                              |
| ------------------------------- | ----------------------------------------------------- |
| `server/internal/test/`         | Service-level unit tests using mock repositories      |
| `server/internal/test/mocks/`   | Mock implementations generated with `testify/mock`    |
| `server/internal/querybuilder/` | Unit tests and benchmarks for the SQL builder package |

The `querybuilder` package includes dedicated injection tests validating that patterns such as `DROP TABLE`, `UNION SELECT`, `DELETE FROM`, `--` comments, and semicolon injection are correctly rejected by the whitelist validator.

---

## Docker

### Build

```bash
docker build -t e-commerce-api:latest ./e-commerce-api
```

The `Dockerfile` uses a **multi-stage build**:

1. **Builder stage** — uses `golang:1.25-alpine` to compile the binary with CGO disabled (`CGO_ENABLED=0`), `-trimpath`, and `-ldflags="-s -w"` for a stripped, path-clean binary
2. **Runtime stage** — copies only the compiled binary and CA certificates into a minimal `alpine:3.20.3` image; runs as non-root user `appuser` with `HEALTHCHECK` configured

The resulting image is typically under 20MB and contains no Go toolchain, build tools, or source code.

### Run

```bash
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=yourpassword \
  -e DB_NAME=ecommerce \
  -e JWT_SECRET=your-secret-min-32-chars \
  -e METRICS_TOKEN=your-metrics-token-min-32-chars \
  -e ENVIRONMENT=production \
  e-commerce-api:latest
```

### Pull from GHCR

```bash
docker pull ghcr.io/akbarandriansyah22/backendproject_and_portofolio/e-commerce-api:latest
```

---

## Security

### Implemented Measures

| Category              | Control                   | Implementation Detail                                                                                                              |
| --------------------- | ------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| Authentication        | JWT token signing         | HMAC-SHA256, configurable expiration via `JWT_EXPIRATION_HOURS`                                                                    |
| Password storage      | Bcrypt hashing            | `golang.org/x/crypto/bcrypt` with cost factor 10                                                                                   |
| Authorization         | Role-based access control | `RequireRole()` middleware with multi-role support, `roleID != 0` validation, and `RoleAdmin`/`RoleCustomer` constants             |
| Input validation      | Request validation        | Explicit field validation in handler and service layers                                                                            |
| SQL injection         | Parameterized queries     | `database/sql` `$N` placeholders for all data queries                                                                              |
| SQL injection         | Dynamic query builder     | Whitelist-based `SQLBuilder` for ORDER BY and WHERE clauses                                                                        |
| Secret scanning       | Gitleaks                  | Automated on every commit in CI, fails the build on detection                                                                      |
| SAST                  | GoSec                     | Scans Go source for known vulnerability patterns, SARIF output                                                                     |
| Dependency scanning   | Trivy filesystem          | Scans Go module graph for known CVEs                                                                                               |
| Container scanning    | Trivy image scan          | Scans the built Docker image layer by layer for CVEs                                                                               |
| SBOM                  | CycloneDX                 | Generated on every CI run, uploaded as a build artifact                                                                            |
| Cross-origin requests | Fiber CORS middleware     | `CORS_ALLOWED_ORIGINS` read from env var — not hardcoded wildcard                                                                  |
| Panic recovery        | Fiber recover middleware  | Prevents server crashes from unhandled runtime panics                                                                              |
| Rate limiting         | Token bucket algorithm    | `NewRateLimiterWithCleanup()` active on `/auth/*` (10 req/min per IP) and `/api/*` (60 req/min per userID) with background cleanup |
| Request timeouts      | Fiber config              | `ReadTimeout: 10s`, `WriteTimeout: 10s`, `IdleTimeout: 60s` — prevents Slowloris attack                                            |
| Request size limit    | Fiber BodyLimit           | 2MB hard limit on all request bodies — prevents OOM DoS                                                                            |
| Metrics auth          | Bearer token              | `/metrics` protected with bearer token and timing-safe `secureCompare()` to prevent timing attacks                                 |
| Non-root container    | Docker USER directive     | Container runs as `appuser` (non-root), binary has `chmod 550` and `chown appuser`                                                 |
| Startup validation    | config.MustLoad()         | Application refuses to start if `JWT_SECRET`, `METRICS_TOKEN`, or `DB_PASSWORD` are missing or use known-weak values               |

### Production Hardening Checklist

**Already implemented:**

- ✅ Non-root container user (`appuser`) with `chmod 550`
- ✅ Alpine image pinned to specific version (`alpine:3.20.3`)
- ✅ `JWT_SECRET` validated at startup — minimum 32 characters, rejects known-weak values
- ✅ `METRICS_TOKEN` required at startup — `/metrics` protected with bearer token
- ✅ Rate limiter active on all public endpoints (`/auth/*` and `/api/*`)
- ✅ Request timeouts and body limit explicitly configured in Fiber
- ✅ CORS configured from env var, not hardcoded wildcard
- ✅ `/health` endpoint does not expose environment name or internal details
- ✅ `HEALTHCHECK` directive in Dockerfile for container orchestration
- ✅ Build binary stripped with `-trimpath -ldflags="-s -w"` to remove path information

**Required before production deployment:**

- Set `DB_SSLMODE=require` and provide valid TLS certificates for all database connections
- Inject all secrets via a secrets manager (HashiCorp Vault, AWS Secrets Manager, GCP Secret Manager) — never via `.env` files
- Configure Alertmanager receivers with production notification channels (PagerDuty, Slack, email)
- Terminate TLS at the load balancer or reverse proxy (NGINX, Caddy, AWS ALB)
- Enable automated PostgreSQL backups with point-in-time recovery (PITR)

---

## License

This project is licensed under the **MIT License**. See [LICENSE](./server/LICENSE) for full terms.

```
MIT License
Copyright (c) 2025 Akbar Andriansyah
```

---

<div align="center">

Developed by **Akbar Andriansyah**

[GitHub](https://github.com/akbarandriansyah22) &nbsp;&nbsp;|&nbsp;&nbsp; [Roadmap.sh Portfolio](https://roadmap.sh/u/akbarandriansyah22)

</div>
