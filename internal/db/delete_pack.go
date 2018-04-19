package db

import (
	"fmt"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	"github.com/toby3d/MyPackBot/internal/models"
	tg "github.com/toby3d/telegram"
)

// DeletePack remove all keys for UserID which contains input SetName
func (db *DataBase) DeletePack(user *tg.User, sticker *tg.Sticker) (bool, error) {
	log.Ln("Trying to remove all", sticker.SetName, "sticker from", user.ID, "user")
	if sticker.SetName == "" {
		sticker.SetName = models.SetUploaded
	}

	var ids []string
	err := db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(
			fmt.Sprint("user:", user.ID, ":set:", sticker.SetName, ":*"),
			func(key, val string) bool {
				keys := strings.Split(key, ":")
				ids = append(ids, keys[5])
				return true
			},
		)
	})

	if len(ids) == 0 {
		return true, nil
	}

	for _, id := range ids {
		var notExist bool
		notExist, err = db.DeleteSticker(user, &tg.Sticker{FileID: id})
		if err != nil {
			return notExist, err
		}
	}

	switch err {
	case buntdb.ErrNotFound:
		log.Ln(user.ID, "not found")
		return true, nil
	}

	return false, err
}
