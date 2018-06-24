package bot

import tg "gitlab.com/toby3d/telegram"

// Bot is a main object of Telegram bot
var Bot *tg.Bot

// New just create new bot by configuration credentials
func New(accessToken string) (bot *tg.Bot, err error) {
	return tg.New(accessToken)
}
