package main

import (
	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	tg "github.com/toby3d/telegram"     // My Telegram bindings
)

func commandAdd(msg *tg.Message, pack bool) {
	bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)

	T, err := switchLocale(msg.From.LanguageCode)
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
	bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	reply := tg.NewMessage(msg.Chat.ID, T("success_add_sticker"))
	reply.ParseMode = tg.ModeMarkdown

	switch {
	case pack && msg.Sticker.SetName == "":
		reply.Text = T("error_empty_add_pack", map[string]interface{}{
			"AddStickerCommand": cmdAddSticker,
		})
	case pack && msg.Sticker.SetName != "":
		set, err := bot.GetStickerSet(msg.Sticker.SetName)
		errCheck(err)

		log.Ln("SetTitle:", set.Title)
		reply.Text = T("success_add_pack", map[string]interface{}{
			"SetTitle": set.Title,
		})

		allExists := true
		for _, sticker := range set.Stickers {
			exists, err := dbAddSticker(
				msg.From.ID,
				sticker.SetName,
				sticker.FileID,
				sticker.Emoji,
			)
			errCheck(err)

			if !exists {
				allExists = false
			}
		}

		log.Ln("All exists?", allExists)

		if allExists {
			reply.Text = T("error_already_add_pack", map[string]interface{}{
				"SetTitle": set.Title,
			})
		} else {
			reply.ReplyMarkup = getCancelButton(T)
		}
	default:
		exists, err := dbAddSticker(
			msg.From.ID,
			msg.Sticker.SetName,
			msg.Sticker.FileID,
			msg.Sticker.Emoji,
		)
		errCheck(err)

		if exists {
			reply.Text = T("error_already_add_sticker")
		}

		if msg.Sticker.Emoji == "" {
			msg.Sticker.Emoji = " "
		}

		reply.ReplyMarkup = getCancelButton(T)
	}

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
