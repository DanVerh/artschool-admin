// db.go
package application

import (
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const dbUri = "mongodb://localhost:27017"

type Database struct {
	client *mongo.Client
}

func NewDbClient() *Database {
	db := &Database{
		client: createClient(),
	}

	return db
}

func createClient() *mongo.Client {
    client, err := mongo.NewClient(options.Client().ApplyURI(dbUri))
    if err != nil {
        log.Fatal(err)
    }

    err = client.Connect(nil)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(nil, nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Connected to MongoDB!")
    return client
}