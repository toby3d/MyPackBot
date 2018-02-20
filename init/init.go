package init

import (
	log "github.com/kirillDanshin/dlog"
	"github.com/toby3d/MyPackBot/internal/bot"
	"github.com/toby3d/MyPackBot/internal/config"
	"github.com/toby3d/MyPackBot/internal/db"
	"github.com/toby3d/MyPackBot/internal/i18n"
)

// init prepare configuration and other things for successful start
func init() {
	log.Ln("Initializing...")
	i18n.Open("translations/")         // Preload localization strings
	config.Open("configs/config.yaml") // Preload configuration file
	db.Open("stickers.db")             // Open database or create new one
	bot.New()                          // Create bot with credentials from config
}
