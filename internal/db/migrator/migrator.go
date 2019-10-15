package migrator

import (
	"strconv"
	"strings"
	"time"

	bolt "github.com/etcd-io/bbolt"
	bunt "github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/store"
)

func AutoMigrate(srcPath string, dst *bolt.DB) error {
	db, err := bunt.Open(srcPath)
	if err != nil {
		return err
	}

	users := make(map[int]model.User)
	stickers := make(map[string]model.Sticker)
	usersStickers := make(map[string]int)
	if err = db.View(func(tx *bunt.Tx) error {
		return tx.AscendKeys("*", func(k, v string) bool {
			parts := strings.Split(k, ":")
			uid, err := strconv.Atoi(parts[1])
			if err != nil {
				return true
			}

			timeStamp := time.Now().UTC().Unix()
			switch {
			case parts[2] == "state":
				users[uid] = model.User{
					ID:           uid,
					LanguageCode: "en",
					CreatedAt:    timeStamp,
					UpdatedAt:    timeStamp,
					LastSeen:     timeStamp,
				}
			case parts[2] == "set":
				setName := parts[3]
				if setName == "?" {
					setName = "uploaded_by_mypackbot"
				}

				stickers[parts[5]] = model.Sticker{
					ID:         parts[5],
					Emoji:      v,
					SetName:    setName,
					CreatedAt:  timeStamp,
					IsAnimated: false,
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

	newStore := store.NewInMemoryStore()
	for _, u := range users {
		u := u
		if _, err = newStore.Users().GetOrCreate(&u); err != nil {
			return err
		}
	}

	for _, s := range stickers {
		s := s
		if _, err = newStore.Stickers().GetOrCreate(&s); err != nil {
			return err
		}
	}

	for sid, uid := range usersStickers {
		sid, uid := sid, uid
		if err = newStore.AddSticker(&model.User{ID: uid}, &model.Sticker{ID: sid}); err != nil {
			return err
		}
	}

	return nil
}
