package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

func (s *APIServer) Serve() {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// Register Services
	usersService := NewUserService(s.store)
	usersService.RegisterRoutes(subRouter)

	tasksService := NewTasksService(s.store)
	tasksService.RegisterRoutes(subRouter)

	fmt.Println("Starting the server at", s.addr)

	log.Fatal(http.ListenAndServe(s.addr, subRouter))
}
