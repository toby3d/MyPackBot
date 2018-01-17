package main

import (
	// log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	tg "github.com/toby3d/telegram" // My Telegram bindings
)

// message function check Message update on commands, sended stickers or other
// user stuff
func messages(msg *tg.Message) {
	if msg.IsCommand() {
		commands(msg)
		return
	}

	state, err := dbGetUserState(msg.From.ID)
	errCheck(err)

	T, err := switchLocale(msg.From.LanguageCode)
	errCheck(err)

	switch msg.Text {
	case T("button_add_sticker"):
		commandAdd(msg, false)
		return
	case T("button_add_pack"):
		commandAdd(msg, true)
		return
	case T("button_del_sticker"):
		commandDelete(msg, false)
		return
	case T("button_del_pack"):
		commandDelete(msg, true)
		return
	case T("button_reset"):
		commandReset(msg)
		return
	case T("button_cancel"):
		commandCancel(msg)
		return
	}

	switch state {
	case stateAddSticker:
		if msg.Sticker == nil {
			return
		}

		actionAdd(msg, false)
		return
	case stateAddPack:
		if msg.Sticker == nil {
			return
		}

		actionAdd(msg, true)
		return
	case stateDeleteSticker:
		if msg.Sticker == nil {
			return
		}

		actionDelete(msg, false)
		return
	case stateDeletePack:
		if msg.Sticker == nil {
			return
		}

		actionDelete(msg, true)
		return
	case stateReset:
		actionReset(msg)
		return
	case stateNone:
		bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)

		reply := tg.NewMessage(
			msg.Chat.ID,
			T("error_unknown", map[string]interface{}{
				"AddStickerCommand":    cmdAddSticker,
				"AddPackCommand":       cmdAddPack,
				"DeleteStickerCommand": cmdDeleteSticker,
				"DeletePackCommand":    cmdDeletePack,
			}))
		reply.ParseMode = tg.ModeMarkdown
		reply.ReplyMarkup = getMenuKeyboard(T)

		_, err = bot.SendMessage(reply)
		errCheck(err)
		return
	}

	err = dbChangeUserState(msg.From.ID, stateNone)
	errCheck(err)

	messages(msg)
}
