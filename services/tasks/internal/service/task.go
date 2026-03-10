package service

import (
	"errors"
	"sync"
)

var (
	ErrorTaskNotFound = errors.New("task not found")
)

type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
	Done        *bool   `json:"done,omitempty"`
}

type Task struct {
	Id          int64   `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
	Done        bool    `json:"done"`
}

type TaskService struct {
	mu    sync.RWMutex
	it    int64
	tasks map[int64]*Task
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks: make(map[int64]*Task),
	}
}

func (s *TaskService) Create(req CreateTaskRequest) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.it++
	task := &Task{
		Id:          s.it,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Done:        false,
	}

	s.tasks[task.Id] = task
	return task
}

func (s *TaskService) GetAll() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *TaskService) GetByID(id int64) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, ErrorTaskNotFound
	}
	return task, nil
}

func (s *TaskService) Update(id int64, req UpdateTaskRequest) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, ErrorTaskNotFound
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.Done != nil {
		task.Done = *req.Done
	}

	return task, nil
}

func (s *TaskService) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrorTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}
