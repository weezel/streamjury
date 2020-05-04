package gameplay

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddReview(t *testing.T) {
	var game GamePlay

	game.Players = make([]Player, 3)

	game.Players[0] = Player{
		Name: "Todd",
		Uid:  1,
		Song: &Song{
			Url:         "https://mysong.url",
			Description: "Goodish song but not my favourite",
		},
	}
	game.Players[0].AddReview("Nancy", "Terrible song 1/10")
	game.Players[0].AddReview("Alice", "Stunning uplift! 10/10")
	assert.Equal(t, game.Players[0].ReceivedReviews[0].Rating, 1)
	assert.Equal(t, game.Players[0].ReceivedReviews[1].Rating, 10)

	game.Players[1] = Player{
		Name: "Alice",
		Uid:  2,
		Song: &Song{
			Url:         "https://www.youtube.com/my-example-song",
			Description: "I had to do it.",
		},
	}
	game.Players[1].AddReview("Nancy", "Oh gosh, why on earth what why oh stop it! 1/10")
	game.Players[1].AddReview("Todd", "I've seen the limits of torture 2/10")
	assert.Equal(t, game.Players[1].ReceivedReviews[0].Rating, 1)
	assert.Equal(t, game.Players[1].ReceivedReviews[1].Rating, 2)

	game.Players[2] = Player{
		Name: "Nancy",
		Uid:  3,
		Song: &Song{
			Url:         "https://www.vimeo.com/ambient-madness",
			Description: "Calm down and enjoy. Close your eyes and feel the breeze.",
		},
	}
	game.Players[2].AddReview("Todd", "I can smell the summer, it's here! 8/10")
	game.Players[2].AddReview("Alice", "Not my piece of cake but I ate it. 7/10")
	assert.Equal(t, game.Players[2].ReceivedReviews[0].Rating, 8)
	assert.Equal(t, game.Players[2].ReceivedReviews[1].Rating, 7)

	game = GamePlay{}
}
