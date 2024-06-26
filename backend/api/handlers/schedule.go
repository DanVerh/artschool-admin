package handler

import (
	"fmt"
	"net/http"
)

// Create struct (class) for Schedule
type Schedule struct{}

// Define all methods of Schedule as handlers for routes
func (schedule *Schedule) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create a schedule")
}

func (schedule *Schedule) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all schedules")
}

func (schedule *Schedule) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get a schedule by ID")
}

func (schedule *Schedule) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update a schedule by ID")
}

func (schedule *Schedule) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete a schedule by ID")
}
