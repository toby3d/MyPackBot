package db

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	tg "gitlab.com/toby3d/telegram"
)

// ResetUser just drop out all stickers keys for input UserID
func (db *DataBase) ResetUser(user *tg.User) error {
	log.Ln("Trying reset all stickers of", user.ID, "user")
	return db.Update(func(tx *buntdb.Tx) error {
		var keys []string
		if err := tx.AscendKeys(
			fmt.Sprint("user:", user.ID, ":set:*"), // index
			func(key, val string) bool { // iterator
				subKeys := strings.Split(key, ":")
				if subKeys[1] == strconv.Itoa(user.ID) {
					keys = append(keys, key)
				}
				return true
			},
		); err != nil {
			return err
		}

		for i := range keys {
			_, err := tx.Delete(keys[i])
			if err != nil {
				break
			}
		}

		return nil
	})
}
