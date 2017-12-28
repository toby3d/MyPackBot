package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	log "github.com/kirillDanshin/dlog"  // Insert logs only in debug builds
	"github.com/nicksnyder/go-i18n/i18n" // Internationalization and localization
	"github.com/olebedev/config"         // Easy configuration file parsing
)

var (
	// Variables with types from imports
	cfg       *config.Config
	channelID int64

	// Setted variables
	flagWebhook = flag.Bool(
		"webhook",
		false,
		"enable work via webhooks (required valid certificates)",
	)
)

// init prepare configuration and other things for successful start of main
// function.
func init() {
	log.Ln("Initializing...")
	log.Ln("Parse flags...")
	flag.Parse()

	err := filepath.Walk("./i18n/", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".all.yaml") {
			i18n.MustLoadTranslationFile(path)
		}
		return nil
	})
	errCheck(err)

	log.Ln("Loading configuration file...")
	cfg, err = config.ParseYamlFile("config.yaml")
	errCheck(err)

	log.Ln("Checking bot access token in configuration file...")
	_, err = cfg.String("telegram.token")
	errCheck(err)

	if *flagWebhook {
		log.Ln("Enabled webhook mode, check configuration strings...")
		log.Ln("Checking webhook set string...")
		_, err = cfg.String("telegram.webhook.set")
		errCheck(err)

		log.Ln("Checking webhook listen string...")
		_, err = cfg.String("telegram.webhook.listen")
		errCheck(err)

		log.Ln("Checking webhook listen string...")
		_, err = cfg.String("telegram.webhook.serve")
		errCheck(err)
	}

	channelID = int64(cfg.UInt("telegram.channel"))
}
