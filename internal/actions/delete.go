package actions

import (
	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Delete action remove sticker or set from user's pack
func Delete(msg *tg.Message, pack bool) {
	if !msg.IsSticker() {
		return
	}

	p := utils.NewPrinter(msg.From.LanguageCode)

	_, err := bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("The sticker has been successfully removed from your set!"))
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.CancelButton(p)

	var notExist bool
	if pack {
		var set *tg.StickerSet
		set, err = bot.Bot.GetStickerSet(msg.Sticker.SetName)
		errors.Check(err)

		log.Ln("SetName:", set.Title)
		reply.Text = p.Sprintf(
			"The set *%s* was successfully removed from your collection!",
			set.Title,
		)

		notExist, err = db.DB.DeletePack(msg.From.ID, msg.Sticker)
		if notExist {
			reply.Text = p.Sprintf("Probably this set is already removed from yours.")
		}
	} else {
		notExist, err = db.DB.DeleteSticker(msg.From.ID, msg.Sticker)
		if notExist {
			reply.Text = p.Sprintf("Probably this sticker is already removed from your set.")
		}
	}
	errors.Check(err)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
