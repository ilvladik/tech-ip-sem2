package store

import (
	"fmt"
	"sync"
)

type Task struct {
	ID          string
	Title       string
	Description *string
	Done        bool
}

type Store struct {
	mu    sync.RWMutex
	tasks []*Task
	next  int
}

func New() *Store {
	desc1 := "Учебный пример"
	desc2 := "GraphQL API"
	return &Store{
		next: 3,
		tasks: []*Task{
			{ID: "t_001", Title: "Первая задача", Description: &desc1, Done: false},
			{ID: "t_002", Title: "Вторая задача", Description: &desc2, Done: true},
		},
	}
}

func (s *Store) List() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Task, len(s.tasks))
	copy(out, s.tasks)
	return out
}

func (s *Store) Get(id string) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.tasks {
		if t.ID == id {
			return cloneTask(t), true
		}
	}
	return nil, false
}

func (s *Store) Create(title string, description *string) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := fmt.Sprintf("t_%03d", s.next)
	s.next++
	t := &Task{ID: id, Title: title, Description: description, Done: false}
	s.tasks = append(s.tasks, t)
	return cloneTask(t)
}

func (s *Store) Update(id string, title *string, description *string, done *bool) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.tasks {
		if t.ID != id {
			continue
		}
		if title != nil {
			t.Title = *title
		}
		if description != nil {
			t.Description = description
		}
		if done != nil {
			t.Done = *done
		}
		return cloneTask(t), nil
	}
	return nil, fmt.Errorf("task not found")
}

func (s *Store) Delete(id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.tasks {
		if t.ID != id {
			continue
		}
		s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
		return true, nil
	}
	return false, fmt.Errorf("task not found")
}

func cloneTask(t *Task) *Task {
	var desc *string
	if t.Description != nil {
		v := *t.Description
		desc = &v
	}
	return &Task{ID: t.ID, Title: t.Title, Description: desc, Done: t.Done}
}

func StrPtr(s string) *string {
	return &s
}
