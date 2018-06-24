package updates

import (
	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/actions"
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/commands"
	"gitlab.com/toby3d/mypackbot/internal/messages"
	tg "gitlab.com/toby3d/telegram"
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
