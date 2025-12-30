package data

import (
	"time"
)

// informasi status progress tugas

type Status int

const (
	Todo Status = iota
	InProgress
	Done
)
// struktur task pada sistem
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
// pembungkus untuk daftar tugas
type TaskWrapper struct {
	Records []Task `json:"records"`
}