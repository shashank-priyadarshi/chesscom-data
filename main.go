package main

import (
	"ccdata/ccdata"
	"ccdata/mongo"
	"os"

	logger "github.com/rs/zerolog/log"
)

// main is the entry point of the program
func main() {
	// Retrieve game data
	gameData, err := ccdata.GetData()
	if err != nil {
		// Log error and exit program
		logger.Info().Err(err).Msg("error while getting data")
		os.Exit(1)
	}

	// Push game data to MongoDB
	err = pushDataToDB(gameData)
	if err != nil {
		os.Exit(1)
	}

	logger.Info().Msg("Plugin executed successfully")
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
