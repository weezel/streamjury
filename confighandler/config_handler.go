package confighandler

import (
	"log"

	toml "github.com/pelletier/go-toml"
)

type StreamjuryConfig struct {
	SuperUserId    int
	ChannelId      int64
	ApiKey         string
	ResultsAbsPath string
}

type TomlConfig struct {
	StreamjuryConfig StreamjuryConfig
}

func LoadConfig(filedata []byte) TomlConfig {
	config := TomlConfig{}
	if err := toml.Unmarshal(filedata, &config); err != nil {
		log.Panicf("Error parsing config: %s", err)
	}
	return config
}
