# Task CLI Manager

Task CLI Manager is a command-line task management application built with Go, implementing clean architecture principles and a layered design.

## Features

- **Create New Task:** Add a task with a description.
- **View All Tasks:** Display a list of tasks with status, creation time, and update details.
- **Edit Task:** Change the description and status of a task.
- **Mark as Done:** Mark a task as completed (Done).
- **Delete Task:** Remove a task by its ID.
- **Statistics:** View statistics of total tasks, progress, and status distribution.

## Architecture

This application uses a layered architecture to maintain separation of concerns:

- **Data Layer:** Defines data models (`Task`, [`Status`](internal/Data/Data.go)) in [internal/Data/Data.go](internal/Data/Data.go).
- **Repository Layer:** Manages data access and CRUD operations on the JSON file in [internal/Repository/repo.go](internal/Repository/repo.go).
- **Service Layer:** Contains business logic and validation in [internal/service/logicbussines.go](internal/service/logicbussines.go).
- **Presentation Layer:** CLI interface, receives input and displays output in [cmd/main.go](cmd/main.go).

## Project Structure

```
backend-task-cli/
├── cmd/
│   └── main.go              # CLI application entry point
├── internal/
│   ├── Data/
│   │   └── Data.go          # Task & Status data models
│   ├── Repository/
│   │   └── repo.go          # CRUD & JSON file access
│   └── service/
│       └── logicbussines.go # Business logic & validation
├── .gitignore
├── go.mod
├── README.md
└── task.json                # Data storage (auto-generated)
```

## How It Works

1. **Run the Application:**  
   Execute the following command in the `backend-task-cli` directory:

   ```sh
   go run cmd/main.go
   ```

2. **CLI Menu:**  
   Users can select menu options to create, view, edit, complete, delete tasks, or view statistics.

3. **Data Storage:**  
   All task data is automatically stored in [task.json](task.json) in JSON format.

## Main File Explanations

- [cmd/main.go](cmd/main.go):  
  Application entry point, displays menu, receives input, and calls the service layer.

- [internal/Data/Data.go](internal/Data/Data.go):  
  Defines the `Task`, `Status` data structures, and data wrappers for JSON serialization.

- [internal/Repository/repo.go](internal/Repository/repo.go):  
  Implements the repository interface for CRUD operations on the JSON file.

- [internal/service/logicbussines.go](internal/service/logicbussines.go):  
  Contains business logic such as input validation, status updates, and statistics calculation.

## Usage Examples

- **Create a task:**  
  Select menu `1`, enter the task description.

- **View tasks:**  
  Select menu `2` to see all tasks with their status and timestamps.

- **Edit a task:**  
  Select menu `3`, enter the task ID, then update the description/status.

- **Mark as done:**  
  Select menu `4`, enter the ID of the task to mark as done.

- **Delete a task:**  
  Select menu `5`, enter the task ID, then confirm deletion.

- **Statistics:**  
  Select menu `6` to view task count and progress statistics.

## License

MIT License © 2025 Akbar Andriansyah

---

For implementation details, please see the source code in each file in [backend-task-cli](.).
