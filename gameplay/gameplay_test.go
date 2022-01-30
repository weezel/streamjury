package gameplay

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAppendPlayer(t *testing.T) {
	var game GamePlay
	p1 := AddPlayer("Todd", 1)
	p2 := AddPlayer("Nancy", 2)
	p3 := AddPlayer("Alice", 3)
	game.AppendPlayer(p1)
	game.AppendPlayer(p2)
	game.AppendPlayer(p3)

	for i, v := range game.Players {
		switch i {
		case 0:
			assert.Equal(t, "Todd", v.Name)
			assert.Equal(t, 1, v.Uid)
		case 1:
			assert.Equal(t, "Nancy", v.Name)
			assert.Equal(t, 2, v.Uid)
		case 2:
			assert.Equal(t, "Alice", v.Name)
			assert.Equal(t, 3, v.Uid)
		}
	}

	game = GamePlay{}
}

func TestIsThereAnySongsLeft(t *testing.T) {
	var game GamePlay
	p1 := AddPlayer("Todd", 1)
	p2 := AddPlayer("Nancy", 2)
	p3 := AddPlayer("Alice", 3)
	game.AppendPlayer(p1)
	game.AppendPlayer(p2)
	game.AppendPlayer(p3)

	assert.True(t, game.IsThereAnySongsLeft(), "There must be songs left")

	game.Players[0].SongPresented = true
	game.Players[1].SongPresented = true
	game.Players[2].SongPresented = true
	assert.False(t, game.IsThereAnySongsLeft(), "There must not be any songs left")

	game = GamePlay{}
}

func TestAllReviewsGiven(t *testing.T) {
	var game GamePlay
	p1 := AddPlayer("Todd", 1)
	p2 := AddPlayer("Nancy", 2)
	p3 := AddPlayer("Alice", 3)
	game.AppendPlayer(p1)
	game.AppendPlayer(p2)
	game.AppendPlayer(p3)

	assert.False(t, game.AllReviewsGiven(), "There's must be reviews left")

	game.Players[0].ReviewGiven = true
	game.Players[1].ReviewGiven = true
	game.Players[2].ReviewGiven = true
	assert.True(t, game.IsThereAnySongsLeft(), "There must not be reviews left")

	game = GamePlay{}
}

func TestAppendSong(t *testing.T) {
	var game GamePlay
	game.Players = make([]Player, 3)
	game.Players[0] = Player{Name: "Todd", Uid: 1}
	game.Players[1] = Player{Name: "Nancy", Uid: 2}
	game.Players[2] = Player{Name: "Alice", Uid: 3}

	game.Players[0].AddSong("Absolutely the bestest song", "https://example.com")
	game.Players[1].AddSong("Here", "https://www.youtube.com/mysong")
	game.Players[2].AddSong("My absolute favourite song", "http://www.vimeo.com/fav")

	assert.Equal(t, "Todd", game.Players[0].Name)
	assert.Equal(
		t,
		"Absolutely the bestest song",
		game.Players[0].Song.Description)
	assert.Equal(
		t,
		"https://example.com",
		game.Players[0].Song.Url)

	assert.Equal(t, "Nancy", game.Players[1].Name)
	assert.Equal(
		t,
		"Here",
		game.Players[1].Song.Description)
	assert.Equal(
		t,
		"https://www.youtube.com/mysong",
		game.Players[1].Song.Url)

	assert.Equal(t, "Alice", game.Players[2].Name)
	assert.Equal(
		t,
		"My absolute favourite song",
		game.Players[2].Song.Description)
	assert.Equal(
		t,
		"http://www.vimeo.com/fav",
		game.Players[2].Song.Url)

	game = GamePlay{}

}

