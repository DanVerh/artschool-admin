package handler

import (
	"fmt"
	"net/http"
)

// Create struct (class) for Student
type Student struct{}

// Define all methods of Strudent as handlers for routes
func (student *Student) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create a student")
}

func (student *Student) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all students")
}

func (student *Student) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get a student by ID")
}

func (student *Student) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update a student by ID")
}

func (student *Student) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete a student by ID")
}
