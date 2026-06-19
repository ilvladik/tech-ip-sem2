package task

import (
	"fmt"
	"sync"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Store struct {
	mu    sync.Mutex
	tasks map[string]Task
	next  int
}

func NewStore() *Store {
	return &Store{
		tasks: make(map[string]Task),
		next:  1,
	}
}

func (s *Store) Create(title, description string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := fmt.Sprintf("t_%03d", s.next)
	s.next++
	t := Task{ID: id, Title: title, Description: description}
	s.tasks[id] = t
	return t
}
