package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DanVerh/artschool-admin/backend/api/db"
)

// Create struct (class) for Student
type Student struct{
    Fullname     string    `json:"fullname" bson:"fullname"`
    Phone        string    `json:"phone" bson:"phone"`
    Subscription *int       `json:"subscription" bson:"subscription"`
    StartDate    *time.Time `json:"startDate" bson:"startDate"`
    LastDate     *time.Time `json:"lastDate" bson:"lastDate"`
    Comments     *string    `json:"comments" bson:"comments"`
}

// Define all methods of Student as handlers for routes

// POST for student creation
func (student *Student) Create(w http.ResponseWriter, r *http.Request) {
	newStudent := &Student{}
	// Check if the method is POST; return 405 in case of error
	if r.Method != http.MethodPost {
		errorMessage := "Invalid request method. Needs to be POST"
		log.Println(errorMessage)
        http.Error(w, errorMessage, http.StatusMethodNotAllowed)
        return
    }

	// Parse JSON request body to Student struct
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(newStudent)
	// Check if parsing is correct; return 400 in case of error
	if err != nil {
		errorMessage := "Invalid JSON"
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}
	// Check if fullname and phone fields are passed in request
	if newStudent.Fullname == "" || newStudent.Phone == "" {
		errorMessage := "Missing fullname or phone field"
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}
	// Define default properties of new student
	newStudent.Subscription, newStudent.StartDate, newStudent.LastDate, newStudent.Comments = nil, nil, nil, nil

	// Connect to DB
	db := db.DbConnect()
	collection := db.Client.Database("artschool-admin").Collection("students")
	
	_, err = collection.InsertOne(nil, newStudent)
	if err != nil {
		log.Printf("Failed to insert document: %v", err)
	    http.Error(w, "Failed to insert the student into the database", http.StatusInternalServerError)
	    return
	}

	// Log the created student
	log.Printf("Created student: %v, %v\n", newStudent.Fullname, newStudent.Phone)

	// Disconnect from the DB
	db.DbDisconnect()

	// Respond with the created student data
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newStudent)
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
