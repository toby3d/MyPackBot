package db

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

// GetUserStickers return array of saved stickers for input UserID and his total count
func GetUserStickers(userID, offset int, query string) ([]string, error) {
	log.Ln("Trying to get", userID, "stickers")
	var i int
	var stickers []string
	offset *= 50

	err := DB.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(
			fmt.Sprint("user:", userID, ":set:*"), // index
			func(key, val string) bool { // iterator
				subKeys := strings.Split(key, ":")
				if subKeys[1] != strconv.Itoa(userID) {
					return true
				}

				if len(stickers) == 50 {
					return false
				}

				i++
				if i < offset {
					return true
				}

				if query != "" && !strings.Contains(query, val) {
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
