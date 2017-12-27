package main

import (
	"strings"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	tg "github.com/toby3d/telegram"     // My Telegram bindings
)

const (
	cmdAddPack    = "addPack"
	cmdAddSticker = "addSticker"
	cmdCancel     = "cancel"
	cmdHelp       = "help"
	cmdDelete     = "del"
	cmdReset      = "reset"
	cmdStart      = "start"
)

func commands(msg *telegram.Message) {
	log.Ln("Received a", msg.Command(), "command")
	switch strings.ToLower(msg.Command()) {
	case strings.ToLower(cmdStart):
		commandStart(msg)
	case strings.ToLower(cmdHelp):
		commandHelp(msg)
	case strings.ToLower(cmdAddSticker):
		commandAdd(msg, false)
	case strings.ToLower(cmdAddPack):
		commandAdd(msg, true)
	case strings.ToLower(cmdDelete):
		commandDelete(msg)
	case strings.ToLower(cmdReset):
		commandReset(msg)
	case strings.ToLower(cmdCancel):
		commandCancel(msg)
	}
}
