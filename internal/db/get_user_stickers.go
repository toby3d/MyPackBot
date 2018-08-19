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
	page, _ := strconv.Atoi(query.Offset)
	from := page * 50

	var i int
	err = db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(
			fmt.Sprint("user:", uid, ":set:*"), // index
			func(key, val string) bool { // iterator
				if len(stickers) >= 50 {
					return false
				}

				subKeys := strings.Split(key, ":")
				if !strings.EqualFold(subKeys[1], strconv.Itoa(uid)) {
					return true
				}

				if query.Query != "" && !strings.ContainsAny(query.Query, val) {
					return true
				}

				i++
				switch {
				case i <= from:
					return true
				default:
					stickers = append(stickers, subKeys[5])
					return true
				}
			},
		)
	})
	if err == buntdb.ErrNotFound {
		log.Ln("Not found stickers")
		return nil, nil
	}

	return stickers, err
}
