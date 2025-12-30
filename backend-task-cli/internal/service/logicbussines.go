package service

import (
	"errors"
	"fmt"
	"time"

	data "task-cli/internal/Data"
	repo "task-cli/internal/Repository"
)

// Interface - Kontrak method yang harus ada
type TaskService interface {
	
	CreateTask(description string) (*data.Task, error)
	UpdateTask(id int, description string, status data.Status) (*data.Task, error)
	DeleteTask(id int) error
	CompleteTask(taskID int) error
	GetTaskStats() (*TaskStats, error)
	ReadAll() ([]data.Task, error) 
}

// Struct untuk statistik
type TaskStats struct {
	TotalTasks         int
	TodoCount          int
	DoneCount          int
	InProgressCount    int
	ProgressPercentage float64
}

// Struct implementasi
type taskService struct {
	repo repo.TaskRepository  
}

// Constructor
func NewTaskService(repository repo.TaskRepository) TaskService {
	return &taskService{
		repo: repository,
	}
}
// implementasi method CreateTask
func (s *taskService) CreateTask(description string) (*data.Task, error) {
	if description == "" {
		return nil, errors.New("deskripsi task tidak boleh kosong")
	}

	task := &data.Task{
		Description: description,
		Status:      data.Todo,
	}

	return s.repo.Create(task)
}
// implementasi method UpdateTask
func (s *taskService) UpdateTask(id int, description string, status data.Status) (*data.Task, error) {
	if description == "" {
		return nil, errors.New("deskripsi task tidak boleh kosong")
	}

	tasks, err := s.repo.ReadAll()  // ✅ Sekarang sudah ada!
	if err != nil {
		return nil, err
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = description
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()

			return s.repo.Update(&tasks[i])
		}
	}

	return nil, fmt.Errorf("task dengan ID %d tidak ditemukan", id)
}
// implementasi method DeleteTask
func (s *taskService) DeleteTask(id int) error {
	tasks, err := s.repo.ReadAll()  // ✅ Sekarang sudah ada!
	if err != nil {
		return err
	}
	for _, task := range tasks {
		if task.ID == id {
			return s.repo.Delete(&task)
		}
	}

	return fmt.Errorf("task dengan ID %d tidak ditemukan", id)
}
// implementasi method CompleteTask
func (s *taskService) CompleteTask(taskID int) error {
	tasks, err := s.repo.ReadAll()  // ✅ Sekarang sudah ada!
	if err != nil {
		return err
	}

	for i, task := range tasks {
		if task.ID == taskID {
			tasks[i].Status = data.Done
			tasks[i].UpdatedAt = time.Now()

			_, err := s.repo.Update(&tasks[i])
			return err
		}
	}

	return fmt.Errorf("task dengan ID %d tidak ditemukan", taskID)
}
// implementasi method GetTaskStats
func (s *taskService) GetTaskStats() (*TaskStats, error) {
	tasks, err := s.repo.ReadAll()  // ✅ Sekarang sudah ada!
	if err != nil {
		return nil, err
	}

	stats := &TaskStats{}
	stats.TotalTasks = len(tasks)

	for _, task := range tasks {
		switch task.Status {
		case data.Todo:
			stats.TodoCount++
		case data.InProgress:
			stats.InProgressCount++
		case data.Done:
			stats.DoneCount++
		}
	}

	if stats.TotalTasks > 0 {
		stats.ProgressPercentage = (float64(stats.DoneCount) / float64(stats.TotalTasks)) * 100
	}

	return stats, nil
}

// ReadAll - ambil semua tasks
func (s *taskService) ReadAll() ([]data.Task, error) {
	return s.repo.ReadAll()
}