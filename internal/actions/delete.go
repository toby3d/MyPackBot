package actions

import (
	log "github.com/kirillDanshin/dlog"
	"github.com/toby3d/MyPackBot/internal/bot"
	"github.com/toby3d/MyPackBot/internal/db"
	"github.com/toby3d/MyPackBot/internal/errors"
	"github.com/toby3d/MyPackBot/internal/helpers"
	"github.com/toby3d/MyPackBot/internal/i18n"
	tg "github.com/toby3d/telegram"
)

// Delete action remove sticker or set from user's pack
func Delete(msg *tg.Message, pack bool) {
	if !msg.IsSticker() {
		return
	}

	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, T("success_del_sticker"))
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = helpers.CancelButton(T)

	var notExist bool
	if pack {
		var set *tg.StickerSet
		set, err = bot.Bot.GetStickerSet(msg.Sticker.SetName)
		errors.Check(err)

		log.Ln("SetName:", set.Title)
		reply.Text = T("success_del_pack", map[string]interface{}{
			"SetTitle": set.Title,
		})

		notExist, err = db.DeletePack(msg.From.ID, msg.Sticker.SetName)
		if notExist {
			reply.Text = T("error_already_del_pack")
		}
	} else {
		notExist, err = db.DeleteSticker(
			msg.From.ID,
			msg.Sticker.SetName,
			msg.Sticker.FileID,
		)
		if notExist {
			reply.Text = T("error_already_del_sticker")
		}
	}
	errors.Check(err)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
