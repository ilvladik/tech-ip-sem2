package student

import "errors"

var ErrStudentNotFound = errors.New("student not found")

type Repo struct {
	data map[int64]Student
}

func NewRepo() *Repo {
	return &Repo{
		data: map[int64]Student{
			1: {ID: 1, FullName: "Иванов Иван Иванович", Group: "ИТТ-01-25", Email: "ivanov@example.com"},
			2: {ID: 2, FullName: "Петрова Мария Сергеевна", Group: "ИТТ-02-25", Email: "petrova@example.com"},
			3: {ID: 3, FullName: "Сидоров Алексей Дмитриевич", Group: "ИТТ-03-25", Email: "sidorov@example.com"},
		},
	}
}

func (r *Repo) GetByID(id int64) (Student, error) {
	st, ok := r.data[id]
	if !ok {
		return Student{}, ErrStudentNotFound
	}
	return st, nil
}