func TestPlayerShuffle(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var game GamePlay
	game.Players = make([]Player, 7)
	game.Players[0] = Player{Name: "Todd", Uid: 1}
	game.Players[1] = Player{Name: "Nancy", Uid: 2}
	game.Players[2] = Player{Name: "Alice", Uid: 3}
	game.Players[3] = Player{Name: "Bob", Uid: 4}
	game.Players[4] = Player{Name: "Maurice", Uid: 5}
	game.Players[5] = Player{Name: "Astrid", Uid: 6}
	game.Players[6] = Player{Name: "Thor", Uid: 7}

	game.ShufflePlayingOrder()

	// Test whether the first and last values are shuffled.
	// This is really a bad test since at times this will fail.
	assert.NotEmpty(t, game.Players[0].Name, "Empty name")
	assert.NotEmpty(t, game.Players[6].Name, "Empty name")
	assert.NotEqual(t, 1, game.Players[0].Uid, "Possibly not shuffled")
	assert.NotEqual(t, 7, game.Players[6].Uid, "Possibly not shuffled")

	game = GamePlay{}
}

func TestNextSong(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var game GamePlay

	game.Players = make([]Player, 3)

	game.Players[0] = Player{
		Name: "Todd",
		Uid:  1,
		Song: &Song{
			Url:         "https://example.com",
			Description: "Absolutely the bestest song"},
	}
	game.Players[1] = Player{
		Name: "Nancy",
		Uid:  2,
		Song: &Song{
			Url:         "https://www.youtube.com/mysong",
			Description: "Here"},
	}
	game.Players[2] = Player{
		Name: "Alice",
		Uid:  3,
		Song: &Song{
			Url:         "http://www.vimeo.com/fav",
			Description: "My absolute favourite song"},
	}

	game.ShufflePlayingOrder()

	// This is purposely duplicate. If two our of three has presented
	// their songs, only ony must be left after that.
	game.NextSongFrom()
	nextSongFrom := game.NextSongFrom()

	for i, p := range game.Players {
		t.Logf("%s %s - %s\n",
			nextSongFrom.Name,
			nextSongFrom.Song.Url,
			nextSongFrom.Song.Description)
		if i < len(game.Players)-1 {
			assert.True(t, p.SongPresented, "Player indexes 0 and 1 must have songs presented")
		} else {
			assert.False(t, p.SongPresented, "The last player must have the song unpresented")
		}

	}

	game = GamePlay{}
}

func TestPublishResults(t *testing.T) {
	var game GamePlay

	game.Players = make([]Player, 3)

	// Oh dear I feel dirty
	pwd, _ := os.Getwd()
	os.Chdir(filepath.Join(pwd, ".."))

	game.Players[0] = Player{
		Name: "Todd",
		Uid:  1,
		Song: &Song{
			Url:         "https://mysong.url",
			Description: "Goodish song but not my favourite",
		},
		ReceivedReviews: []Review{
			Review{
				FromPlayer: "Nancy",
				UserReview: "Terrible song",
				Rating:     5,
			},
			Review{
				FromPlayer: "Alice",
				UserReview: "Stunning uplift!",
				Rating:     10,
			},
		},
	}

	game.Players[1] = Player{
		Name: "Alice",
		Uid:  2,
		Song: &Song{
			Url:         "https://www.youtube.com/my-example-song",
			Description: "I had to do it.",
		},
		ReceivedReviews: []Review{
			Review{
				FromPlayer: "Nancy",
				UserReview: "Oh gosh, why on earth what why oh stop it!",
				Rating:     1,
			},
			Review{
				FromPlayer: "Todd",
				UserReview: "I've seen the limits of torture.",
				Rating:     2,
			},
		},
	}

	game.Players[2] = Player{
		Name: "Nancy",
		Uid:  3,
		Song: &Song{
			Url:         "https://www.vimeo.com/ambient-madness",
			Description: "Calm down and enjoy. Close your eyes and feel the breeze.",
		},
		ReceivedReviews: []Review{
			Review{
				FromPlayer: "Tood",
				UserReview: "I can smell the summer, it's here!",
				Rating:     8,
			},
			Review{
				FromPlayer: "Alice",
				UserReview: "Not my piece of cake but I i still enjoyed it.",
				Rating:     7,
			},
		},
	}
	// check := func(err error) {
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// var stringWriter bytes.Buffer
	// game.PublishResults(&stringWriter)
	// err := ioutil.WriteFile("mytest.dump", stringWriter.Bytes(), 0644)
	// check(err)

	game = GamePlay{}
}
