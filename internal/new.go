package internal

import (
	json "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"gitlab.com/toby3d/mypackbot/internal/config"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/model/photos"
	"gitlab.com/toby3d/mypackbot/internal/model/stickers"
	"gitlab.com/toby3d/mypackbot/internal/model/users"
	usersphotos "gitlab.com/toby3d/mypackbot/internal/model/users/photos"
	usersstickers "gitlab.com/toby3d/mypackbot/internal/model/users/stickers"
	"gitlab.com/toby3d/mypackbot/internal/store"
	tg "gitlab.com/toby3d/telegram"
)

type MyPackBot struct {
	bot           *tg.Bot
	config        *viper.Viper
	photos        photos.Manager
	stickers      stickers.Manager
	users         users.Manager
	usersPhotos   usersphotos.Manager
	usersStickers usersstickers.Manager
}

func New(path string) (mpb *MyPackBot, err error) {
	mpb = new(MyPackBot)

	if mpb.config, err = config.Open(path); err != nil {
		return nil, err
	}

	conn, err := db.Open(mpb.config.GetString("database.filepath"))
	if err != nil {
		return nil, err
	}

	marshler := json.ConfigFastest
	mpb.photos = store.NewPhotosStore(conn, marshler)
	mpb.stickers = store.NewStickersStore(conn, marshler)
	mpb.users = store.NewUsersStore(conn, marshler)
	mpb.usersPhotos = store.NewUsersPhotosStore(conn, mpb.users, mpb.photos, marshler)
	mpb.usersStickers = store.NewUsersStickersStore(conn, mpb.users, mpb.stickers, marshler)

	if mpb.bot, err = tg.New(mpb.config.GetString("telegram.token")); err != nil {
		return nil, err
	}

	return mpb, nil
}
