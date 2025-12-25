package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	data "task-cli/internal/Data"
	repository "task-cli/internal/Repository"
	"task-cli/internal/service"
)

var (
	taskService service.TaskService
	scanner     *bufio.Scanner
)

func main() {
	// Setup
	repo := repository.NewFileTaskRepository("task.json")
	taskService = service.NewTaskService(repo)
	scanner = bufio.NewScanner(os.Stdin)

	fmt.Println("╔════════════════════════════════════╗")
	fmt.Println("║     TASK MANAGER CLI               ║")
	fmt.Println("╚════════════════════════════════════╝")

	// Main loop
	for {
		showMenu()
		choice := readInput("\nPilih menu (0-6): ")

		switch choice {
		case "1":
			createTask()
		case "2":
			listTasks()
		case "3":
			updateTask()
		case "4":
			completeTask()
		case "5":
			deleteTask()
		case "6":
			showStats()
		case "0":
			fmt.Println("\n Terima kasih! Goodbye!")
			return
		default:
			fmt.Println("Pilihan tidak valid!")
		}
	}
}

func showMenu() {
	fmt.Println("\n" + strings.Repeat("=", 40))
	fmt.Println("MENU:")
	fmt.Println("1. Create Task")
	fmt.Println("2. ist All Tasks")
	fmt.Println("3. Edit Task")
	fmt.Println("4. Mark as Done")
	fmt.Println("5. Delete Task")
	fmt.Println("6. Show Statistics")
	fmt.Println("0. Exit")
	fmt.Println(strings.Repeat("=", 40))
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}


//  CREATE TASK
func createTask() {
	fmt.Println("\nCREATE NEW TASK")
	fmt.Println(strings.Repeat("-", 40))
	
	description := readInput("Task description: ")
	if description == "" {
		fmt.Println("Description cannot be empty!")
		return
	}

	task, err := taskService.CreateTask(description)
	if err != nil {
		fmt.Println(" Error:", err)
		return
	}

	fmt.Printf("\nTask created successfully!\n")
	fmt.Printf("   ID: %d\n", task.ID)
	fmt.Printf("   Description: %s\n", task.Description)
	fmt.Printf("   Status: %s\n", statusToString(task.Status))
}


// 2. LIST ALL TASKS

func listTasks() {
	fmt.Println("\nALL TASKS")
	fmt.Println(strings.Repeat("-", 40))

	tasks, err := taskService.ReadAll()
	if err != nil {
		fmt.Println(" Error:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found. Create one first!")
		return
	}

	for _, task := range tasks {
		statusIcon := getStatusIcon(task.Status)
		fmt.Printf("\n[%d] %s %s\n", task.ID, statusIcon, task.Description)
		fmt.Printf("    Status: %s\n", statusToString(task.Status))
		fmt.Printf("    Created: %s\n", task.CreatedAt.Format("2006-01-02 15:04"))
		fmt.Printf("    Updated: %s\n", task.UpdatedAt.Format("2006-01-02 15:04"))
	}
}


// 3. UPDATE/EDIT TASK
func updateTask() {
	fmt.Println("\nEDIT TASK")
	fmt.Println(strings.Repeat("-", 40))

	// Show current tasks
	tasks, err := taskService.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println(" No tasks to edit!")
		return
	}

	// Show list
	fmt.Println("\nAvailable tasks:")
	for _, t := range tasks {
		fmt.Printf("  [%d] %s (%s)\n", t.ID, t.Description, statusToString(t.Status))
	}

	// Input task ID
	idStr := readInput("\nTask ID to edit: ")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid ID!")
		return
	}

	// Find task
	var currentTask *data.Task
	for _, t := range tasks {
		if t.ID == id {
			currentTask = &t
			break
		}
	}

	if currentTask == nil {
		fmt.Printf(" Task with ID %d not found!\n", id)
		return
	}

	// Show current data
	fmt.Println("\nCurrent task:")
	fmt.Printf("  Description: %s\n", currentTask.Description)
	fmt.Printf("  Status: %s\n", statusToString(currentTask.Status))

	// Input new data
	fmt.Println("\nEnter new data (press Enter to keep current):")
	
	newDesc := readInput(fmt.Sprintf("New description [%s]: ", currentTask.Description))
	if newDesc == "" {
		newDesc = currentTask.Description
	}

	fmt.Println("\nStatus options:")
	fmt.Println("  0 = Todo")
	fmt.Println("  1 = In Progress")
	fmt.Println("  2 = Done")
	statusStr := readInput(fmt.Sprintf("New status [%d]: ", currentTask.Status))
	
	var newStatus data.Status
	if statusStr == "" {
		newStatus = currentTask.Status
	} else {
		statusInt, err := strconv.Atoi(statusStr)
		if err != nil || statusInt < 0 || statusInt > 2 {
			fmt.Println("Invalid status!")
			return
		}
		newStatus = data.Status(statusInt)
	}

	// Update
	updated, err := taskService.UpdateTask(id, newDesc, newStatus)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("\nTask updated successfully!")
	fmt.Printf("   ID: %d\n", updated.ID)
	fmt.Printf("   Description: %s\n", updated.Description)
	fmt.Printf("   Status: %s\n", statusToString(updated.Status))
}


