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

func dbConnection() *Database {
	db := &Database{
		client: connectClient(),
	}

	return db
}

func connectClient() *mongo.Client {
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


// Close disconnects the MongoDB client and drops the connection
/*func (db *Database) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.client.Disconnect(ctx)
	if err != nil {
		log.Fatalf("Failed to disconnect MongoDB client: %v", err)
	}

	log.Println("Disconnected from MongoDB")
}*/