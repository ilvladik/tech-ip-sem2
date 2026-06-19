package task

import (
	"errors"
	"sort"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")

type Repo struct {
	data map[int64]Task
}

func NewRepo() *Repo {
	return &Repo{
		data: map[int64]Task{
			1: {
				ID:          1,
				Title:       "Изучить Redis",
				Description: "Освоить cache-aside",
				DueDate:     time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC),
			},
			2: {
				ID:          2,
				Title:       "Написать кэш",
				Description: "Инвалидация кэша по id",
				DueDate:     time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC),
			},
		},
	}
}

func (r *Repo) GetByID(id int64) (Task, error) {
	t, ok := r.data[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}
	return t, nil
}

func (r *Repo) List(page, limit int) []Task {
	ids := make([]int, 0, len(r.data))
	for id := range r.data {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	start := (page - 1) * limit
	if start >= len(ids) {
		return []Task{}
	}
	end := start + limit
	if end > len(ids) {
		end = len(ids)
	}

	out := make([]Task, 0, end-start)
	for _, id := range ids[start:end] {
		out = append(out, r.data[int64(id)])
	}
	return out
}

func (r *Repo) Update(t Task) error {
	if _, ok := r.data[t.ID]; !ok {
		return ErrTaskNotFound
	}
	r.data[t.ID] = t
	return nil
}

func (r *Repo) Delete(id int64) error {
	if _, ok := r.data[id]; !ok {
		return ErrTaskNotFound
	}
	delete(r.data, id)
	return nil
}
