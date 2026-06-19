package cache

import (
	"encoding/json"
	"errors"

	"example.com/pz9-redis-cache/internal/task"
)

const notFoundMarker = "__NOT_FOUND__"

type Serializer struct{}

func NewSerializer() Serializer {
	return Serializer{}
}

func (Serializer) MarshalTask(t task.Task) ([]byte, error) {
	return json.Marshal(t)
}

func (Serializer) UnmarshalTask(data []byte) (task.Task, error) {
	var t task.Task
	if err := json.Unmarshal(data, &t); err != nil {
		return task.Task{}, err
	}
	return t, nil
}

func (Serializer) MarshalTaskList(items []task.Task) ([]byte, error) {
	return json.Marshal(items)
}

func (Serializer) UnmarshalTaskList(data []byte) ([]task.Task, error) {
	var items []task.Task
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (Serializer) MarshalNotFound() []byte {
	return []byte(notFoundMarker)
}

func (Serializer) IsNotFound(data []byte) bool {
	return string(data) == notFoundMarker
}

var ErrNegativeCacheHit = errors.New("negative cache hit")
