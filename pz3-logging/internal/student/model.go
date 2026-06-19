package student

type Student struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Group    string `json:"group"`
	Email    string `json:"email"`
}

type CreateInput struct {
	FullName string `json:"full_name"`
	Group    string `json:"group"`
	Email    string `json:"email"`
}
