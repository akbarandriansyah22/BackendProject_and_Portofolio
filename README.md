<div align="center">

# Backend & DevSecOps Portfolio

**Three production-minded Golang backend projects demonstrating backend engineering, clean architecture, and DevSecOps practices**

[![Go](https://img.shields.io/badge/Go-1.25.2-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![GitHub Actions](https://img.shields.io/badge/CI%2FCD-GitHub%20Actions-2088FF?logo=github-actions&logoColor=white)](https://github.com/akbarandriansyah22/BackendProject_and_Portofolio/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

</div>

---

## Overview

This repository contains three Golang-based backend projects built with a **production-ready mindset** — emphasizing maintainability, security, and automation rather than tutorial-style implementations. Each project addresses a different scope and complexity level, collectively demonstrating a range of backend and DevOps competencies.

---

## Projects

| #   | Project                                                         | Type                          | Path                     |
| --- | --------------------------------------------------------------- | ----------------------------- | ------------------------ |
| 1   | [E-Commerce API](#1-e-commerce-api--devsecops-oriented-backend) | REST API + DevSecOps Pipeline | `e-commerce-api/`        |
| 2   | [Blogging Platform API](#2-blogging-platform-api)               | REST API                      | `blogging-platform-api/` |
| 3   | [Task CLI Manager](#3-task-cli-manager)                         | Command-Line Application      | `backend-task-cli/`      |

---

## 1. E-Commerce API — DevSecOps-Oriented Backend

**Reference:** [roadmap.sh/projects/ecommerce-api](https://roadmap.sh/projects/ecommerce-api)

**Tech Stack:** Go (Fiber), PostgreSQL, Docker, GitHub Actions, golangci-lint, GoSec, Trivy

### Overview

The E-Commerce API is a REST backend built with **Clean Architecture** principles and extended with a comprehensive **DevSecOps pipeline**. The project prioritizes code quality, application security, and lifecycle automation — from development to deployment.

Core domain responsibilities include product catalog management, cart and order processing, stock validation, and transaction handling.

### Architecture

The application is structured in three distinct layers with clear separation of concerns:

- **Handler** — HTTP routing, request parsing, and response formatting
- **Service** — business logic and domain rules, independent of transport and persistence concerns
- **Repository** — database access, SQL execution, and data mapping

Configuration is environment-based, loaded from `.env` at startup, with no hardcoded values in application code.

### DevSecOps Pipeline

The CI/CD pipeline is implemented in **GitHub Actions** and runs automatically on every push and pull request. The pipeline executes the following steps in order:

1. Dependency validation (`go mod tidy`, `go mod verify`)
2. Static analysis with `golangci-lint`
3. SAST with GoSec — output uploaded as SARIF to the GitHub Security tab
4. Software Composition Analysis (SCA) with Trivy — filesystem and container image scans, also uploaded as SARIF
5. Docker image build followed by image-level vulnerability scan
6. Unit tests with race condition detection and coverage reporting
7. Binary build with artifact upload (7-day retention)

All security scan results are centralized in the **GitHub Security Dashboard** for continuous visibility.

**Path:** `e-commerce-api/`

---

## 2. Blogging Platform API

**Reference:** [roadmap.sh/projects/blogging-platform-api](https://roadmap.sh/projects/blogging-platform-api)

**Tech Stack:** Go, PostgreSQL

### Overview

The Blogging Platform API is a RESTful backend service demonstrating **core backend development fundamentals** — clean API design, consistent project structure, and straightforward maintainability. It is built around real-world CRUD use cases and serves as a foundation-level companion to the more infrastructure-heavy E-Commerce project.

### Features

- Full CRUD operations for blog posts
- Search and filtering via query parameters
- Consistent response structure and error handling

**Path:** `blogging-platform-api/`

---

## 3. Task CLI Manager

**Reference:** [roadmap.sh/projects/task-tracker](https://roadmap.sh/projects/task-tracker)

**Tech Stack:** Go

### Overview

Task CLI Manager is a command-line application for local task management, built using Go with a **layered design approach**. The project demonstrates the ability to apply clean architecture principles outside of HTTP contexts — maintaining separation between business logic and data persistence in a non-web application.

### Features

- Create, view, update, and delete tasks
- Mark tasks as completed
- Task statistics summary

**Path:** `backend-task-cli/`

---

## Skills Demonstrated

- Backend development with Go
- Clean Architecture and layered design patterns
- RESTful API design and consistent response structures
- CI/CD automation with GitHub Actions
- DevSecOps practices — SAST, SCA, and container image security scanning
- Docker-based build and deployment workflows
- Secure coding principles and production hardening

---

## Contact

- **GitHub:** [github.com/akbarandriansyah22](https://github.com/akbarandriansyah22)
- **LinkedIn:** [linkedin.com/in/akbar-andriansyah-b3907322b](https://www.linkedin.com/in/akbar-andriansyah-b3907322b/)

---

<div align="center">

This repository serves as a backend and DevSecOps portfolio, reflecting readiness for production-grade engineering roles.

</div>
