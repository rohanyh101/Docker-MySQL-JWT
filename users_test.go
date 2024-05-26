package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	ms := &MockStore{}
	service := NewUserService(ms)

	t.Run("should validate if email is empty", func(t *testing.T) {
		payload := &RegisterUserPayload{
			FirstName: "bob",
			LastName:  "cj",
			Email:     "",
			Password:  "5Vi64w^&",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("POST /users/register", service.HandleUserRegister)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		var response ErrorResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		if response.Error != errEmailRequired.Error() {
			t.Errorf("expected error message %s, got %s", response.Error, errEmailRequired.Error())
		}
	})

	t.Run("should create a user", func(t *testing.T) {
		payload := &RegisterUserPayload{
			FirstName: "bob",
			LastName:  "cj",
			Email:     "bob@gmail.com",
			Password:  "5Vi64w^&",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := http.NewServeMux()

		router.HandleFunc("POST /users/register", service.HandleUserRegister)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})

}

func TestValidateUserPayload(t *testing.T) {
	type args struct {
		user *User
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "should return error if email is empty",
			args: args{
				user: &User{
					FirstName: "manoj",
					LastName:  "g",
				},
			},
			want: errEmailRequired,
		},
		{
			name: "should return error if first name is empty",
			args: args{
				user: &User{
					LastName: "g",
					Email:    "manoj@gmail.com",
				},
			},
			want: errFirstNameRequired,
		},
		{
			name: "should return error if last name is empty",
			args: args{
				user: &User{
					FirstName: "manoj",
					Email:     "manoj@gmail.com",
				},
			},
			want: errLastNameRequired,
		},
		{
			name: "should return error if password is empty",
			args: args{
				user: &User{
					FirstName: "manoj",
					LastName:  "g",
					Email:     "manoj@gmail.com",
				},
			},
			want: errPasswordRequired,
		},
		{
			name: "should return nil if all the fields are present",
			args: args{
				user: &User{
					FirstName: "manoj",
					LastName:  "g",
					Email:     "manoj@gmail.com",
					Password:  "wf#$w7",
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateUserPayload(tt.args.user); got != tt.want {
				t.Errorf("validateUserPayload() = %v, want %v", got, tt.want)
			}
		})
	}

}
