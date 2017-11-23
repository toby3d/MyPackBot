package main

import (
	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/go-telegram"     // My Telegram bindings
)

// message function check Message update on commands, sended stickers or other
// user stuff
func messages(msg *telegram.Message) {
	if msg.IsCommand() {
		log.Ln("Received a command message")
		commands(msg)
		return
	}

	state, err := dbGetUserState(msg.From.ID)
	errCheck(err)

	switch state {
	case stateNone:
		bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

		T, err := switchLocale(msg.From.LanguageCode)
		errCheck(err)

		reply := telegram.NewMessage(
			msg.Chat.ID,
			T("error_unknown", map[string]interface{}{
				"AddStickerCommand": cmdAddSticker,
				"AddPackCommand":    cmdAddPack,
				"DeleteCommand":     cmdDelete,
			}))
		reply.ParseMode = telegram.ModeMarkdown

		_, err = bot.SendMessage(reply)
		errCheck(err)
	case stateAddSticker:
		if msg.Sticker != nil {
			log.D(msg.Sticker)
			log.D(msg.Sticker.Emoji)

			actionAdd(msg, false)
		}
	case stateAddPack:
		if msg.Sticker != nil {
			log.D(msg.Sticker)
			log.D(msg.Sticker.Emoji)

			actionAdd(msg, true)
		}
	case stateDelete:
		if msg.Sticker != nil {
			log.D(msg.Sticker)
			log.D(msg.Sticker.Emoji)

			actionDelete(msg)
		}
	case stateReset:
		actionReset(msg)
	default:
		err = dbChangeUserState(msg.From.ID, stateNone)
		errCheck(err)

		messages(msg)
	}
}
