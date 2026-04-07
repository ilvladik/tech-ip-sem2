package usecases

import "time"

type CreateTaskInput struct {
	Title       string
	Description *string
	DueDate     *time.Time
}

type UpdateTaskInput struct {
	Title       *string
	Description *string
	DueDate     *time.Time
	Done        *bool
}
