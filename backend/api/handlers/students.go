// TODO add id return in GET request

package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"strings"

	"github.com/DanVerh/artschool-admin/backend/api/db"
	"github.com/DanVerh/artschool-admin/backend/api/errorHandling"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create struct (class) for StudentHandler to handle requests
type StudentHandler struct{}

// Create struct (class) for Student
type Student struct{
	Id 			 primitive.ObjectID	`json:"id" bson:"_id"`
    Fullname     string    `json:"fullname" bson:"fullname"`
    Phone        string    `json:"phone" bson:"phone"`
    Subscription *int       `json:"subscription" bson:"subscription"`
    StartDate    *time.Time `json:"startDate" bson:"startDate"`
    LastDate     *time.Time `json:"lastDate" bson:"lastDate"`
    Comments     *string    `json:"comments" bson:"comments"`
}

// Define all methods of Student as handlers for routes

// POST for student creation
func (studentHandler *StudentHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Create Student object
	student := &Student{}

	// Check if the method is POST; return 405 in case of error
	if r.Method != http.MethodPost {
		errorHandling.ThrowError(w, http.StatusMethodNotAllowed, "Invalid request method. Needs to be POST", nil)
		return
	}

	// Parse JSON request body to Student struct
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(student)
	// Check if parsing is correct; return 400 in case of error
	if err != nil {
		errorHandling.ThrowError(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}
	// Check if fullname and phone fields are passed in request
	if student.Fullname == "" || student.Phone == "" {
		errorHandling.ThrowError(w, http.StatusBadRequest, "Missing fullname or phone field", nil)
		return
	}

	// Define default properties of new student
	student.Id, student.Subscription, student.StartDate, student.LastDate, student.Comments = primitive.NewObjectID(),nil, nil, nil, nil

	// Connect to DB
	db := db.DbConnect()
	// Disconnect from the DB
	defer db.DbDisconnect()
	collection := db.Client.Database("artschool-admin").Collection("students")
	
	_, err = collection.InsertOne(nil, student)
	if err != nil {
		errorHandling.ThrowError(w, http.StatusInternalServerError, "Failed to insert the student into the database", &err)
		return
	}

	// Log the created student
	log.Printf("Created student: %v, %v\n", student.Fullname, student.Phone)

	// Respond with the created student data
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}


// GET for students list
func (studentHandler *StudentHandler) List(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET; return 405 in case of error
	if r.Method != http.MethodGet {
		errorHandling.ThrowError(w, http.StatusMethodNotAllowed, "Invalid request method. Needs to be GET", nil)
		return
	}

	// Connect to DB
	db := db.DbConnect()
	defer db.DbDisconnect()
	collection := db.Client.Database("artschool-admin").Collection("students")

	// Retrieve all documents without context
	cursor, err := collection.Find(nil, bson.M{})
	if err != nil {
		errorHandling.ThrowError(w, http.StatusInternalServerError, "Failed to retrieve documents from the database", &err)
		return
	}
	defer cursor.Close(nil)

	// Prepare a slice to hold the documents
	var students []Student

	// Iterate through the cursor and decode each document into a Student struct
	for cursor.Next(nil) {
		var student Student
		if err := cursor.Decode(&student); err != nil {
			errorHandling.ThrowError(w, http.StatusInternalServerError, "Failed to decode document", &err)
			return
		}
		students = append(students, student)
	}

	if err := cursor.Err(); err != nil {
		errorHandling.ThrowError(w, http.StatusInternalServerError, "Error occurred during cursor iteration", &err)
		return
	}

	// Respond with the list of students as JSON
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}


// GET for one student by ID
func (studentHandler *StudentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Check if the method is GET; return 405 in case of error
	if r.Method != http.MethodGet {
		errorHandling.ThrowError(w, http.StatusMethodNotAllowed, "Invalid request method. Needs to be GET", nil)
		return
	}

	// Extract the ObjectId from the URL path
    id := strings.TrimPrefix(r.URL.Path, "/students/")
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
	collection := db.Client.Database("artschool-admin").Collection("students")

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


// PUT for one student by ID
func (studentHandler *StudentHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	// Check if the method is PUT; return 405 in case of error
	if r.Method != http.MethodPut {
		errorHandling.ThrowError(w, http.StatusMethodNotAllowed, "Invalid request method. Needs to be PUT", nil)
		return
	}

	// Extract the ObjectId from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/students/")
	// Convert the string ID to a MongoDB ObjectId type
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errorHandling.ThrowError(w, http.StatusBadRequest, "Invalid ObjectId format", nil)
		return
	}

	var updateBody bson.M
    jsonDecoder := json.NewDecoder(r.Body)
    err = jsonDecoder.Decode(&updateBody)
    if err != nil {
		errorHandling.ThrowError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
    }

	// Connect to DB
	db := db.DbConnect()
	// Disconnect from the DB
	defer db.DbDisconnect()
	// Define collection
	collection := db.Client.Database("artschool-admin").Collection("students")

	// Find the record with required id
	updateResult, err := collection.UpdateByID(nil, objectID, bson.M{"$set": updateBody})
	if err != nil {
		errorHandling.ThrowError(w, http.StatusInternalServerError, "Failed to update student", &err)
		return
	}
	if updateResult.MatchedCount == 0 {
		errorHandling.ThrowError(w, http.StatusNotFound, "No record found with the provided ID", nil)
		return
	}

    // Check if any student field is updated and save these fields to slice
	var updateKeys []string
    for updateKey := range updateBody {
        if _, found := map[string]bool{"fullname": true, "phone": true, "subscription": true, "startDate": true, "lastDate": true, "comments": true}[updateKey]; !found {
			errorHandling.ThrowError(w, http.StatusBadRequest, "No student field is updated", nil)
			return
		}
		updateKeys = append(updateKeys, updateKey)
    }

	// Write the response with updated keys
	response := fmt.Sprintf("Student with id %v fields updated successfully: %v",id, updateKeys)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(response))
}


// DELETE for one student by ID
func (studentHandler *StudentHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	// Check if the method is Delete; return 405 in case of error
	if r.Method != http.MethodDelete {
		errorHandling.ThrowError(w, http.StatusMethodNotAllowed, "Invalid request method. Needs to be Delete", nil)
		return
	}

	// Extract the ObjectId from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/students/")
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
	collection := db.Client.Database("artschool-admin").Collection("students")

	// Delete record with mentioned id
	deleteResult, err := collection.DeleteOne(nil, bson.M{"_id": objectID})
	if err != nil {
		errorHandling.ThrowError(w, http.StatusInternalServerError, "Failed to delete schedule", &err)
		return
	} 
	if deleteResult.DeletedCount == 0 {
		errorHandling.ThrowError(w, http.StatusInternalServerError, fmt.Sprintf("No student found with the provided ID: %v", id), nil)
		return
	}

	// Write the response with deleted student id
	response := fmt.Sprintf("Deleted student by mentioned id: %v", id)
	w.WriteHeader(http.StatusOK)
    w.Write([]byte(response))
}
