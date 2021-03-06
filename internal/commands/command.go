package commands

import (
	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/models"
	tg "gitlab.com/toby3d/telegram"
)

// Command check's got user command
func Command(msg *tg.Message) {
	log.Ln("command:", msg.Command())
	switch {
	case msg.IsCommandEqual(tg.CommandStart):
		Start(msg)
	case msg.IsCommandEqual(tg.CommandHelp):
		Help(msg)
	case msg.IsCommandEqual(models.CommandAddSticker):
		Add(msg, false)
	case msg.IsCommandEqual(models.CommandAddPack):
		Add(msg, true)
	case msg.IsCommandEqual(models.CommandDeleteSticker):
		Delete(msg, false)
	case msg.IsCommandEqual(models.CommandDeletePack):
		Delete(msg, true)
	case msg.IsCommandEqual(models.CommandReset):
		Reset(msg)
	case msg.IsCommandEqual(models.CommandCancel):
		Cancel(msg)
	}
}
