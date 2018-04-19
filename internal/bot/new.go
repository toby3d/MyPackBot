package bot

import tg "github.com/toby3d/telegram"

// Bot is a main object of Telegram bot
var Bot *tg.Bot

// New just create new bot by configuration credentials
func New(accessToken string) (*tg.Bot, error) {
	return tg.New(accessToken)
}
