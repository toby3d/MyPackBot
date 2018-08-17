package actions

import (
	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/i18n"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Delete action remove sticker or set from user's pack
func Delete(msg *tg.Message, pack bool) {
	if !msg.IsSticker() {
		return
	}

	t, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, t("success_del_sticker"))
	reply.ParseMode = tg.StyleMarkdown
	reply.ReplyMarkup = utils.CancelButton(t)

	var notExist bool
	if pack {
		var set *tg.StickerSet
		set, err = bot.Bot.GetStickerSet(msg.Sticker.SetName)
		errors.Check(err)

		log.Ln("SetName:", set.Title)
		reply.Text = t("success_del_pack", map[string]interface{}{
			"SetTitle": set.Title,
		})

		notExist, err = db.DB.DeletePack(msg.From.ID, msg.Sticker)
		if notExist {
			reply.Text = t("error_already_del_pack")
		}
	} else {
		notExist, err = db.DB.DeleteSticker(msg.From.ID, msg.Sticker)
		if notExist {
			reply.Text = t("error_already_del_sticker")
		}
	}
	errors.Check(err)

	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
