package db

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	tg "gitlab.com/toby3d/telegram"
)

// ChangeUserState change current user state on input state.
func (db *DataBase) ChangeUserState(user *tg.User, state string) error {
	log.Ln("Trying to change", user.ID, "state to", state)
	return db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(fmt.Sprint("user:", user.ID, ":state"), state, nil)
		return err
	})
}
