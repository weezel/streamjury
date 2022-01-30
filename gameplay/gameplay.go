package gameplay

import (
	_ "embed"
	"math/rand"
	"time"
)

const (
	GameIdleTimeout time.Duration = 15 * 60 * time.Second
)

type GamePlay struct {
	// TODO won't be needed sync.Mutex
	StartedAt        time.Time `json:"time"`
	GameStarterUID   int64     `json:"game_starter_uid"`
	CurrentPresenter *Player   `json:"curren_presenter"`
	Players          []Player  `json:"players"`
}

func (gameplay GamePlay) IsInGame(userId int64) bool {
	for _, player := range gameplay.Players {
		if player.Uid == userId {
			return true
		}
	}
	return false
}

func (g *GamePlay) ShufflePlayingOrder() {
	rand.Shuffle(len(g.Players), func(i int, j int) {
		g.Players[i], g.Players[j] = g.Players[j], g.Players[i]
	})
}

func (g GamePlay) HasIdleTimedOut() (bool, time.Duration) {
	idleTime := time.Now().Sub(g.StartedAt)
	if idleTime > GameIdleTimeout {
		return true, idleTime
	}
	return false, idleTime
}

func (g *GamePlay) Reset() {
	for i := range g.Players {
		g.Players[i].Song = nil
		for r := range g.Players[i].ReceivedReviews {
			g.Players[i].ReceivedReviews[r] = Review{}
		}
		g.Players[i].ReceivedReviews = nil
		g.Players[i] = Player{}
	}
	g.Players = []Player{}
	g = &GamePlay{}
}

func (g *GamePlay) AppendPlayer(p *Player) {
	g.Players = append(g.Players, *p)
}

func (g *GamePlay) IsThereAnySongsLeft() bool {
	for _, p := range g.Players {
		if p.SongPresented == false {
			return true
		}
	}
	return false
}

func (g *GamePlay) AllReviewsGiven() bool {
	for _, p := range g.Players {
		if p.ReviewGiven == false {
			return false
		}
	}
	return true
}

func (g *GamePlay) NextSongFrom() *Player {
	for i, p := range g.Players {
		// Song not given yet,
		// though this shouldn't be possible
		// at this point in game
		if g.Players[i].Song == nil {
			return nil
		}

		if p.SongPresented == false {
			g.Players[i].SongPresented = true
			return &g.Players[i]
		}
	}
	return nil
}
