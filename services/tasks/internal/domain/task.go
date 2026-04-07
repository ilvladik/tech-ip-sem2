package domain

import "time"

type Task struct {
	Id          int
	Title       string
	Description *string
	DueDate     *time.Time
	Done        bool
}
