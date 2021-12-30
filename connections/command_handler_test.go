package connections

import (
	"net/http"
	"reflect"
	"streamjury/gameplay"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type myBotIF interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type User struct {
	ID           int
	FirstName    string
	LastName     string
	UserName     string
	LanguageCode string
	IsBot        bool
}

type MockBot struct {
	Client *http.Client
}

func (b *MockBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return tgbotapi.Message{}, nil
}

func Test_handleBegin(t *testing.T) {
	// a := tgbotapi.NewBotAPIWithClient()
	var testbot myBotIF = &MockBot{}
	type args struct {
		bot    BotIF
		player *gameplay.Player
	}
	tests := []struct {
		name string
		args args
		want gameplay.GameState
	}{
		{
			"All songs given",
			args{
				bot: testbot,
				player: &gameplay.Player{
					Name:          "a",
					Uid:           1,
					ReviewGiven:   false,
					SongSubmitted: false,
					SongPresented: false,
					Song: &gameplay.Song{
						Description: "b",
						Url:         "c",
					},
					ReceivedReviews: []gameplay.Review{
						{
							Rating:     0,
							FromPlayer: "jorma",
							UserReview: "kiva",
						},
						{
							Rating:     0,
							FromPlayer: "luke",
							UserReview: "omsakas",
						},
					},
				},
			},
			gameplay.WaitingForPlayers,
		},
		{
			"Missing song",
			args{
				bot: testbot,
				player: &gameplay.Player{
					Name:          "a",
					Uid:           1,
					ReviewGiven:   false,
					SongSubmitted: true,
					SongPresented: false,
					Song: &gameplay.Song{
						Description: "b",
						Url:         "c",
					},
					ReceivedReviews: []gameplay.Review{},
				},
			},
			gameplay.WaitingForPlayers,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleBegin(tt.args.bot, tt.args.player); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleBegin() = %v, want %v", got, tt.want)
			}
		})
	}
}
