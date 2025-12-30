package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	data "task-cli/internal/Data"
	"time"
)

// interface untuk kontrak repository
type TaskRepository interface {
	Create(task *data.Task) (*data.Task, error)
	Update(task *data.Task) (*data.Task, error)
	Delete(task *data.Task) error
	ReadAll() ([]data.Task, error)
}

// struct implementasi
type FileTaskRepository struct {
	filePath string
	mu      sync.RWMutex
}


// CONSTRUCTOR
func NewFileTaskRepository(filePath string) TaskRepository {
	return &FileTaskRepository{
		filePath: filePath,
	}
}

// Implementasi read untuk FileTaskRepository
func (r *FileTaskRepository) readTasks() ([]data.Task, error) {
	
	file, err := os.Open(r.filePath)
	if err != nil {
		
		if errors.Is(err, os.ErrNotExist) {
			return []data.Task{}, nil  
		}
		return nil, err  
	}
	defer file.Close()  

	// json dibungkus memakai wrapper
	var wrapper data.TaskWrapper
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&wrapper); err != nil {
		return nil, err
	}

	return wrapper.Records, nil
}

//  tulis semua tasks ke file
func (r *FileTaskRepository) writeTasks(tasks []data.Task) error {
	
	wrapper := data.TaskWrapper{Records: tasks}

	file, err := os.Create(r.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode dengan indentasi untuk keterbacaan
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")  
	
	if err := encoder.Encode(&wrapper); err != nil {
		return err
	}

	return nil
}
// Implementasi Create untuk FileTaskRepositry
func (r *FileTaskRepository) Create(task *data.Task) (*data.Task, error) {
	
	r.mu.Lock()
	defer r.mu.Unlock()

	// untuk membaca semua tasks
	tasks, err := r.readTasks()
	if err != nil {
		return nil, err
	}

	// batas maksimal ukuran ID
	var maxID int
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}

	// informasi task baru
	task.ID = maxID + 1
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	tasks = append(tasks, *task)

	
	if err := r.writeTasks(tasks); err != nil {
		return nil, err
	}

	
	return task, nil
}
// Implementasi Update untuk FileTaskRepositry
func (r *FileTaskRepository) Update(task *data.Task) (*data.Task, error) {
	
	r.mu.Lock()
	defer r.mu.Unlock()

	
	tasks, err := r.readTasks()
	if err != nil {
		return nil, err
	}
	for i, t := range tasks {
		if t.ID == task.ID {
			task.CreatedAt = t.CreatedAt
			task.UpdatedAt = time.Now()
			tasks[i] = *task
			if err := r.writeTasks(tasks); err != nil {
				return nil, err
			}
			return task, nil
		}
	}

	return nil, fmt.Errorf("task dengan ID %d tidak ditemukan", task.ID)
}
// Implementasi Delete untuk FileTaskRepositry
func (r *FileTaskRepository) Delete(task *data.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	tasks, err := r.readTasks()
	if err != nil {
		return err
	}

	for i, t := range tasks {
		if t.ID == task.ID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			if err := r.writeTasks(tasks); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("task dengan ID %d tidak ditemukan", task.ID)
}

//  Implementasi ReadAll untuk FileTaskRepositry
func (r *FileTaskRepository) ReadAll() ([]data.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return r.readTasks()
}

