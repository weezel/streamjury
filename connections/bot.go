package connections

import (
	"log"
	"os"
	"streamjury/confighandler"
	"streamjury/gameplay"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	game              gameplay.GamePlay
	currentGameState  gameplay.GameState = gameplay.InitState
	curSongFrom       *gameplay.Player
	superUserId       int64
	channelId         int64
	telegramApiKey    string
	resultsAbsPath    string
	allSongsPresented = false
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func SetConfigValues(config confighandler.TomlConfig) {
	superUserId = config.StreamjuryConfig.SuperUserId
	channelId = config.StreamjuryConfig.ChannelId
	telegramApiKey = config.StreamjuryConfig.ApiKey
	resultsAbsPath = config.StreamjuryConfig.ResultsAbsPath
}

func ConnectionHandler() {
	bot, err := tgbotapi.NewBotAPI(telegramApiKey)
	if err != nil {
		log.Panicf("Possible error in config file: %s", err)
	}

	bot.Debug = func() bool {
		return strings.ToLower(os.Getenv("DEBUG")) == "true"
	}()

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		splitted := strings.Split(update.Message.Text, " ")
		if len(splitted) < 2 {
			continue
		}
		if strings.HasPrefix(splitted[0], "@") {
			splitted = append(splitted, splitted[1:]...)
		}
		if strings.ToLower(splitted[0]) != "levyraati" {
			continue
		}

		var playerName string
		if update.Message.From.UserName == "" {
			playerName = update.Message.From.FirstName
		} else {
			playerName = update.Message.From.UserName
		}
		player := gameplay.AddPlayer(
			playerName,
			update.Message.From.ID,
		)

		log.Printf("Message from %s(%d): %s",
			player.Name,
			player.Uid,
			update.Message.Text)

		// Only certain commands are feasible in certains states
		var command string = strings.ToLower(splitted[1])
		if gameplay.IsCommandFeasibleInState(currentGameState, command) == false {
			log.Printf("Player %s(%d) initiated command '%s' in wrong state",
				player.Name,
				player.Uid,
				command)
			continue
		}

		// Someone initiated the command while not joined in game?
		playerIsInTheGame := isPlayerInTheGame(player)
		if !playerIsInTheGame &&
			currentGameState > gameplay.WaitingForPlayers {
			log.Printf("Player %s(%d) initiated esitys without being in the game",
				player.Name,
				player.Uid)
			continue
		}

		// At first, handle users's input
		log.Printf("State before handling command: %+v\n",
			currentGameState)
		switch command {
		case "aloita":
			currentGameState = handleBegin(bot, player)
		case "lopeta":
			currentGameState = handleStop(bot, player, currentGameState)
		case "jatka":
			currentGameState = handleContinue(bot, player, currentGameState)
		case "arvio", "arvioi", "arvostele":
			currentGameState = handleReview(
				bot,
				player,
				update.Message.Chat.ID,
				currentGameState,
				splitted)
		case "esitä", "esitys":
			currentGameState = handleSubmittedSong(
				bot,
				player,
				update.Message.Chat.ID,
				currentGameState,
				splitted)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"Parametri uupuu, pitää olla jokin näistä: [aloita|arvio tai arvioi tai arvostele|esitä tai esitys|lopeta]")
			if _, err = bot.Send(msg); err != nil {
				log.Printf("ERROR: %+v", err)
			}
		}
		log.Printf("State after handling command: %+v\n",
			currentGameState)

		// And then handle state regarding the previous input.
		switch currentGameState {
		case gameplay.InitState:
			game = gameplay.GamePlay{}
			game.Players = []gameplay.Player{}
		case gameplay.PublishingSong:
			// Clear ReviewGiven flags before the next song
			// Mark review given since on this round
			// presenter won't review anything
			currentGameState, replyMsg := handlePublishing(bot)
			switch currentGameState {
			case gameplay.StopGame:
				allSongsPresented = true
			case gameplay.WaitingForReviews:
				curSongFrom.ReviewGiven = true
			}

			msg := tgbotapi.NewMessage(channelId, replyMsg)
			if _, err = bot.Send(msg); err != nil {
				log.Printf("ERROR: %+v", err)
			}
			log.Printf("State: %+v\n", currentGameState)
		}

		// Game over man, game over
		if currentGameState == gameplay.StopGame {
			handlePublishResults(game)
			currentGameState = handleQuit(game)
		} // for update
	}
}

func isPlayerInTheGame(player *gameplay.Player) bool {
	for _, p := range game.Players {
		if player.Uid == p.Uid {
			return true
		}
	}
	return false
}
