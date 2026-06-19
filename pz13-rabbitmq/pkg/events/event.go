package events

import "time"

type TaskEvent struct {
	Event     string `json:"event"`
	TaskID    string `json:"task_id"`
	TS        string `json:"ts"`
	RequestID string `json:"request_id,omitempty"`
	Producer  string `json:"producer,omitempty"`
	Version   string `json:"version,omitempty"`
}

func NewTaskCreated(taskID, requestID, producer string) TaskEvent {
	return TaskEvent{
		Event:     "task.created",
		TaskID:    taskID,
		TS:        time.Now().UTC().Format(time.RFC3339),
		RequestID: requestID,
		Producer:  producer,
		Version:   "1",
	}
}
