package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/DanVerh/artschool-admin/backend/api/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create struct (class) for StudentHandler to handle requests
type ScheduleHandler struct{}

// Create struct (class) for Classes that will be added to Schedule
type Class struct {
	StudentId  primitive.ObjectID `json:"studentId" bson:"studentId"`
	Time       string             `json:"time" bson:"time"`
	Type       string             `json:"type" bson:"type"`
	Attendence *bool              `json:"attendance" bson:"attendance"`
}

// Create struct (class) for Schedule
type Schedule struct {
	Id      primitive.ObjectID `json:"id" bson:"_id"`
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

	// Check if classes array is not empty; return 400 in case of error
	if len(schedule.Classes) == 0 {
		errorMessage := "No classes found for schedule creation"
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	// Create primitive object id in mongo for schedule
	schedule.Id = primitive.NewObjectID()

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

// GET for schedules list
func (scheduleHandler *ScheduleHandler) List(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET; return 405 in case of error
	if r.Method != http.MethodGet {
		errorMessage := "Invalid request method. Needs to be GET"
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusMethodNotAllowed)
		return
	}

	// Connect to DB
	db := db.DbConnect()
	defer db.DbDisconnect()
	collection := db.Client.Database("artschool-admin").Collection("schedule")

	// Retrieve all documents without context
	cursor, err := collection.Find(nil, bson.M{})
	if err != nil {
		log.Printf("Failed to retrieve documents: %v", err)
		http.Error(w, "Failed to retrieve documents from the database", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(nil)

	// Prepare a slice to hold the documents
	var schedules []Schedule

	// Iterate through the cursor and decode each document into a Schedule struct
	for cursor.Next(nil) {
		var schedule Schedule
		if err := cursor.Decode(&schedule); err != nil {
			log.Printf("Failed to decode document: %v", err)
			http.Error(w, "Failed to decode document", http.StatusInternalServerError)
			return
		}
		schedules = append(schedules, schedule)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		http.Error(w, "Error occurred during cursor iteration", http.StatusInternalServerError)
		return
	}

	// Respond with the list of schedules as JSON
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedules)
}

// GET for one schedule by ID
func (scheduleHandler *ScheduleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET; return 405 in case of error
	if r.Method != http.MethodGet {
		errorMessage := "Invalid request method. Needs to be GET"
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusMethodNotAllowed)
		return
	}

	// Extract the ObjectId from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/schedule/")
	// Convert the string ID to a MongoDB ObjectId type
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ObjectId format", http.StatusBadRequest)
		return
	}
	// Create a filter to search for the document with this ObjectId
	filter := bson.M{"_id": objectID}
	var result bson.M

	// Connect to DB
	db := db.DbConnect()
	// Disconnect from the DB
	defer db.DbDisconnect()
	// Define collection
	collection := db.Client.Database("artschool-admin").Collection("schedule")

	// Find the record with required id
	err = collection.FindOne(nil, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No document found with the given ObjectId", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve document", http.StatusInternalServerError)
		}
		return
	}

	// Set the response header to JSON and encode the result
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (scheduleHandler *ScheduleHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update a schedule by ID")
}

func (scheduleHandler *ScheduleHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete a schedule by ID")
}
