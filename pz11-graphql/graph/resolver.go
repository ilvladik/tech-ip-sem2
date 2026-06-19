package graph

import "example.com/pz11-graphql/internal/service"

type Resolver struct {
	TaskSvc *service.TaskService
}
