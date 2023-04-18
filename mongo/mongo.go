package mongo

import (
	"context"
	"errors"
	"os"

	// "github.com/docker/docker/daemon/logger"
	logger "github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	done = make(chan bool)
	db   *mongo.Database
)

// setMongoConnection creates a new MongoDB client and sets the global "db" variable to
// the specified database. A goroutine is started to close the client when the program finishes
func setMongoConnection() {
	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Info().Err(err).Msg("error while creating mongo client:")
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logger.Info().Err(err).Msg("error while pinging mongo: ")
	}

	// Start a goroutine to close the MongoDB connection when the program finishes
	go func() {
		// Wait for the task to finish
		<-done

		// Task finished, close the MongoDB connection
		err := client.Disconnect(context.TODO())
		if err != nil {
			logger.Info().Err(err).Msg("Error while closing connection to mongo: ")
		}
	}()

	// Set the global "db" variable to the specified database
	// NOTE: This should be done after the connection is established and checked for errors
	db = client.Database(os.Getenv("GAME_DB"))
}

// WriteDataToCollection writes data to a MongoDB collection.
// It takes a collection name and the data to be inserted as parameters.
// Returns an error if the operation was unsuccessful.
func WriteDataToCollection(collectionName string, data interface{}) error {
	// setMongoConnection() function connects to MongoDB and initializes the db variable.
	setMongoConnection()

	// Get the collection handle for the given collection name.
	collection := db.Collection(collectionName)

	// Insert the data into the collection.
	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		// If there was an error during insertion, log the error.
		logger.Info().Err(errors.New(err.Error())).Msg("Error while inserting data to MongoDB")
	}

	// Signal completion of the operation.
	done <- true

	return err
}
