package ccdata

const (
//		GAMERESULTCODES = `{
//	  "win": "Win",
//	  "checkmated": "Checkmated",
//	  "agreed": "Draw agreed",
//	  "repetition": "Draw by repetition",
//	  "timeout": "Timeout",
//	  "resigned": "Resigned",
//	  "stalemate": "Stalemate",
//	  "lose": "Lose",
//	  "insufficient": "Insufficient material",
//	  "50move": "Draw by 50-move rule",
//	  "abandoned": "Abandoned",
//	  "kingofthehill": "Opponent king reached the hill",
//	  "threecheck": "Checked for the 3rd time",
//	  "timevsinsufficient": "Draw by timeout vs insufficient material",
//	  "bughousepartnerlose": "Bughouse partner lost"
//	}`
)

var (
	archive = "https://api.chess.com/pub/player/%v/games/archives"
	// resultcode = parseGameData(GAMERESULTCODES)
)

type AssortedGamePGN struct {
	GameDetails string `json:"game_details"`
	PGN         string `json:"pgn"`
	Result      string `json:"result"`
	GameURL     string `json:"url"`
}

type Archive struct {
	Archives []string `json:"archives"`
}

type GameList struct {
	Games []Game `json:"games"`
}

// Generated by https://quicktype.io
type Game struct {
	URL string `json:"url"`
	Pgn string `json:"pgn"`
	// TimeControl  string     `json:"time_control"`
	// EndTime      int64      `json:"end_time"`
	// Rated        bool       `json:"rated"`
	// Accuracies   Accuracies `json:"accuracies"`
	// Tcn          string     `json:"tcn"`
	// UUID         string     `json:"uuid"`
	// InitialSetup string     `json:"initial_setup"`
	// Fen          string     `json:"fen"`
	// TimeClass    string     `json:"time_class"`
	// Rules        string     `json:"rules"`
	White Player `json:"white"`
	Black Player `json:"black"`
}

// type Accuracies struct {
// 	White float64 `json:"white"`
// 	Black float64 `json:"black"`
// }

type Player struct {
	// Rating   int64  `json:"rating"`
	Result string `json:"result"`
	// ID       string `json:"@id"`
	Username string `json:"username"`
	// UUID     string `json:"uuid"`
}
