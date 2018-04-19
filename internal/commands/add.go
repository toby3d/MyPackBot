package commands

import (
	log "github.com/kirillDanshin/dlog"
	"github.com/toby3d/MyPackBot/internal/bot"
	"github.com/toby3d/MyPackBot/internal/db"
	"github.com/toby3d/MyPackBot/internal/errors"
	"github.com/toby3d/MyPackBot/internal/helpers"
	"github.com/toby3d/MyPackBot/internal/i18n"
	"github.com/toby3d/MyPackBot/internal/models"
	tg "github.com/toby3d/telegram"
)

// Add command prepare user for adding some stickers or sets to his pack
func Add(msg *tg.Message, pack bool) {
	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, T("reply_add_sticker"))
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = helpers.CancelButton(T)

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
