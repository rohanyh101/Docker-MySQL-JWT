package main

import (
	"database/sql"
)

type Store interface {
	// Users
	CreateUser(u *User) (*User, error)
	GetUserByID(id string) (*User, error)

	// Tasks
	CreateTask(t *Task) (*Task, error)
	GetTask(id string) (*Task, error)

	// Project
	CreateProject(p *Project) (*Project, error)
	GetProjectByID(id string) (*Project, error)
	DeleteProject(id string) (int64, error)
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(u *User) (*User, error) {
	rows, err := s.db.Exec("INSERT INTO users (email, first_name, last_name, password) VALUES (?, ?, ?, ?)",
		u.Email, u.FirstName, u.LastName, u.Password)

	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.ID = id
	return u, nil
}

func (s *Storage) GetUserByID(id string) (*User, error) {
	var u User
	err := s.db.QueryRow("SELECT id, email, first_name, last_name, created_at FROM users WHERE id = ?", id).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt,
	)

	return &u, err
}

func (s *Storage) CreateTask(t *Task) (*Task, error) {
	rows, err := s.db.Exec("INSERT INTO tasks (name, project_id, assigned_to) VALUES (?, ?, ?)",
		t.Name, t.ProjectID, t.AssignedTo)

	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	t.ID = id
	return t, nil
}

func (s *Storage) GetTask(id string) (*Task, error) {
	var t Task
	err := s.db.QueryRow("SELECT id, name, status, project_id, assigned_to, created_at FROM tasks WHERE id = ?", id).Scan(
		&t.ID, &t.Name, &t.Status, &t.ProjectID, &t.AssignedTo, &t.CreatedAt,
	)
	return &t, err
}

// CreateProject implements Store.
func (s *Storage) CreateProject(p *Project) (*Project, error) {
	rows, err := s.db.Exec("INSERT INTO projects (name) VALUES (?)", p.Name)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	p.ID = id
	return p, nil
}

// DeleteProject implements Store.
func (s *Storage) DeleteProject(id string) (int64, error) {
	rows, err := s.db.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// GetProjectByID implements Store.
func (s *Storage) GetProjectByID(id string) (*Project, error) {
	var p Project
	err := s.db.QueryRow("SELECT id, name, created_at FROM projects WHERE id = ?", id).Scan(
		&p.ID, &p.Name, &p.CreatedAt,
	)

	return &p, err
}
