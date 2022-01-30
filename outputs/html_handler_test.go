package outputs

import (
	"os"
	"streamjury/gameplay"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGamePlay_PublishResults(t *testing.T) {
	var game gameplay.GamePlay

	game.Players[0] = gameplay.Player{
		Name: "Todd",
		Uid:  1,
		Song: &gameplay.Song{
			Url:         "https://mysong.url",
			Description: "Goodish song but not my favourite",
		},
		ReceivedReviews: []gameplay.Review{
			{
				FromPlayer: "Nancy",
				UserReview: "Terrible song",
				Rating:     5,
			},
			{
				FromPlayer: "Alice",
				UserReview: "Stunning uplift!",
				Rating:     10,
			},
		},
	}

	game.Players[1] = gameplay.Player{
		Name: "Alice",
		Uid:  2,
		Song: &gameplay.Song{
			Url:         "https://www.youtube.com/my-example-song",
			Description: "I had to do it.",
		},
		ReceivedReviews: []gameplay.Review{
			{
				FromPlayer: "Nancy",
				UserReview: "Oh gosh, why on earth what why oh stop it!",
				Rating:     1,
			},
			{
				FromPlayer: "Todd",
				UserReview: "I've seen the limits of torture.",
				Rating:     2,
			},
		},
	}

	game.Players[2] = gameplay.Player{
		Name: "Nancy",
		Uid:  3,
		Song: &gameplay.Song{
			Url:         "https://www.vimeo.com/ambient-madness",
			Description: "Calm down and enjoy. Close your eyes and feel the breeze.",
		},
		ReceivedReviews: []gameplay.Review{
			{
				FromPlayer: "Tood",
				UserReview: "I can smell the summer, it's here!",
				Rating:     8,
			},
			{
				FromPlayer: "Alice",
				UserReview: "Not my piece of cake but I i still enjoyed it.",
				Rating:     7,
			},
		},
	}

}

func TestPublishResultsInHTML(t *testing.T) {
	var game gameplay.GamePlay = gameplay.GamePlay{}
	game.Players = make([]gameplay.Player, 3)

	game.StartedAt = time.Date(2022, 1, 12, 14, 39, 19, 0, time.UTC)
	game.Players[0] = gameplay.Player{
		Name: "Todd",
		Uid:  1,
		Song: &gameplay.Song{
			Url:         "https://mysong.url",
			Description: "Goodish song but not my favourite",
		},
		ReceivedReviews: []gameplay.Review{
			{
				FromPlayer: "Nancy",
				UserReview: "Terrible song",
				Rating:     5,
			},
			{
				FromPlayer: "Alice",
				UserReview: "Stunning #>S$! uplift!",
				Rating:     10,
			},
		},
	}

	game.Players[1] = gameplay.Player{
		Name: "Alice",
		Uid:  2,
		Song: &gameplay.Song{
			Url:         "https://www.youtube.com/my-example-song",
			Description: "I had to do it.",
		},
		ReceivedReviews: []gameplay.Review{
			{
				FromPlayer: "Nancy",
				UserReview: "Oh gosh, why on earth what why oh stop it!",
				Rating:     1,
			},
			{
				FromPlayer: "Todd",
				UserReview: "I've seen the limits of torture.",
				Rating:     2,
			},
		},
	}

	game.Players[2] = gameplay.Player{
		Name: "Nancy",
		Uid:  3,
		Song: &gameplay.Song{
			Url:         "https://www.vimeo.com/ambient-madness",
			Description: "Calm down and enjoy. Close your eyes and feel the breeze.",
		},
		ReceivedReviews: []gameplay.Review{
			{
				FromPlayer: "Tood",
				UserReview: "I can smell the summer, it's here!",
				Rating:     8,
			},
			{
				FromPlayer: "Alice",
				UserReview: "Not my piece of cake but I i still enjoyed it.",
				Rating:     7,
			},
		},
	}

	type args struct {
		g gameplay.GamePlay
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "aa",
			args:    args{g: game},
			want:    []byte{'a'},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PublishResultsInHTML(tt.args.g)
			if (err != nil) != tt.wantErr {
				t.Errorf("PublishResultsInHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Exppected output differs:\n%s", diff)
				if err := os.WriteFile("dingdong.html", got, 0600); err != nil {
					panic(err)
				}
			}
		})
	}
}
