package main

import (
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func commandDelete(msg *tg.Message, pack bool) {
	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	_, total, err := dbGetUserStickers(msg.From.ID, 0, "")
	errCheck(err)

	_, err = bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	if total <= 0 {
		err = dbChangeUserState(msg.From.ID, stateNone)
		errCheck(err)

		reply := tg.NewMessage(msg.Chat.ID, T("error_empty_del"))
		reply.ReplyMarkup = getMenuKeyboard(T)
		_, err = bot.SendMessage(reply)
		errCheck(err)
		return
	}

	reply := tg.NewMessage(msg.Chat.ID, T("reply_del_sticker"))
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = getCancelButton(T)

	err = dbChangeUserState(msg.From.ID, stateDeleteSticker)
	errCheck(err)

	if pack {
		err = dbChangeUserState(msg.From.ID, stateDeletePack)
		errCheck(err)

		reply.Text = T("reply_del_pack")
	}

	_, err = bot.SendMessage(reply)
	errCheck(err)

	_, err = bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	reply = tg.NewMessage(msg.Chat.ID, T("reply_switch_button"))
	reply.ReplyMarkup = getSwitchButton(T)
	_, err = bot.SendMessage(reply)
	errCheck(err)
}

func actionDelete(msg *tg.Message, pack bool) {
	if msg.Sticker == nil {
		return
	}

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	_, err = bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errCheck(err)

	reply := tg.NewMessage(msg.Chat.ID, T("success_del_sticker"))
	reply.ParseMode = tg.ModeMarkdown
	reply.ReplyMarkup = getCancelButton(T)

	var notExist bool
	if pack {
		var set *tg.StickerSet
		set, err = bot.GetStickerSet(msg.Sticker.SetName)
		errCheck(err)

		log.Ln("SetName:", set.Title)
		reply.Text = T("success_del_pack", map[string]interface{}{
			"SetTitle": set.Title,
		})

		notExist, err = dbDeletePack(msg.From.ID, msg.Sticker.SetName)
	} else {
		notExist, err = dbDeleteSticker(msg.From.ID, msg.Sticker.SetName, msg.Sticker.FileID)
	}
	errCheck(err)

	if notExist {
		reply.Text = T("error_already_del")
	}

	_, err = bot.SendMessage(reply)
	errCheck(err)
}
