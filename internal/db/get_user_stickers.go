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
func (db *DataBase) GetUserStickers(uid int, query *tg.InlineQuery) (stickers []string, err error) {
	log.Ln("Trying to get", uid, "stickers")
	var i int
	offset, _ := strconv.Atoi(query.Offset)
	offset *= 50

	err = db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(
			fmt.Sprint("user:", uid, ":set:*"), // index
			func(key, val string) bool { // iterator
				subKeys := strings.Split(key, ":")
				if subKeys[1] != strconv.Itoa(uid) {
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
