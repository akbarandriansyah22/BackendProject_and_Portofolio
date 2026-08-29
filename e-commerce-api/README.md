<div align="center">

# E-Commerce API

RESTful API for a small e-commerce catalog, cart, and orders.

Built with Go · PostgreSQL · Fiber

[![CI Pipeline](https://img.shields.io/github/actions/workflow/status/akbarandriansyah22/BackendProject_and_Portofolio/ci.yml?branch=main&label=CI&logo=github&logoColor=white)](https://github.com/akbarandriansyah22/BackendProject_and_Portofolio/actions)
[![GHCR](https://img.shields.io/github/actions/workflow/status/akbarandriansyah22/BackendProject_and_Portofolio/cd.yml?branch=main&label=GHCR%20publish&logo=docker&logoColor=white)](https://github.com/akbarandriansyah22/BackendProject_and_Portofolio/actions)
[![Go Version](https://img.shields.io/badge/Go-1.25.2-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./server/LICENSE)

</div>

---

## Overview

This folder contains a Go HTTP API for user auth, products, categories, cart, and orders. Layers are split into handlers, services, and PostgreSQL repositories behind interfaces in `server/internal/ports/`.

CI lives at the **repository root** (`.github/workflows/`), not inside this folder. Successful CI on `main` publishes a Docker image to GitHub Container Registry. There is no automated deploy to a staging cluster.

---

## Technology Stack

| Technology | Purpose |
| ---------- | ------- |
| Go 1.25.2 | Language |
| Fiber v2 | HTTP |
| PostgreSQL 12+ | Data store |
| golang-jwt/jwt v5 | JWT |
| bcrypt (`golang.org/x/crypto`) | Password hashing |
| zap | Logging |
| Prometheus client | `/metrics` |

---

## Project Structure

```
BackendProject_and_Portofolio/
|-- .github/workflows/
|   |-- ci.yml                 # test, SAST, image scan (root)
|   |-- cd.yml                 # build+push GHCR after CI success
|-- e-commerce-api/
    |-- migrations/001_init.sql
    |-- Dockerfile
    |-- .env.example
    |-- server/cmd/main.go
    |-- server/internal/...
```

---

## Database Schema

Tables: `roles`, `users`, `products`, `categories`, `product_categories`, `carts`, `cart_items`, `orders`, `order_items`, `payments`.

Default roles seeded by the schema: Admin (`id=1`), Customer (`id=2`).

---

## Getting Started

### Prerequisites

- Go 1.25.2 or later
- PostgreSQL 12 or later
- Git

### 1. Clone

```bash
git clone https://github.com/akbarandriansyah22/BackendProject_and_Portofolio.git
cd BackendProject_and_Portofolio/e-commerce-api
```

### 2. Environment

```bash
cp .env.example .env
# JWT_SECRET and METRICS_TOKEN must be at least 32 characters:
#   openssl rand -hex 32
# Set DB_PASSWORD to your local Postgres password.
```

The process refuses to start if `JWT_SECRET`, `METRICS_TOKEN`, or `DB_PASSWORD` are missing or too short.

### 3. Database + schema

```bash
createdb ecommerce   # or: psql -U postgres -c "CREATE DATABASE ecommerce;"
psql -U postgres -d ecommerce -f migrations/001_init.sql
```

Docker Postgres example (mount the migration as init SQL):

```bash
docker run --name ecommerce-pg -e POSTGRES_PASSWORD=change-me \
  -e POSTGRES_DB=ecommerce -p 5432:5432 -d postgres:16
# wait until ready, then:
psql -h localhost -U postgres -d ecommerce -f migrations/001_init.sql
```

### 4. Run

```bash
go mod download
go run ./server/cmd
```

### 5. Health checks

```bash
curl -s http://localhost:8080/live
curl -s http://localhost:8080/ready
curl -s http://localhost:8080/health
```

`/live` is process-up only (used by the Docker `HEALTHCHECK`). `/ready` pings PostgreSQL and returns 503 if the database is down. `/health` returns `{"status","db","version"}`.

### 6. Example requests

Register and login:

```bash
curl -s -X POST http://localhost:8080/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"Secret123","full_name":"Demo User"}'

curl -s -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"Secret123"}'
```

List products (public):

```bash
curl -s http://localhost:8080/api/products
```

---

## Environment Configuration

See `.env.example` for every variable the app reads (`DB_*`, `SERVER_*`, `JWT_*`, `METRICS_TOKEN`, CORS, `APP_*`, `LOG_LEVEL`). Do not commit `.env`.

---

## API Reference

Auth routes are also available under `/auth/register` and `/auth/login`.

### Authentication

| Method | Endpoint | Access |
| ------ | -------- | ------ |
| POST | `/api/auth/register` | Public |
| POST | `/api/auth/login` | Public |
| GET | `/api/auth/profile` | JWT |
| PUT | `/api/auth/profile` | JWT |
| PUT | `/api/auth/change-password` | JWT |

### Products / categories

| Method | Endpoint | Access |
| ------ | -------- | ------ |
| GET | `/api/products` | Public |
| GET | `/api/products/:id` | Public |
| GET | `/api/products/slug/:slug` | Public |
| GET | `/api/products/search?q=` | Public |
| GET | `/api/products/category/:id` | Public |
| GET | `/api/categories` | Public |
| GET | `/api/categories/:id` | Public |
| GET | `/api/categories/:id/products` | Public |
| GET | `/api/categories/:id/subcategories` | Public |
| POST/PUT/DELETE | `/api/admin/products...` | Admin JWT |
| POST/PUT/DELETE | `/api/admin/categories...` | Admin JWT |

### Cart / orders

| Method | Endpoint | Access |
| ------ | -------- | ------ |
| GET | `/api/cart` | JWT |
| POST | `/api/cart/items` | JWT |
| DELETE | `/api/cart/items/:id` | JWT |
| DELETE | `/api/cart` | JWT |
| POST | `/api/orders` | JWT |
| GET | `/api/orders` | JWT |
| GET | `/api/orders/:id` | JWT |
| POST | `/api/orders/:id/cancel` | JWT |
| GET/PUT | `/api/admin/orders...` | Admin JWT |

### System

| Method | Endpoint |
| ------ | -------- |
| GET | `/live` |
| GET | `/ready` |
| GET | `/health` |
| GET | `/metrics` (Bearer `METRICS_TOKEN`) |

Send JWT as `Authorization: Bearer <token>`. Role `1` is Admin, role `2` is Customer.

---

## CI/CD

Workflows are in the **monorepo root**: `.github/workflows/ci.yml` and `.github/workflows/cd.yml`.

- **CI** runs when `e-commerce-api/**` or those workflow files change. Steps: modules, Gitleaks (`gitleaks detect` on the working tree), golangci-lint, GoSec (fails on high severity; SARIF upload is `continue-on-error`), Trivy, tests with coverage printed (no hard coverage gate), `go build ./server/cmd`.
- **CD** runs only after that CI workflow completes **successfully**. It builds and pushes the image to GHCR. It does not SSH, apply Kubernetes manifests, or deploy to a staging host.

Image: `ghcr.io/akbarandriansyah22/backendproject_and_portofolio/e-commerce-api`

---

## Testing

```bash
go test -v -race ./...
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Mocks live in `server/internal/test/mocks/`.

---

## Docker

```bash
docker build -t e-commerce-api:latest .
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=yourpassword \
  -e DB_NAME=ecommerce \
  -e JWT_SECRET="$(openssl rand -hex 32)" \
  -e METRICS_TOKEN="$(openssl rand -hex 32)" \
  -e ENVIRONMENT=development \
  e-commerce-api:latest
```

The image `HEALTHCHECK` hits `/live` so a brief database outage does not restart the container.

---

## License

MIT. See [LICENSE](./server/LICENSE).
