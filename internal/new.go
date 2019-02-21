package internal

import (
	"gitlab.com/toby3d/mypackbot/internal/config"
	"gitlab.com/toby3d/mypackbot/internal/db"
	tg "gitlab.com/toby3d/telegram"
)

type MyPackBot struct {
	bot    *tg.Bot
	db     db.Reader
	config config.Reader
}

func New(path string) (*MyPackBot, error) {
	var mpb MyPackBot

	var err error
	if mpb.config = config.Open(path); err != nil {
		return nil, err
	}

	if mpb.db, err = db.Open("./stickers.db"); err != nil {
		return nil, err
	}

	if mpb.bot, err = tg.New(mpb.GetString("telegram.token")); err != nil {
		return nil, err
	}

	return &mpb, nil
}
