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
	var err error

	// Preload configuration file
	config.Config, err = config.Open("./configs")
	errors.Check(err)

	// Preload localization strings
	err = i18n.Open("i18n/")
	errors.Check(err)

	// Open database or create new one
	db.DB, err = db.Open("stickers.db")
	errors.Check(err)

	// Create bot with credentials from config
	bot.Bot, err = bot.New(config.Config.GetString("telegram.token"))
	errors.Check(err)
}
