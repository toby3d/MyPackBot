package db

import (
	"fmt"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/models"
	tg "gitlab.com/toby3d/telegram"
)

// DeletePack remove all keys for UserID which contains input SetName
func (db *DataBase) DeletePack(uid int, sticker *tg.Sticker) (bool, error) {
	log.Ln("Trying to remove all", sticker.SetName, "sticker from", uid, "user")
	if sticker.SetName == "" {
		sticker.SetName = models.SetUploaded
	}

	var ids []string
	err := db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(
			fmt.Sprint("user:", uid, ":set:", sticker.SetName, ":*"),
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
		notExist, err = db.DeleteSticker(uid, &tg.Sticker{FileID: id})
		if err != nil {
			return notExist, err
		}
	}

	if err == buntdb.ErrNotFound {
		log.Ln(uid, "not found")
		return true, nil
	}

	return false, err
}
