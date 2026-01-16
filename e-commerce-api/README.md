# E-Commerce API Documentation

## Table of Contents

- Overview
- Technology Stack
- Database Schema
- Project Structure
- Installation
- Configuration
- API Endpoints
- Authentication
- Error Handling
- Development
- Testing
- Deployment
- Security Considerations
- License
- Changelog

---

## Overview

This document describes a production-ready RESTful API for an e-commerce platform built with Go (Golang) and PostgreSQL. The application is designed following Clean Architecture principles, ensuring a clear separation of concerns between handlers, services, and repositories.

The API provides end-to-end e-commerce functionality, including user authentication, product and category management, shopping cart operations, order processing, and payment handling. The system is designed to be scalable, maintainable, and secure.

### Key Features

- JWT-based authentication and authorization
- Role-based access control (Admin, Customer)
- Product and category management
- Shopping cart lifecycle management
- Order creation, tracking, and cancellation
- Payment transaction handling
- RESTful API design with consistent response formats
- Centralized error handling and structured logging
- Database transaction management
- Input validation and sanitization

---

## Technology Stack

### Backend

- **Go 1.25.0** – Core programming language
- **Fiber v2.52.10** – High-performance HTTP web framework
- **PostgreSQL** – Relational database management system

### Key Dependencies

- `github.com/gofiber/fiber/v2` – HTTP framework
- `github.com/lib/pq` – PostgreSQL driver
- `github.com/golang-jwt/jwt/v5` – JWT authentication
- `golang.org/x/crypto` – Secure password hashing (bcrypt)
- `github.com/joho/godotenv` – Environment variable loader

### Development Tooling

- Custom structured logging with log levels
- Database connection pooling
- CORS middleware
- Recovery middleware for panic handling

---

## Database Schema

The database schema is normalized and designed to support transactional consistency, extensibility, and reporting.

### User and Role Management

#### users

Stores user account information.

Key attributes include email-based authentication, role association, account status, and audit timestamps.

#### roles

Defines system roles used for access control (e.g., Admin, Customer).

---

### Product Management

#### product

Stores product catalog data, including pricing, inventory, and activation status.

#### category

Defines hierarchical product categories with optional parent-child relationships.

#### product_category

Implements a many-to-many relationship between products and categories.

---

### Shopping Cart

#### cart

Represents a single active cart per user.

#### cart_item

Stores individual items within a cart, including quantity and price at the time of addition.

---

### Order Management

#### orders

Stores customer orders, including status, payment method, totals, and shipping information.

#### order_item

Stores line items associated with an order, capturing price and quantity at checkout time.

---

### Payment

#### payment

Tracks payment transactions associated with orders, including external transaction references and payment status.

---

## Project Structure

The project follows a modular Clean Architecture layout:

```
e-commerce-api/
├── server/
│   ├── cmd/                # Application entry point
│   └── internal/
│       ├── config/         # Configuration management
│       ├── database/       # Database initialization and helpers
│       ├── handler/        # HTTP handlers (controllers)
│       ├── middleware/     # HTTP middleware (auth, CORS, logging)
│       ├── models/         # Domain models and DTOs
│       ├── repository/     # Data access layer
│       ├── service/        # Business logic layer
│       └── utils/          # Shared utilities
├── go.mod
├── go.sum
├── Dockerfile
├── .env.example
└── README.md
```

---

## Installation

### Prerequisites

- Go 1.25.0 or later
- PostgreSQL 12 or later
- Git

### Setup Steps

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd e-commerce-api
   ```

2. Download dependencies:

   ```bash
   go mod download
   ```

3. Create the database:

   ```sql
   CREATE DATABASE ecommerce;
   ```

4. Apply database schema and migrations.

5. Create environment configuration:

   ```bash
   cp .env.example .env
   ```

6. Start the application:

   ```bash
   cd server/cmd
   go run main.go
   ```

---

## Configuration

All configuration is managed via environment variables.

Key configuration domains include:

- Server settings
- Database connection parameters
- JWT authentication settings
- CORS configuration
- Application metadata and logging levels

Environment-specific configuration is supported for development and production deployments.

---

## API Endpoints

The API is organized by functional domains:

- Authentication
- Products (Public and Admin)
- Categories (Public and Admin)
- Shopping Cart
- Orders (User and Admin)
- Health Check

All endpoints follow RESTful conventions and return consistent JSON response formats.

---

## Authentication

Authentication is implemented using JSON Web Tokens (JWT).

Clients must include a valid token in the `Authorization` header:

```
Authorization: Bearer <token>
```

### Token Claims

- User ID
- Email
- Role ID
- Full name
- Issued at (iat)
- Expiration time (exp)

Token expiration is configurable via environment variables.

---

## Error Handling

The API uses standard HTTP status codes and structured JSON responses.

### Error Response Principles

- Clear, user-safe error messages
- Detailed validation feedback for client-side correction
- No leakage of internal implementation details

---

## Development

### Local Development

Recommended tools:

- `air` for hot reloading
- `golangci-lint` for static analysis

### Coding Standards

- Follow Go formatting and idioms
- Keep functions small and cohesive
- Handle errors explicitly
- Use dependency injection
- Write tests for business logic

---

## Testing

The project uses Go’s built-in testing framework.

- Unit tests for services and utilities
- Repository tests against a test database
- Coverage reporting supported

---

## Deployment

### Production Build

```bash
go build -o ecommerce-api ./server/cmd/main.go
```

### Docker Deployment

The application can be deployed using Docker or Docker Compose with PostgreSQL as a service dependency.

Environment variables must be configured appropriately for production use.

---

## Security Considerations

### Implemented Measures

- Bcrypt password hashing
- JWT-based authentication
- Parameterized SQL queries
- Input validation and sanitization
- Configurable CORS policy

### Recommended Enhancements

- Rate limiting
- HTTPS enforcement
- Secrets management solution
- Automated backups
- Monitoring and alerting

---

## License

This project is proprietary software. All rights reserved.

---

## Changelog

### Version 1.0.0 – 2026-01-16

- Initial release
- Authentication and authorization
- Product and category management
- Shopping cart functionality
- Order processing
- Payment support
- Administrative endpoints
