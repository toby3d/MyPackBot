package actions

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

// Add action add sticker or set to user's pack
func Add(msg *tg.Message, pack bool) {
	if !msg.IsSticker() {
		return
	}

	t, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, t("success_add_sticker"))
	reply.ParseMode = tg.StyleMarkdown

	if !pack {
		var exist bool
		sticker := msg.Sticker
		exist, err = db.DB.AddSticker(msg.From, sticker)
		errors.Check(err)

		if exist {
			reply.Text = t("error_already_add_sticker")
		}

		reply.ReplyMarkup = utils.CancelButton(t)
		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
		return
	}

	reply.Text = t("error_empty_add_pack", map[string]interface{}{
		"AddStickerCommand": models.CommandAddSticker,
	})

	if msg.Sticker.SetName != "" {
		var set *tg.StickerSet
		set, err = bot.Bot.GetStickerSet(msg.Sticker.SetName)
		errors.Check(err)

		log.Ln("SetTitle:", set.Title)
		reply.Text = t("success_add_pack", map[string]interface{}{
			"SetTitle": set.Title,
		})

		allExists := true
		for i := range set.Stickers {
			var exist bool
			exist, err = db.DB.AddSticker(msg.From, &set.Stickers[i])
			errors.Check(err)

			if !exist {
				allExists = false
			}
		}

		log.Ln("All exists?", allExists)
		if allExists {
			reply.Text = t("error_already_add_pack", map[string]interface{}{
				"SetTitle": set.Title,
			})
		}
	}

	reply.ReplyMarkup = utils.CancelButton(t)
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
