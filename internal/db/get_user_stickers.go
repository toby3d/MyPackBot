package db

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	tg "gitlab.com/toby3d/telegram"
)

// GetUserStickers return array of saved stickers for input UserID and his total count
func (db *DataBase) GetUserStickers(user *tg.User, query *tg.InlineQuery) ([]string, error) {
	log.Ln("Trying to get", user.ID, "stickers")
	var i int
	var stickers []string
	offset, _ := strconv.Atoi(query.Offset)
	offset *= 50

	err := db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(
			fmt.Sprint("user:", user.ID, ":set:*"), // index
			func(key, val string) bool { // iterator
				subKeys := strings.Split(key, ":")
				if subKeys[1] != strconv.Itoa(user.ID) {
					return true
				}

				if len(stickers) == 50 {
					return false
				}

				i++
				if i < offset {
					return true
				}

				if query.Query != "" && !strings.Contains(query.Query, val) {
					return true
				}

				stickers = append(stickers, subKeys[5])
				return true
			},
		)
	})

	if err == buntdb.ErrNotFound {
		log.Ln("Not found stickers")
		return nil, nil
	}

	return stickers, err
}
