package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/DanVerh/artschool-admin/backend/api/db"
	"github.com/DanVerh/artschool-admin/backend/api/errorHandling"
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
		errorHandling.ThrowError(w, http.StatusMethodNotAllowed, "Invalid request method. Needs to be POST", nil)
		return
	}

	// Parse JSON request body to Schedule struct
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&schedule)
	// Check if parsing is correct; return 400 in case of error
	if err != nil {
		errorHandling.ThrowError(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	// Check if classes array is not empty; return 400 in case of error
	if len(schedule.Classes) == 0 {
		errorHandling.ThrowError(w, http.StatusBadRequest, "No classes found for schedule creation", nil)
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
		errorHandling.ThrowError(w, http.StatusInternalServerError, "Failed to insert the schedule into the database", &err)
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
		errorHandling.ThrowError(w, http.StatusMethodNotAllowed, "Invalid request method. Needs to be GET", nil)
		return
	}

	// Connect to DB
	db := db.DbConnect()
	defer db.DbDisconnect()
	collection := db.Client.Database("artschool-admin").Collection("schedule")

	// Retrieve all documents without context
	cursor, err := collection.Find(nil, bson.M{})
	if err != nil {
		errorHandling.ThrowError(w, http.StatusInternalServerError, "Failed to retrieve documents from the database", &err)
		return
	}
	defer cursor.Close(nil)

	// Prepare a slice to hold the documents
	var schedules []Schedule

	// Iterate through the cursor and decode each document into a Schedule struct
	for cursor.Next(nil) {
		var schedule Schedule
		if err := cursor.Decode(&schedule); err != nil {
			errorHandling.ThrowError(w, http.StatusInternalServerError, "Failed to decode document", &err)
			return
		}
		schedules = append(schedules, schedule)
	}

	if err := cursor.Err(); err != nil {
		errorHandling.ThrowError(w, http.StatusInternalServerError, "Error occurred during cursor iteration", &err)
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
		errorHandling.ThrowError(w, http.StatusMethodNotAllowed, "Invalid request method. Needs to be GET", nil)
		return
	}

	// Extract the ObjectId from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/schedule/")
	// Convert the string ID to a MongoDB ObjectId type
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorHandling.ThrowError(w, http.StatusBadRequest, "Invalid ObjectId format", nil)
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
			errorHandling.ThrowError(w, http.StatusNotFound, "No document found with the given ObjectId", nil)
		} else {
			errorHandling.ThrowError(w, http.StatusInternalServerError, "Failed to retrieve document", &err)
		}
		return
	}

	// Set the response header to JSON and encode the result
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// PUT for schedule classes update
func (scheduleHandler *ScheduleHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	// Check if the method is PUT; return 405 in case of error
	if r.Method != http.MethodPut {
		errorHandling.ThrowError(w, http.StatusMethodNotAllowed, "Invalid request method. Needs to be PUT", nil)
		return
	}

	// Extract the ObjectId from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/schedule/")
	// Convert the string ID to a MongoDB ObjectId type
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorHandling.ThrowError(w, http.StatusBadRequest, "Invalid ObjectId format", nil)
		return
	}

	// GET CURRENT SCHEDULE

	// Create a filter to search for the document with this ObjectId
	filter := bson.M{"_id": objectID}
	var currentSchedule Schedule

	// Connect to DB
	db := db.DbConnect()
	// Disconnect from the DB
	defer db.DbDisconnect()
	// Define collection
	collection := db.Client.Database("artschool-admin").Collection("schedule")

	// Find the record with required id
	err = collection.FindOne(nil, filter).Decode(&currentSchedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			errorHandling.ThrowError(w, http.StatusNotFound, "No document found with the given ObjectId", nil)
		} else {
			errorHandling.ThrowError(w, http.StatusInternalServerError, "Failed to retrieve document", nil)
		}
		return
	}

	// Get current schedule student ids
	var currentStudentIds []primitive.ObjectID
	for classIndex := range currentSchedule.Classes {
		currentStudentIds = append(currentStudentIds, currentSchedule.Classes[classIndex].StudentId)
	}

	// UPDATE THE CURRENT SCHEDULE

	// Decode request body to schedule object
	updatedClass := Class{}
	jsonDecoder := json.NewDecoder(r.Body)
	err = jsonDecoder.Decode(&updatedClass)
	if err != nil {
		errorHandling.ThrowError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Check if class is already booked for this student
	var studentClassExists bool
	var updatedClassIndex int
	for index, scheduledStudentId := range currentStudentIds {
		if scheduledStudentId == updatedClass.StudentId {
			studentClassExists = true
			updatedClassIndex = index
			break
		}
		studentClassExists = false
	}
	log.Println(studentClassExists, updatedClassIndex)

	if !(studentClassExists) {
		currentSchedule.Classes = append(currentSchedule.Classes, updatedClass)
	} else {
		currentSchedule.Classes[updatedClassIndex] = updatedClass
	}

	log.Println(currentSchedule)

	// Find the record with required id
	updateResult, err := collection.UpdateByID(nil, objectID, bson.M{"$set": currentSchedule})
	if err != nil {
		errorHandling.ThrowError(w, http.StatusInternalServerError, "Failed to update schedule", &err)
		return
	}
	if updateResult.MatchedCount == 0 {
		errorHandling.ThrowError(w, http.StatusNotFound, "No record found with the provided ID", nil)
		return
	}

	// Write the response with updated keys
	response := fmt.Sprintf("Schedule updated successfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// DELETE for deleting schedule
func (scheduleHandler *ScheduleHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	// Check if the method is Delete; return 405 in case of error
	if r.Method != http.MethodDelete {
		errorHandling.ThrowError(w, http.StatusMethodNotAllowed, "Invalid request method. Needs to be Delete", nil)
		return
	}

	// Extract the ObjectId from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/schedule/")
	// Convert the string ID to a MongoDB ObjectId type
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorHandling.ThrowError(w, http.StatusBadRequest, "Invalid ObjectId format", nil)
		return
	}

	// Connect to DB
	db := db.DbConnect()
	// Disconnect from the DB
	defer db.DbDisconnect()
	// Define collection
	collection := db.Client.Database("artschool-admin").Collection("schedule")

	// Delete record with mentioned id
	deleteResult, err := collection.DeleteOne(nil, bson.M{"_id": objectID})
	if err != nil {
		errorHandling.ThrowError(w, http.StatusInternalServerError, "Failed to delete schedule", &err)
		return
	}
	if deleteResult.DeletedCount == 0 {
		errorHandling.ThrowError(w, http.StatusInternalServerError, fmt.Sprintf("No schedule found with the provided ID: %v", id), nil)
		return
	}

	// Write the response with deleted schedule id
	response := fmt.Sprintf("Deleted schedule by mentioned id: %v", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
