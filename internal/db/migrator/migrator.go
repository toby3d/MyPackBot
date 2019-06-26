package migrator

import (
	"strconv"
	"strings"

	bolt "github.com/etcd-io/bbolt"
	bunt "github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/store"
)

func AutoMigrate(srcPath string, dst *bolt.DB) error {
	db, err := bunt.Open(srcPath)
	if err != nil {
		return err
	}

	users := make(map[int]models.User)
	stickers := make(map[string]models.Sticker)
	usersStickers := make(map[string]int)
	if err = db.View(func(tx *bunt.Tx) error {
		return tx.AscendKeys("*", func(k, v string) bool {
			parts := strings.Split(k, ":")
			uid, err := strconv.Atoi(parts[1])
			if err != nil {
				return true
			}

			switch {
			case parts[2] == "state":
				users[uid] = models.User{
					ID:           uid,
					LanguageCode: "en",
					AutoSaving:   true,
				}
			case parts[2] == "set":
				setName := parts[3]
				if setName == "?" {
					setName = ""
				}

				stickers[parts[5]] = models.Sticker{
					Model:   models.Model{ID: parts[5]},
					Emoji:   v,
					SetName: setName,
				}

				usersStickers[parts[5]] = uid
			}

			return true
		})
	}); err != nil {
		return err
	}
	if err = db.Close(); err != nil {
		return err
	}

	newStore, err := store.New(dst)
	if err != nil {
		return err
	}

	for _, u := range users {
		if err = newStore.CreateUser(&u); err != nil {
			return err
		}
	}

	for _, s := range stickers {
		if err = newStore.CreateSticker(&s); err != nil {
			return err
		}
	}

	for sid, uid := range usersStickers {
		if err = newStore.AddSticker(uid, sid); err != nil {
			return err
		}
	}

	return nil
}
