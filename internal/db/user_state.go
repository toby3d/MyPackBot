package db

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	"gitlab.com/toby3d/mypackbot/internal/models"
	tg "gitlab.com/toby3d/telegram"
)

// UserState return current state for UserID
func (db *DataBase) UserState(usr *tg.User) (string, error) {
	log.Ln("Trying to get", usr.ID, "state")
	var state string
	err := DB.View(func(tx *buntdb.Tx) error {
		var err error
		state, err = tx.Get(fmt.Sprint("user:", usr.ID, ":state"))
		return err
	})

	switch err {
	case buntdb.ErrNotFound:
		log.Ln(usr.ID, "not found, create new one")
		if err = db.ChangeUserState(usr, models.StateNone); err != nil {
			return state, err
		}
	}

	return state, err
}
