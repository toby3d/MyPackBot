package main

import (
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

// actions function check Message update on commands, sended stickers or other user stuff
func actions(msg *tg.Message) {
	state, err := dbGetUserState(msg.From.ID)
	errCheck(err)

	log.Ln("state:", state)
	switch state {
	case stateAddSticker:
		actionAdd(msg, false)
	case stateAddPack:
		actionAdd(msg, true)
	case stateDeleteSticker:
		actionDelete(msg, false)
	case stateDeletePack:
		actionDelete(msg, true)
	case stateReset:
		actionReset(msg)
	default:
		err = dbChangeUserState(msg.From.ID, stateNone)
		errCheck(err)

		actionError(msg)
	}
}
