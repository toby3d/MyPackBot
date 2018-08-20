package commands

import (
	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Add command prepare user for adding some stickers or sets to his pack
func Add(msg *tg.Message, pack bool) {
	p := utils.NewPrinter(msg.From.LanguageCode)

	_, err := bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("Send stickers from any other sets to add them one by one."))
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.CancelButton(p)

	err = db.DB.ChangeUserState(msg.From.ID, models.StateAddSticker)
	errors.Check(err)

	if pack {
		reply.Text = p.Sprintf("Send stickers from any other sets to completely add their sets to your.")

		err = db.DB.ChangeUserState(msg.From.ID, models.StateAddPack)
		errors.Check(err)
	}

	log.Ln("Sending add reply...")
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
