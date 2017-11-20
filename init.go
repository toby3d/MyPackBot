package main

import (
	"flag"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/olebedev/config"        // Easy configuration file parsing
)

var (
	tgHookSet, tgHookListen, tgHookServe string

	cfg *config.Config

	flagWebhook = flag.Bool("webhook", false, "enable work via webhooks (required valid certificates)")
)

// init prepare configuration and other things for successful start of main
// function.
func init() {
	log.Ln("[init] Initializing...")
	log.Ln("[init] Parse flags...")
	flag.Parse()

	var err error
	log.Ln("[init] Loading configuration file...")
	cfg, err = config.ParseYamlFile("config.yaml")
	errCheck(err)

	if *flagWebhook {
		log.Ln("[init] Enabled webhook mode, check configuration strings...")
		log.Ln("[init] Checking webhook set string...")
		tgHookSet, err = cfg.String("telegram.webhook.set")
		errCheck(err)

		log.Ln("[init] Checking webhook listen string...")
		tgHookListen, err = cfg.String("telegram.webhook.listen")
		errCheck(err)

		log.Ln("[init] Checking webhook listen string...")
		tgHookServe, err = cfg.String("telegram.webhook.serve")
		errCheck(err)
	}
}
