package service

import (
	"fmt"

	"example.com/pz11-graphql/internal/store"
)

// TaskService — общий слой данных для GraphQL (тот же домен, что tasks в pz7/pz10).
type TaskService struct {
	store *store.Store
}

func New(st *store.Store) *TaskService {
	return &TaskService{store: st}
}

func (s *TaskService) List() []*store.Task {
	return s.store.List()
}

func (s *TaskService) Get(id string) (*store.Task, error) {
	t, ok := s.store.Get(id)
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (s *TaskService) Create(title string, description *string) *store.Task {
	return s.store.Create(title, description)
}

func (s *TaskService) Update(id string, title, description *string, done *bool) (*store.Task, error) {
	t, err := s.store.Update(id, title, description, done)
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}
	return t, nil
}

func (s *TaskService) Delete(id string) (bool, error) {
	return s.store.Delete(id)
}
