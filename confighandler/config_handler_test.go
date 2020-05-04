package confighandler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigHandler(t *testing.T) {
	data := []byte(`
[streamjuryconfig]
SuperUserId = 123456
ChannelId = -987654
ApiKey = "abcdefg:1234"
ResultsAbsPath = "/var/www/blargh/"
`)
	tomlConfig := LoadConfig(data)
	assert.Equal(t, 123456, tomlConfig.StreamjuryConfig.SuperUserId)
	assert.Equal(t, int64(-987654), tomlConfig.StreamjuryConfig.ChannelId)
	assert.Equal(t, "abcdefg:1234", tomlConfig.StreamjuryConfig.ApiKey)
	assert.Equal(t, "/var/www/blargh/", tomlConfig.StreamjuryConfig.ResultsAbsPath)
}
