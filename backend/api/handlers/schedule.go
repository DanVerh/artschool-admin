package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DanVerh/artschool-admin/backend/api/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create struct (class) for StudentHandler to handle requests
type ScheduleHandler struct{}

// Create struct (class) for Classes that will be added to Schedule
type Class struct {
	StudentId  primitive.ObjectID `json:"studentId" bson:"studentId"`
	Time       string             `json:"time" bson:"time"`
	Type       string             `json:"type" bson:"type"`
	Attendence *bool              `json:"attendence" bson:"attendence"`
}

// Create struct (class) for Schedule
type Schedule struct {
	Date    primitive.DateTime `bson:"date" json:"date"`
	Classes []Class            `bson:"classes" json:"classes"`
}

// Define all methods of Schedule as handlers for routes

// POST for schedule creation
func (scheduleHandler *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Create Schedule object
	schedule := &Schedule{}

	// Check if the method is POST; return 405 in case of error
	if r.Method != http.MethodPost {
		errorMessage := "Invalid request method. Needs to be POST"
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body to Schedule struct
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&schedule)
	// Check if parsing is correct; return 400 in case of error
	if err != nil {
		errorMessage := "Invalid JSON"
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	// Connect to DB
	db := db.DbConnect()
	// Disconnect from the DB
	defer db.DbDisconnect()
	collection := db.Client.Database("artschool-admin").Collection("schedule")

	// Insert schedule object to schedule collection in mongo
	_, err = collection.InsertOne(nil, schedule)
	if err != nil {
		log.Printf("Failed to insert document: %v", err)
	    http.Error(w, "Failed to insert the schedule into the database", http.StatusInternalServerError)
	    return
	}

	// Log the created schedule
	log.Printf("Created schedule")

	// Respond with the created student data
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
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
