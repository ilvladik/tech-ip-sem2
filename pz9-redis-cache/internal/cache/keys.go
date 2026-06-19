package cache

import "fmt"

type KeyBuilder struct{}

func NewKeyBuilder() KeyBuilder {
	return KeyBuilder{}
}

func (KeyBuilder) TaskByID(id int64) string {
	return fmt.Sprintf("tasks:task:%d", id)
}

func (KeyBuilder) TasksList() string {
	return "tasks:list"
}

func (KeyBuilder) TasksListPaged(page, limit int) string {
	return fmt.Sprintf("tasks:list:page=%d:limit=%d", page, limit)
}

// Backward-compatible helpers.
func TaskByIDKey(id int64) string {
	return NewKeyBuilder().TaskByID(id)
}

func TasksListKey() string {
	return NewKeyBuilder().TasksList()
}

func TasksListPagedKey(page, limit int) string {
	return NewKeyBuilder().TasksListPaged(page, limit)
}
