package main

import (
	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	tg "github.com/toby3d/telegram"     // My Telegram bindings
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

	switch state {
	case stateNone:
		bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)

		T, err := switchLocale(msg.From.LanguageCode)
		errCheck(err)

		reply := tg.NewMessage(
			msg.Chat.ID,
			T("error_unknown", map[string]interface{}{
				"AddStickerCommand":    cmdAddSticker,
				"AddPackCommand":       cmdAddPack,
				"DeleteStickerCommand": cmdDeleteSticker,
				"DeletePackCommand":    cmdDeletePack,
			}))
		reply.ParseMode = tg.ModeMarkdown

		_, err = bot.SendMessage(reply)
		errCheck(err)
		return
	case stateAddSticker:
		if msg.Sticker == nil {
			return
		}

		log.D(msg.Sticker)
		log.D(msg.Sticker.Emoji)

		actionAdd(msg, false)
		return
	case stateAddPack:
		if msg.Sticker == nil {
			return
		}

		log.D(msg.Sticker)
		log.D(msg.Sticker.Emoji)

		actionAdd(msg, true)
		return
	case stateDeleteSticker:
		if msg.Sticker == nil {
			return
		}

		log.D(msg.Sticker)
		log.D(msg.Sticker.Emoji)

		actionDelete(msg, false)
		return
	case stateDeletePack:
		if msg.Sticker == nil {
			return
		}

		log.D(msg.Sticker)
		log.D(msg.Sticker.Emoji)

		actionDelete(msg, true)
		return
	case stateReset:
		actionReset(msg)
		return
	}

	err = dbChangeUserState(msg.From.ID, stateNone)
	errCheck(err)

	messages(msg)
}
