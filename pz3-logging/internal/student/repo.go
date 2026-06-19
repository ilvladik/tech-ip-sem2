package student

import (
	"errors"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var (
	ErrStudentNotFound = errors.New("student not found")
	ErrEmailExists     = errors.New("email already exists")
)

type Repo struct {
	log    *zap.Logger
	mu     sync.RWMutex
	data   map[int64]Student
	nextID int64
}

func NewRepo(log *zap.Logger) *Repo {
	return &Repo{
		log:    log,
		nextID: 4,
		data: map[int64]Student{
			1: {ID: 1, FullName: "Иванов Иван Иванович", Group: "ИТТ-01-25", Email: "ivanov@example.com"},
			2: {ID: 2, FullName: "Петрова Мария Сергеевна", Group: "ИТТ-02-25", Email: "petrova@example.com"},
			3: {ID: 3, FullName: "Сидоров Алексей Дмитриевич", Group: "ИТТ-03-25", Email: "sidorov@example.com"},
		},
	}
}

func (r *Repo) GetByID(id int64) (Student, error) {
	r.log.Debug("repo: lookup student started", zap.Int64("student_id", id))

	r.mu.RLock()
	defer r.mu.RUnlock()

	st, ok := r.data[id]
	if !ok {
		return Student{}, ErrStudentNotFound
	}
	return st, nil
}

func (r *Repo) Create(in CreateInput) (Student, error) {
	r.log.Debug("repo: create student started",
		zap.String("full_name", in.FullName),
		zap.String("group", in.Group),
		zap.String("email", in.Email),
	)

	r.mu.Lock()
	defer r.mu.Unlock()

	email := strings.TrimSpace(strings.ToLower(in.Email))
	for _, st := range r.data {
		if strings.ToLower(st.Email) == email {
			return Student{}, ErrEmailExists
		}
	}

	st := Student{
		ID:       r.nextID,
		FullName: strings.TrimSpace(in.FullName),
		Group:    strings.TrimSpace(in.Group),
		Email:    email,
	}
	r.data[r.nextID] = st
	r.nextID++

	return st, nil
}
