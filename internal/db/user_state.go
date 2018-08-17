package db

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/models"
)

// UserState return current state for UserID
func (db *DataBase) UserState(uid int) (string, error) {
	log.Ln("Trying to get", uid, "state")
	var state string
	err := DB.View(func(tx *buntdb.Tx) error {
		var err error
		state, err = tx.Get(fmt.Sprint("user:", uid, ":state"))
		return err
	})

	if err == buntdb.ErrNotFound {
		log.Ln(uid, "not found, create new one")
		if err = db.ChangeUserState(uid, models.StateNone); err != nil {
			return state, err
		}
	}

	return state, err
}
