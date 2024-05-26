package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var errEmailRequired = errors.New("email is required")
var errFirstNameRequired = errors.New("first name is required")
var errLastNameRequired = errors.New("last name is required")
var errPasswordRequired = errors.New("password is required")

type UserService struct {
	store Store
}

func NewUserService(s Store) *UserService {
	return &UserService{
		store: s,
	}
}

func (s *UserService) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("POST /users/register", s.HandleUserRegister)
	r.HandleFunc("GET /users/{user_id}", WithJWTAuth(s.HandleUserGet, s.store))
}

func (s *UserService) HandleUserGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("user_id")
	if id == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "missing field user id"})
		return
	}

	u, err := s.store.GetUserByID(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "error getting user by id"})
		return
	}

	WriteJSON(w, http.StatusOK, u)
}

func (s *UserService) HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var payload *User
	err = json.Unmarshal(body, &payload)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := validateUserPayload(payload); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	hashedPWD, err := HashPassword(payload.Password)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "error hashing password"})
		return
	}

	payload.Password = hashedPWD
	u, err := s.store.CreateUser(payload)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "error creating user"})
		return
	}

	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "error while setting cookie"})
		return
	}

	response := map[string]string{"token": token}
	WriteJSON(w, http.StatusCreated, response)
}

func validateUserPayload(user *User) error {
	if user.Email == "" {
		return errEmailRequired
	}

	// both first_name and last_name, you can make them optional....
	if user.FirstName == "" {
		return errFirstNameRequired
	}

	if user.LastName == "" {
		return errLastNameRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	secret := []byte(Envs.JWTSecret)
	token, err := CreateJWT(secret, id)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}
