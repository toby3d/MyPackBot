package bot

import (
	"github.com/toby3d/MyPackBot/internal/config"
	"github.com/toby3d/MyPackBot/internal/errors"
	tg "github.com/toby3d/telegram"
)

// Bot is a main object of Telegram bot
var Bot *tg.Bot

// New just create new bot by configuration credentials
func New() {
	accessToken, err := config.Config.String("telegram.token")
	errors.Check(err)

	Bot, err = tg.NewBot(accessToken)
	errors.Check(err)
}
