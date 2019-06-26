package internal

import (
	"github.com/spf13/viper"
	"gitlab.com/toby3d/mypackbot/internal/config"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/store"
	tg "gitlab.com/toby3d/telegram"
)

type MyPackBot struct {
	bot     *tg.Bot
	store   *store.Store
	config  *viper.Viper
	updates tg.UpdatesChannel
}

func New(path string) (*MyPackBot, error) {
	var mpb MyPackBot

	var err error
	if mpb.config, err = config.Open(path); err != nil {
		return nil, err
	}

	dataBase, err := db.Open(mpb.config.GetString("database.filepath"))
	if err != nil {
		return nil, err
	}

	if mpb.store, err = store.New(dataBase); err != nil {
		dataBase.Close()
		return nil, err
	}

	if mpb.bot, err = tg.New(mpb.config.GetString("telegram.token")); err != nil {
		dataBase.Close()
		return nil, err
	}

	return &mpb, nil
}