// 4. COMPLETE TASK
func completeTask() {
	fmt.Println("\nMARK TASK AS DONE")
	fmt.Println(strings.Repeat("-", 40))

	// Show current tasks
	tasks, err := taskService.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println(" No tasks available!")
		return
	}

	// Show incomplete tasks only
	fmt.Println("\nIncomplete tasks:")
	hasIncomplete := false
	for _, t := range tasks {
		if t.Status != data.Done {
			fmt.Printf("  [%d] %s (%s)\n", t.ID, t.Description, statusToString(t.Status))
			hasIncomplete = true
		}
	}

	if !hasIncomplete {
		fmt.Println(" All tasks are already done!")
		return
	}

	// Input task ID
	idStr := readInput("\nTask ID to mark as done: ")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(" Invalid ID!")
		return
	}

	// Complete
	if err := taskService.CompleteTask(id); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("\n Task #%d marked as done!\n", id)
}


// 5. DELETE TASK

func deleteTask() {
	fmt.Println("\n DELETE TASK")
	fmt.Println(strings.Repeat("-", 40))

	// Show current tasks
	tasks, err := taskService.ReadAll()
	if err != nil {
		fmt.Println(" Error:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println(" No tasks to delete!")
		return
	}

	// Show list
	fmt.Println("\nAvailable tasks:")
	for _, t := range tasks {
		fmt.Printf("  [%d] %s (%s)\n", t.ID, t.Description, statusToString(t.Status))
	}

	// Input task ID
	idStr := readInput("\nTask ID to delete: ")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(" Invalid ID!")
		return
	}

	// Confirm
	confirm := readInput(fmt.Sprintf(" Delete task #%d? (y/n): ", id))
	if strings.ToLower(confirm) != "y" {
		fmt.Println("Deletion cancelled")
		return
	}

	// Delete
	if err := taskService.DeleteTask(id); err != nil {
		fmt.Println(" Error:", err)
		return
	}

	fmt.Printf("\n Task #%d deleted successfully!\n", id)
}

// 6. SHOW STATISTICS

func showStats() {
	fmt.Println("\n STATISTICS")
	fmt.Println(strings.Repeat("-", 40))

	stats, err := taskService.GetTaskStats()
	if err != nil {
		fmt.Println(" Error:", err)
		return
	}

	fmt.Printf("\nTotal Tasks:    %d\n", stats.TotalTasks)
	fmt.Printf("Todo:           %d\n", stats.TodoCount)
	fmt.Printf(" In Progress:    %d\n", stats.InProgressCount)
	fmt.Printf("Done:           %d\n", stats.DoneCount)
	fmt.Printf("Progress:       %.1f%%\n", stats.ProgressPercentage)
}


// HELPER FUNCTIONS


func statusToString(status data.Status) string {
	switch status {
	case data.Todo:
		return "Todo"
	case data.InProgress:
		return "In Progress"
	case data.Done:
		return "Done"
	default:
		return "Unknown"
	}
}

func getStatusIcon(status data.Status) string {
	switch status {
	case data.Todo:
		return ""
	case data.InProgress:
		return ""
	case data.Done:
		return ""
	default:
		return ""
	}
}