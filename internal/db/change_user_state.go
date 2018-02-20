package db

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

// ChangeUserState change current user state on input state.
func ChangeUserState(userID int, state string) error {
	log.Ln("Trying to change", userID, "state to", state)
	return DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(fmt.Sprint("user:", userID, ":state"), state, nil)
		return err
	})
}
