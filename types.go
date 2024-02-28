package main

import "time"

type ErrorResponse struct {
	Error string `json:"error"`
}

type Task struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	ProjectID    int       `json:"projectID"`
	AssignedToID int       `json:"assignedTo"`
	CreatedAt    time.Time `json:"createdAt"`
}
