package models

type AssignedTask struct {
	Email string `json:"email"`
	Task  Task   `json:"task"`
}
