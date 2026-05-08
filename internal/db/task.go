package db

import "time"

// Task represents the domain model for a task.
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" validate:"required,min=1,max=100"`
	Description string    `json:"description"`
	Status      string    `json:"status" validate:"oneof=todo doing done blocked"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
