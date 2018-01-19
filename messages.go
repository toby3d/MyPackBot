package main

import tg "github.com/toby3d/telegram"

// message function check Message update on commands, sended stickers or other
// user stuff
func messages(msg *tg.Message) {
	state, err := dbGetUserState(msg.From.ID)
	errCheck(err)

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
	case stateNone:
		actionError(msg)
	default:
		err = dbChangeUserState(msg.From.ID, stateNone)
		errCheck(err)

		messages(msg)
	}
}
