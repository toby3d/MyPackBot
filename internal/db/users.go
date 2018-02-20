package db

import (
	"strconv"
	"strings"

	// log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

// Users return array of all avaliable UserID in database
func Users() ([]int, error) {
	var users []int
	err := DB.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(
			"user:*:state",
			func(key, val string) bool {
				subKeys := strings.Split(key, ":")
				id, err := strconv.Atoi(subKeys[1])
				if err == nil {
					users = append(users, id)
				}

				return true
			},
		)
	})
	if err == buntdb.ErrNotFound {
		return nil, nil
	}

	return users, err
}
