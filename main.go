package main

import (
	"ccdata/ccdata"
	"ccdata/mongo"
	"fmt"
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

	// Count the number of games found
	gameCount := len(gameData)
	logger.Info().Msg(fmt.Sprintf("Games found: %v", gameCount))

	// Write game data to MongoDB "games" collection in batches of 5000 games
	for i := 0; i < gameCount; i += 5000 {
		// Determine the index of the last game in the batch
		var index int
		if i+5000 > gameCount {
			index = gameCount
		} else {
			index = i + 5000
		}

		// Write the batch of games to the "games" collection
		err = mongo.WriteDataToCollection(os.Getenv("GAME_COLLECTION"), struct{ Games []ccdata.AssortedGamePGN }{Games: gameData[i:index]})
		if err != nil {
			// Log error and exit program
			logger.Info().Err(err).Msg("error while writing data to collection")
			os.Exit(1)
		}
	}
}
