package main

import (
	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/botan"
	"github.com/toby3d/go-telegram" // My Telegram bindings
)

// message function check Message update on commands, sended stickers or other
// user stuff
func messages(msg *telegram.Message) {
	if msg.IsCommand() {
		commands(msg)
		return
	}

	state, err := dbGetUserState(msg.From.ID)
	errCheck(err)

	switch state {
	case stateNone:
		appMetrika.TrackAsync(
			"Nothing", msg.From.ID, *msg,
			func(answer *botan.Answer, err error) {
				log.Ln("Asynchonous:", answer.Status)
				metrika <- true
			},
		)

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

		<-metrika
	case stateAddSticker:
		if msg.Sticker == nil {
			appMetrika.Track("Message", msg.From.ID, *msg)
			return
		}

		log.D(msg.Sticker)
		log.D(msg.Sticker.Emoji)

		appMetrika.TrackAsync(
			"Add single sticker", msg.From.ID, *msg,
			func(answer *botan.Answer, err error) {
				log.Ln("Asynchonous:", answer.Status)
				metrika <- true
			},
		)

		actionAdd(msg, false)

		<-metrika
	case stateAddPack:
		if msg.Sticker == nil {
			appMetrika.Track("Message", msg.From.ID, *msg)
			return
		}

		appMetrika.TrackAsync(
			"Add pack", msg.From.ID, *msg,
			func(answer *botan.Answer, err error) {
				log.Ln("Asynchonous:", answer.Status)
				metrika <- true
			},
		)

		log.D(msg.Sticker)
		log.D(msg.Sticker.Emoji)

		actionAdd(msg, true)

		<-metrika
	case stateDelete:
		if msg.Sticker == nil {
			appMetrika.Track("Message", msg.From.ID, *msg)
			return
		}

		appMetrika.TrackAsync(
			"Delete sticker", msg.From.ID, *msg,
			func(answer *botan.Answer, err error) {
				log.Ln("Asynchonous:", answer.Status)
				metrika <- true
			},
		)

		log.D(msg.Sticker)
		log.D(msg.Sticker.Emoji)

		actionDelete(msg)

		<-metrika
	case stateReset:
		appMetrika.TrackAsync(
			"Reset pack", msg.From.ID, *msg,
			func(answer *botan.Answer, err error) {
				log.Ln("Asynchonous:", answer.Status)
				metrika <- true
			},
		)

		actionReset(msg)

		<-metrika
	default:
		appMetrika.TrackAsync(
			"Message", msg.From.ID, *msg,
			func(answer *botan.Answer, err error) {
				log.Ln("Asynchonous:", answer.Status)
				metrika <- true
			},
		)

		err = dbChangeUserState(msg.From.ID, stateNone)
		errCheck(err)

		messages(msg)

		<-metrika
	}
}
