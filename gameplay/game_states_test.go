package gameplay

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCommandFeasibleInState(t *testing.T) {
	assert.True(t, IsCommandFeasibleInState(InitState, "aloita"))
	assert.False(t, IsCommandFeasibleInState(InitState, "jatka"))
	assert.False(t, IsCommandFeasibleInState(InitState, "lopeta"))
	assert.True(t, IsCommandFeasibleInState(WaitingForPlayers, "aloita"))
	assert.True(t, IsCommandFeasibleInState(WaitingForSongs, "lopeta"))
	assert.False(t, IsCommandFeasibleInState(PublishingSong, "jatka"))
	assert.True(t, IsCommandFeasibleInState(WaitingForReviews, "arvioi"))
	assert.True(t, IsCommandFeasibleInState(WaitingForReviews, "arvio"))
	assert.True(t, IsCommandFeasibleInState(WaitingForReviews, "arvostele"))
	assert.False(t, IsCommandFeasibleInState(WaitingForReviews, "arvost"))
	assert.False(t, IsCommandFeasibleInState(StopGame, "jatka"))
}
