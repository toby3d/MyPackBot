package config

import (
	"github.com/olebedev/config"
	"github.com/toby3d/MyPackBot/internal/errors"
)

var (
	// Config is a main object of preloaded configuration YAML
	Config *config.Config

	// ChannelID is a announcements channel ID
	ChannelID int64
)

// Open just open configuration file for parsing some data in other functions
func Open(path string) {
	var err error
	Config, err = config.ParseYamlFile(path)
	errors.Check(err)

	ChannelID = int64(Config.UInt("telegram.channel"))
}
