package db

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	"github.com/toby3d/MyPackBot/internal/models"
)

// UserState return current state for UserID
func UserState(userID int) (string, error) {
	log.Ln("Trying to get", userID, "state")
	var state string
	err := DB.View(func(tx *buntdb.Tx) error {
		var err error
		state, err = tx.Get(fmt.Sprint("user:", userID, ":state"))
		return err
	})
	if err == buntdb.ErrNotFound {
		log.Ln(userID, "not found, create new one")
		if err = ChangeUserState(userID, models.StateNone); err != nil {
			return state, err
		}
	}

	return state, err
}
