package actions

import (
	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/models"
	tg "gitlab.com/toby3d/telegram"
)

// Action function check Message update on commands, sended stickers or other
// user stuff if user state is not 'none'
func Action(msg *tg.Message) {
	state, err := db.DB.GetUserState(msg.From.ID)
	errors.Check(err)

	log.Ln("state:", state)
	switch state {
	case models.StateAddSticker:
		Add(msg, false)
	case models.StateAddPack:
		Add(msg, true)
	case models.StateDeleteSticker:
		Delete(msg, false)
	case models.StateDeletePack:
		Delete(msg, true)
	case models.StateReset:
		Reset(msg)
	default:
		err = db.DB.ChangeUserState(msg.From.ID, models.StateNone)
		errors.Check(err)

		Error(msg)
	}
}
