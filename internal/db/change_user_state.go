package db

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

// ChangeUserState change current user state on input state.
func (db *DataBase) ChangeUserState(uid int, state string) error {
	log.Ln("Trying to change", uid, "state to", state)
	return db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(fmt.Sprint("user:", uid, ":state"), state, nil)
		return err
	})
}
