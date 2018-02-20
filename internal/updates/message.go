package updates

import (
	log "github.com/kirillDanshin/dlog"
	"github.com/toby3d/MyPackBot/internal/actions"
	"github.com/toby3d/MyPackBot/internal/bot"
	"github.com/toby3d/MyPackBot/internal/commands"
	"github.com/toby3d/MyPackBot/internal/messages"
	tg "github.com/toby3d/telegram"
)

// Message checks Message updates for answer to user commands, replies or sended
// stickers
func Message(msg *tg.Message) {
	if bot.Bot.IsMessageFromMe(msg) ||
		bot.Bot.IsForwardFromMe(msg) {
		log.Ln("Ignore message update")
		return
	}

	switch {
	case bot.Bot.IsCommandToMe(msg):
		commands.Command(msg)
	case msg.IsText():
		messages.Message(msg)
	default:
		actions.Action(msg)
	}
}
