package db

import (
	"context"
	"errors"
	"fmt"

	"task-cli/ent"
	"task-cli/ent/task"

	"github.com/go-playground/validator/v10"
)

var (
	// ErrNotFound is returned when a task is not found.
	ErrNotFound = errors.New("task not found")
)

// Repository handles database operations for tasks.
type Repository struct {
	client   *ent.Client
	validate *validator.Validate
}

// NewRepository creates a new Repository.
func NewRepository(client *ent.Client) *Repository {
	return &Repository{
		client:   client,
		validate: validator.New(),
	}
}

// Map ent.Task to db.Task
func mapEntToDomain(et *ent.Task) *Task {
	return &Task{
		ID:          et.ID,
		Title:       et.Title,
		Description: et.Description,
		Status:      string(et.Status),
		CreatedAt:   et.CreatedAt,
		UpdatedAt:   et.UpdatedAt,
	}
}

// CreateTask inserts a new task.
func (r *Repository) CreateTask(ctx context.Context, title string) (*Task, error) {
	// Validation
	input := struct {
		Title string `validate:"required,min=1,max=100"`
	}{Title: title}
	if err := r.validate.Struct(input); err != nil {
		return nil, fmt.Errorf("repository.CreateTask validation: %w", err)
	}

	et, err := r.client.Task.
		Create().
		SetTitle(title).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository.CreateTask: %w", err)
	}

	return mapEntToDomain(et), nil
}

// GetAllTasks returns all tasks ordered by created_at descending.
func (r *Repository) GetAllTasks(ctx context.Context) ([]*Task, error) {
	ets, err := r.client.Task.
		Query().
		Order(ent.Desc(task.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAllTasks: %w", err)
	}

	tasks := make([]*Task, len(ets))
	for i, et := range ets {
		tasks[i] = mapEntToDomain(et)
	}
	return tasks, nil
}

// GetTaskByID returns a single task by its integer ID.
func (r *Repository) GetTaskByID(ctx context.Context, id int) (*Task, error) {
	et, err := r.client.Task.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repository.GetTaskByID: %w", err)
	}

	return mapEntToDomain(et), nil
}

// UpdateStatus changes the status of a task.
func (r *Repository) UpdateStatus(ctx context.Context, id int, status string) (*Task, error) {
	// Validation
	input := struct {
		Status string `validate:"oneof=todo doing done blocked"`
	}{Status: status}
	if err := r.validate.Struct(input); err != nil {
		return nil, fmt.Errorf("repository.UpdateStatus validation: %w", err)
	}

	et, err := r.client.Task.
		UpdateOneID(id).
		SetStatus(task.Status(status)).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repository.UpdateStatus: %w", err)
	}

	return mapEntToDomain(et), nil
}

// UpdateDescription sets the Markdown description on a task.
func (r *Repository) UpdateDescription(ctx context.Context, id int, description string) (*Task, error) {
	// Validation
	input := struct {
		Description string `validate:"max=10000"`
	}{Description: description}
	if err := r.validate.Struct(input); err != nil {
		return nil, fmt.Errorf("repository.UpdateDescription validation: %w", err)
	}

	et, err := r.client.Task.
		UpdateOneID(id).
		SetDescription(description).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repository.UpdateDescription: %w", err)
	}

	return mapEntToDomain(et), nil
}
