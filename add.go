package main

import (
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func commandAdd(msg *tg.Message, pack bool) {
	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	_, err = bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	reply := tg.NewMessage(msg.Chat.ID, T("reply_add_sticker"))
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = getCancelButton(T)

	err = dbChangeUserState(msg.From.ID, stateAddSticker)
	errCheck(err)

	if pack {
		reply.Text = T("reply_add_pack")

		err = dbChangeUserState(msg.From.ID, stateAddPack)
		errCheck(err)
	}

	log.Ln("Sending add reply...")
	_, err = bot.SendMessage(reply)
	errCheck(err)
}

func actionAdd(msg *tg.Message, pack bool) {
	if !msg.IsSticker() {
		return
	}

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	_, err = bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	reply := tg.NewMessage(msg.Chat.ID, T("success_add_sticker"))
	reply.ParseMode = tg.ModeMarkdown

	if !pack {
		var exist bool
		sticker := msg.Sticker
		exist, err = dbAddSticker(
			msg.From.ID, sticker.SetName, sticker.FileID, sticker.Emoji,
		)
		errCheck(err)

		if exist {
			reply.Text = T("error_already_add_sticker")
		}

		reply.ReplyMarkup = getCancelButton(T)
		_, err = bot.SendMessage(reply)
		errCheck(err)
		return
	}

	reply.Text = T("error_empty_add_pack", map[string]interface{}{
		"AddStickerCommand": cmdAddSticker,
	})

	if msg.Sticker.SetName != "" {
		var set *tg.StickerSet
		set, err = bot.GetStickerSet(msg.Sticker.SetName)
		errCheck(err)

		log.Ln("SetTitle:", set.Title)
		reply.Text = T("success_add_pack", map[string]interface{}{
			"SetTitle": set.Title,
		})

		allExists := true
		for _, sticker := range set.Stickers {
			var exist bool
			exist, err = dbAddSticker(
				msg.From.ID, sticker.SetName, sticker.FileID, sticker.Emoji,
			)
			errCheck(err)

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

	reply.ReplyMarkup = getCancelButton(T)
	_, err = bot.SendMessage(reply)
	errCheck(err)
}
