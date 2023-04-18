package ccdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	logger "github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

// GetData retrieves user game data and logs the number of games found
// and returns the marshalled JSON data.
func GetData() ([]AssortedGamePGN, error) {
	// log that we're getting data
	logger.Info().Msg("Getting data...")

	// get userwise game data
	gameList, err := userWiseGameData()
	if err != nil {
		// log error if failed to get userwise game data
		logger.Info().Err(err).Msg("Failed to get userwise game data")
		return nil, err
	}

	// return marshalled JSON data
	return gameList, nil
}

// parseUserName parses a comma-separated string of users from the USERNAME environment variable
// and returns them as a slice of strings.
//
// If the USERNAME environment variable is not set or empty, the function logs an error and exits with a non-zero code.
func parseUserName() (users []string) {
	// Get the raw user string from the environment variable
	rawUserStr := os.Getenv("USERNAME")
	if rawUserStr == "" {
		// Log an error if the environment variable is not set or empty
		logger.Info().Msg("No username provided")
		os.Exit(1)
	}
	// Split the user string into a slice of strings
	users = strings.Split(rawUserStr, ",")
	// Print the users to the console (for debugging purposes)
	return
}

// userWiseGameData retrieves a list of games played by all users and returns them in an array.
// It makes a HTTP request to the archive API to get the user data, parses it into an Archive struct,
// and then calls parseGameData to extract the game data from the parsed Archive struct.
// If any errors occur during the process, they are logged and the function continues to the next user.
func userWiseGameData() (gameList []AssortedGamePGN, err error) {
	// retrieve list of user names
	users := parseUserName()

	// iterate through each user
	for _, user := range users {
		// make HTTP request to archive API
		data, statusCode, err := httpcall(fmt.Sprintf(archive, user))
		if statusCode != 200 {
			// log error if request was unsuccessful and continue to next user
			logger.Info().Err(err).Msg(fmt.Sprintf("error while making request to %v: %v", archive, statusCode))
			continue
		}

		data = []byte(strings.ReplaceAll(string(data), "\n", ""))

		// parse JSON response into Archive struct
		var archiveList Archive
		err = json.Unmarshal(data, &archiveList)
		if err != nil {
			// log error if failed to parse archive data and continue to next user
			logger.Info().Err(err).Msg("Failed to parse archive data")
			continue
		}

		// parse game data from Archive struct
		parsedGameList, err := parseGameData(archiveList)
		if err != nil {
			// log error if failed to parse game data and continue to next user
			logger.Info().Err(err).Msg("Failed to parse game data")
		}

		// add parsed game data to gameList array
		gameList = append(gameList, parsedGameList...)
	}
	return
}

/*
parseGameData takes an Archive object and returns an array of AssortedGamePGN objects and an error.

It makes an HTTP call to each URL in the archiveList and unmarshals the response into a GameList object. It then loops through each Game object in GameList and splits each game's PGN into GameDetails and PGN properties. The function creates an AssortedGamePGN object with GameURL, GameDetails, PGN, and Result properties. If the game was played by "k_heerathakur", the Result property is set to the result of the white player, otherwise it is set to the result of the black player. The function returns an array of all created AssortedGamePGN objects and any errors that may have occurred.

Parameters:
- archiveList (Archive): an object containing a list of archive URLs

Returns:
- gameList ([]AssortedGamePGN): an array of AssortedGamePGN objects
- err (error): any errors that may have occurred
*/
func parseGameData(archiveList Archive) (gameList []AssortedGamePGN, err error) {
	// Loop through each URL in the archiveList
	for _, url := range archiveList.Archives {
		var games GameList

		// Make an HTTP call to the URL and check the status code
		data, statusCode, err := httpcall(url)
		if statusCode != 200 {
			logger.Info().Err(err).Msg(fmt.Sprintf("error while making request to %v: %v", archive, statusCode))
			continue
		}

		data = []byte(strings.ReplaceAll(string(data), "\n", ""))

		// Unmarshal the response into a GameList object
		err = json.Unmarshal(data, &games)
		if err != nil {
			logger.Info().Err(err).Msg("Failed to parse game data")
			continue
		}

		// Loop through each Game object in GameList
		lo.ForEach(games.Games, func(game Game, index int) {
			// Split each game's PGN into GameDetails and PGN properties
			split := strings.Split(game.Pgn, "\n\n")
			// Create an AssortedGamePGN object with GameURL, GameDetails, PGN, and Result properties
			assortedGame := AssortedGamePGN{
				GameURL:     game.URL,
				GameDetails: split[0],
				PGN:         split[1],
			}
			// If the game was played by "k_heerathakur", set the Result property to the result of the white player, otherwise set it to the result of the black player
			if strings.EqualFold(game.White.Username, "k_heerathakur") {
				assortedGame.Result = game.White.Result
			} else {
				assortedGame.Result = game.Black.Result
			}
			// Add the created AssortedGamePGN object to the gameList array
			gameList = append(gameList, assortedGame)
		})
	}
	// Return the gameList array and any errors that may have occurred
	return
}

// httpcall sends an HTTP GET request to the specified URL and returns the response body and status code.
// If an error occurs, an empty byte slice and an appropriate status code are returned.
// reqURL: the URL to send the GET request to
// returns: the response body and status code
func httpcall(reqURL string) ([]byte, int, error) {
	// create an HTTP client
	client := http.Client{}

	// create a new GET request with an empty request body
	request, err := http.NewRequest("GET", reqURL, bytes.NewBuffer([]byte("")))
	if err != nil {
		logger.Info().Err(err).Msg("err in creating new request: ")
		return []byte{}, http.StatusServiceUnavailable, err
	}

	// set the request header to specify that the request body is in JSON format
	request.Header.Set("Content-Type", "application/json")

	// send the request and get the response
	resp, err := client.Do(request)
	if err != nil {
		logger.Info().Err(err).Msg("err in making request: ")
		return []byte{}, http.StatusServiceUnavailable, err
	}

	// if the response status code indicates an error, log it and return the appropriate status code
	if resp.StatusCode != http.StatusOK {
		logger.Info().Err(fmt.Errorf("error '%v' while making request to %v: %v", err, reqURL, resp.StatusCode))
		return []byte{}, resp.StatusCode, err
	}

	// read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Info().Err(err).Msg("err in reading req response: ")
		return []byte{}, http.StatusServiceUnavailable, err
	}

	// return the response body and status code
	return respBody, resp.StatusCode, nil
}

// func parseGameData(data string) map[string]string {
// 	var code map[string]string
// 	err := json.Unmarshal([]byte(data), &code)
// 	if err != nil {
// 		logger.Info().Err(err).Msg("Failed to parse game result code data")
// 		return nil
// 	}
// 	return code
// }
