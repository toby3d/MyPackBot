package db

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	"github.com/toby3d/MyPackBot/internal/models"
	tg "github.com/toby3d/telegram"
)

// AddSticker add sticker FileID, Emoji and SetName meta for UserID
func (db *DataBase) AddSticker(user *tg.User, sticker *tg.Sticker) (bool, error) {
	log.Ln("Trying to add", sticker.FileID, "sticker from", user.ID, "user")
	if sticker.SetName == "" {
		sticker.SetName = models.SetUploaded
	}

	var exists bool
	err := db.Update(func(tx *buntdb.Tx) error {
		var err error
		_, exists, err = tx.Set(
			fmt.Sprint("user:", user.ID, ":set:", sticker.SetName, ":sticker:", sticker.FileID), // key
			sticker.Emoji, // value
			nil,           // options
		)
		if err == buntdb.ErrIndexExists {
			exists = true
			return nil
		}

		return err
	})

	return exists, err
}
