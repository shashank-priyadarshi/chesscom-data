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

func GetData() {
	logger.Info().Msg("Getting data...")
	gameList := userWiseGameData()
	marshalled, _ := json.Marshal(gameList)
	logger.Info().Msg(fmt.Sprintf("Games found: %v", len(gameList)))
	fmt.Println(string(marshalled))
}

func parseUserName() (users []string) {
	rawUserStr := os.Getenv("USERNAME")
	if rawUserStr == "" {
		logger.Info().Msg("No username provided")
		os.Exit(1)
	}
	users = strings.Split(rawUserStr, ",")
	fmt.Println(users)
	return
}

func userWiseGameData() (gameList []AssortedGamePGN) {
	users := parseUserName()
	for _, user := range users {
		data, statusCode := httpcall(fmt.Sprintf(archive, user))
		if statusCode != 200 {
			logger.Info().Err(fmt.Errorf("error while making request to %v: %v", archive, statusCode)).Msg("")
			continue
		}

		var archiveList Archive
		err := json.Unmarshal(data, &archiveList)
		if err != nil {
			logger.Info().Err(err).Msg("Failed to parse archive data")
			continue
		}

		gameList, err = parseGameData(archiveList)
		if err != nil {
			logger.Info().Err(err).Msg("Failed to parse game data")
		}
	}
	return
}

func parseGameData(archiveList Archive) (gameList []AssortedGamePGN, err error) {
	for _, url := range archiveList.Archives {
		var games GameList
		data, statusCode := httpcall(url)
		if statusCode != 200 {
			logger.Info().Err(fmt.Errorf("error while making request to %v: %v", url, statusCode)).Msg("")
			continue
		}

		err = json.Unmarshal(data, &games)
		if err != nil {
			logger.Info().Err(err).Msg("Failed to parse game data")
			continue
		}

		lo.ForEach(games.Games, func(game Game, index int) {
			split := strings.Split(game.Pgn, "\n\n")
			assortedGame := AssortedGamePGN{
				GameURL:     game.URL,
				GameDetails: split[0],
				PGN:         split[1],
			}
			if strings.EqualFold(game.White.Username, "k_heerathakur") {
				assortedGame.Result = game.White.Result
			} else {
				assortedGame.Result = game.Black.Result
			}
			gameList = append(gameList, assortedGame)
		})
	}
	return
}

func httpcall(reqURL string) ([]byte, int) {
	client := http.Client{}

	request, err := http.NewRequest("GET", reqURL, bytes.NewBuffer([]byte("")))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		logger.Info().Err(err).Msg("err in creating new request: ")
		return []byte{}, 503
	}

	resp, err := client.Do(request)

	if resp.StatusCode != http.StatusOK {
		logger.Info().Err(fmt.Errorf("error '%v' while making request to %v: %v", err, reqURL, resp.StatusCode))
		return []byte{}, resp.StatusCode
	}

	if err != nil {
		logger.Info().Err(err).Msg("err in reading req response: ")
		return []byte{}, resp.StatusCode
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Info().Err(err).Msg("err in reading req response: ")
		return []byte{}, 503
	}

	return respBody, resp.StatusCode
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
