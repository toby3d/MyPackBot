package actions

import (
	log "github.com/kirillDanshin/dlog"
	"github.com/toby3d/MyPackBot/internal/db"
	"github.com/toby3d/MyPackBot/internal/errors"
	"github.com/toby3d/MyPackBot/internal/models"
	tg "github.com/toby3d/telegram"
)

// Action function check Message update on commands, sended stickers or other
// user stuff if user state is not 'none'
func Action(msg *tg.Message) {
	state, err := db.DB.GetUserState(msg.From)
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
		err = db.DB.ChangeUserState(msg.From, models.StateNone)
		errors.Check(err)

		Error(msg)
	}
}
