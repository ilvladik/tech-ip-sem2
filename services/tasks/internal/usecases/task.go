package usecases

import (
	"context"
	"errors"
	"sync"
	"tech-ip-sem2/services/tasks/internal/domain"
)

var (
	ErrorTaskNotFound = errors.New("Task not found")
)

type TaskUsecase struct {
	mu    sync.RWMutex
	it    int
	tasks map[int]*domain.Task
}

func NewTaskUsecase() *TaskUsecase {
	return &TaskUsecase{
		tasks: make(map[int]*domain.Task),
	}
}

func (s *TaskUsecase) Add(ctx context.Context, in CreateTaskInput) (*domain.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.it++
	task := &domain.Task{
		Id:          s.it,
		Title:       in.Title,
		Description: in.Description,
		DueDate:     in.DueDate,
		Done:        false,
	}

	s.tasks[task.Id] = task
	return task, nil
}

func (s *TaskUsecase) GetAll(ctx context.Context) []domain.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]domain.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, *task)
	}
	return tasks
}

func (s *TaskUsecase) Get(ctx context.Context, id int) (*domain.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, ErrorTaskNotFound
	}
	return task, nil
}

func (s *TaskUsecase) Update(ctx context.Context, id int, in UpdateTaskInput) (*domain.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, ErrorTaskNotFound
	}

	if in.Title != nil {
		task.Title = *in.Title
	}
	if in.Description != nil {
		task.Description = in.Description
	}
	if in.DueDate != nil {
		task.DueDate = in.DueDate
	}
	if in.Done != nil {
		task.Done = *in.Done
	}

	return task, nil
}

func (s *TaskUsecase) Delete(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrorTaskNotFound
	}

	delete(s.tasks, id)
	return nil
}
