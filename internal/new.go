package internal

import (
	bolt "github.com/etcd-io/bbolt"
	"github.com/spf13/viper"
	"gitlab.com/toby3d/mypackbot/internal/config"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/models/sticker"
	"gitlab.com/toby3d/mypackbot/internal/models/user"
	"gitlab.com/toby3d/mypackbot/internal/store"
	tg "gitlab.com/toby3d/telegram"
)

type MyPackBot struct {
	bot          *tg.Bot
	db           *bolt.DB
	userStore    user.Store
	stickerStore sticker.Store
	config       *viper.Viper
	updates      tg.UpdatesChannel
}

func New(path string) (*MyPackBot, error) {
	var mpb MyPackBot

	var err error
	if mpb.config, err = config.Open(path); err != nil {
		return nil, err
	}

	if mpb.db, err = db.Open(mpb.config.GetString("database.filepath")); err != nil {
		return nil, err
	}

	mpb.userStore = store.NewUserStore(mpb.db)
	mpb.stickerStore = store.NewStickerStore(mpb.db)

	if mpb.bot, err = tg.New(mpb.config.GetString("telegram.token")); err != nil {
		return nil, err
	}

	return &mpb, nil
}
