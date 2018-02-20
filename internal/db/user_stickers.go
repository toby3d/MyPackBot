package db

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

// UserStickers return array of saved stickers for input UserID and his total count
func UserStickers(userID, offset int, query string) ([]string, int, error) {
	log.Ln("Trying to get", userID, "stickers")
	var total, count int
	var stickers []string
	offset = offset * 50

	err := DB.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(
			fmt.Sprint("user:", userID, ":set:*"), // index
			func(key, val string) bool { // iterator
				subKeys := strings.Split(key, ":")
				if subKeys[1] != strconv.Itoa(userID) {
					return true
				}

				total++
				if count >= 51 {
					return true
				}

				if total < offset {
					return true
				}

				if query != "" && !strings.Contains(query, val) {
					return true
				}

				count++
				stickers = append(stickers, subKeys[5])
				return true
			},
		)
	})

	if err == buntdb.ErrNotFound {
		log.Ln("Not found stickers")
		return nil, total, nil
	}

	return stickers, total, err
}
