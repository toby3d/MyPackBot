package internal

import (
	"github.com/spf13/viper"
	"gitlab.com/toby3d/mypackbot/internal/config"
	storemodel "gitlab.com/toby3d/mypackbot/internal/model/store"
	"gitlab.com/toby3d/mypackbot/internal/store"
	tg "gitlab.com/toby3d/telegram"
)

type MyPackBot struct {
	bot    *tg.Bot
	config *viper.Viper
	store  storemodel.Manager
}

func New(path string) (*MyPackBot, error) {
	var mpb MyPackBot
	mpb.store = store.NewInMemoryStore()

	var err error
	if mpb.config, err = config.Open(path); err != nil {
		return nil, err
	}

	if mpb.bot, err = tg.New(mpb.config.GetString("telegram.token")); err != nil {
		return nil, err
	}

	return &mpb, nil
}

func (mpb *MyPackBot) Bot() *tg.Bot {
	return mpb.bot
}

func (mpb *MyPackBot) Store() storemodel.Manager {
	return mpb.store
}
