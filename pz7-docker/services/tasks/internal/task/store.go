package task

type Item struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func List() []Item {
	return []Item{
		{ID: 1, Title: "Read Dockerfile guide", Status: "done"},
		{ID: 2, Title: "Build multi-stage image", Status: "in_progress"},
		{ID: 3, Title: "Run docker compose", Status: "todo"},
	}
}
