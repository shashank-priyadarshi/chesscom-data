package main

import (
	"ccdata/ccdata"
	"ccdata/mongo"
	"fmt"
	"os"

	logger "github.com/rs/zerolog/log"
)

func main() {
	// Run the main logic of the application
	if err := runApp(); err != nil {
		// Log any errors encountered during runtime
		logger.Info().Err(err).Msg("Application error")
		os.Exit(1)
	}

	// Log successful execution of the plugin
	logger.Info().Msg("Plugin executed successfully")
}

// runApp retrieves game data from ccdata package and pushes it to MongoDB database.
//
// Returns an error if there's an issue retrieving or pushing data.
func runApp() error {
	gameData, err := ccdata.GetData()
	if err != nil {
		return fmt.Errorf("error retrieving game data: %v", err)
	}

	err = pushDataToDB(gameData)
	if err != nil {
		return fmt.Errorf("error pushing game data to MongoDB: %v", err)
	}

	return nil
}

// pushDataToDB writes AssortedGamePGN data to a MongoDB collection in batches of 5000 games.
//
// Params:
// gameData: A slice of AssortedGamePGN containing the game data to be written.
//
// Returns:
// An error if the function fails to write the data to the collection.
func pushDataToDB(gameData []ccdata.AssortedGamePGN) (err error) {
	// Count the number of games found
	gameCount := len(gameData)
	logger.Info().Msgf("Games found: %v", gameCount)

	// Write game data to MongoDB "games" collection in batches of 5000 games
	for i := 0; i < gameCount; i += 5000 {
		// Determine the index of the last game in the batch
		index := i + 5000
		if index > gameCount {
			index = gameCount
		}

		// Write the batch of games to the "games" collection
		if err = mongo.WriteDataToCollection(os.Getenv("GAME_COLLECTION"), struct{ Games []ccdata.AssortedGamePGN }{Games: gameData[i:index]}); err != nil {
			// Log error and exit program
			logger.Error().Err(err).Msgf("error while writing objects %v to %v", i, index)
		}
	}
	return
}
