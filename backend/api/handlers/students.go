package handler

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"
	"log"

	//"go.mongodb.org/mongo-driver/bson"
    //"go.mongodb.org/mongo-driver/mongo"
    //"go.mongodb.org/mongo-driver/mongo/options"
    //"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Create struct (class) for Student
type Student struct{
    Fullname     string    `json:"fullname" bson:"fullname"`
    Phone        string    `json:"phone" bson:"phone"`
    Subscription int       `json:"subscription" bson:"subscription"`
    StartDate    time.Time `json:"startDate" bson:"startDate"`
    LastDate     time.Time `json:"lastDate" bson:"lastDate"`
    Comments     string    `json:"comments" bson:"comments"`
}

// Define all methods of Student as handlers for routes

// POST for student creation
func (student *Student) Create(w http.ResponseWriter, r *http.Request) {
	newStudent := &Student{}
	// Check if the method is POST; return 405 in case of error
	if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method. Needs to be POST", http.StatusMethodNotAllowed)
        return
    }

	// Parse JSON request body to Student struct
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(newStudent)
	// Check if parsing is correct; return 400 in case of error
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Check if fullname and phone fields are passed in request
	if newStudent.Fullname == "" || newStudent.Phone == "" {
		http.Error(w, "Missing fullname or phone field", http.StatusBadRequest)
		return
	}

	// Log the created student
	log.Printf("Created student: %v, %v\n", newStudent.Fullname, newStudent.Phone)
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
