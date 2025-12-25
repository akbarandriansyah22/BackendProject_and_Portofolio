# Task CLI Manager

A command-line task management application built with Go, implementing clean architecture principles and layered design patterns.

## Features

- Create new tasks
- List all tasks with detailed information
- Edit task description and status
- Mark tasks as completed
- Delete tasks
- View task statistics and progress

## Tech Stack

- **Language:** Go 1.21+
- **Storage:** JSON file-based persistence
- **Architecture:** Clean layered architecture
  - Data Layer (Models)
  - Repository Layer (Data Access)
  - Service Layer (Business Logic)
  - Presentation Layer (CLI Interface)

## Project Structure

```
backend-task-cli/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── Data/
│   │   └── Data.go          # Task and Status models
│   ├── Repository/
│   │   └── repo.go          # Data access and CRUD operations
│   └── service/
│       └── logicbussines.go # Business logic and validation
├── .gitignore
├── go.mod
├── README.md
└── task.json                # Data storage (auto-generated)
```
