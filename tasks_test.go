package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTask(t *testing.T) {

	ms := &MockStore{}
	service := NewTasksService(ms)

	t.Run("should return error if task name is missing", func(t *testing.T) {
		payload := &CreateTaskPayload{
			Name: "",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/tasks", service.HandleCreateTask)

		router.ServeHTTP(rr, req)

		// it should be failing case cuz, name is empty...
		if rr.Code != http.StatusBadRequest {
			t.Error("invalid status code, it should pass")
		}

		// var response ErrorResponse
		// err = json.NewDecoder(rr.Body).Decode(&response)
		// if err != nil {
		// 	t.Fatal(err)
		// }

		// if response.Error != errTaskNameRequired.Error() {
		// 	t.Errorf("expected error message %s, got %s", errTaskNameRequired, response.Error)
		// }
	})

	t.Run("should create a task", func(t *testing.T) {
		payload := &CreateTaskPayload{
			Name:       "REST API IN GO",
			ProjectID:  3,
			AssignedTo: 26,
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/tasks", service.HandleCreateTask)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}
