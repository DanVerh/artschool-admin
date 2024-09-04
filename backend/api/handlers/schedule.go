package handler

import (
	"fmt"
	"net/http"
)

// Create struct (class) for StudentHandler to handle requests
type ScheduleHandler struct{}

// Create struct (class) for Schedule
type Schedule struct{}

// Define all methods of Schedule as handlers for routes
func (scheduleHandler *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create a schedule")
}

func (scheduleHandler *ScheduleHandler) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all schedules")
}

func (scheduleHandler *ScheduleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get a schedule by ID")
}

func (scheduleHandler *ScheduleHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update a schedule by ID")
}

func (scheduleHandler *ScheduleHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete a schedule by ID")
}
