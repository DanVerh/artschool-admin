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
	student := &Student{}
	// Check if the method is POST; return 405 in case of error
	if r.Method != http.MethodPost {
		errorMessage := "Invalid request method. Needs to be POST"
		log.Println(errorMessage)
        http.Error(w, errorMessage, http.StatusMethodNotAllowed)
        return
    }

	// Parse JSON request body to Student struct
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(student)
	// Check if parsing is correct; return 400 in case of error
	if err != nil {
		errorMessage := "Invalid JSON"
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}
	// Check if fullname and phone fields are passed in request
	if student.Fullname == "" || student.Phone == "" {
		errorMessage := "Missing fullname or phone field"
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
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
		log.Printf("Failed to insert document: %v", err)
	    http.Error(w, "Failed to insert the student into the database", http.StatusInternalServerError)
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
		errorMessage := "Invalid request method. Needs to be GET"
		log.Println(errorMessage)
        http.Error(w, errorMessage, http.StatusMethodNotAllowed)
        return
    }

	// Connect to DB
	db := db.DbConnect()
	collection := db.Client.Database("artschool-admin").Collection("students")

	// Retrieve all documents without context
	cursor, err := collection.Find(nil, bson.M{})
	if err != nil {
		log.Printf("Failed to retrieve documents: %v", err)
		http.Error(w, "Failed to retrieve documents from the database", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(nil)

	// Prepare a slice to hold the documents
	var students []Student

	// Iterate through the cursor and decode each document into a Student struct
	for cursor.Next(nil) {
		var student Student
		if err := cursor.Decode(&student); err != nil {
			log.Printf("Failed to decode document: %v", err)
			http.Error(w, "Failed to decode document", http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		http.Error(w, "Error occurred during cursor iteration", http.StatusInternalServerError)
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
		errorMessage := "Invalid request method. Needs to be GET"
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusMethodNotAllowed)
		return
	}

	// Extract the ObjectId from the URL path
    id := strings.TrimPrefix(r.URL.Path, "/students/")
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
	collection := db.Client.Database("artschool-admin").Collection("students")

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


// PUT for one student by ID
func (studentHandler *StudentHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	// Check if the method is PUT; return 405 in case of error
	if r.Method != http.MethodPut {
		errorMessage := "Invalid request method. Needs to be PUT"
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusMethodNotAllowed)
		return
	}

	// Extract the ObjectId from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/students/")
	// Convert the string ID to a MongoDB ObjectId type
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ObjectId format", http.StatusBadRequest)
		return
	}

	var updateBody bson.M
    jsonDecoder := json.NewDecoder(r.Body)
    err = jsonDecoder.Decode(&updateBody)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
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
		log.Printf("Failed to update student: %v", err)
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
	}
	if updateResult.MatchedCount == 0 {
		log.Println("No record found with the provided ID")
		http.Error(w, "No record found with the provided ID", http.StatusNotFound)
		return
	}

    // Check if any student field is updated and save these fields to slice
	var updateKeys []string
    for updateKey := range updateBody {
        if _, found := map[string]bool{"fullname": true, "phone": true, "subscription": true, "startDate": true, "lastDate": true, "comments": true}[updateKey]; !found {
			log.Println("No student field is updated")
			http.Error(w, "No student field is updated", http.StatusBadRequest)
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
	// Check if the method is DELETE; return 405 in case of error
	if r.Method != http.MethodDelete {
		errorMessage := "Invalid request method. Needs to be DELETE"
		log.Println(errorMessage)
		http.Error(w, errorMessage, http.StatusMethodNotAllowed)
		return
	}

	// Extract the ObjectId from the URL path
	id := strings.TrimPrefix(r.URL.Path, "/students/")
	// Convert the string ID to a MongoDB ObjectId type
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ObjectId format", http.StatusBadRequest)
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
		log.Printf("Failed to update student: %v", err)
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
	} 
	if deleteResult.DeletedCount == 0 {
		log.Printf("No record found with the provided ID: %v", id)
		http.Error(w, fmt.Sprintf("No record found with the provided ID: %v", id), http.StatusNotFound)
		return
	}

	// Write the response with deleted student id
	response := fmt.Sprintf("Deleted student by mentioned id: %v", id)
	w.WriteHeader(http.StatusOK)
    w.Write([]byte(response))
}
