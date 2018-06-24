package db

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/models"
	tg "gitlab.com/toby3d/telegram"
)

// DeleteSticker just remove specified sticker key from database.
func (db *DataBase) DeleteSticker(user *tg.User, sticker *tg.Sticker) (bool, error) {
	log.Ln("Trying to remove", sticker.FileID, "sticker from", user.ID, "user")
	if sticker.SetName == "" {
		sticker.SetName = models.SetUploaded
	}

	err := db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(
			fmt.Sprint("user:", user.ID, ":set:", sticker.SetName, ":sticker:", sticker.FileID),
		)
		return err
	})
	if err == buntdb.ErrNotFound {
		log.Ln(user.ID, "not found, create new one")
		return true, nil
	}

	return false, err
}
