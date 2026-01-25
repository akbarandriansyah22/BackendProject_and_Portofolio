# Backend & DevSecOps Portfolio

This repository contains **three Golang-based backend projects** designed to demonstrate my capabilities in **backend engineering**, **clean architecture**, and **DevSecOps / CI/CD practices**.

All projects are structured with a **production-ready mindset**, focusing on maintainability, security, and automation rather than tutorial-style implementations.

---

## Project List

1. **E-Commerce API (DevSecOps Project)**https://roadmap.sh/projects/ecommerce-api
2. **Blogging Platform API**https://roadmap.sh/projects/blogging-platform-api
3. **Task CLI Manager**https://roadmap.sh/projects/task-tracker

---

## 1. E-Commerce API – DevSecOps-Oriented Backend

**Type**: Backend REST API with CI/CD and Security Automation
**Tech Stack**:

- Golang (Fiber)
- PostgreSQL
- Docker and Docker Compose
- GitHub Actions
- golangci-lint, GoSec, Trivy

### Overview

The E-Commerce API is a backend service built using **Clean Architecture principles** and enhanced with a comprehensive **DevSecOps pipeline**. The project emphasizes code quality, application security, and automation throughout the development lifecycle.

Core responsibilities of the API include:

- Product management
- Cart and order processing
- Stock validation
- Transaction handling

### Architecture

- Clean / Layered Architecture
  - Handler (HTTP layer)
  - Service (business logic)
  - Repository (database access)

- Environment-based configuration
- Clear separation of concerns

### DevSecOps Implementation

The CI/CD pipeline is implemented using **GitHub Actions** and runs automatically on every push and pull request.

The pipeline includes:

- Dependency validation (`go mod tidy`, `go mod verify`)
- Static code analysis using `golangci-lint`
- SAST with GoSec (SARIF output)
- SCA using Trivy (filesystem and container image scanning)
- Docker image build and vulnerability scanning
- Unit testing with race detection and coverage reporting
- Binary build and artifact generation

All security scan results are integrated into the **GitHub Security Dashboard** for centralized visibility.

**Path**: `e-commerce-api/`

---

## 2. Blogging Platform API

**Type**: Backend REST API
**Tech Stack**:

- Golang
- PostgreSQL

### Overview

The Blogging Platform API is a RESTful backend service created to demonstrate **core backend development fundamentals** using Go and PostgreSQL.

The project focuses on clean API design, consistent structure, and ease of maintenance. It is inspired by backend practice challenges and real-world CRUD use cases.

### Features

- CRUD operations for blog posts
- Search and filtering via query parameters
- Simple and maintainable project structure

**Path**: `blogging-platform-api/`

---

## 3. Task CLI Manager

**Type**: Command Line Application
**Tech Stack**:

- Golang

### Overview

Task CLI Manager is a command-line application for task management, built using Go with a **layered and clean design approach**.

This project demonstrates the ability to design non-HTTP applications while maintaining clear separation between business logic and data handling.

### Features

- Create tasks
- View all tasks
- Update existing tasks
- Mark tasks as completed
- Delete tasks
- Task statistics

**Path**: `backend-task-cli/`

---

## Skills Demonstrated

- Backend development with Golang
- Clean Architecture and layered design
- RESTful API design
- CI/CD automation using GitHub Actions
- DevSecOps practices (SAST, SCA, container security)
- Docker-based workflows
- Secure coding principles

---

## Contact

For further discussion or project review:

- GitHub: [https://github.com/](https://github.com/)<username>
- LinkedIn: (optional)

---

This repository serves as a backend and DevSecOps portfolio, highlighting readiness for production-gra
