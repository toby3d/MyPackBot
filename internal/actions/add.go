package actions

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

// Add action add sticker or set to user's pack
func Add(msg *tg.Message, pack bool) {
	if !msg.IsSticker() {
		return
	}

	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	_, err = bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, T("success_add_sticker"))
	reply.ParseMode = tg.ModeMarkdown

	if !pack {
		var exist bool
		sticker := msg.Sticker
		exist, err = db.DB.AddSticker(msg.From, sticker)
		errors.Check(err)

		if exist {
			reply.Text = T("error_already_add_sticker")
		}

		reply.ReplyMarkup = helpers.CancelButton(T)
		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
		return
	}

	reply.Text = T("error_empty_add_pack", map[string]interface{}{
		"AddStickerCommand": models.CommandAddSticker,
	})

	if msg.Sticker.SetName != "" {
		var set *tg.StickerSet
		set, err = bot.Bot.GetStickerSet(msg.Sticker.SetName)
		errors.Check(err)

		log.Ln("SetTitle:", set.Title)
		reply.Text = T("success_add_pack", map[string]interface{}{
			"SetTitle": set.Title,
		})

		allExists := true
		for _, sticker := range set.Stickers {
			var exist bool
			exist, err = db.DB.AddSticker(msg.From, &sticker)
			errors.Check(err)

			if !exist {
				allExists = false
			}
		}

		log.Ln("All exists?", allExists)
		if allExists {
			reply.Text = T("error_already_add_pack", map[string]interface{}{
				"SetTitle": set.Title,
			})
		}
	}

	reply.ReplyMarkup = helpers.CancelButton(T)
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
