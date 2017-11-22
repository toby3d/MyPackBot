package main

import (
	log "github.com/kirillDanshin/dlog"  // Insert logs only in debug builds
	"github.com/nicksnyder/go-i18n/i18n" // Internationalization and localization
	"github.com/toby3d/go-telegram"      // My Telegram bindings
)

func commandAdd(msg *telegram.Message, pack bool) {
	log.Ln("Received a /add command")
	log.Ln("Check", msg.From.LanguageCode, "localization")
	T, err := i18n.Tfunc(msg.From.LanguageCode)
	if err != nil {
		log.Ln("Unsupported language, change to 'en-us' by default")
		T, err = i18n.Tfunc(langDefault)
		errCheck(err)
	}

	reply := telegram.NewMessage(msg.Chat.ID, T("reply_add_sticker"))
	reply.ParseMode = telegram.ModeMarkdown

	log.Ln("Change", msg.From.ID, "state to", stateAddSticker)
	err = dbChangeUserState(msg.From.ID, stateAddSticker)
	errCheck(err)

	if pack {
		reply.Text = T("reply_add_pack")

		log.Ln("Change", msg.From.ID, "state to", stateAddPack)
		err = dbChangeUserState(msg.From.ID, stateAddPack)
		errCheck(err)
	}

	log.Ln("Sending add reply...")
	_, err = bot.SendMessage(reply)
	errCheck(err)
}

func actionAdd(msg *telegram.Message, pack bool) {
	log.Ln("Received a /add action")
	bot.SendChatAction(msg.Chat.ID, telegram.ActionTyping)

	log.Ln("Check", msg.From.LanguageCode, "localization")
	T, err := i18n.Tfunc(msg.From.LanguageCode)
	if err != nil {
		log.Ln("Unsupported language, change to 'en-us' by default")
		T, err = i18n.Tfunc(langDefault)
		errCheck(err)
	}

	reply := telegram.NewMessage(msg.Chat.ID, T("success_add_sticker"))
	reply.ParseMode = telegram.ModeMarkdown

	switch {
	case pack &&
		msg.Sticker.SetName == "":
		reply.Text = T("error_empty_add_pack")
	case pack &&
		msg.Sticker.SetName != "":

		set, err := bot.GetStickerSet(msg.Sticker.SetName)
		errCheck(err)

		log.Ln("SetTitle:", set.Title)

		reply.Text = T("success_add_pack", map[string]interface{}{
			"SetTitle": set.Title,
		})

		allExists := true
		for _, sticker := range set.Stickers {
			exists, err := dbAddSticker(msg.From.ID, sticker.FileID, sticker.Emoji)
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
		}
	default:
		exists, err := dbAddSticker(msg.From.ID, msg.Sticker.FileID, msg.Sticker.Emoji)
		errCheck(err)

		if exists {
			reply.Text = T("error_already_add_sticker")
		}

		markup := telegram.NewInlineKeyboardMarkup(
			telegram.NewInlineKeyboardRow(
				telegram.NewInlineKeyboardButtonSwitch(
					T("button_share"),
					msg.Sticker.Emoji,
				),
			),
		)
		reply.ReplyMarkup = &markup
	}

	_, err = bot.SendMessage(reply)
	errCheck(err)
}