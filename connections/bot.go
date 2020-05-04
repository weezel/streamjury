package connections

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"streamjury/confighandler"
	"streamjury/gameplay"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	game              gameplay.GamePlay
	curGameState      int = gameplay.InitState
	curSongFrom       *gameplay.Player
	superUserId       int
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

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if !game.StartedAt.IsZero() {
			timedOut, timeInIdle := game.HasIdleTimedOut()
			if timedOut {
				log.Printf("Timeout reached: %v", timeInIdle)
			} else {
				log.Printf("Currently idled: %v / %v", timeInIdle,
					gameplay.GameIdleTimeout)
			}
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
		if gameplay.IsCommandFeasibleInState(curGameState, command) == false {
			log.Printf("Player %s(%d) initiated command '%s' in wrong state",
				player.Name,
				player.Uid,
				command)
			continue
		}
		// Someone initiated the command while not joined in game?
		var playerIsInTheGame bool = false
		for _, p := range game.Players {
			if player.Uid == p.Uid {
				playerIsInTheGame = true
			}
		}
		if playerIsInTheGame == false &&
			curGameState > gameplay.WaitingForPlayers {
			log.Printf("Player %s(%d) initiated esitys without being in the game",
				player.Name,
				player.Uid)
			continue
		}

		switch command {
		case "aloita":
			var replyMsg string
			var outMsg tgbotapi.MessageConfig

			if game.IsInGame(player.Uid) {
				replyMsg = "Olet jo mukana pelissä ;)"
				outMsg = tgbotapi.NewMessage(channelId, replyMsg)
				bot.Send(outMsg)
				log.Printf("%s(%d) is already joined\n",
					player.Name,
					player.Uid)
			} else {
				// Game starter has super powers
				if len(game.Players) == 0 {
					game.GameStarterUid = player.Uid
					replyMsg = fmt.Sprintf("%s haluaa aloittaa levyraadin, "+
						"odotetaan muita pelaajia. "+
						"Liity peliin kirjoittamalla: levyraati aloita. "+
						"Kun kaikki ovat liittyneet, aloita-komennon ensimmäinen "+
						"käskijä voi käynnistää pelin kirjoittamalla 'levyraati jatka'.",
						player.Name)
					outMsg = tgbotapi.NewMessage(channelId, replyMsg)
					bot.Send(outMsg)
				} else {
					replyMsg = fmt.Sprintf("Pelaaja %s liittyi peliin!",
						player.Name)
					outMsg = tgbotapi.NewMessage(channelId, replyMsg)
					bot.Send(outMsg)

					// Show joined players
					joinedPlayers := func() []string {
						names := make([]string, len(game.Players))
						for i, p := range game.Players {
							names[i] = p.Name
						}
						return names
					}
					replyMsg = fmt.Sprintf("Tähän mennessä liittyneet: %s",
						strings.Join(joinedPlayers(), ", "))
					outMsg = tgbotapi.NewMessage(channelId, replyMsg)
					bot.Send(outMsg)
				}
				game.AppendPlayer(player)
				log.Printf("Player %s(%d) joined", player.Name, player.Uid)
				log.Printf("Joined players: %v\n", game.Players)
			}

			if game.StartedAt.IsZero() {
				game.StartedAt = time.Now()
			}

			curGameState = gameplay.WaitingForPlayers
			log.Printf("State: %v\n", curGameState)
		case "lopeta":
			var replyMsg string

			if player.Uid == game.GameStarterUid ||
				player.Uid == superUserId {
				if curGameState > 0 {
					replyMsg = "Kierros lopetettu"
				} else {
					replyMsg = "Ei peliä käynnissä"
				}
				curGameState = gameplay.StopGame
			} else {
				log.Printf("Player %s(%d) tried to stop the game, ignored",
					player.Name,
					player.Uid)
				outMsg := tgbotapi.NewMessage(channelId,
					"Et ole pelin aloittaja, joten teikäläisen natsat ei riitä")
				bot.Send(outMsg)
				continue
			}

			outMsg := tgbotapi.NewMessage(channelId, replyMsg)
			bot.Send(outMsg)

			_, timeInIdle := game.HasIdleTimedOut()
			log.Printf("Gameplay stopped. Idle time was %v\n",
				timeInIdle)
			game.Reset()
			curGameState = gameplay.InitState
			log.Printf("State: %v\n", curGameState)
		case "jatka":
			var replyMsg string

			log.Printf("Player %s(%d) commanded jatka",
				player.Name,
				player.Uid)

			if curGameState == gameplay.WaitingForPlayers &&
				len(game.Players) == 1 {
				replyMsg = "Yksin ei voi pelata :´( Lopetetaan peli"
				curGameState = gameplay.StopGame
				log.Printf("Player %s(%d) tried to play alone",
					player.Name,
					player.Uid)
				msg := tgbotapi.NewMessage(channelId,
					replyMsg)
				bot.Send(msg)
				break
			}

			if player.Uid == game.GameStarterUid ||
				player.Uid == superUserId {
				replyMsg = "Jatketaan..."
				gameplay.NextGameState(&curGameState)
				replyMsg = gameplay.GetGameStateIntro(curGameState)
			} else {
				replyMsg = "Et ole pelin aloittaja, joten oikeutesi ei riitä komentelemaan"
			}

			outMsg := tgbotapi.NewMessage(channelId, replyMsg)
			bot.Send(outMsg)
			log.Printf("State: %v\n", curGameState)
		case "arvio", "arvioi", "arvostele":
			var replyMsg string

			if player.Uid == curSongFrom.Uid {
				log.Printf("Player %s(%d) tried to rate his/hers own song",
					player.Name,
					player.Uid)
				replyMsg = "Eh, yritit sitten arvioida oman kappaleesi"
			} else {
				var review string = strings.Join(splitted[2:], " ")
				// Add review for the current song presenter
				if err := curSongFrom.AddReview(player.Name, review); err != nil {
					errMsg := tgbotapi.NewMessage(
						update.Message.Chat.ID,
						err.Error())
					bot.Send(errMsg)
					continue
				}
				for i := range game.Players {
					if player.Uid == game.Players[i].Uid {
						game.Players[i].ReviewGiven = true
						log.Printf("Player %s(%d) reviewed the song",
							player.Name,
							player.Uid)
						break
					}
				}

				replyMsg = fmt.Sprintf("Arviointisi on rekisteröity, %s",
					player.Name)
				log.Printf("Player %s(%d) submitted review: %s",
					player.Name,
					player.Uid,
					review)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				replyMsg)
			bot.Send(msg)

			if game.AllReviewsGiven() {
				curGameState = gameplay.PublishingSong
				log.Print("All reviews given")
			}

			log.Printf("State: %v\n", curGameState)
		case "esitä", "esitys":
			curGameState = gameplay.WaitingForSongs
			var replyMsg string

			// Get song
			var possibleUrl string = splitted[len(splitted)-1]
			songUrl, err := url.Parse(possibleUrl)
			if err != nil {
				replyMsg = "VIRHE: Viimeisenä ei ollut linkkiä"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID,
					replyMsg)
				bot.Send(msg)
				continue
			}

			// Get song description
			// description := strings.Join(splitted[2:len(splitted)-1], " ")
			thirdSpace := func(msg string, cmp rune, nThMatch int) int {
				spaceCount := 0
				for i, r := range update.Message.Text {
					if r == cmp {
						spaceCount++
					}
					if spaceCount == nThMatch {
						return i
					}
				}
				return -1
			}
			startIdx := thirdSpace(update.Message.Text, ' ', 2)
			endIdx := strings.LastIndex(update.Message.Text, " ")
			if startIdx == -1 || endIdx == -1 {
				log.Printf("startIdx=%d, endIdx=%d message: %s",
					startIdx,
					endIdx,
					update.Message.Text)
				continue
			}
			description := update.Message.Text[startIdx+1 : endIdx]

			for i, p := range game.Players {
				// Continue until the submitter is found
				if player.Uid != p.Uid {
					continue
				}

				if game.Players[i].Song != nil {
					replyMsg = fmt.Sprintf("Olet jo lähettänyt kappaleen: %s",
						game.Players[i].Song.Url)
					log.Printf("Player %s tried to send a new song: %s\n",
						player.Name, songUrl.String())
					continue
				} else {
					game.Players[i].AddSong(description, songUrl.String())
					game.Players[i].SongSubmitted = true
					log.Printf("Player %s added song %s with description: %s\n",
						p.Name,
						songUrl.String(),
						description)
					log.Printf("%v\n", game.Players)
					break
				}
			}

			// Don't go to next state until all players have submitted
			// their song
			var allPlayersSubmittedSong bool = true
			for _, p := range game.Players {
				if p.SongSubmitted == false {
					allPlayersSubmittedSong = false
					break
				}
			}
			if allPlayersSubmittedSong {
				curGameState = gameplay.PublishingSong
				replyMsg = gameplay.GetGameStateIntro(curGameState)
				msg := tgbotapi.NewMessage(channelId, replyMsg)
				bot.Send(msg)
			} else {
				privMsg := tgbotapi.NewMessage(update.Message.Chat.ID,
					fmt.Sprintf("Kiitos kappaleesta %s", player.Name))
				chanReply := tgbotapi.NewMessage(channelId,
					fmt.Sprintf("%s lähetti kappaleen", player.Name))

				bot.Send(privMsg)
				bot.Send(chanReply)
			}

			log.Printf("State: %v\n", curGameState)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"Parametri uupuu, pitää olla jokin näistä: [aloita|arvio tai arvioi tai arvostele|esitä tai esitys|lopeta]")
			bot.Send(msg)
		}

		switch curGameState {
		case gameplay.InitState:
			game = gameplay.GamePlay{}
			game.Players = []gameplay.Player{}
		case gameplay.PublishingSong:
			var replyMsg string

			curSongFrom = game.NextSongFrom()
			if curSongFrom == nil {
				replyMsg = "Kaikki Muumit ovat esittäneet kappaleet. Peli loppuu ja muumitalo lukitaan"
				allSongsPresented = true
				curGameState = gameplay.StopGame
			} else {
				// Clear ReviewGiven flags before the next song
				for i := range game.Players {
					game.Players[i].ReviewGiven = false
				}
				replyMsg = fmt.Sprintf("Seuraava kappale tulee käyttäjältä %s. Linkki kappaleeseen %s ja kuvaus: %s",
					curSongFrom.Name,
					curSongFrom.Song.Url,
					curSongFrom.Song.Description)
				// Mark review given since on this round
				// presenter won't review anything
				curSongFrom.ReviewGiven = true
				curGameState = gameplay.WaitingForReviews
				infoMsg := tgbotapi.NewMessage(channelId, "Odotetaan arvioita. Arvioi kappale näin: levyraati arvioi selitys-tähän 10/10")
				bot.Send(infoMsg)
			}
			msg := tgbotapi.NewMessage(channelId, replyMsg)
			bot.Send(msg)
			log.Printf("State: %v\n", curGameState)
		}

		if curGameState == gameplay.StopGame {
			if allSongsPresented {
				dateTimeNow := time.Now().Format("2006-01-02_150405")
				reviewsFilename := fmt.Sprintf("gameplay-%s.html", dateTimeNow)
				reviewsFullPath := filepath.Join(resultsAbsPath, reviewsFilename)
				var stringWriter bytes.Buffer
				game.PublishResults(&stringWriter)
				err := ioutil.WriteFile(
					reviewsFullPath,
					stringWriter.Bytes(),
					0644)
				check(err)
				log.Printf("Game written to file %s", reviewsFullPath)
				allSongsPresented = false
			}
			game.Reset()
			log.Printf("After reset: %v", game)
			curGameState = gameplay.InitState
		}
	} // for update
}
