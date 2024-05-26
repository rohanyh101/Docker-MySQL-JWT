package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type APIServer struct {
	addr  string
	store Store
}

func NewAPIServer(addr string, store Store) *APIServer {
	return &APIServer{
		addr:  addr,
		store: store,
	}
}

func (s *APIServer) Run() {
	router := http.NewServeMux()
	server := http.Server{
		Addr:         s.addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	subRouter := http.NewServeMux()
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", subRouter))

	// registering services...

	// task service...
	tasksService := NewTasksService(s.store)
	tasksService.RegisterRoutes(subRouter)

	// user service...
	userService := NewUserService(s.store)
	userService.RegisterRoutes(subRouter)

	// project service...
	projectService := NewProjectService(s.store)
	projectService.RegisterRoutes(subRouter)

	// health check route...
	// route "GET /" is not working !!!
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{"message": "server is up and running..."}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})

	log.Println("Server is listening on port", s.addr)
	log.Fatal(server.ListenAndServe())
}
