package connections

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"streamjury/gameplay"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotIF interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func handleBegin(bot BotIF, player *gameplay.Player) gameplay.GameState {
	var replyMsg string
	var outMsg tgbotapi.MessageConfig
	var err error

	if game.IsInGame(player.Uid) {
		replyMsg = "Olet jo mukana pelissä ;)"
		outMsg = tgbotapi.NewMessage(channelId, replyMsg)
		if _, err = bot.Send(outMsg); err != nil {
			log.Printf("ERROR (handleBegin1): %+v", err)
		}
		log.Printf("%s(%d) is already joined\n",
			player.Name,
			player.Uid)
	} else {
		// Game starter has super powers
		if len(game.Players) == 0 {
			game.GameStarterUID = player.Uid
			replyMsg = fmt.Sprintf("%s haluaa aloittaa levyraadin, "+
				"odotetaan muita pelaajia. "+
				"Liity peliin kirjoittamalla: levyraati aloita. "+
				"Kun kaikki ovat liittyneet, aloita-komennon ensimmäinen "+
				"käskijä voi käynnistää pelin kirjoittamalla 'levyraati jatka'.",
				player.Name)
			outMsg = tgbotapi.NewMessage(channelId, replyMsg)
			if _, err = bot.Send(outMsg); err != nil {
				log.Printf("ERROR (handleBegin2): %+v", err)
			}
		} else {
			replyMsg = fmt.Sprintf("Pelaaja %s liittyi peliin!",
				player.Name)
			outMsg = tgbotapi.NewMessage(channelId, replyMsg)
			if _, err = bot.Send(outMsg); err != nil {
				log.Printf("ERROR (handleBegin3): %+v", err)
			}

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
			if _, err = bot.Send(outMsg); err != nil {
				log.Printf("ERROR (handleBegin4): %+v", err)
			}
		}
		game.AppendPlayer(player)
		log.Printf("Player %s(%d) joined", player.Name, player.Uid)
		log.Printf("Joined players: %+v\n", game.Players)
	}

	if game.StartedAt.IsZero() {
		game.StartedAt = time.Now()
	}

	return gameplay.WaitingForPlayers
}

func handleStop(bot *tgbotapi.BotAPI, player *gameplay.Player, curGameState gameplay.GameState) gameplay.GameState {
	var replyMsg string
	var err error

	if player.Uid == game.GameStarterUID ||
		player.Uid == superUserId {
		if curGameState > 0 {
			replyMsg = "Kierros lopetettu"
		} else {
			replyMsg = "Ei peliä käynnissä"
		}
		return gameplay.StopGame
	} else {
		log.Printf("Player %s(%d) tried to stop the game, ignored",
			player.Name,
			player.Uid)
		outMsg := tgbotapi.NewMessage(channelId,
			"Et ole pelin aloittaja, joten teikäläisen natsat ei riitä")
		if _, err = bot.Send(outMsg); err != nil {
			log.Printf("ERROR (handleStop1): %+v", err)
		}
		return curGameState
	}

	outMsg := tgbotapi.NewMessage(channelId, replyMsg)
	if _, err = bot.Send(outMsg); err != nil {
		log.Printf("ERROR (handleStop2): %+v", err)
	}

	_, timeInIdle := game.HasIdleTimedOut()
	log.Printf("Gameplay stopped. Idle time was %+v\n",
		timeInIdle)
	game.Reset()

	return gameplay.InitState
}

func handleContinue(
	bot *tgbotapi.BotAPI,
	player *gameplay.Player,
	curGameState gameplay.GameState,
) gameplay.GameState {
	var replyMsg string
	var err error

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
		if _, err := bot.Send(msg); err != nil {
			log.Printf("ERROR (handleContinue1): %+v", err)
		}
		return curGameState
	}

	if player.Uid == game.GameStarterUID ||
		player.Uid == superUserId {
		replyMsg = "Jatketaan..."
		curGameState = gameplay.NextGameState(curGameState)
		replyMsg = gameplay.GetGameStateIntro(curGameState)
	} else {
		replyMsg = "Et ole pelin aloittaja, joten oikeutesi ei riitä komentelemaan"
	}

	outMsg := tgbotapi.NewMessage(channelId, replyMsg)
	if _, err = bot.Send(outMsg); err != nil {
		log.Printf("ERROR (handleContinue2): %+v", err)
	}

	return curGameState
}

