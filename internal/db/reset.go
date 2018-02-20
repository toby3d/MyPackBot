package db

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

// ResetUser just drop out all stickers keys for input UserID
func ResetUser(userID int) error {
	log.Ln("Trying reset all stickers of", userID, "user")
	return DB.Update(func(tx *buntdb.Tx) error {
		var keys []string
		if err := tx.AscendKeys(
			fmt.Sprint("user:", userID, ":set:*"), // index
			func(key, val string) bool { // iterator
				subKeys := strings.Split(key, ":")
				if subKeys[1] == strconv.Itoa(userID) {
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
