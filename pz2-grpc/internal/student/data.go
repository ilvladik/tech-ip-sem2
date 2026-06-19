package student

import (
	"errors"
	"strings"
	"sync"

	"example.com/pz2-grpc/gen/studentpb"
)

var (
	ErrStudentNotFound = errors.New("student not found")
	ErrEmailExists     = errors.New("email already exists")
)

type Repository struct {
	mu     sync.RWMutex
	data   map[int64]*studentpb.Student
	nextID int64
}

func NewRepository() *Repository {
	return &Repository{
		nextID: 4,
		data: map[int64]*studentpb.Student{
			1: {
				Id:              1,
				FullName:        "Иванов Иван Иванович",
				Group:           "ИТТ-01-25",
				Email:           "ivanov@example.com",
				Specialization:  "Программная инженерия",
			},
			2: {
				Id:              2,
				FullName:        "Петрова Мария Сергеевна",
				Group:           "ИТТ-02-25",
				Email:           "petrova@example.com",
				Specialization:  "Информационная безопасность",
			},
			3: {
				Id:              3,
				FullName:        "Сидоров Алексей Дмитриевич",
				Group:           "ИТТ-03-25",
				Email:           "sidorov@example.com",
				Specialization:  "Прикладная информатика",
			},
		},
	}
}

func (r *Repository) GetByID(id int64) (*studentpb.Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	st, ok := r.data[id]
	if !ok {
		return nil, ErrStudentNotFound
	}
	return st, nil
}

func (r *Repository) ListAll() []*studentpb.Student {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*studentpb.Student, 0, len(r.data))
	for _, st := range r.data {
		list = append(list, st)
	}
	return list
}

func (r *Repository) Create(fullName, group, email, specialization string) (*studentpb.Student, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	email = strings.TrimSpace(strings.ToLower(email))
	for _, st := range r.data {
		if strings.ToLower(st.GetEmail()) == email {
			return nil, ErrEmailExists
		}
	}

	st := &studentpb.Student{
		Id:             r.nextID,
		FullName:       fullName,
		Group:          group,
		Email:          email,
		Specialization: specialization,
	}
	r.data[r.nextID] = st
	r.nextID++

	return st, nil
}
