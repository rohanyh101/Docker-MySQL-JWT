package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var errProjectNameRequired = errors.New("project name is required")

type ProjectService struct {
	store Store
}

func NewProjectService(s Store) *ProjectService {
	return &ProjectService{
		store: s,
	}
}

func (s *ProjectService) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("POST /projects", WithJWTAuth(s.HandleProjectCreate, s.store))
	r.HandleFunc("GET /projects/{project_id}", WithJWTAuth(s.HandleProjectGet, s.store))
	r.HandleFunc("DELETE /projects/{project_id}", WithJWTAuth(s.HandleProjectDelete, s.store))
}

func (s *ProjectService) HandleProjectCreate(w http.ResponseWriter, r *http.Request) {
	// read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "error reading request body"})
		return
	}

	defer r.Body.Close()

	// unmarshal request body into project struct
	var payload *Project
	err = json.Unmarshal(body, &payload)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "invalid request payload"})
		return
	}

	// validate project payload
	if err := validateProjectPayload(payload); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "error validating project payload"})
		return
	}

	// call store.CreateProject
	p, err := s.store.CreateProject(payload)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "error creating project"})
		return
	}

	// write response
	WriteJSON(w, http.StatusCreated, p)
}

func (s *ProjectService) HandleProjectGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("project_id")
	if id == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "field id is missing"})
		return
	}

	p, err := s.store.GetProjectByID(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "error getting project"})
		return
	}

	WriteJSON(w, http.StatusOK, p)
}

func (s *ProjectService) HandleProjectDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("project_id")
	if id == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "field id is missing"})
		return
	}

	d, err := s.store.DeleteProject(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "error deleting project"})
		return
	}

	WriteJSON(w, http.StatusNoContent, d)
}

func validateProjectPayload(p *Project) error {
	if p.Name == "" {
		return errProjectNameRequired
	}

	return nil
}
