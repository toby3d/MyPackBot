package main

import (
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

const (
	cmdAddPack       = "addPack"
	cmdAddSticker    = "addSticker"
	cmdCancel        = "cancel"
	cmdHelp          = "help"
	cmdDeleteSticker = "delSticker"
	cmdDeletePack    = "delPack"
	cmdReset         = "reset"
	cmdStart         = "start"
)

func commands(msg *tg.Message) {
	log.Ln("command:", msg.Command())
	switch {
	case msg.IsCommand(cmdStart):
		commandStart(msg)
	case msg.IsCommand(cmdHelp):
		commandHelp(msg)
	case msg.IsCommand(cmdAddSticker):
		commandAdd(msg, false)
	case msg.IsCommand(cmdAddPack):
		commandAdd(msg, true)
	case msg.IsCommand(cmdDeleteSticker):
		commandDelete(msg, false)
	case msg.IsCommand(cmdDeletePack):
		commandDelete(msg, true)
	case msg.IsCommand(cmdReset):
		commandReset(msg)
	case msg.IsCommand(cmdCancel):
		commandCancel(msg)
	}
}
