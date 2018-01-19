package main

import (
	"strings"

	// log "github.com/kirillDanshin/dlog"
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
	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	switch {
	case strings.ToLower(msg.Command()) == strings.ToLower(cmdStart):
		commandStart(msg)
	case strings.ToLower(msg.Command()) == strings.ToLower(cmdHelp):
		commandHelp(msg)
	case msg.Text == T("button_add_sticker"),
		strings.ToLower(msg.Command()) == strings.ToLower(cmdAddSticker):
		commandAdd(msg, false)
	case msg.Text == T("button_add_pack"),
		strings.ToLower(msg.Command()) == strings.ToLower(cmdAddPack):
		commandAdd(msg, true)
	case msg.Text == T("button_del_sticker"),
		strings.ToLower(msg.Command()) == strings.ToLower(cmdDeleteSticker):
		commandDelete(msg, false)
	case msg.Text == T("button_del_pack"),
		strings.ToLower(msg.Command()) == strings.ToLower(cmdDeletePack):
		commandDelete(msg, true)
	case msg.Text == T("button_reset"),
		strings.ToLower(msg.Command()) == strings.ToLower(cmdReset):
		commandReset(msg)
	case msg.Text == T("button_cancel"),
		strings.ToLower(msg.Command()) == strings.ToLower(cmdCancel):
		commandCancel(msg)
	}
}
