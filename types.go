package main

import "time"

type ErrorResponse struct {
	Error string `json:"error"`
}

type Task struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	ProjectID  int64     `json:"project_id"`
	AssignedTo int64     `json:"assigned_to"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateTaskPayload struct {
	Name       string `json:"name"`
	ProjectID  int64  `json:"project_id"`
	AssignedTo int64  `json:"assigned_to"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type RegisterUserPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Project struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateProjectPayload struct {
	Name string `json:"name"`
}
