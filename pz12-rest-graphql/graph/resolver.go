package graph

import "example.com/pz12-rest-graphql/internal/service"

type Resolver struct {
	TaskSvc *service.TaskService
}
