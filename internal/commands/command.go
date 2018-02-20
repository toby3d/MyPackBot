package commands

import (
	log "github.com/kirillDanshin/dlog"
	"github.com/toby3d/MyPackBot/internal/models"
	tg "github.com/toby3d/telegram"
)

// Command check's got user command
func Command(msg *tg.Message) {
	log.Ln("command:", msg.Command())
	switch {
	case msg.IsCommand(models.CommandStart):
		Start(msg)
	case msg.IsCommand(models.CommandHelp):
		Help(msg)
	case msg.IsCommand(models.CommandAddSticker):
		Add(msg, false)
	case msg.IsCommand(models.CommandAddPack):
		Add(msg, true)
	case msg.IsCommand(models.CommandDeleteSticker):
		Delete(msg, false)
	case msg.IsCommand(models.CommandDeletePack):
		Delete(msg, true)
	case msg.IsCommand(models.CommandReset):
		Reset(msg)
	case msg.IsCommand(models.CommandCancel):
		Cancel(msg)
	}
}
