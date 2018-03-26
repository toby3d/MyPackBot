package init

import (
	log "github.com/kirillDanshin/dlog"
	"github.com/toby3d/MyPackBot/internal/bot"
	"github.com/toby3d/MyPackBot/internal/config"
	"github.com/toby3d/MyPackBot/internal/db"
	"github.com/toby3d/MyPackBot/internal/errors"
	"github.com/toby3d/MyPackBot/internal/i18n"
)

// init prepare configuration and other things for successful start
func init() {
	log.Ln("Initializing...")

	// Preload localization strings
	err := i18n.Open("i18n/")
	errors.Check(err)

	// Preload configuration file
	config.Open("configs/config.yaml")

	// Open database or create new one
	db.Open("stickers.db")

	// Create bot with credentials from config
	bot.New()
}
