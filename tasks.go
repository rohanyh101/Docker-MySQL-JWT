package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var errTaskNameRequired = errors.New("name is required")
var errProjectIDRequired = errors.New("project id is required")
var errUserIDRequired = errors.New("user id is required")

type TasksService struct {
	store Store
}

func NewTasksService(s Store) *TasksService {
	return &TasksService{
		store: s,
	}
}

func (s *TasksService) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("POST /tasks", WithJWTAuth(s.HandleCreateTask, s.store))
	r.HandleFunc("GET /tasks/{task_id}", WithJWTAuth(s.HandleGetTask, s.store))
}

func (s *TasksService) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "error reading request body"})
		return
	}

	defer r.Body.Close()

	var task *Task
	if err := json.Unmarshal(body, &task); err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "error unmarshalling task payload"})
		return
	}

	if err := validateTaskPayload(task); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "error validating task payload" + err.Error()})
		return
	}

	t, err := s.store.CreateTask(task)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "error creating task: " + err.Error()})
		return
	}

	WriteJSON(w, http.StatusCreated, t)
}

func (s *TasksService) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("task_id")
	if id == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "task id is required"})
		return
	}

	t, err := s.store.GetTask(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "error getting task: " + err.Error()})
		return
	}

	WriteJSON(w, http.StatusOK, t)
}

func validateTaskPayload(task *Task) error {
	if task.Name == "" {
		return errTaskNameRequired
	}

	if task.ProjectID == 0 {
		return errProjectIDRequired
	}

	if task.AssignedTo == 0 {
		return errUserIDRequired
	}

	return nil
}