func handleReview(
	bot *tgbotapi.BotAPI,
	player *gameplay.Player,
	chatId int64,
	curGameState gameplay.GameState,
	review []string,
) gameplay.GameState {
	var replyMsg string
	var err error

	if player.Uid == curSongFrom.Uid {
		log.Printf("Player %s(%d) tried to rate his/hers own song",
			player.Name,
			player.Uid)
		replyMsg = "Eh, yritit sitten arvioida oman kappaleesi"
	} else {
		var review string = strings.Join(review[2:], " ")
		// Add review for the current song presenter
		if err := curSongFrom.AddReview(player.Name, review); err != nil {
			errMsg := tgbotapi.NewMessage(
				chatId,
				err.Error())
			if _, err = bot.Send(errMsg); err != nil {
				log.Printf("ERROR (handleReview1): %+v", err)
			}
			return curGameState
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
	msg := tgbotapi.NewMessage(chatId, replyMsg)
	if _, err = bot.Send(msg); err != nil {
		log.Printf("ERROR (handleReview2): %+v", err)
	}

	if game.AllReviewsGiven() {
		curGameState = gameplay.PublishingSong
		log.Print("All reviews given")
	}

	return curGameState
}

func handlePresent(
	bot *tgbotapi.BotAPI,
	player *gameplay.Player,
	chatId int64,
	curGameState gameplay.GameState,
	review []string,
	origMsg string,
) gameplay.GameState {
	curGameState = gameplay.WaitingForSongs
	var replyMsg string

	// Get song
	var possibleUrl string = review[len(review)-1]
	songUrl, err := url.Parse(possibleUrl)
	if err != nil {
		replyMsg = "VIRHE: Viimeisenä ei ollut linkkiä"
		msg := tgbotapi.NewMessage(chatId, replyMsg)
		if _, err = bot.Send(msg); err != nil {
			log.Printf("ERROR (handlePresent1): %+v", err)
		}
		return curGameState
	}

	// Get song description
	// description := strings.Join(splitted[2:len(splitted)-1], " ")
	thirdSpace := func(msg string, cmp rune, nThMatch int) int {
		spaceCount := 0
		for i, r := range origMsg {
			if r == cmp {
				spaceCount++
			}
			if spaceCount == nThMatch {
				return i
			}
		}
		return -1
	}
	startIdx := thirdSpace(origMsg, ' ', 2)
	endIdx := strings.LastIndex(origMsg, " ")
	if startIdx == -1 || endIdx == -1 {
		log.Printf("startIdx=%d, endIdx=%d message: %s",
			startIdx,
			endIdx,
			origMsg)
		return curGameState
	}
	description := origMsg[startIdx+1 : endIdx]

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
			log.Printf("%+v\n", game.Players)
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
		if replyMsg != "" {
			msg := tgbotapi.NewMessage(channelId, replyMsg)
			if _, err = bot.Send(msg); err != nil {
				log.Printf("ERROR (handlePresent2): %+v", err)
			}
		}
	} else {
		privMsg := tgbotapi.NewMessage(chatId,
			fmt.Sprintf("Kiitos kappaleesta %s", player.Name))
		chanReply := tgbotapi.NewMessage(channelId,
			fmt.Sprintf("%s lähetti kappaleen", player.Name))

		if _, err = bot.Send(privMsg); err != nil {
			log.Printf("ERROR (handlePresent3): %+v", err)
		}
		if _, err = bot.Send(chanReply); err != nil {
			log.Printf("ERROR (handlePresent4): %+v", err)
		}
	}

	return curGameState
}

func handleQuit(gamePlay gameplay.GamePlay) gameplay.GameState {
	if allSongsPresented {
		log.Print("Publishing end results")
		resultsInJSON, err := json.Marshal(game)
		if err != nil {
			log.Printf("JSON marshalling failed: %v", err)
		} else {
			log.Printf("Results in JSON: %s\n", resultsInJSON)
		}
		dateTimeNow := time.Now().Format("2006-01-02_150405")
		reviewsFilename := fmt.Sprintf("gameplay-%s.html", dateTimeNow)
		reviewsFullPath := filepath.Join(resultsAbsPath, reviewsFilename)
		htmlData := game.PublishResults()
		err = ioutil.WriteFile(
			reviewsFullPath,
			htmlData,
			0644)
		check(err)
		log.Printf("Results written to file %s", reviewsFullPath)
		allSongsPresented = false
	}
	game.Reset()
	log.Printf("Game finished: %+v", game)

	return gameplay.InitState
}
