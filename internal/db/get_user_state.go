package db

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	"github.com/toby3d/MyPackBot/internal/models"
	tg "github.com/toby3d/telegram"
)

// UserState return current state for UserID
func (db *DataBase) GetUserState(user *tg.User) (string, error) {
	log.Ln("Trying to get", user.ID, "state")
	var state string
	err := db.View(func(tx *buntdb.Tx) error {
		var err error
		state, err = tx.Get(fmt.Sprint("user:", user.ID, ":state"))
		return err
	})
	if err == buntdb.ErrNotFound {
		log.Ln(user.ID, "not found, create new one")
		if err = db.ChangeUserState(user, models.StateNone); err != nil {
			return state, err
		}
	}

	return state, err
}
