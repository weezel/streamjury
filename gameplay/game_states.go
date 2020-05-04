package gameplay

const (
	InitState         = iota
	WaitingForPlayers = iota
	WaitingForSongs   = iota
	PublishingSong    = iota
	WaitingForReviews = iota
	StopGame          = iota
)

func NextGameState(gameState *int) {
	switch *gameState {
	case InitState:
		*gameState = WaitingForPlayers
	case WaitingForPlayers:
		*gameState = WaitingForSongs
	case WaitingForSongs:
		*gameState = PublishingSong
	case PublishingSong:
		*gameState = WaitingForReviews
	case WaitingForReviews:
		*gameState = StopGame
	case StopGame:
		*gameState = InitState
	}
}

func IsCommandFeasibleInState(currentStatee int, command string) bool {
	switch currentStatee {
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

func GetGameStateIntro(gameState int) string {
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
