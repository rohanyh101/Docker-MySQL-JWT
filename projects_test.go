package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateProject(t *testing.T) {
	ms := &MockStore{}
	service := NewProjectService(ms)

	t.Run("shoud create a project", func(t *testing.T) {
		payload := &CreateProjectPayload{
			Name: "NO_NAME",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/projects", service.HandleProjectCreate)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Error("invalid status code, it should pass")
		}
	})

	t.Run("should validate if name is a empty field", func(t *testing.T) {
		payload := &CreateProjectPayload{
			Name: "",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/projects", service.HandleProjectCreate)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		// var response ErrorResponse
		// err = json.NewDecoder(rr.Body).Decode(&response)
		// if err != nil {
		// 	t.Fatal(err)
		// }

		// if response.Error != errProjectNameRequired.Error() {
		// 	t.Errorf("expected error message %s, got %s", errProjectNameRequired, response.Error)
		// }
	})
}

func TestGetProject(t *testing.T) {

	ms := &MockStore{}
	service := NewProjectService(ms)

	t.Run("should return a project", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/projects/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/projects/{project_id}", service.HandleProjectGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Error("invalid status code, it should pass")
		}

		// var response Project
		// err = json.NewDecoder(rr.Body).Decode(&response)
		// if err != nil {
		// 	t.Fatal(err)
		// }

		// if response.Name != "NO_NAME" {
		// 	t.Errorf("expected project name %s, got %s", "my new project", response.Name)
		// }
	})
}

func TestDeleteProject(t *testing.T) {

	ms := &MockStore{}
	service := NewProjectService(ms)

	t.Run("should delete the project", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, "/projects", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("/projects", service.HandleProjectDelete)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusNoContent, rr.Code)
		}
	})
}
