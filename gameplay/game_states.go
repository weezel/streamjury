package gameplay

type GameState int

const (
	InitState         GameState = iota // 0
	WaitingForPlayers GameState = iota // 1
	WaitingForSongs   GameState = iota // 2
	PublishingSong    GameState = iota // 3
	WaitingForReviews GameState = iota // 4
	StopGame          GameState = iota // 5
)

func NextGameState(gameState GameState) GameState {
	switch gameState {
	case InitState:
		return WaitingForPlayers
	case WaitingForPlayers:
		return WaitingForSongs
	case WaitingForSongs:
		return PublishingSong
	case PublishingSong:
		return WaitingForReviews
	case WaitingForReviews:
		return StopGame
	case StopGame:
		return InitState
	}
	return StopGame
}

// IsCommandFeasibleInState function returns boolean whether
// a given command is valid on a current state.
func IsCommandFeasibleInState(currentState GameState, command string) bool {
	switch currentState {
	case InitState:
		if command == "aloita" {
			return true
		}
	case WaitingForPlayers:
		switch command {
		case "aloita", "jatka", "lopeta":
			return true
		}
	case WaitingForSongs:
		switch command {
		case "esitys", "esitä", "jatka", "lopeta":
			return true
		}
	case WaitingForReviews:
		switch command {
		case "arvio", "arvioi", "arvostele", "jatka", "lopeta":
			return true
		}
	case PublishingSong, StopGame:
		return false
	}
	return false
}

func GetGameStateIntro(gameState GameState) string {
	switch gameState {
	case InitState:
		return ""
	case WaitingForPlayers:
		return "Odotetaan muiden pelaajien liittymistä. Liity komentamalla: levyraati aloita"
	case WaitingForSongs:
		return "Odotetaan kappaleita. Esitä kappale botille privaattiviestillä: levyraati esitä selitys-tähän url-linkki"
	case PublishingSong:
		return ""
	case WaitingForReviews:
		return ""
	case StopGame:
		return ""
	}
	return ""
}
