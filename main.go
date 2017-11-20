package main

import (
	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

var bot *telegram.Bot

// main function is a general function for work of this bot
func main() {
	log.Ln("[main] Let'g Get It Started...")
	var err error
	log.Ln("[main] Initializing new bot via checking access_token")
	bot, err = telegram.NewBot(cfg.UString("telegram.token"))
	errCheck(err)
}
