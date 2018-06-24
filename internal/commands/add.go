package commands

import (
	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/i18n"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Add command prepare user for adding some stickers or sets to his pack
func Add(msg *tg.Message, pack bool) {
	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, T("reply_add_sticker"))
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.CancelButton(T)

	err = db.DB.ChangeUserState(msg.From, models.StateAddSticker)
	errors.Check(err)

	if pack {
		reply.Text = T("reply_add_pack")

		err = db.DB.ChangeUserState(msg.From, models.StateAddPack)
		errors.Check(err)
	}

	log.Ln("Sending add reply...")
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
